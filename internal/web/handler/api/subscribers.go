package api

import (
	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/dto"
	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// Get active subscribers
// @Summary Get active subscribers
// @Description Get subscribers
// @Tags settings
// @Produce json
// @Success 200 {object} dto.ResponseHTTP{data=[]model.Subscriber}
// @Failure 503 {object} dto.ResponseHTTP{}
// @Router /settings/subscribers [get]
// @Security ApiKeyAuth
func ApiGetSubscribers(ctx *app.Ctx) fiber.Handler {
	return func(f *fiber.Ctx) error {
		return f.JSON(dto.ResponseHTTP{
			Success: true,
			Message: "Successfully fetched subscribers",
			Data:    ctx.Config.GetCfg().Subscribers,
		})
	}
}

// Post subscriber
// @Summary Create a new subscriber
// @Description Create a subscriber
// @Tags settings
// @Produce json
// @Success 200 {object} dto.ResponseHTTP{data=[]model.Subscriber}
// @Failure 503 {object} dto.ResponseHTTP{}
// @Param dto.SubscriberInput body dto.SubscriberInput true "Subscriber"
// @Router /settings/subscribers [post]
// @Security ApiKeyAuth
func ApiCreateSubscriber(ctx *app.Ctx) fiber.Handler {
	return func(f *fiber.Ctx) error {
		newSubscriber := new(dto.SubscriberInput)

		if err := f.BodyParser(newSubscriber); err != nil {
			return f.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
		}

		newConfig := ctx.Config.GetCfg()

		_, alreadyExists := lo.Find(newConfig.Subscribers, func(s model.Subscriber) bool {
			return s.ApiaryId == newSubscriber.ApiaryId && s.HiveId == newSubscriber.HiveId
		})

		if !alreadyExists {
			subscriber := model.Subscriber{ApiaryId: newSubscriber.ApiaryId, HiveId: newSubscriber.HiveId}
			newConfig.Subscribers = append(newConfig.Subscribers, subscriber)
			ctx.Config.Update(newConfig)
		}

		return f.JSON(dto.ResponseHTTP{
			Success: true,
			Message: "Registered subscriber for creation",
			Data:    ctx.Config.GetCfg().Subscribers,
		})
	}
}
