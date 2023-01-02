package cmd

import (
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
)

type WorkersCmd struct {
	app *core.App
}

func NewWorkersRunner(app *core.App) core.CmdRunner {
	cmd := &WorkersCmd{app}
	return cmd
}

func (cmd *WorkersCmd) Run() {
	if token := cmd.app.MqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (cmd *WorkersCmd) CleanUp() {
	log.Println("Gracefully closing mqtt workers...")

	go func() {
		cmd.app.MqttClient.Disconnect(250)
	}()
}
