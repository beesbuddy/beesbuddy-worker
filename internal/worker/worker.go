package worker

import (
	"fmt"

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
	queue   chan int64
}

func NewWorkersRunner(appCtx *app.Ctx) internal.Ctx {
	storage, err := tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithDataPath(appCtx.Pref.GetConfig().StoragePath),
	)

	if err != nil {
		panic("unable to create storage")
	}

	queue := make(chan int64, 100)

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

			log.Info.Println("[Re]configuring MQTT:", cfg.BrokerTCPUrl)

			if !w.appCtx.MqttClient.IsConnectionOpen() || !client.IsConnected() {
				NewConnection(w.appCtx.MqttClient)
			}

			w.initializeSubscribers()

			<-pref.GetSubscriber(internal.WorkerKey)

			w.cleanUpSubscribers()
		}
	}(w)
}

func (w *workerCtx) CleanUp() {
	log.Info.Println("Gracefully closing mqtt workers...")

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

func (m *workerCtx) Subscribe(c MQTT.Client, topic string) {
	if token := c.Subscribe(topic, 0, m.DefaultMessageHandler()); token.Wait() && token.Error() != nil {
		log.Error.Fatal(token.Error())
		panic(token.Error())
	}
}
