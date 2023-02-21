package worker

import (
	"encoding/json"

	"github.com/beesbuddy/beesbuddy-worker/internal/log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func (w *workerComponent) DefaultMessageHandler() MQTT.MessageHandler {
	return func(client MQTT.Client, msg MQTT.Message) {
		log.Info.Println("TOPIC: ", msg.Topic())
		log.Info.Println("MSG: ", msg.Payload())
		log.Info.Println("########################################################")
	}
}

func (w *workerComponent) PersistMessageHandler() MQTT.MessageHandler {
	return func(client MQTT.Client, msg MQTT.Message) {
		log.Info.Println("received message from topic: ", msg.Topic())

		var m metrics
		json.Unmarshal(msg.Payload(), &m)

		select {
		case w.queue <- m:
			log.Info.Println("sending message to queue: ", m)
		default:
			log.Error.Println("unable to send message")
		}
	}
}
