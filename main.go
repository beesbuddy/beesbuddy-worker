package main

import (
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/beesbuddy/beesbuddy-worker/cmd"
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//go:embed public/*.*
//go:embed public/*/*.*
var embedFS embed.FS

func main() {
	core.InitializeCfg()

	app := core.NewApplication()
	appRunner := cmd.NewApplicationRunner(app)
	appRunner.Run()

	ui := cmd.NewUIRunner(app, &embedFS)
	ui.Run()

	opts := MQTT.NewClientOptions().AddBroker(core.GetCfgModel().BrokerTCPUrl)
	client := MQTT.NewClient(opts)
	mqttClientRunner := cmd.NewMqttClientRunner(client)
	mqttClientRunner.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	mqttClientRunner.CleanUp()
	appRunner.CleanUp()
}
