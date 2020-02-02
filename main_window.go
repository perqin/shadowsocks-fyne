package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/perqin/go-shadowsocks2"
	customWidget "github.com/perqin/shadowsocks-fyne/widget"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var mainWindow fyne.Window

func showMainWindow() {
	if mainWindow != nil {
		mainWindow.RequestFocus()
		return
	}
	mainWindow = application.NewWindow(appName)
	mainWindow.SetOnClosed(func() {
		mainWindow = nil
	})
	mainWindow.Resize(fyne.Size{
		Width:  1200,
		Height: 800,
	})
	mainWindow.CenterOnScreen()
	subscriptionsTabs = widget.NewTabContainer(buildTabs()...)
	subscriptionsTabs.SetTabLocation(widget.TabLocationLeading)
	mainWindow.SetContent(fyne.NewContainerWithLayout(customWidget.NewVWeightLayout(),
		buildToolbar(),
		customWidget.NewWeightedItem(subscriptionsTabs, 1)))
	mainWindow.Show()
}

// TODO: Drop direct ssd:// support and follow the subscription protocol by official Shadowsocks project
//  Maybe the support for other protocols can be met with the help of subconverter project.
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
var subscriptionsTabs *widget.TabContainer
var currentProfileWidget *ProfileWidget

func buildServerList(index int, subscription Subscription) fyne.CanvasObject {
	list := make([]fyne.CanvasObject, 0)
	for indexP, profile := range subscription.Profiles {
		indexP := indexP
		profileWidget := NewProfileWidget()
		profileWidget.OnTapped = func() {
			if currentProfileWidget != nil {
				currentProfileWidget.Selected = false
				currentProfileWidget.Refresh()
			}
			currentProfileWidget = profileWidget
			selectCurrentProfile(indexP, index)
		}
		profileWidget.Profile = profile
		if index == config.CurrentProfileSubscription && indexP == config.CurrentProfile {
			profileWidget.Selected = true
		}
		profileWidget.Refresh()
		list = append(list, profileWidget)
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
	}, mainWindow)

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
	// TODO
	si := config.CurrentProfileSubscription
	if si < 0 || si >= len(config.Subscriptions) {
		log.Printf("Invalid index %d\n", si)
		return
	}
	pi := config.CurrentProfile
	if pi < 0 || pi >= len(config.Subscriptions[si].Profiles) {
		log.Printf("Invalid index %d %d\n", si, pi)
		return
	}
	profile := config.Subscriptions[si].Profiles[pi]
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

func showEditProfileDialog(edit bool) {
	var si, pi int
	var profile Profile
	if edit {
		si = config.CurrentProfileSubscription
		pi = config.CurrentProfile
		if si != subscriptionsTabs.CurrentTabIndex() {
			log.Println("Please switch to the tab containing the profile you want to edit.")
			return
		}
		if si < 0 || si >= len(config.Subscriptions) {
			log.Printf("Invalid subscription index %d\n", si)
			return
		}
		if pi < 0 || pi >= len(config.Subscriptions[si].Profiles) {
			log.Printf("Invalid profile index %d\n", pi)
			return
		}
		profile = config.Subscriptions[si].Profiles[pi]
	} else {
		si = subscriptionsTabs.CurrentTabIndex()
	}
	// Can show dialog now
	nameEntry := widget.NewEntry()
	serverEntry := widget.NewEntry()
	portEntry := widget.NewEntry()
	methodEntry := widget.NewEntry()
	passwordEntry := widget.NewEntry()
	if edit {
		nameEntry.Text = profile.Name
		serverEntry.Text = profile.Server
		portEntry.Text = fmt.Sprint(profile.ServerPort)
		methodEntry.Text = profile.Method
		passwordEntry.Text = profile.Password
	}
	form := widget.NewForm(
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Server", serverEntry),
		widget.NewFormItem("Port", portEntry),
		widget.NewFormItem("Method", methodEntry),
		widget.NewFormItem("Password", passwordEntry))
	dialog.ShowCustomConfirm("Edit Profile", "Save", "Cancel", form, func(confirmed bool) {
		if confirmed {
			profile.Name = nameEntry.Text
			profile.Server = serverEntry.Text
			// TODO: Validate
			profile.ServerPort, _ = strconv.Atoi(portEntry.Text)
			profile.Method = methodEntry.Text
			profile.Password = passwordEntry.Text
			var err error
			if edit {
				if si != 0 {
					log.Println("Only custom profile can be edited.")
					return
				}
				err = saveProfile(si, pi, profile)
			} else {
				err = addProfile(si, profile)
			}
			if err != nil {
				log.Println(err)
			} else {
				// TODO: Performance issue
				updateTabs()
			}
		}
	}, mainWindow)
}

func onAddProfileAction() {
	showEditProfileDialog(false)
}

func onEditProfileAction() {
	showEditProfileDialog(true)
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
		widget.NewToolbarAction(addProfileIcon, onAddProfileAction),
		widget.NewToolbarAction(editProfileIcon, onEditProfileAction))
}
