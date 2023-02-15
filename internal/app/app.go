package app

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/config"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/leonidasdeim/goconfig"
)

type Ctx struct {
	Fiber          *fiber.App
	MqttClient     MQTT.Client
	Config         *goconfig.Config[config.AppPreferences]
	InlfuxDbClient influxdb2.Client
}

func NewContext(config *goconfig.Config[config.AppPreferences]) *Ctx {
	router := fiber.New(fiber.Config{Prefork: config.GetCfg().IsPrefork})

	opts := MQTT.NewClientOptions().AddBroker(config.GetCfg().BrokerTCPUrl).SetAutoReconnect(true)
	mqttClient := MQTT.NewClient(opts)

	influxDbClient := influxdb2.NewClient(config.GetCfg().InfluxDbURL, config.GetCfg().InfluxDbAccessToken)

	ctx := &Ctx{
		Fiber:          router,
		MqttClient:     mqttClient,
		Config:         config,
		InlfuxDbClient: influxDbClient,
	}

	return ctx
}
