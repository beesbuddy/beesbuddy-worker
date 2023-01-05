package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/beesbuddy/beesbuddy-worker/cmd"
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
)

// @title BeesBuddy worker
// @version 1.0
// @description This is an API for Worker Module

// @contact.name Viktor Nareiko
// @contact.email vnareiko.lt@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app := core.NewApplication()

	webRunner := cmd.NewWebRunner(app)
	webRunner.Run()

	workersRunner := cmd.NewWorkersRunner(app)
	workersRunner.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	workersRunner.CleanUp()
	webRunner.CleanUp()
}
