package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Version                    int            `json:"version"`
	Subscriptions              []Subscription `json:"subscriptions"`
	LocalAddress               string         `json:"local_address"`
	LocalPort                  int            `json:"local_port"`
	CurrentProfile             int            `json:"current_profile"`
	CurrentProfileSubscription int            `json:"current_profile_subscription"`
}

var config Config

func LoadConfig() error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Fail to read file %s, creating default config", path)
		config.Version = 1
		config.Subscriptions = []Subscription{
			{
				Url:      "",
				Name:     "Custom",
				Profiles: nil,
			},
		}
		config.LocalAddress = "127.0.0.1"
		config.LocalPort = 1180
		config.CurrentProfile = -1
		config.CurrentProfileSubscription = -1
		return SaveConfig()
	}
	return json.Unmarshal(bytes, &config)
}

func SaveConfig() error {
	bytes, err := json.MarshalIndent(&config, "", "  ")
	if err != nil {
		return err
	}
	file, err := getConfigFile()
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}

func getConfigFile() (*os.File, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
}

func getConfigFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(userHomeDir, ".ssd-go")
	if err = os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}
