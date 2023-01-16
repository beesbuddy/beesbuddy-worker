package handlers

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/core"
	"github.com/beesbuddy/beesbuddy-worker/internal/dto"
	"github.com/beesbuddy/beesbuddy-worker/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// Get active subscribers
// @Summary Get active subscribers
// @Description Get subscribers
// @Tags settings
// @Produce json
// @Success 200 {object} core.ResponseHTTP{data=[]models.Subscriber}
// @Failure 503 {object} core.ResponseHTTP{}
// @Router /settings/subscribers [get]
// @Security ApiKeyAuth
func ApiGetSubscribers(app *core.App) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.JSON(core.ResponseHTTP{
			Success: true,
			Message: "Successfully fetched subscribers",
			Data:    core.GetCfg().Subscribers,
		})
	}
}

// Post subscriber
// @Summary Create a new subscriber
// @Description Create a subscriber
// @Tags settings
// @Produce json
// @Success 200 {object} core.ResponseHTTP{data=[]models.Subscriber}
// @Failure 503 {object} core.ResponseHTTP{}
// @Param dto.SubscriberInput body dto.SubscriberInput true "Subscriber"
// @Router /settings/subscribers [post]
// @Security ApiKeyAuth
func ApiCreateSubscriber(app *core.App) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newSubscriber := new(dto.SubscriberInput)

		if err := ctx.BodyParser(newSubscriber); err != nil {
			return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
		}

		newConfig := core.GetCfg()

		_, alreadyExists := lo.Find(newConfig.Subscribers, func(s models.Subscriber) bool {
			return s.ApiaryId == newSubscriber.ApiaryId && s.HiveId == newSubscriber.HiveId
		})

		if !alreadyExists {
			subscriber := models.Subscriber{ApiaryId: newSubscriber.ApiaryId, HiveId: newSubscriber.HiveId}
			newConfig.Subscribers = append(newConfig.Subscribers, subscriber)
			core.GetCfgObject().Update(newConfig)
		}

		return ctx.JSON(core.ResponseHTTP{
			Success: true,
			Message: "Registered subscriber for creation",
			Data:    core.GetCfg().Subscribers,
		})
	}
}
