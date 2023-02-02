package modules

import (
	"fmt"
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/mqtt"
	"github.com/samber/lo"
)

type workersModule struct {
	ctx    *core.Ctx
	topics []string
}

func NewWorkersRunner(ctx *core.Ctx) core.Module {
	ctx.WebModuleSync.Add(1)
	m := &workersModule{ctx: ctx}
	return m
}

func (m *workersModule) Run() {
	mqtt.NewConnection(m.ctx.MqttClient)

	go func(m *workersModule) {
		defer m.CleanUp()

		for {
			log.Println("[Re]configuring MQTT:", m.ctx.Config.GetCfg().BrokerTCPUrl)

			if !m.ctx.MqttClient.IsConnectionOpen() || !m.ctx.MqttClient.IsConnected() {
				mqtt.NewConnection(m.ctx.MqttClient)
			}

			m.initializeSubscribers()

			<-m.ctx.Config.GetSubscriber(core.WorkerKey)

			m.cleanUpSubscribers()
		}
	}(m)

}

func (m *workersModule) CleanUp() {
	log.Println("Gracefully closing mqtt workers...")

	if m.ctx.MqttClient.IsConnectionOpen() && m.ctx.MqttClient.IsConnected() {
		m.cleanUpSubscribers()
		mqtt.Disconnect(m.ctx.MqttClient)
	}

	m.ctx.WebModuleSync.Done()
}

func (m *workersModule) cleanUpSubscribers() {
	for _, s := range m.ctx.Config.GetCfg().Subscribers {
		topic := fmt.Sprintf(core.TopicPath, s.ApiaryId, s.HiveId)
		topicToDelete, alreadyExists := lo.Find(m.topics, func(t string) bool {
			return t == topic
		})

		if alreadyExists {
			go func(topic string) {
				mqtt.Unsubscribe(m.ctx.MqttClient, topicToDelete)
			}(topic)
		}

	}
}

func (m *workersModule) initializeSubscribers() {
	for _, s := range m.ctx.Config.GetCfg().Subscribers {
		topic := fmt.Sprintf(core.TopicPath, s.ApiaryId, s.HiveId)
		_, alreadyExists := lo.Find(m.topics, func(t string) bool {
			return t == topic
		})

		if !alreadyExists {
			go func(topic string) {
				mqtt.Subscribe(m.ctx.MqttClient, topic)
			}(topic)
		}
	}

}
