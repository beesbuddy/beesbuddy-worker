package cmd

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	mod "github.com/beesbuddy/beesbuddy-worker/internal/modules"
	"github.com/petaki/support-go/cli"
)

func WebServe(group *cli.Group, command *cli.Command, arguments []string) int {
	_, err := command.Parse(arguments)
	if err != nil {
		return command.PrintHelp(group)
	}

	app := core.NewApp()
	var cmd *exec.Cmd

	if !core.GetCfg().IsProd {
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

	tstorage := mod.NewTstorageRunner(app)
	tstorage.Run()
	workersRunner := mod.NewWorkersRunner(app)
	workersRunner.Run()
	webRunner := mod.NewWebRunner(app)
	webRunner.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
	}

	workersRunner.CleanUp()
	tstorage.CleanUp()
	webRunner.CleanUp()

	return cli.Success
}
