package modules

import "github.com/beesbuddy/beesbuddy-worker/internal/core"

type tstorageModule struct {
	app *core.App
}

func NewTstorageRunner(app *core.App) core.Module {
	m := &tstorageModule{app: app}
	return m
}

func (m *tstorageModule) Run() {

}

func (m *tstorageModule) CleanUp() {

}
