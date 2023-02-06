package cmd

import (
	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	m "github.com/beesbuddy/beesbuddy-worker/internal/module"
	"github.com/petaki/support-go/cli"
)

func Migrate(ctx *c.Ctx) func(*cli.Group, *cli.Command, []string) int {
	return func(group *cli.Group, command *cli.Command, arguments []string) int {
		_, err := command.Parse(arguments)

		if err != nil {
			return command.PrintHelp(group)
		}

		migration := m.NewMigrationRunner(ctx)
		migration.Run()

		return cli.Success
	}
}
