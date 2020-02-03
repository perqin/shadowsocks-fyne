package main

import (
	"os"
	"path/filepath"
)

func getConfigDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	desktopConfig := filepath.Join(filepath.Join(homeDir, "AppData"), "Roaming")
	return filepath.Join(desktopConfig, "fyne", applicationId), nil
}
