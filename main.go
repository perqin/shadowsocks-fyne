package main

import (
	customWidget "SSD-Go/widget"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/perqin/go-shadowsocks2"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
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
var subscriptionsTabs *widget.TabContainer

func buildServerList(index int, subscription Subscription) fyne.CanvasObject {
	list := make([]fyne.CanvasObject, 0)
	itemHeight := 64
	for indexP, profile := range subscription.Profiles {
		indexP := indexP
		selectionIndicator := customWidget.NewColoredBox()
		selectionIndicator.Resize(fyne.Size{Width: 8, Height: itemHeight})
		selectionIndicator.SetBackgroundColor(color.Transparent)
		if index == config.CurrentProfileSubscription && indexP == config.CurrentProfile {
			selectionIndicator.SetBackgroundColor(fyne.CurrentApp().Settings().Theme().PrimaryColor())
		}
		listItemContent := customWidget.NewTwoLineListItem(profile.Name, fmt.Sprintf("%s:%d", profile.Server, profile.ServerPort))
		list = append(list, customWidget.NewTappableWidget(
			fyne.NewContainerWithLayout(customWidget.NewStackLayout(), selectionIndicator, listItemContent),
			func() {
				selectProfile(index, indexP, selectionIndicator)
			}))
	}
	return widget.NewScrollContainer(widget.NewVBox(list...))
}

func buildTabs() []*widget.TabItem {
	var item []*widget.TabItem
	for indexS, subscription := range config.Subscriptions {
		item = append(item, widget.NewTabItem(subscription.Name, buildServerList(indexS, subscription)))
	}
	return item
}

func updateTabs() {
	// Ensure refresh not crash
	subscriptionsTabs.SelectTabIndex(0)
	// Ensure same item count
	for len(subscriptionsTabs.Items) != len(config.Subscriptions) {
		if len(subscriptionsTabs.Items) < len(config.Subscriptions) {
			subscriptionsTabs.Append(widget.NewTabItem("", widget.NewHBox()))
		} else {
			subscriptionsTabs.RemoveIndex(0)
		}
	}
	// Ensure correct UI
	for i, s := range config.Subscriptions {
		tab := subscriptionsTabs.Items[i]
		tab.Text = s.Name
		tab.Content = buildServerList(i, s)
	}
	// And refresh!
	subscriptionsTabs.Refresh()
}

var profile Profile
var selectedIndicator *customWidget.ColoredBox

func selectProfile(si, pi int, indicator *customWidget.ColoredBox) {
	if selectedIndicator != nil {
		selectedIndicator.SetBackgroundColor(color.Transparent)
	}
	selectCurrentProfile(si, pi)
	profile = config.Subscriptions[si].Profiles[pi]
	selectedIndicator = indicator
	selectedIndicator.SetBackgroundColor(fyne.CurrentApp().Settings().Theme().PrimaryColor())
}

func showEditSubscriptionDialog(index int) {
	if index == 0 {
		log.Println("The Subscription for custom profiles cannot be edited.")
		return
	}
	edit := index != -1
	if edit && (index < 0 || index >= len(config.Subscriptions)) {
		log.Println(fmt.Sprintf("Invalid index %d", index))
		return
	}
	var subscription Subscription
	if edit {
		subscription = config.Subscriptions[index]
	}
	// Show dialog for new subscription
	urlEntry := widget.NewEntry()
	if edit {
		urlEntry.Text = subscription.Url
	}
	form := widget.NewForm(
		widget.NewFormItem("Url", urlEntry))
	dialog.ShowCustomConfirm("Add subscription", "Add", "Cancel", form, func(confirmed bool) {
		if confirmed {
			if edit {
				// Update subscription
				subscription.Url = urlEntry.Text
				SaveSubscription(index, subscription)
			} else {
				// Add new one
				AddSubscription(Subscription{
					Name: "(Untitled)",
					Url:  urlEntry.Text,
				})
			}
			updateTabs()
		}
	}, window)

}

func onAddSubscriptionAction() {
	showEditSubscriptionDialog(-1)
}

func onRefreshAction() {
	// Update
	index := subscriptionsTabs.CurrentTabIndex()
	size := len(config.Subscriptions)
	if index < 0 || index >= size {
		customWidget.Toast("Index out of bound")
		return
	}
	url := config.Subscriptions[index].Url
	go func() {
		directClient := http.Client{Transport: &http.Transport{Proxy: nil}}
		response, err := directClient.Get(url)
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
		protocolPrefix := "ssd://"
		if !strings.HasPrefix(content, protocolPrefix) {
			log.Println("Unsupported protocol")
			return
		}
		content = content[len(protocolPrefix):]
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
		newSubscription := Subscription{
			Url:      url,
			Name:     ssdData.Airport,
			Profiles: nil,
		}
		for _, server := range ssdData.Servers {
			var method string
			if server.Encryption != "" {
				method = server.Encryption
			} else {
				method = ssdData.Encryption
			}
			var password string
			if server.Password != "" {
				password = server.Password
			} else {
				password = ssdData.Password
			}
			newSubscription.Profiles = append(newSubscription.Profiles, Profile{
				Name:       server.Remarks,
				Server:     server.Server,
				ServerPort: server.Port,
				Method:     method,
				Password:   password,
			})
		}
		// TODO: Should I update data and UI on dedicated thread?
		if err = UpdateSubscription(newSubscription); err != nil {
			log.Println(err)
			return
		}
		log.Println("Successfully saved")
		updateTabs()
	}()
}

func onEditSubscriptionAction() {
	showEditSubscriptionDialog(config.CurrentProfileSubscription)
}

func onRemoveSubscriptionAction() {
	subscriptionIndex := subscriptionsTabs.CurrentTabIndex()
	err := removeSubscription(subscriptionIndex)
	if err == nil {
		subscriptionsTabs.SelectTabIndex(0)
		updateTabs()
	} else {
		log.Println(err)
	}
}

func onRunAction() {
	client := fmt.Sprintf("ss://%s:%s@%s:%d", profile.Method, profile.Password, profile.Server, profile.ServerPort)
	log.Printf("onRunAction client:%s\n", client)
	if err := runShadowsocks(shadowsocks2.Flags{
		Client: client,
		Socks:  "127.0.0.1:2080",
	}); err != nil {
		customWidget.Toast(fmt.Sprintf("Fail to stop: %v\n", err))
	}
}

func onStopAction() {
	if err := stopShadowsocks(); err != nil {
		customWidget.Toast(fmt.Sprintf("Fail to stop: %v\n", err))
	}
}

func onEditProfileAction() {
	// TODO
}

func buildToolbar() fyne.CanvasObject {
	return widget.NewToolbar(
		widget.NewToolbarAction(addIcon, onAddSubscriptionAction),
		widget.NewToolbarAction(refreshIcon, onRefreshAction),
		widget.NewToolbarAction(editSubscriptionIcon, onEditSubscriptionAction),
		widget.NewToolbarAction(deleteIcon, onRemoveSubscriptionAction),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(playIcon, onRunAction),
		widget.NewToolbarAction(stopIcon, onStopAction),
		widget.NewToolbarAction(editProfileIcon, onEditProfileAction))
}

func main() {
	// TODO: Use Preference to r/w config
	if err := LoadConfig(); err != nil {
		log.Fatalln(err)
	}
	shadowsocks2.SetConfig(shadowsocks2.Config{
		Verbose:    true,
		UDPTimeout: time.Minute * 5,
	})

	application := app.New()

	window = application.NewWindow("SSD Go")
	window.Resize(fyne.Size{
		Width:  1200,
		Height: 800,
	})
	window.CenterOnScreen()

	// At least 1 tab is required, or index out of range is thrown
	subscriptionsTabs = widget.NewTabContainer(buildTabs()...)
	subscriptionsTabs.SetTabLocation(widget.TabLocationLeading)
	window.SetContent(fyne.NewContainerWithLayout(customWidget.NewVWeightLayout(),
		buildToolbar(),
		customWidget.NewWeightedItem(subscriptionsTabs, 1)))

	window.ShowAndRun()

	// Cleanup here
	log.Println("Exiting")
	_ = stopShadowsocks()
}
