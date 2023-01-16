package core

import (
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/chmike/securecookie"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/mix"
)

type App struct {
	Router         *fiber.App
	MixManager     *mix.Mix
	InertiaManager *inertia.Inertia
	SessionManager *scs.SessionManager
	RememberCookie *securecookie.Obj
	MqttClient     MQTT.Client
	Pool           chan int64
}

func NewApplication() *App {
	router := fiber.New(fiber.Config{Prefork: GetCfg().IsPrefork})

	debug := !GetCfg().IsProd
	url := ""

	if debug {
		url = GetCfg().UiHotReloadUrl
	}

	mixManager, inertiaManager, err := newMixAndInertiaManager(
		debug,
		url,
	)

	if err != nil {
		log.Fatal(err)
	}

	opts := MQTT.NewClientOptions().AddBroker(GetCfg().BrokerTCPUrl)
	client := MQTT.NewClient(opts)

	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		Router:         router,
		MixManager:     mixManager,
		InertiaManager: inertiaManager,
		MqttClient:     client,
	}

	return app
}
