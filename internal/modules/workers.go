package modules

import (
	"fmt"
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/mqtt"
	"github.com/samber/lo"
)

type WorkersModule struct {
	app    *core.App
	topics []string
}

func NewWorkersRunner(app *core.App) core.Module {
	m := &WorkersModule{app: app}
	return m
}

func (m *WorkersModule) Run() {
	mqtt.NewConnection(m.app.MqttClient)

	go func() {
		for {
			log.Println("[Re]configuring MQTT:", core.GetCfg().BrokerTCPUrl)

			if !m.app.MqttClient.IsConnectionOpen() || !m.app.MqttClient.IsConnected() {
				mqtt.NewConnection(m.app.MqttClient)
			}

			m.initializeSubscribers()
			<-core.GetCfgObject().GetSubscriber(core.WorkerKey)
			m.cleanUpSubscribers()
		}
	}()

}

func (m *WorkersModule) CleanUp() {
	log.Println("Gracefully closing mqtt workers...")
	if m.app.MqttClient.IsConnectionOpen() && m.app.MqttClient.IsConnected() {
		m.cleanUpSubscribers()
		mqtt.Disconnect(m.app.MqttClient)
	}
}

func (m *WorkersModule) cleanUpSubscribers() {
	for _, s := range core.GetCfg().Subscribers {
		topic := fmt.Sprintf("apiary/%s/hive/%s", s.ApiaryId, s.HiveId)
		topicToDelete, alreadyExists := lo.Find(m.topics, func(t string) bool {
			return t == topic
		})

		if alreadyExists {
			go func(topic string) {
				mqtt.Unsubscribe(m.app.MqttClient, topicToDelete)
			}(topic)
		}

	}
}

func (m *WorkersModule) initializeSubscribers() {
	for _, s := range core.GetCfg().Subscribers {
		topic := fmt.Sprintf("apiary/%s/hive/%s", s.ApiaryId, s.HiveId)
		_, alreadyExists := lo.Find(m.topics, func(t string) bool {
			return t == topic
		})

		if !alreadyExists {
			go func(topic string) {
				mqtt.Subscribe(m.app.MqttClient, topic)
			}(topic)
		}
	}

}
