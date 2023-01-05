package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/beesbuddy/beesbuddy-worker/docs"
	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	h "github.com/beesbuddy/beesbuddy-worker/internal/handlers"
	s "github.com/beesbuddy/beesbuddy-worker/static"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/redirect/v2"
)

type webCmd struct {
	app *c.App
}

func NewWebRunner(app *c.App) c.CmdRunner {
	module := &webCmd{app: app}
	return module
}

func (cmd *webCmd) Run() {
	cfg := c.GetCfg()

	router := cmd.app.Router

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
	auth.Post("/login", h.ApiLogin())

	// Set up static file serving
	var fileServer http.Handler
	var docsServer http.Handler

	if c.GetCfg().IsProd {
		staticFS := http.FS(s.Files)
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

	// Set up ui and inertia for handling vue serving
	ui := router.Group("/")
	ui.Use(adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return cmd.app.InertiaManager.Middleware(next)
	}))
	ui.Get("/", h.HomeHandler(cmd.app))

	go func() {
		if err := cmd.app.Router.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)); err != nil {
			log.Panic(err)
		}
	}()
}

func (cmd *webCmd) CleanUp() {
	log.Println("Gracefully closing web...")
	go func() {
		err := cmd.app.Router.Shutdown()

		if err != nil {
			log.Panic(err)
			os.Exit(1)
		}
	}()
}
