package cmd

import (
	"fmt"
	"log"

	"github.com/beesbuddy/beesbuddy-worker/app"
	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/beesbuddy/beesbuddy-worker/internal/pref"
	"github.com/petaki/support-go/cli"
	"github.com/samber/lo"
)

func Token(ctx *app.Ctx) func(*cli.Group, *cli.Command, []string) int {
	return func(group *cli.Group, command *cli.Command, arguments []string) int {
		parsed, err := command.Parse(arguments)

		config := ctx.Pref.GetConfig()

		if err != nil {
			return command.PrintHelp(group)
		}

		clients := config.Clients

		client, ok := lo.Find(clients, func(client pref.Client) bool {
			return client.AppKey == parsed[0]
		})

		if !ok {
			log.Fatal("unable to get client")
			return cli.Failure
		}

		secret := config.Secret

		token, err := internal.GenerateToken(client.AppKey, secret)

		if err != nil {
			return cli.Failure
		}

		fmt.Printf("succesfully generated token: %v\n", token)

		return cli.Success
	}
}
