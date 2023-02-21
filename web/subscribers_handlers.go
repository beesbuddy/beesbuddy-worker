package web

import (
	"time"

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
	pref := w.appCtx.Pref

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
	newConfig := w.appCtx.Pref.GetConfig()
	subscribers := newConfig.Subscribers

	newConfig.Subscribers = lo.Filter(subscribers, func(item p.Subscriber, _ int) bool {
		return item.ApiaryId != apiaryId
	})

	pref.UpdateConfig(newConfig)

	return f.JSON(dto.ResponseHTTP{
		Success: true,
		Message: "Registered subscribers after deletion",
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
	newConfig := w.appCtx.Pref.GetConfig()
	subscribers := newConfig.Subscribers

	for index, item := range subscribers {
		if item.ApiaryId == apiaryId && item.HiveId == hiveId {
			subscribers = append(subscribers[:index], subscribers[index+1:]...)
		}
	}

	newConfig.Subscribers = subscribers
	pref.UpdateConfig(newConfig)

	return f.JSON(dto.ResponseHTTP{
		Success: true,
		Message: "Registered subscribers after deletion",
		Data: lo.Map(pref.GetConfig().Subscribers, func(s p.Subscriber, _ int) dto.SubscriberOutput {
			return dto.SubscriberOutput(s)
		}),
	})
}
