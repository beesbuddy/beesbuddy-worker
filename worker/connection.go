package worker

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func NewConnection(mqttClient MQTT.Client) {
	if connection := mqttClient.Connect(); connection.Wait() && connection.Error() != nil {
		panic(connection.Error())
	}
}

func Disconnect(mqttClient MQTT.Client) {
	go func() {
		mqttClient.Disconnect(250)
	}()
}
