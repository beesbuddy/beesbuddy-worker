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
		log.Info.Println("retrying the MQTT connection in", int(t.Abs().Seconds()), "seconds...")

		time.Sleep(t)

		count *= 2
		if count > 32 {
			count = 1
		}
	}
	log.Info.Println("successfully connected to MQTT:", w.appCtx.Pref.GetConfig().BrokerTCPUrl)
}

func (w *workerComponent) Disconnect(mqttClient MQTT.Client) {
	go func() {
		mqttClient.Disconnect(250)
	}()
}
