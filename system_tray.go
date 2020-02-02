package main

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

func setupSystemTray() {
	// systray.Run will block a goroutine
	go func() {
		// The system tray has the same lifecycle of the application
		systray.Run(func() {
			// onReady
			configureSystemTray()
		}, nil)
	}()
}

func configureSystemTray() {
	systray.SetTitle(appName)
	// TODO: Icon for app
	systray.SetIcon(icon.Data)
	showMenu := systray.AddMenuItem("Show", "Show main window")
	exitMenu := systray.AddMenuItem("Exit", "Exit application")
	// Menu handlers
	go func() {
		for {
			select {
			case <-showMenu.ClickedCh:
				showMainWindow()
			case <-exitMenu.ClickedCh:
				systray.Quit()
				exitApp()
			}
		}
	}()
}
