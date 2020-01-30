package main

import "log"

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
