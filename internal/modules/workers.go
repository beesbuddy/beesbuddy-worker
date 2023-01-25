package modules

import (
	"fmt"
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/mqtt"
	"github.com/samber/lo"
)

type WorkersMod struct {
	app    *core.App
	topics []string
}

func NewWorkersRunner(app *core.App) core.Mod {
	mod := &WorkersMod{app: app}
	return mod
}

func (mod *WorkersMod) Run() {
	mqtt.NewConnection(mod.app.MqttClient)

	for {
		log.Println("[Re]configuring MQTT:", core.GetCfg().BrokerTCPUrl)

		if !mod.app.MqttClient.IsConnectionOpen() || !mod.app.MqttClient.IsConnected() {
			mqtt.NewConnection(mod.app.MqttClient)
		}

		mod.initializeSubscribers()
		<-core.GetCfgObject().GetSubscriber(core.WorkerKey)
		mod.cleanUpSubscribers()
	}
}

func (mod *WorkersMod) CleanUp() {
	log.Println("Gracefully closing mqtt workers...")
	if mod.app.MqttClient.IsConnectionOpen() && mod.app.MqttClient.IsConnected() {
		mod.cleanUpSubscribers()
		mqtt.Disconnect(mod.app.MqttClient)
	}
}

func (mod *WorkersMod) cleanUpSubscribers() {
	for _, s := range core.GetCfg().Subscribers {
		topic := fmt.Sprintf("apiary/%s/hive/%s", s.ApiaryId, s.HiveId)
		topicToDelete, alreadyExists := lo.Find(mod.topics, func(t string) bool {
			return t == topic
		})

		if alreadyExists {
			go func(topic string) {
				mqtt.Unsubscribe(mod.app.MqttClient, topicToDelete)
			}(topic)
		}

	}
}

func (mod *WorkersMod) initializeSubscribers() {
	for _, s := range core.GetCfg().Subscribers {
		topic := fmt.Sprintf("apiary/%s/hive/%s", s.ApiaryId, s.HiveId)
		_, alreadyExists := lo.Find(mod.topics, func(t string) bool {
			return t == topic
		})

		if !alreadyExists {
			go func(topic string) {
				mqtt.Subscribe(mod.app.MqttClient, topic)
			}(topic)
		}
	}

}
