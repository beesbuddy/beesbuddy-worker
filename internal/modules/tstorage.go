package modules

import "github.com/beesbuddy/beesbuddy-worker/internal/core"

type tstorageMod struct {
	app *core.App
}

func NewTstorageRunner(app *core.App) core.Mod {
	mod := &tstorageMod{app: app}
	return mod
}

func (mod *tstorageMod) Run() {

}

func (mod *tstorageMod) CleanUp() {

}
