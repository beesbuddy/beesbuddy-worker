package app

import (
	p "github.com/beesbuddy/beesbuddy-worker/pref"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Ctx struct {
	Fiber          *fiber.App
	MqttClient     MQTT.Client
	Pref           p.Preferences[p.AppPreferences]
	InlfuxDbClient influxdb2.Client
}

func NewContext(pref p.Preferences[p.AppPreferences]) *Ctx {
	router := fiber.New(fiber.Config{Prefork: pref.GetConfig().IsPrefork})

	opts := MQTT.NewClientOptions().AddBroker(pref.GetConfig().BrokerTCPUrl).SetAutoReconnect(true)
	mqttClient := MQTT.NewClient(opts)

	influxDbClient := influxdb2.NewClient(pref.GetConfig().InfluxDbURL, pref.GetConfig().InfluxDbAccessToken)

	ctx := &Ctx{
		Fiber:          router,
		MqttClient:     mqttClient,
		Pref:           pref,
		InlfuxDbClient: influxDbClient,
	}

	return ctx
}
