package core

import (
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/beesbuddy/beesbuddy-worker/internal/models"
	"github.com/chmike/securecookie"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	c "github.com/leonidasdeim/goconfig"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/mix"
	"gorm.io/gorm"
)

type Ctx struct {
	Router         *fiber.App
	MixManager     *mix.Mix
	InertiaManager *inertia.Inertia
	SessionManager *scs.SessionManager
	RememberCookie *securecookie.Obj
	MqttClient     MQTT.Client
	Services       map[string]interface{}
	Repositories   map[string]interface{}
	Config         *c.Config[models.Config]
	Orm            *gorm.DB
}

func NewContext(config *c.Config[models.Config]) *Ctx {
	router := fiber.New(fiber.Config{Prefork: config.GetCfg().IsPrefork})

	debug := !config.GetCfg().IsProd
	url := ""

	if debug {
		url = config.GetCfg().UiHotReloadUrl
	}

	mixManager, inertiaManager, err := newMixAndInertiaManager(
		debug,
		url,
		config.GetCfg().AppName,
	)

	if err != nil {
		log.Fatal(err)
	}

	opts := MQTT.NewClientOptions().AddBroker(config.GetCfg().BrokerTCPUrl).SetAutoReconnect(true)
	mqttClient := MQTT.NewClient(opts)

	if err != nil {
		log.Fatal(err)
	}

	gorm := NewDatabase(config.GetCfg())

	ctx := &Ctx{
		Router:         router,
		MixManager:     mixManager,
		InertiaManager: inertiaManager,
		MqttClient:     mqttClient,
		Config:         config,
		Orm:            gorm,
	}

	return ctx
}
