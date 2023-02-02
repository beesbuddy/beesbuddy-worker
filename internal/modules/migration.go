package modules

import (
	"log"

	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/models"
)

type migrationModule struct {
	ctx *core.Ctx
}

func NewMigrationRunner(ctx *core.Ctx) core.Module {
	m := &migrationModule{ctx}
	return m
}

func (m *migrationModule) Run() {
	log.Println("Runing database migtations")
	var err = m.ctx.Orm.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalln("Failed to migrate `User` entity")
	}
	log.Println("☑️ - migrated `User` entity")

}

func (m *migrationModule) CleanUp() {
	// Nothing to clean up after migrations
}
