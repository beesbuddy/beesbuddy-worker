package cmd

import (
	"fmt"
	"log"

	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	"github.com/beesbuddy/beesbuddy-worker/internal/util"
	"github.com/petaki/support-go/cli"
	"github.com/samber/lo"
)

func Token(ctx *c.Ctx) func(*cli.Group, *cli.Command, []string) int {
	return func(group *cli.Group, command *cli.Command, arguments []string) int {
		parsed, err := command.Parse(arguments)

		if err != nil {
			return command.PrintHelp(group)
		}

		clients := ctx.Config.GetCfg().Clients

		client, ok := lo.Find(clients, func(client model.Client) bool {
			return client.AppKey == parsed[0]
		})

		if !ok {
			log.Fatal("unable to get client")
			return cli.Failure
		}

		secret := ctx.Config.GetCfg().Secret

		token, err := util.GenerateJWTToken(client.AppKey, secret)

		if err != nil {
			return cli.Failure
		}

		fmt.Printf("succesfully generated token: %v\n", token)

		return cli.Success
	}
}
