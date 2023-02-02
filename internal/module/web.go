package module

import (
	"fmt"
	"log"
	"net/http"
	"os"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/beesbuddy/beesbuddy-worker/docs"
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/handler"
	"github.com/beesbuddy/beesbuddy-worker/static"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/redirect/v2"
)

type webModule struct {
	ctx *core.Ctx
}

func NewWebRunner(ctx *core.Ctx) core.Module {
	m := &webModule{ctx}
	return m
}

func (m *webModule) Run() {
	cfg := m.ctx.Config.GetCfg()

	router := m.ctx.Router

	// Set up base handlers / middleware
	router.Use(recover.New())
	router.Use(logger.New())
	router.Get("/dashboard", monitor.New())
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	apiV1 := router.Group("/api/v1", limiter.New(limiter.Config{Max: 100}))

	// Redirect rules for /api/v1
	router.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/docs/": "/swagger/index.html",
			"/docs":  "/swagger/index.html",
		},
		StatusCode: 301,
	}))

	// Docs
	router.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger.json",
		DeepLinking: true,
	}))

	// Auth
	auth := apiV1.Group("/auth")
	auth.Post("/token", handler.ApiGenerateToken(m.ctx))

	// Settings
	settings := apiV1.Group("/settings")
	settings.Use(jwtware.New(jwtware.Config{
		SigningKey:   []byte(m.ctx.Config.GetCfg().Secret),
		ErrorHandler: core.AuthError,
	}))
	settings.Get("/subscribers", handler.ApiGetSubscribers(m.ctx))
	settings.Post("/subscribers", handler.ApiCreateSubscriber(m.ctx))

	// Set up static file serving
	var fileServer http.Handler
	var docsServer http.Handler

	if m.ctx.Config.GetCfg().IsProd {
		staticFS := http.FS(static.Files)
		docsFS := http.FS(docs.Files)
		fileServer = http.FileServer(staticFS)
		docsServer = http.FileServer(docsFS)
	} else {
		fileServer = http.FileServer(http.Dir("./static/"))
		docsServer = http.FileServer(http.Dir("./docs/"))
	}

	router.Use("/css/", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return fileServer
	}))
	router.Use("/js/", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return fileServer
	}))
	router.Use("/images/", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return fileServer
	}))
	router.Use("/favicon.ico", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return fileServer
	}))
	router.Use("/swagger.json", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return docsServer
	}))

	// Set up ui and inertia for handling vue
	ui := router.Group("/")
	ui.Use(adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return m.ctx.InertiaManager.Middleware(next)
	}))
	// Pages
	ui.Get("/", handler.WebHomeHandler(m.ctx))

	go func(m *webModule) {
		defer m.CleanUp()
		if err := m.ctx.Router.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)); err != nil {
			log.Panic(err)
		}
	}(m)
}

func (m *webModule) CleanUp() {
	log.Println("Gracefully closing web...")

	go func() {
		err := m.ctx.Router.Shutdown()

		if err != nil {
			log.Panic(err)
			os.Exit(1)
		}
	}()
}
