package worker

import (
	"fmt"
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/samber/lo"
)

type workersModule struct {
	ctx    *app.Ctx
	topics []string
}

func NewWorkersRunner(ctx *app.Ctx) app.Module {
	m := &workersModule{ctx: ctx}
	return m
}

func (m *workersModule) Run() {
	NewConnection(m.ctx.MqttClient)

	go func(m *workersModule) {
		defer m.CleanUp()

		for {
			log.Println("[Re]configuring MQTT:", m.ctx.Config.GetCfg().BrokerTCPUrl)

			if !m.ctx.MqttClient.IsConnectionOpen() || !m.ctx.MqttClient.IsConnected() {
				NewConnection(m.ctx.MqttClient)
			}

			m.initializeSubscribers()

			<-m.ctx.Config.GetSubscriber(internal.WorkerKey)

			m.cleanUpSubscribers()
		}
	}(m)
}

func (m *workersModule) CleanUp() {
	log.Println("Gracefully closing mqtt workers...")

	if m.ctx.MqttClient.IsConnectionOpen() && m.ctx.MqttClient.IsConnected() {
		m.cleanUpSubscribers()
		Disconnect(m.ctx.MqttClient)
	}
}

func (m *workersModule) cleanUpSubscribers() {
	for _, s := range m.ctx.Config.GetCfg().Subscribers {
		topic := fmt.Sprintf(internal.TopicPath, s.ApiaryId, s.HiveId)
		topicToDelete, alreadyExists := lo.Find(m.topics, func(t string) bool {
			return t == topic
		})

		if alreadyExists {
			go func(topic string) {
				Unsubscribe(m.ctx.MqttClient, topicToDelete)
			}(topic)
		}

	}
}

func (m *workersModule) initializeSubscribers() {
	for _, s := range m.ctx.Config.GetCfg().Subscribers {
		topic := fmt.Sprintf(internal.TopicPath, s.ApiaryId, s.HiveId)
		_, alreadyExists := lo.Find(m.topics, func(t string) bool {
			return t == topic
		})

		if !alreadyExists {
			go func(topic string) {
				Subscribe(m.ctx.MqttClient, topic)
			}(topic)
		}
	}

}
