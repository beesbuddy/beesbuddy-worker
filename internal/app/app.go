package app

import (
	"github.com/alexedwards/scs/v2"
	"github.com/beesbuddy/beesbuddy-worker/internal/app/settings"
	"github.com/chmike/securecookie"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/leonidasdeim/goconfig"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/mix"
)

type Ctx struct {
	Fiber          *fiber.App
	MixManager     *mix.Mix
	InertiaManager *inertia.Inertia
	SessionManager *scs.SessionManager
	RememberCookie *securecookie.Obj
	MqttClient     MQTT.Client
	Config         *goconfig.Config[settings.AppSettings]
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
