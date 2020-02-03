package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type ssProfileData struct {
	Server     string `json:"server"`
	ServerPort int    `json:"server_port"`
	Method     string `json:"method"`
	Password   string `json:"password"`
	Remarks    string `json:"remarks"`
	Plugin     string `json:"plugin"`
	PluginOpts string `json:"plugin_opts"`
}

func UpdateSubscription(index int) error {
	if index < 0 || index >= len(config.Subscriptions) {
		return errors.New("index out of bound")
	}
	oldSubscription := config.Subscriptions[index]
	directClient := http.Client{Transport: &http.Transport{Proxy: nil}}
	response, err := directClient.Get(oldSubscription.Url)
	if err != nil {
		return err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	//content := string(responseBody)
	//log.Println(content)
	var dataList []ssProfileData
	if err = json.Unmarshal(responseBody, &dataList); err != nil {
		return err
	}
	log.Printf("Successfully updated [%s] and found %d servers\n", oldSubscription.Name, len(dataList))
	//// Now convert into config
	newSubscription := Subscription{
		Url:      oldSubscription.Url,
		Name:     oldSubscription.Name,
		Profiles: nil,
	}
	for _, server := range dataList {
		newSubscription.Profiles = append(newSubscription.Profiles, Profile{
			Name:       server.Remarks,
			Server:     server.Server,
			ServerPort: server.ServerPort,
			Method:     server.Method,
			Password:   server.Password,
		})
	}
	config.Subscriptions[index] = newSubscription
	if err = SaveConfig(); err != nil {
		return err
	}
	log.Println("Successfully saved")
	return nil
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
