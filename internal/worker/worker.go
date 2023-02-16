package worker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/log"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/nakabonne/tstorage"
	"github.com/samber/lo"
)

type workerCtx struct {
	appCtx  *app.Ctx
	storage tstorage.Storage
	topics  []string
	queue   chan metrics
}

func NewWorkersRunner(appCtx *app.Ctx) internal.Ctx {
	config := appCtx.Pref.GetConfig()
	duration := config.PartitionDuration

	storage, err := tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithPartitionDuration(time.Duration(duration)*time.Hour),
		tstorage.WithDataPath(appCtx.Pref.GetConfig().StoragePath),
	)

	if err != nil {
		panic("unable to create storage")
	}

	queue := make(chan metrics, internal.WorkerChanBuffer)

	m := &workerCtx{appCtx: appCtx, storage: storage, queue: queue}
	NewConnection(m.appCtx.MqttClient)
	return m
}

func (w *workerCtx) Run() {
	go func(w *workerCtx) {
		for {
			cfg := w.appCtx.Pref.GetConfig()
			client := w.appCtx.MqttClient
			pref := w.appCtx.Pref

			log.Info.Println("[re]configuring MQTT:", cfg.BrokerTCPUrl)

			if !w.appCtx.MqttClient.IsConnectionOpen() || !client.IsConnected() {
				NewConnection(w.appCtx.MqttClient)
			}

			w.initializeSubscribers()

			<-pref.GetSubscriber(internal.WorkerKey)

			w.cleanUpSubscribers()
		}
	}(w)

	for i := 0; i < w.appCtx.Pref.GetConfig().StorageWorkersCount; i++ {
		go w.storageWorker()
	}
}

func (w *workerCtx) CleanUp() {
	log.Info.Println("gracefully closing mqtt workers...")

	if w.appCtx.MqttClient.IsConnectionOpen() && w.appCtx.MqttClient.IsConnected() {
		w.cleanUpSubscribers()
		Disconnect(w.appCtx.MqttClient)
	}

	w.storage.Close()
}

func (w *workerCtx) cleanUpSubscribers() {
	for _, s := range w.appCtx.Pref.GetConfig().Subscribers {
		topic := fmt.Sprintf(internal.TopicPath, s.ApiaryId, s.HiveId)
		topicToDelete, alreadyExists := lo.Find(w.topics, func(t string) bool {
			return t == topic
		})

		if alreadyExists {
			go func(topic string) {
				w.Unsubscribe(w.appCtx.MqttClient, topicToDelete)
			}(topic)
		}

	}
}

func (m *workerCtx) initializeSubscribers() {
	for _, s := range m.appCtx.Pref.GetConfig().Subscribers {
		topic := fmt.Sprintf(internal.TopicPath, s.ApiaryId, s.HiveId)
		_, alreadyExists := lo.Find(m.topics, func(t string) bool {
			return t == topic
		})

		if !alreadyExists {
			go func(topic string) {
				m.Subscribe(m.appCtx.MqttClient, topic)
			}(topic)
		}
	}

}

func (m *workerCtx) Unsubscribe(c MQTT.Client, topic string) {
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Error.Fatal(token.Error())
		panic(token.Error())
	}
}

func (w *workerCtx) Subscribe(c MQTT.Client, topic string) {
	if token := c.Subscribe(topic, 0, w.PersistMessageHandler()); token.Wait() && token.Error() != nil {
		log.Error.Fatal(token.Error())
		panic(token.Error())
	}
}

func (w *workerCtx) persist(m metrics) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(internal.WorkerTimeout)*time.Second)
	defer cancel()

	labels := []tstorage.Label{
		{Name: "ApriaryId", Value: m.ApiaryId},
		{Name: "HiveId", Value: m.HiveId},
	}

	w.storeMetric("temperature", m.Temperature, labels)
	w.storeMetric("humidity", m.Humidity, labels)
	w.storeMetric("weight", m.Weight, labels)

	ctx.Done()
	return nil
}

func (w *workerCtx) storageWorker() {
	for msq := range w.queue {
		log.Debug.Println("try to persist metrics: ", msq)
		err := w.persist(msq)
		if err != nil {
			log.Error.Print(err)
		}
	}
}

func (w *workerCtx) storeMetric(name, value string, labels []tstorage.Label) {
	v, err := strconv.ParseFloat(value, 64)

	if err != nil {
		log.Error.Println("unable to parse temperature for ", labels)
	} else {
		err = w.storage.InsertRows([]tstorage.Row{
			{
				Metric:    name,
				Labels:    labels,
				DataPoint: tstorage.DataPoint{Timestamp: time.Now().Unix(), Value: v},
			},
		})

		if err != nil {
			log.Error.Println("unable to store temperature for ", labels)
		}
	}
}
