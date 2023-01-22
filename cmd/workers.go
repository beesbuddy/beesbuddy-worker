package cmd

import (
	"fmt"
	"log"

	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	q "github.com/beesbuddy/beesbuddy-worker/internal/queue"
	"github.com/samber/lo"
)

type WorkersCmd struct {
	app    *c.App
	topics []string
}

func NewWorkersRunner(app *c.App) c.CmdRunner {
	cmd := &WorkersCmd{app: app}
	return cmd
}

func (cmd *WorkersCmd) Run() {
	q.NewConnection(cmd.app.MqttClient)

	for {
		log.Println("[Re]configuring MQTT:", c.GetCfg().BrokerTCPUrl)

		if !cmd.app.MqttClient.IsConnectionOpen() || !cmd.app.MqttClient.IsConnected() {
			q.NewConnection(cmd.app.MqttClient)
		}

		cmd.initializeSubscribers()
		<-c.GetCfgObject().GetSubscriber(c.WorkerKey)
		cmd.cleanUpSubscribers()
	}
}

func (cmd *WorkersCmd) CleanUp() {
	log.Println("Gracefully closing mqtt workers...")
	if cmd.app.MqttClient.IsConnectionOpen() && cmd.app.MqttClient.IsConnected() {
		cmd.cleanUpSubscribers()
		q.Disconnect(cmd.app.MqttClient)
	}
}

func (cmd *WorkersCmd) cleanUpSubscribers() {
	for _, s := range c.GetCfg().Subscribers {
		topic := fmt.Sprintf("apiary/%s/hive/%s", s.ApiaryId, s.HiveId)
		topicToDelete, alreadyExists := lo.Find(cmd.topics, func(t string) bool {
			return t == topic
		})

		if alreadyExists {
			go func(topic string) {
				q.Unsubscribe(cmd.app.MqttClient, topicToDelete)
			}(topic)
		}

	}
}

func (cmd *WorkersCmd) initializeSubscribers() {
	for _, s := range c.GetCfg().Subscribers {
		topic := fmt.Sprintf("apiary/%s/hive/%s", s.ApiaryId, s.HiveId)
		_, alreadyExists := lo.Find(cmd.topics, func(t string) bool {
			return t == topic
		})

		if !alreadyExists {
			go func(topic string) {
				q.Subscribe(cmd.app.MqttClient, topic)
			}(topic)
		}
	}

}
