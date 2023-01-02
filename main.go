package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/beesbuddy/beesbuddy-worker/cmd"
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
)

func main() {
	app := core.NewApplication()

	appRunner := cmd.NewWebRunner(app)
	appRunner.Run()

	mqttClientRunner := cmd.NewMqttClientRunner(app)
	mqttClientRunner.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	mqttClientRunner.CleanUp()
	appRunner.CleanUp()
}
