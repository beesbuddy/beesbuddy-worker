package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/gofiber/fiber/v2"
)

type appModule struct {
	app *fiber.App
}

func NewApplicationRunner(app *fiber.App) core.ModuleRunner {
	module := &appModule{app: app}
	return module
}

func (mod *appModule) Run() {
	cfg := core.GetCfgModel()

	go func() {
		if err := mod.app.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)); err != nil {
			log.Panic(err)
		}
	}()
}

func (mod *appModule) CleanUp() {
	log.Println("Gracefully shutting down fiber...")
	err := mod.app.Shutdown()

	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
}
