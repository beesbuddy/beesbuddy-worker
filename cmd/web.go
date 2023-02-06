package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/web"
	"github.com/beesbuddy/beesbuddy-worker/internal/worker"
	"github.com/petaki/support-go/cli"
)

func WebServe(ctx *app.Ctx) func(*cli.Group, *cli.Command, []string) int {
	return func(group *cli.Group, command *cli.Command, arguments []string) int {
		_, err := command.Parse(arguments)
		if err != nil {
			return command.PrintHelp(group)
		}

		workersRunner := worker.NewWorkersRunner(ctx)
		workersRunner.Run()
		webRunner := web.NewWebRunner(ctx)
		webRunner.Run()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		<-interrupt

		workersRunner.CleanUp()
		webRunner.CleanUp()

		return cli.Success
	}
}
