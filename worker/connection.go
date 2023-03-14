package worker

import (
	"time"

	"github.com/beesbuddy/beesbuddy-worker/internal/log"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func (w *workerComponent) NewConnection(mqttClient MQTT.Client) {
	count := 1
	for connection := mqttClient.Connect(); connection.Wait() && connection.Error() != nil; {
		log.Error.Println(connection.Error())
		t := time.Duration(count * int(time.Second))
		log.Info.Printf("retrying the MQTT connection in %d seconds...\n", int(t.Abs().Seconds()))

		time.Sleep(t)

		count *= 2
		if count > 32 {
			count = 1
		}
	}
}

func (w *workerComponent) Disconnect(mqttClient MQTT.Client) {
	go func() {
		mqttClient.Disconnect(250)
	}()
}
