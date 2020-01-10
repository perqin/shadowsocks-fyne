package main

import (
	"errors"
	"fmt"
	"log"
)

type Subscription struct {
	Url      string    `json:"url"`
	Name     string    `json:"name"`
	Profiles []Profile `json:"profiles"`
}

func AddSubscription(subscription Subscription) {
	config.Subscriptions = append(config.Subscriptions, subscription)
	if err := SaveConfig(); err != nil {
		log.Println(err)
	}
}

func UpdateSubscription(subscription Subscription) error {
	for i, s := range config.Subscriptions {
		if s.Url == subscription.Url {
			config.Subscriptions[i] = subscription
			return SaveConfig()
		}
	}
	return errors.New(fmt.Sprintf("fail to update Subscription with Url:%s", subscription.Url))
}
