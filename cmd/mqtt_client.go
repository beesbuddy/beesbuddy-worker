package cmd

import (
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type mqttClientModule struct {
	client MQTT.Client
}

func NewMqttClientRunner(client MQTT.Client) core.ModuleRunner {
	module := &mqttClientModule{client: client}
	return module
}

func (mod *mqttClientModule) Run() {
	if token := mod.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (mod *mqttClientModule) CleanUp() {
	log.Println("Gracefully closing mqtt connections...")
	mod.client.Disconnect(250)
}
