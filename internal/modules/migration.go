package modules

import c "github.com/beesbuddy/beesbuddy-worker/internal/core"

type migrationModule struct {
	app *c.Ctx
}

func NewMigrationRunner(app *c.Ctx) c.Module {
	m := &migrationModule{app: app}
	return m
}

func (m *migrationModule) Run() {

}

func (m *migrationModule) CleanUp() {

}
