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

func SaveSubscription(index int, subscription Subscription) {
	if index >= 0 && index < len(config.Subscriptions) {
		config.Subscriptions[index] = subscription
	}
	if err := SaveConfig(); err != nil {
		log.Println(err)
	}
}

func removeSubscription(index int) error {
	if index == 0 {
		return errors.New("the Subscription for custom profiles cannot be removed")
	}
	if index < 0 || index >= len(config.Subscriptions) {
		return errors.New(fmt.Sprintf("invalid index %d", index))
	}
	config.Subscriptions = append(config.Subscriptions[:index], config.Subscriptions[index+1:]...)
	return SaveConfig()
}
