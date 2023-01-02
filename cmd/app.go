package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
)

type appCmd struct {
	app *core.App
}

func NewApplicationRunner(app *core.App) core.CmdRunner {
	module := &appCmd{app: app}
	return module
}

func (cmd *appCmd) Run() {
	cfg := core.GetCfg()

	go func() {
		if err := cmd.app.Router.Listen(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)); err != nil {
			log.Panic(err)
		}
	}()
}

func (cmd *appCmd) CleanUp() {
	log.Println("Gracefully closing app...")
	go func() {
		err := cmd.app.Router.Shutdown()

		if err != nil {
			log.Panic(err)
			os.Exit(1)
		}
	}()
}
