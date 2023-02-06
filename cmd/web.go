package cmd

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
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

		var cmd *exec.Cmd

		if ctx.Config.GetCfg().HotReload {
			name := "/bin/sh"
			arg := "-c"
			command := "npm run hot"

			if runtime.GOOS == "windows" {
				name = "cmd.exe"
				arg = "/C"
			}

			cmd = exec.Command(name, arg, command)
			cmd.Stdout = os.Stdout
			cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
		}

		workersRunner := worker.NewWorkersRunner(ctx)
		workersRunner.Run()
		webRunner := web.NewWebRunner(ctx)
		webRunner.Run()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		<-interrupt

		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}

		workersRunner.CleanUp()
		webRunner.CleanUp()

		return cli.Success
	}
}
