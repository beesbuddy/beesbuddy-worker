package worker

import (
	"fmt"
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/nakabonne/tstorage"
	"github.com/samber/lo"
)

type workerCtx struct {
	appCtx  *app.Ctx
	storage tstorage.Storage
	topics  []string
}

func NewWorkersRunner(appCtx *app.Ctx) internal.ModuleCtx {
	storage, err := tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithDataPath(appCtx.Config.GetCfg().StoragePath),
	)

	if err != nil {
		panic("unable to create storage")
	}

	m := &workerCtx{appCtx: appCtx, storage: storage}
	NewConnection(m.appCtx.MqttClient)
	return m
}

func (m *workerCtx) Run() {
	go func(m *workerCtx) {
		defer m.CleanUp()

		for {
			log.Println("[Re]configuring MQTT:", m.appCtx.Config.GetCfg().BrokerTCPUrl)

			if !m.appCtx.MqttClient.IsConnectionOpen() || !m.appCtx.MqttClient.IsConnected() {
				NewConnection(m.appCtx.MqttClient)
			}

			m.initializeSubscribers()

			<-m.appCtx.Config.GetSubscriber(internal.WorkerKey)

			m.cleanUpSubscribers()
		}
	}(m)
}

func (m *workerCtx) CleanUp() {
	log.Println("Gracefully closing mqtt workers...")

	if m.appCtx.MqttClient.IsConnectionOpen() && m.appCtx.MqttClient.IsConnected() {
		m.cleanUpSubscribers()
		Disconnect(m.appCtx.MqttClient)
	}

	m.storage.Close()
}

func (m *workerCtx) cleanUpSubscribers() {
	for _, s := range m.appCtx.Config.GetCfg().Subscribers {
		topic := fmt.Sprintf(internal.TopicPath, s.ApiaryId, s.HiveId)
		topicToDelete, alreadyExists := lo.Find(m.topics, func(t string) bool {
			return t == topic
		})

		if alreadyExists {
			go func(topic string) {
				m.Unsubscribe(m.appCtx.MqttClient, topicToDelete)
			}(topic)
		}

	}
}

func (m *workerCtx) initializeSubscribers() {
	for _, s := range m.appCtx.Config.GetCfg().Subscribers {
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
		log.Fatal(token.Error())
		panic(token.Error())
	}
}

func (m *workerCtx) Subscribe(c MQTT.Client, topic string) {
	if token := c.Subscribe(topic, 0, m.DefaultMessageHandler()); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
		panic(token.Error())
	}
}
