package queue

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Unsubscribe(c MQTT.Client, topic string) {
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

func Subscribe(c MQTT.Client, topic string) {
	if token := c.Subscribe(topic, 0, DefaultMessageHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}
