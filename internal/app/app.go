package app

import (
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	"github.com/chmike/securecookie"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	c "github.com/leonidasdeim/goconfig"
	"github.com/petaki/inertia-go"
	"github.com/petaki/support-go/mix"
	"gorm.io/gorm"
)

type Ctx struct {
	Fiber          *fiber.App
	MixManager     *mix.Mix
	InertiaManager *inertia.Inertia
	SessionManager *session.Store
	RememberCookie *securecookie.Obj
	MqttClient     MQTT.Client
	UserRepository *model.UserRepository
	Config         *c.Config[model.Config]
	Orm            *gorm.DB
}

func NewContext(config *c.Config[model.Config]) *Ctx {
	sessionManager := session.New()
	rememberCookie, err := securecookie.New(internal.RememberCookieNameKey, []byte(config.GetCfg().Secret), securecookie.Params{
		Path:     "/",
		MaxAge:   157680000, // Five years
		HTTPOnly: true,
	})

	if err != nil {
		panic("unable to set up remember cookie")
	}

	router := fiber.New(fiber.Config{Prefork: config.GetCfg().IsPrefork})

	debug := !config.GetCfg().IsProd
	url := ""

	if debug {
		url = config.GetCfg().UiHotReloadUrl
	}

	mixManager, inertiaManager, err := core.NewMixAndInertiaManager(
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

	gorm := core.NewDatabase(config.GetCfg())

	ctx := &Ctx{
		Fiber:          router,
		MixManager:     mixManager,
		InertiaManager: inertiaManager,
		MqttClient:     mqttClient,
		Config:         config,
		Orm:            gorm,
		UserRepository: model.NewUserRepository(gorm),
		RememberCookie: rememberCookie,
		SessionManager: sessionManager,
	}

	return ctx
}
