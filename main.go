package main

import (
	"flag"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"github.com/perqin/go-shadowsocks2"
	"github.com/perqin/shadowsocks-fyne/material"
	"github.com/perqin/shadowsocks-fyne/resources"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	setupSignals()
	setupSystemTray()
	startupFyneGui()
}

var appName = "Shadowsocks Fyne"

func setupSignals() {
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Received SIGINT or SIGTERM")
		exitApp()
	}()
}

var application fyne.App
var applicationId = "shadowsocks-fyne"

func startupFyneGui() {
	// Setup GUI
	application = app.NewWithID(applicationId)
	application.SetIcon(resources.IconPng)
	application.Settings().SetTheme(material.NewLightTheme())
	// A invisible window is required, otherwise the Fyne will exit after all windows are closed
	application.NewWindow("")
	// Show main window if needed
	if !flags.Minimize {
		showMainWindow()
	}

	application.Driver().Run()
}

// exitApp performs necessary cleanups and then exit the application
func exitApp() {
	if mainWindow != nil {
		mainWindow.Close()
	}
	_ = stopShadowsocks()
	// Quit must be called on the last, because it will break the blocking main thread
	application.Quit()
}
