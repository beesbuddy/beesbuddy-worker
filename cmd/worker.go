package cmd

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/log"
	"github.com/beesbuddy/beesbuddy-worker/internal/starter"
	"github.com/beesbuddy/beesbuddy-worker/internal/web"
	"github.com/beesbuddy/beesbuddy-worker/internal/worker"
	"github.com/petaki/support-go/cli"
)

func WrokerServe(ctx *app.Ctx) func(*cli.Group, *cli.Command, []string) int {
	return func(group *cli.Group, command *cli.Command, arguments []string) int {
		_, err := command.Parse(arguments)
		if err != nil {
			return command.PrintHelp(group)
		}

		workersRunner := worker.NewWorkersRunner(ctx)
		workersRunner.Run()
		webRunner := web.NewWebRunner(ctx)
		webRunner.Run()

		// Add shutdown handlers
		starter.Handle(webRunner.Flush)
		starter.Handle(workersRunner.Flush)
		starter.Handle(log.Flush)

		// Init starter. It will also make application running till interrupt signal will be received.
		starter.Ignite(cli.Success)

		return cli.Success
	}
}
