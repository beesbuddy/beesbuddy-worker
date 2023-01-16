package cmd

import (
	"fmt"
	"log"

	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	q "github.com/beesbuddy/beesbuddy-worker/internal/queue"
)

type WorkersCmd struct {
	app *c.App
}

func NewWorkersRunner(app *c.App) c.CmdRunner {
	cmd := &WorkersCmd{app: app}
	return cmd
}

func (cmd *WorkersCmd) Run() {
	q.NewConnection(cmd.app.MqttClient)

	for {
		log.Println("[Re]configuring MQTT:", c.GetCfg().BrokerTCPUrl)
		cmd.initSubscribers()
		<-c.GetCfgObject().GetSubscriber(c.WorkerKey)
		cmd.cleanUpSubscribers()
	}
}

func (cmd *WorkersCmd) CleanUp() {
	log.Println("Gracefully closing mqtt workers...")
	cmd.cleanUpSubscribers()
	q.Disconnect(cmd.app.MqttClient)
}

func (cmd *WorkersCmd) cleanUpSubscribers() {
	for _, s := range c.GetCfg().Subscribers {
		topic := fmt.Sprintf("apiary/%s/hive/%s", s.ApiaryId, s.HiveId)
		go func(topic string) {
			q.Unsubscribe(cmd.app.MqttClient, topic)
		}(topic)
	}
}

func (cmd *WorkersCmd) initSubscribers() {
	for _, s := range c.GetCfg().Subscribers {
		topic := fmt.Sprintf("apiary/%s/hive/%s", s.ApiaryId, s.HiveId)
		go func(topic string) {
			q.Subscribe(cmd.app.MqttClient, topic)
		}(topic)
	}
}
