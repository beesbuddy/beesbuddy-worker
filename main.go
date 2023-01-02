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

	appRunner := cmd.NewApplicationRunner(app)
	appRunner.Run()

	webRunner := cmd.NewWebRunner(app)
	webRunner.Run()

	mqttClientRunner := cmd.NewMqttClientRunner(app)
	mqttClientRunner.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	webRunner.CleanUp()
	mqttClientRunner.CleanUp()
	appRunner.CleanUp()
}
