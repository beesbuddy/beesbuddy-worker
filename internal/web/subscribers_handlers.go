package web

import (
	"time"

	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/dto"
	p "github.com/beesbuddy/beesbuddy-worker/internal/pref"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// Get active subscribers
// @Summary Get active subscribers
// @Description Get subscribers
// @Tags settings
// @Produce json
// @Success 200 {object} dto.ResponseHTTP{data=[]dto.SubscriberOutput}
// @Failure 503 {object} dto.ResponseHTTP{}
// @Router /settings/subscribers [get]
// @Security ApiKeyAuth
func ApiGetSubscribers(ctx *app.Ctx) fiber.Handler {
	return func(f *fiber.Ctx) error {
		return f.JSON(dto.ResponseHTTP{
			Success: true,
			Message: "Successfully fetched subscribers",
			Data:    ctx.Pref.GetConfig().Subscribers,
		})
	}
}

// Post subscriber
// @Summary Create a new subscriber
// @Description Create a subscriber
// @Tags settings
// @Produce json
// @Success 200 {object} dto.ResponseHTTP{data=[]dto.SubscriberOutput}
// @Failure 503 {object} dto.ResponseHTTP{}
// @Param dto.SubscriberInput body dto.SubscriberInput true "Subscriber"
// @Router /settings/subscribers [post]
// @Security ApiKeyAuth
func ApiCreateSubscriber(ctx *app.Ctx) fiber.Handler {
	return func(f *fiber.Ctx) error {
		newSubscriber := new(dto.SubscriberInput)
		pref := ctx.Pref

		if err := f.BodyParser(newSubscriber); err != nil {
			return f.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
		}

		newConfig := pref.GetConfig()

		_, alreadyExists := lo.Find(newConfig.Subscribers, func(s p.Subscriber) bool {
			return s.ApiaryId == newSubscriber.ApiaryId && s.HiveId == newSubscriber.HiveId
		})

		if !alreadyExists {
			subscriber := p.Subscriber{ApiaryId: newSubscriber.ApiaryId, HiveId: newSubscriber.HiveId, CreatedAt: time.Now()}
			newConfig.Subscribers = append(newConfig.Subscribers, subscriber)
			pref.UpdateConfig(newConfig)
		}

		return f.JSON(dto.ResponseHTTP{
			Success: true,
			Message: "Registered subscriber for creation",
			Data: lo.Map(pref.GetConfig().Subscribers, func(s p.Subscriber, _ int) dto.SubscriberOutput {
				return dto.SubscriberOutput(s)
			}),
		})
	}
}
