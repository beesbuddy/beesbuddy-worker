package cmd

import (
	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/petaki/support-go/cli"
)

func Migrate(ctx *c.Ctx) func(*cli.Group, *cli.Command, []string) int {
	return func(group *cli.Group, command *cli.Command, arguments []string) int {
		_, err := command.Parse(arguments)
		if err != nil {
			return command.PrintHelp(group)
		}

		return cli.Success
	}
}
