package main

import "log"

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
