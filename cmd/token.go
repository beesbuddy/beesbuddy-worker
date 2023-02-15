package cmd

import (
	"fmt"
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/config"
	"github.com/petaki/support-go/cli"
	"github.com/samber/lo"
)

func Token(ctx *app.Ctx) func(*cli.Group, *cli.Command, []string) int {
	return func(group *cli.Group, command *cli.Command, arguments []string) int {
		parsed, err := command.Parse(arguments)

		if err != nil {
			return command.PrintHelp(group)
		}

		clients := ctx.Config.GetCfg().Clients

		client, ok := lo.Find(clients, func(client config.Client) bool {
			return client.AppKey == parsed[0]
		})

		if !ok {
			log.Fatal("unable to get client")
			return cli.Failure
		}

		secret := ctx.Config.GetCfg().Secret

		token, err := internal.GenerateJWTToken(client.AppKey, secret)

		if err != nil {
			return cli.Failure
		}

		fmt.Printf("succesfully generated token: %v\n", token)

		return cli.Success
	}
}
