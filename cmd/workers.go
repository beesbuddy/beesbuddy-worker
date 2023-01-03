package cmd

import (
	"log"

	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	m "github.com/beesbuddy/beesbuddy-worker/internal/messaging"
)

type WorkersCmd struct {
	app *c.App
}

func NewWorkersRunner(app *c.App) c.CmdRunner {
	cmd := &WorkersCmd{app: app}
	return cmd
}

func (cmd *WorkersCmd) Run() {
	m.NewConnection(cmd.app.MqttClient)

	for {
		log.Println("[Re]configuring RabbitMQ:", c.GetCfg().BrokerTCPUrl)

		<-c.GetCfgObject().GetSubscriber(c.WorkerKey)
	}
}

func (cmd *WorkersCmd) CleanUp() {
	log.Println("Gracefully closing mqtt workers...")
	m.Disconnect(cmd.app.MqttClient)
}
