package cmd

import (
	"fmt"

	c "github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	"github.com/petaki/support-go/cli"
	"github.com/petaki/support-go/forms"
)

func User(ctx *c.Ctx) func(*cli.Group, *cli.Command, []string) int {
	return func(group *cli.Group, command *cli.Command, arguments []string) int {
		enabled := command.FlagSet().Bool("enabled", true, "User Enabled")

		parsed, err := command.Parse(arguments)

		if err != nil {
			return command.PrintHelp(group)
		}

		form := forms.New(map[string]interface{}{
			"username":   parsed[0],
			"email":      parsed[1],
			"password":   parsed[2],
			"is_enabled": *enabled,
		})

		model.UserCreateRules(form)

		if !form.IsValid() {
			return command.PrintError(fmt.Errorf("make user: invalid arguments: %v", form.Errors))
		}

		user, err := (&model.User{}).Fill(form)

		if err != nil {
			cli.ErrorLog.Panicln("unable to generate password hash")
			return cli.Failure
		}

		ctx.UserRepository.Create(user)

		return cli.Success
	}
}
