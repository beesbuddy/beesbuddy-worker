package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func NewConnection(mqttClient MQTT.Client) {
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func Disconnect(mqttClient MQTT.Client) {
	go func() {
		mqttClient.Disconnect(250)
	}()
}
