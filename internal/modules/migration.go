package modules

import "github.com/beesbuddy/beesbuddy-worker/internal/core"

type migrationModule struct {
	app *core.Ctx
}

func NewMigrationRunner(app *core.Ctx) core.Module {
	m := &migrationModule{app: app}
	return m
}

func (m *migrationModule) Run() {

}

func (m *migrationModule) CleanUp() {

}
