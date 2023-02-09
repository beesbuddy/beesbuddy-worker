package handler

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var DefaultMessageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf("TOPIC: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
	log.Println("########################################################")
}

var PersistMessageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf("Received message from topic: %s\n", msg.Topic())
	// TODO: Implement logic to store payload in tstorage engine and to send message to influxdb
}
