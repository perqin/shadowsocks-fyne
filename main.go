package main

import (
	customWidget "SSD-Go/widget"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// SSD Data
type SsdData struct {
	// Required fields
	Airport    string          `json:"airport"`
	Port       int             `json:"port"`
	Encryption string          `json:"encryption"`
	Password   string          `json:"password"`
	Servers    []SsdDataServer `json:"servers"`
	// Extended fields
	Plugin        string `json:"plugin"`
	PluginOptions string `json:"plugin_options"`
	// Optional fields
	TrafficUsed  float64 `json:"traffic_used"`
	TrafficTotal float64 `json:"traffic_total"`
	Expiry       string  `json:"expiry"`
	Url          string  `json:"url"`
}

type SsdDataServer struct {
	// Required fields
	Server string `json:"server"`
	// Extended fields
	Port          int    `json:"port"`
	Encryption    string `json:"encryption"`
	Password      string `json:"password"`
	Plugin        string `json:"plugin"`
	PluginOptions string `json:"plugin_options"`
	// Optional fields
	Id      int     `json:"id"`
	Remarks string  `json:"remarks"`
	Ratio   float64 `json:"ratio"`
}

// UI widgets that should be updated on data changed
var window fyne.Window
var emptyTab = widget.NewTabItem("", fyne.NewContainer())
var subscriptionsTabs *widget.TabContainer
var noSubscriptionsObject fyne.CanvasObject

func updateTabs() {
	if len(config.Subscriptions) == 0 {
		subscriptionsTabs.Hide()
		noSubscriptionsObject.Show()
		//window.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("No subscriptions.")))
	} else {
		// Clear old tabs
		for len(subscriptionsTabs.Items) != 0 {
			subscriptionsTabs.RemoveIndex(0)
		}
		for _, subscription := range config.Subscriptions {
			list := make([]fyne.CanvasObject, 0)
			for _, profile := range subscription.Profiles {
				list = append(list, widget.NewLabel(fmt.Sprintf("%s (%s)", profile.Server, profile.Remark)))
			}
			//subscriptionsTabs.Append(widget.NewTabItem(subscription.Name, widget.NewLabel(fmt.Sprintf("URL: %s", subscription.Url))))
			subscriptionsTabs.Append(widget.NewTabItem(subscription.Name, fyne.NewContainerWithLayout(customWidget.NewVWeightLayout(),
				widget.NewLabel(fmt.Sprintf("URL: %s", subscription.Url)),
				// TODO: Add list
				customWidget.NewWeightedItem(widget.NewScrollContainer(), 1),
			)))
		}
		//window.SetContent(subscriptionsTabs)
		subscriptionsTabs.Show()
		noSubscriptionsObject.Hide()
	}
}

func newMainMenu(window fyne.Window) *fyne.MainMenu {
	subscriptionsMenu := fyne.NewMenu("Subscriptions",
		fyne.NewMenuItem("Add", func() {
			// Show dialog for new subscription
			urlEntry := widget.NewEntry()
			form := widget.NewForm(
				widget.NewFormItem("Url", urlEntry))
			dialog.ShowCustomConfirm("Add subscription", "Add", "Cancel", form, func(confirmed bool) {
				if confirmed {
					AddSubscription(Subscription{Url: urlEntry.Text})
				}
			}, window)
		}))
	return fyne.NewMainMenu(
		subscriptionsMenu)
}

func buildToolbar() fyne.CanvasObject {
	addIcon, err := fyne.LoadResourceFromPath("./res/add.png")
	if err != nil {
		log.Fatalln(err)
	}
	return widget.NewToolbar(
		widget.NewToolbarAction(addIcon, func() {
			// Show dialog for new subscription
			urlEntry := widget.NewEntry()
			form := widget.NewForm(
				widget.NewFormItem("Url", urlEntry))
			dialog.ShowCustomConfirm("Add subscription", "Add", "Cancel", form, func(confirmed bool) {
				if confirmed {
					AddSubscription(Subscription{Url: urlEntry.Text})
				}
			}, window)
		}),
		widget.NewToolbarAction(addIcon, func() {
			// Update
			currentSubscription := CurrentSubscription()
			go func() {
				response, err := http.Get(currentSubscription.Url)
				if err != nil {
					log.Println(err)
					return
				}
				responseBody, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Println(err)
					return
				}
				content := string(responseBody)
				//log.Println(content)
				if !strings.HasPrefix(content, "ssd://") {
					log.Println("Unsupported protocol")
					return
				}
				content = content[len("ssd://"):]
				//log.Println(content[24640])
				//log.Println(len(content))
				decoded, err := base64.RawStdEncoding.DecodeString(content)
				if err != nil {
					log.Println(err)
					return
				}
				ssdData := SsdData{}
				err = json.Unmarshal(decoded, &ssdData)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("Successfully updated [%s] and found %d servers\n", ssdData.Airport, len(ssdData.Servers))
				// Now convert into config
				profiles := make([]Profile, 0)
				for _, server := range ssdData.Servers {
					profiles = append(profiles, Profile{
						Server: server.Server,
						Remark: server.Remarks,
					})
				}
				// TODO: Fix!
				config.Subscriptions[0].Profiles = profiles
				err = SaveConfig()
				if err != nil {
					log.Println(err)
					return
				}
				log.Println("Successfully saved")
				updateTabs()
			}()
		}))
}

func main() {
	// TODO: Use Preference to r/w config
	if err := LoadConfig(); err != nil {
		log.Fatalln(err)
	}
	if len(config.Subscriptions) > 0 {
		SetCurrentSubscription(config.Subscriptions[0])
	}

	application := app.New()

	window = application.NewWindow("SSD Go")
	window.SetMainMenu(newMainMenu(window))
	window.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})
	window.CenterOnScreen()

	// At least 1 tab is required, or index out of range is thrown
	subscriptionsTabs = widget.NewTabContainer(emptyTab)
	subscriptionsTabs.SetTabLocation(widget.TabLocationLeading)
	noSubscriptionsObject = fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("No subscriptions."))
	subscriptionsOrNothing := fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		subscriptionsTabs,
		noSubscriptionsObject)
	window.SetContent(fyne.NewContainerWithLayout(customWidget.NewVWeightLayout(),
		buildToolbar(),
		customWidget.NewWeightedItem(subscriptionsOrNothing, 1)))

	// Update UI
	updateTabs()

	window.ShowAndRun()
}
