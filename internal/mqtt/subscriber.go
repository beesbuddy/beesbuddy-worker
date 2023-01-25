package mqtt

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Unsubscribe(c MQTT.Client, topic string) {
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
		panic(token.Error())
	}
}

func Subscribe(c MQTT.Client, topic string) {
	if token := c.Subscribe(topic, 0, DefaultMessageHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
		panic(token.Error())
	}
}
