package web

import (
	"fmt"
	"time"

	"github.com/beesbuddy/beesbuddy-worker/dto"
	p "github.com/beesbuddy/beesbuddy-worker/pref"
	"github.com/samber/lo"
)

func (w webComponent) createSubscriberForHive(subscriberInput dto.SubscriberInput) error {
	pref := w.appCtx.Pref
	newConfig := pref.GetConfig()

	_, alreadyExists := lo.Find(newConfig.Subscribers, func(s p.Subscriber) bool {
		return s.ApiaryId == subscriberInput.ApiaryId && s.HiveId == subscriberInput.HiveId
	})

	if !alreadyExists {
		subscriber := p.Subscriber{ApiaryId: subscriberInput.ApiaryId, HiveId: subscriberInput.HiveId, CreatedAt: time.Now()}
		newConfig.Subscribers = append(newConfig.Subscribers, subscriber)
		pref.UpdateConfig(newConfig)

		return nil
	}

	return fmt.Errorf("subscriber already exist")
}

func (w webComponent) deleteSubscriberForApiary(apiaryId string) error {
	pref := w.appCtx.Pref
	newConfig := pref.GetConfig()
	subscribers := newConfig.Subscribers

	newConfig.Subscribers = lo.Filter(subscribers, func(item p.Subscriber, _ int) bool {
		return item.ApiaryId != apiaryId
	})

	pref.UpdateConfig(newConfig)

	return nil
}

func (w webComponent) deleteSubscriberForHive(apiaryId, hiveId string) error {
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

	return nil
}
