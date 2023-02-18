package worker

import (
	"fmt"
	"time"

	"github.com/beesbuddy/beesbuddy-worker/app"
	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/component"
	"github.com/beesbuddy/beesbuddy-worker/internal/log"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nakabonne/tstorage"
	"github.com/samber/lo"
)

type workerCtx struct {
	appCtx         *app.Ctx
	storage        tstorage.Storage
	topics         []string
	queue          chan metrics
	influxDbClient influxdb2.Client
}

func NewWorkersRunner(appCtx *app.Ctx) component.Component {
	config := appCtx.Pref.GetConfig()
	duration := config.PartitionDuration

	storage, err := tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithPartitionDuration(time.Duration(duration)*time.Hour),
		tstorage.WithDataPath(appCtx.Pref.GetConfig().StoragePath),
	)

	influxDbClient := influxdb2.NewClient(config.InfluxDbURL, config.InfluxDbAccessToken)

	if err != nil {
		panic("unable to create storage")
	}

	queue := make(chan metrics, internal.WorkerChanBuffer)

	w := &workerCtx{appCtx: appCtx, storage: storage, queue: queue, influxDbClient: influxDbClient}
	NewConnection(w.appCtx.MqttClient)
	return w
}

func (w *workerCtx) Init() {
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

func (w *workerCtx) Flush() {
	log.Info.Println("gracefully closing mqtt workers...")

	if w.appCtx.MqttClient.IsConnectionOpen() && w.appCtx.MqttClient.IsConnected() {
		w.cleanUpSubscribers()
		Disconnect(w.appCtx.MqttClient)
	}

	w.storage.Close()
	w.influxDbClient.Close()
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
