package module

import (
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
)

type migrationModule struct {
	ctx *app.Ctx
}

func NewMigrationRunner(app *app.Ctx) core.Module {
	m := &migrationModule{app}
	return m
}

func (m *migrationModule) Run() {
	log.Println("Runing database migtations")
	var err = m.ctx.Orm.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalln("Failed to migrate `User` entity")
	}
	log.Println("☑️ - migrated `User` entity")

}

func (m *migrationModule) CleanUp() {
	// Nothing to clean up after migrations
}
