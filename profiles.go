package main

import (
	"errors"
	"log"
)

type Profile struct {
	Name       string `json:"name"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Method     string `json:"method"`
	Password   string `json:"password"`
	Acl        string `json:"acl"`
}

func selectCurrentProfile(profile, subscription int) {
	config.CurrentProfile = profile
	config.CurrentProfileSubscription = subscription
	if err := SaveConfig(); err != nil {
		log.Println(err)
	}
}

func addProfile(si int, p Profile) error {
	if si != 0 {
		return errors.New("only custom profile can be edit")
	}
	config.Subscriptions[si].Profiles = append(config.Subscriptions[si].Profiles, p)
	return SaveConfig()
}

func saveProfile(si, pi int, p Profile) error {
	if si != 0 {
		return errors.New("only custom profile can be edit")
	}
	config.Subscriptions[si].Profiles[pi] = p
	return SaveConfig()
}
