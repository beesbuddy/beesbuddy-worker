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
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/web/handler/api"
	"github.com/beesbuddy/beesbuddy-worker/internal/web/handler/middleware"
	"github.com/beesbuddy/beesbuddy-worker/internal/web/handler/ui"
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
	ctx *app.Ctx
}

func NewWebRunner(ctx *app.Ctx) core.Module {
	m := &webModule{ctx}
	return m
}

func (m *webModule) Run() {
	cfg := m.ctx.Config.GetCfg()

	// TODO: Move to routes file
	fiber := m.ctx.Fiber

	// Set up base handlers / middleware
	fiber.Use(recover.New())
	fiber.Use(logger.New())
	fiber.Get("/dashboard", monitor.New())
	fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	apiV1 := fiber.Group("/api/v1", limiter.New(limiter.Config{Max: 100}))

	// Redirect rules for /api/v1
	fiber.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/docs/": "/swagger/index.html",
			"/docs":  "/swagger/index.html",
		},
		StatusCode: 301,
	}))

	// Docs
	fiber.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger.json",
		DeepLinking: true,
	}))

	// Auth
	auth := apiV1.Group("/auth")
	auth.Post("/token", api.ApiGenerateToken(m.ctx))

	// Settings
	settings := apiV1.Group("/settings")
	settings.Use(jwtware.New(jwtware.Config{
		SigningKey:   []byte(m.ctx.Config.GetCfg().Secret),
		ErrorHandler: internal.AuthError,
	}))
	settings.Get("/subscribers", api.ApiGetSubscribers(m.ctx))
	settings.Post("/subscribers", api.ApiCreateSubscriber(m.ctx))

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
	fiber.Use("/swagger.json", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return docsServer
	}))

	// Set up ui and inertia for handling vue
	front := fiber.Group("/")
	front.Use(adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return m.ctx.InertiaManager.Middleware(next)
	}))

	// Pages
	front.Get("/login", ui.WebGetLoginHandler(m.ctx))
	front.Post("/login", ui.WebPostLoginHandler(m.ctx))
	front.Get("/", middleware.RedirectIfNotAuthenticated(m.ctx), ui.WebGetHomeHandler(m.ctx))

	go func(m *webModule) {
		defer m.CleanUp()
		if err := m.ctx.Fiber.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)); err != nil {
			log.Panic(err)
		}
	}(m)
}

func (m *webModule) CleanUp() {
	log.Println("Gracefully closing web...")

	go func() {
		err := m.ctx.Fiber.Shutdown()

		if err != nil {
			log.Panic(err)
			os.Exit(1)
		}
	}()
}
