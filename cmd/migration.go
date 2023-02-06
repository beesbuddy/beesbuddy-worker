package cmd

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	migration "github.com/beesbuddy/beesbuddy-worker/internal/migration"
	"github.com/petaki/support-go/cli"
)

func Migrate(ctx *app.Ctx) func(*cli.Group, *cli.Command, []string) int {
	return func(group *cli.Group, command *cli.Command, arguments []string) int {
		_, err := command.Parse(arguments)

		if err != nil {
			return command.PrintHelp(group)
		}

		migration := migration.NewMigrationRunner(ctx)
		migration.Run()

		return cli.Success
	}
}
