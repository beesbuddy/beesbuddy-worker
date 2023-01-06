package handlers

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/gofiber/fiber/v2"
)

// Get active workers
// @Summary Get active workers
// @Description Get workers
// @Tags settings
// @Produce json
// @Success 200 {object} core.ResponseHTTP{data=[]string}
// @Failure 503 {object} core.ResponseHTTP{}
// @Router /settings/workers [get]
// @Security ApiKeyAuth
func ApiGetWorkers(ctx *fiber.Ctx) error {
	return ctx.JSON(core.ResponseHTTP{
		Success: true,
		Message: "Successfully fetched workers",
		Data:    []string{"w1", "w2", "w3", "w4"},
	})
}
