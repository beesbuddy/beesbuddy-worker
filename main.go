package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/beesbuddy/beesbuddy-worker/cmd"
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
)

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
