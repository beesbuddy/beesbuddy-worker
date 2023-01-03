package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	h "github.com/beesbuddy/beesbuddy-worker/internal/handlers"
	"github.com/beesbuddy/beesbuddy-worker/static"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	// Set up static file serving
	var fileServer http.Handler

	if c.GetCfg().IsProd {
		staticFS := http.FS(static.Files)
		fileServer = http.FileServer(staticFS)
	} else {
		fileServer = http.FileServer(http.Dir("./static/"))
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
