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

// getConfigFilePath returns the path of the config file. It ensures that the parent directory is created.
// TODO: Migrate from `getConfigDirPath` to Preference API:
//  > @andy.xyz: The plan is that we will create a way for developers to ask for a new storage file so the library can
//  > track them - that way we will be able to keep your app data in sync across computers eventually
//  > The preferences api was a start and we have more to do in that space
//  See https://gophers.slack.com/archives/CB4QUBXGQ/p1580730440194900.
func getConfigFilePath() (string, error) {
	dir, err := getConfigDirPath()
	if err != nil {
		return "", err
	}
	_, err = os.Stat(dir)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0700); err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, "config.json"), nil
}
