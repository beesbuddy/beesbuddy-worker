package cmd

import (
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
)

type MqttClientCmd struct {
	app *core.App
}

func NewMqttClientRunner(app *core.App) core.CmdRunner {
	cmd := &MqttClientCmd{app}
	return cmd
}

func (cmd *MqttClientCmd) Run() {
	if token := cmd.app.MqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (cmd *MqttClientCmd) CleanUp() {
	log.Println("Gracefully closing mqtt connections...")

	go func() {
		cmd.app.MqttClient.Disconnect(250)
	}()
}
