package main

import (
	"flag"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"github.com/perqin/go-shadowsocks2"
	"log"
	"time"
)

var flags struct {
	Minimize bool
}

func main() {
	flag.BoolVar(&flags.Minimize, "minimize", false, "start and minimize to system tray")
	flag.Parse()

	// TODO: Use Preference to r/w config
	if err := LoadConfig(); err != nil {
		log.Fatalln(err)
	}
	shadowsocks2.SetConfig(shadowsocks2.Config{
		Verbose:    true,
		UDPTimeout: time.Minute * 5,
	})

	setupSystemTray()

	startupFyneGui()
}

var appName = "Shadowsocks Fyne"
var application fyne.App

func startupFyneGui() {
	// Setup GUI
	application = app.New()
	// A invisible window is required, otherwise the Fyne will exit after all windows are closed
	application.NewWindow("")
	// Show main window if needed
	if !flags.Minimize {
		showMainWindow()
	}

	application.Driver().Run()
}

// exitApp performs necessary cleanups before the process terminates.
func exitApp() {
	if mainWindow != nil {
		mainWindow.Close()
	}
	application.Quit()
	_ = stopShadowsocks()
}
