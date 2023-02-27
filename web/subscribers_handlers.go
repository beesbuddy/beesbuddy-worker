package web

import (
	"github.com/beesbuddy/beesbuddy-worker/dto"
	p "github.com/beesbuddy/beesbuddy-worker/pref"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// Get active subscribers
// @Summary Get active subscribers
// @Description Get subscribers
// @Tags preferences
// @Produce json
// @Success 200 {object} dto.ResponseHTTP{data=[]dto.SubscriberOutput}
// @Failure 503 {object} dto.ResponseHTTP{}
// @Router /preferences/subscribers [get]
// @Security ApiKeyAuth
func (w *webComponent) apiGetSubscribers(f *fiber.Ctx) error {
	return f.JSON(dto.ResponseHTTP{
		Success: true,
		Message: "Successfully fetched subscribers",
		Data:    w.appCtx.Pref.GetConfig().Subscribers,
	})
}

// Post subscriber
// @Summary Create a new subscriber
// @Description Create a subscriber
// @Tags preferences
// @Produce json
// @Success 200 {object} dto.ResponseHTTP{data=[]dto.SubscriberOutput}
// @Failure 503 {object} dto.ResponseHTTP{}
// @Param dto.SubscriberInput body dto.SubscriberInput true "Subscriber"
// @Router /preferences/subscribers [post]
// @Security ApiKeyAuth
func (w *webComponent) apiCreateSubscriber(f *fiber.Ctx) error {
	newSubscriber := &dto.SubscriberInput{}

	if err := f.BodyParser(newSubscriber); err != nil {
		return f.Status(fiber.StatusBadRequest).JSON(dto.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := w.createSubscriberForHive(*newSubscriber); err != nil {
		return f.Status(fiber.StatusConflict).JSON(dto.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	pref := w.appCtx.Pref

	return f.Status(fiber.StatusCreated).JSON(dto.ResponseHTTP{
		Success: true,
		Message: "registered subscriber for creation",
		Data: lo.Map(pref.GetConfig().Subscribers, func(s p.Subscriber, _ int) dto.SubscriberOutput {
			return dto.SubscriberOutput(s)
		}),
	})
}

// Delete subscribers for aiary
// @Summary Create a new subscriber
// @Description Create a subscriber
// @Tags preferences
// @Produce json
// @Success 200 {object} dto.ResponseHTTP{data=[]dto.SubscriberOutput}
// @Failure 503 {object} dto.ResponseHTTP{}
// @Param apiary_id path string true "Apiary Id"
// @Router /preferences/subscribers/{apiary_id} [delete]
// @Security ApiKeyAuth
func (w *webComponent) apiDeleteSubscriberForApiary(f *fiber.Ctx) error {
	apiaryId := f.Params("apiary_id")
	pref := w.appCtx.Pref

	w.deleteSubscriberForApiary(apiaryId)

	return f.JSON(dto.ResponseHTTP{
		Success: true,
		Message: "registered subscribers after deletion",
		Data: lo.Map(pref.GetConfig().Subscribers, func(s p.Subscriber, _ int) dto.SubscriberOutput {
			return dto.SubscriberOutput(s)
		}),
	})
}

// Delete subscribers for aiary
// @Summary Create a new subscriber
// @Description Create a subscriber
// @Tags preferences
// @Produce json
// @Success 200 {object} dto.ResponseHTTP{data=[]dto.SubscriberOutput}
// @Failure 503 {object} dto.ResponseHTTP{}
// @Param apiary_id path string true "Apiary Id"
// @Param hive_id path string true "Hive Id"
// @Router /preferences/subscribers/{apiary_id}/{hive_id} [delete]
// @Security ApiKeyAuth
func (w *webComponent) apiDeleteSubscriberForHive(f *fiber.Ctx) error {
	apiaryId := f.Params("apiary_id")
	hiveId := f.Params("hive_id")
	pref := w.appCtx.Pref

	w.deleteSubscriberForHive(apiaryId, hiveId)

	return f.JSON(dto.ResponseHTTP{
		Success: true,
		Message: "registered subscribers after deletion",
		Data: lo.Map(pref.GetConfig().Subscribers, func(s p.Subscriber, _ int) dto.SubscriberOutput {
			return dto.SubscriberOutput(s)
		}),
	})
}
