package core

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/beesbuddy/beesbuddy-worker/static"
	"github.com/chmike/securecookie"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	Queue          chan int64
}

func NewApplication() *App {
	fiber := fiber.New(fiber.Config{Prefork: GetCfg().IsPrefork})
	// Default handlers
	fiber.Use(recover.New())
	fiber.Use(logger.New())
	fiber.Get("/dashboard", monitor.New())
	fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Set up static file serving
	var fileServer http.Handler

	if GetCfg().IsProd {
		staticFS := http.FS(static.Files)
		fileServer = http.FileServer(staticFS)
	} else {
		fileServer = http.FileServer(http.Dir("./static/"))
	}

	fiber.Use("/css/", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return fileServer
	}))
	fiber.Use("/js/", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return fileServer
	}))
	fiber.Use("/images/", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return fileServer
	}))
	fiber.Use("/favicon.ico", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return fileServer
	}))

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
		Router:         fiber,
		MixManager:     mixManager,
		InertiaManager: inertiaManager,
		MqttClient:     client,
	}

	return app
}
