package app

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/app/settings"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/leonidasdeim/goconfig"
)

type Ctx struct {
	Fiber      *fiber.App
	MqttClient MQTT.Client
	Config     *goconfig.Config[settings.AppSettings]
}

func NewContext(config *goconfig.Config[settings.AppSettings]) *Ctx {
	router := fiber.New(fiber.Config{Prefork: config.GetCfg().IsPrefork})

	opts := MQTT.NewClientOptions().AddBroker(config.GetCfg().BrokerTCPUrl).SetAutoReconnect(true)
	mqttClient := MQTT.NewClient(opts)

	ctx := &Ctx{
		Fiber:      router,
		MqttClient: mqttClient,
		Config:     config,
	}

	return ctx
}
