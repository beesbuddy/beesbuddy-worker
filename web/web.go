package web

import (
	"fmt"
	"net/http"
	"os"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/beesbuddy/beesbuddy-worker/app"
	"github.com/beesbuddy/beesbuddy-worker/docs"
	"github.com/beesbuddy/beesbuddy-worker/internal/component"
	"github.com/beesbuddy/beesbuddy-worker/internal/log"
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

func NewWebRunner(ctx *app.Ctx) component.Component {
	w := &webCtx{ctx}
	return w
}

func (w *webCtx) Init() {
	appCtx := w.appCtx
	cfg := appCtx.Pref.GetConfig()
	fiber := appCtx.Fiber

	// set up handlers / middlewares
	fiber.Use(recover.New())
	fiber.Use(logger.New())

	fiber.Get("/metrics", monitor.New())

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
		SigningKey:   []byte(cfg.Secret),
		ErrorHandler: AuthError,
	}))
	settings.Get("/subscribers", ApiGetSubscribers(appCtx))
	settings.Post("/subscribers", ApiCreateSubscriber(appCtx))

	// set up static file serving
	var docsServer http.Handler

	if cfg.IsProd {
		docsFS := http.FS(docs.Files)
		docsServer = http.FileServer(docsFS)
	} else {
		docsServer = http.FileServer(http.Dir("./docs/"))
	}

	fiber.Use("/swagger.json", adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		return docsServer
	}))

	go func(w *webCtx) {
		if err := w.appCtx.Fiber.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)); err != nil {
			log.Error.Println(err)
			panic(err)
		}
	}(w)
}

func (w *webCtx) Flush() {
	log.Info.Println("gracefully closing web...")

	go func() {
		err := w.appCtx.Fiber.Shutdown()

		if err != nil {
			log.Error.Println(err)
			os.Exit(1)
		}
	}()
}
