package worker

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func (w *workerCtx) DefaultMessageHandler() MQTT.MessageHandler {
	return func(client MQTT.Client, msg MQTT.Message) {
		log.Printf("TOPIC: %s\n", msg.Topic())
		log.Printf("MSG: %s\n", msg.Payload())
		log.Println("########################################################")
	}
}

func (w *workerCtx) PersistMessageHandler() MQTT.MessageHandler {
	return func(client MQTT.Client, msg MQTT.Message) {
		log.Printf("Received message from topic: %s\n", msg.Topic())
		// TODO: Implement logic to store payload in tstorage engine and to send message to influxdb
	}
}
