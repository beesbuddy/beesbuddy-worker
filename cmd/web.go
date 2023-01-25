package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	mod "github.com/beesbuddy/beesbuddy-worker/internal/modules"
	"github.com/petaki/support-go/cli"
)

func WebServe(group *cli.Group, command *cli.Command, arguments []string) int {
	command.FlagSet().Bool("debug", false, "Application debug mode")
	command.FlagSet().Bool("config", false, "Application configuration path")

	_, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	app := core.NewApp()

	tstorage := mod.NewTstorageRunner(app)
	tstorage.Run()
	workersRunner := mod.NewWorkersRunner(app)
	workersRunner.Run()
	webRunner := mod.NewWebRunner(app)
	webRunner.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	workersRunner.CleanUp()
	tstorage.CleanUp()
	webRunner.CleanUp()

	return cli.Success
}
