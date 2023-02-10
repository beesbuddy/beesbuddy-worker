package web

import (
	"fmt"
	"log"
	"net/http"
	"os"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/beesbuddy/beesbuddy-worker/docs"
	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/redirect/v2"
)

type webCtx struct {
	appCtx *app.Ctx
}

func NewWebRunner(ctx *app.Ctx) internal.ModuleCtx {
	m := &webCtx{ctx}
	return m
}

func (m *webCtx) Run() {
	appCtx := m.appCtx
	cfg := appCtx.Config.GetCfg()
	fiber := appCtx.Fiber

	// set up handlers / middlewares
	fiber.Use(recover.New())
	fiber.Use(logger.New())

	fiber.Get("/dashboard", monitor.New())
	fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	apiV1 := fiber.Group("/api/v1", limiter.New(limiter.Config{Max: 100}))

	// redirect rules for docs
	fiber.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/":      "/swagger/index.html",
			"/docs/": "/swagger/index.html",
			"/docs":  "/swagger/index.html",
		},
		StatusCode: 301,
	}))

	// docs
	fiber.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger.json",
		DeepLinking: true,
	}))

	// settings
	settings := apiV1.Group("/settings")
	settings.Use(jwtware.New(jwtware.Config{
		SigningKey:   []byte(m.appCtx.Config.GetCfg().Secret),
		ErrorHandler: AuthError,
	}))
	settings.Get("/subscribers", ApiGetSubscribers(appCtx))
	settings.Post("/subscribers", ApiCreateSubscriber(appCtx))

	// set up static file serving
	var docsServer http.Handler

	if appCtx.Config.GetCfg().IsProd {
		docsFS := http.FS(docs.Files)
		docsServer = http.FileServer(docsFS)
	} else {
		docsServer = http.FileServer(http.Dir("./docs/"))
	}

	fiber.Use("/swagger.json", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return docsServer
	}))

	go func(m *webCtx) {
		defer m.CleanUp()
		if err := m.appCtx.Fiber.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)); err != nil {
			log.Panic(err)
		}
	}(m)
}

func (m *webCtx) CleanUp() {
	log.Println("Gracefully closing web...")

	go func() {
		err := m.appCtx.Fiber.Shutdown()

		if err != nil {
			log.Panic(err)
			os.Exit(1)
		}
	}()
}
