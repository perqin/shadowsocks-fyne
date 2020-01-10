package main

import "log"

var currentSubscription Subscription

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

func SetCurrentSubscription(subscription Subscription) {
	currentSubscription = subscription
}

func CurrentSubscription() Subscription {
	return currentSubscription
}
