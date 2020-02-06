package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/perqin/go-shadowsocks2"
	"github.com/perqin/shadowsocks-fyne/material"
	"github.com/perqin/shadowsocks-fyne/resources"
	customWidget "github.com/perqin/shadowsocks-fyne/widget"
	"log"
	"strconv"
)

const toolbarHeight = 56

var mainWindow fyne.Window
var drawer *subscriptionDrawer
var titleLabel *widget.Label
var updateSubscriptionButton *material.ToolbarActionButton
var editSubscriptionButton *material.ToolbarActionButton
var removeSubscriptionButton *material.ToolbarActionButton
var runButton *material.ToolbarActionButton
var stopButton *material.ToolbarActionButton
var addProfileButton *material.ToolbarActionButton
var editProfileButton *material.ToolbarActionButton
var settingsButton *material.ToolbarActionButton
var profileList *widget.Box

func showMainWindow() {
	if mainWindow != nil {
		mainWindow.RequestFocus()
		return
	}
	mainWindow = application.NewWindow(appName)
	mainWindow.SetOnClosed(func() {
		mainWindow = nil
	})
	mainWindow.CenterOnScreen()
	drawer = newSubscriptionDrawer(config.Subscriptions)
	titleLabel = widget.NewLabel("Title")
	updateSubscriptionButton = material.NewToolbarActionButton(resources.RefreshPng, onRefreshAction)
	editSubscriptionButton = material.NewToolbarActionButton(resources.EditsubscriptionPng, onEditSubscriptionAction)
	removeSubscriptionButton = material.NewToolbarActionButton(resources.DeletePng, onRemoveSubscriptionAction)
	runButton = material.NewToolbarActionButton(resources.PlayPng, onRunAction)
	stopButton = material.NewToolbarActionButton(resources.StopPng, onStopAction)
	addProfileButton = material.NewToolbarActionButton(resources.AddprofilePng, onAddProfileAction)
	editProfileButton = material.NewToolbarActionButton(resources.EditprofilePng, onEditProfileAction)
	settingsButton = material.NewToolbarActionButton(resources.SettingsPng, onSettingsAction)
	profileList = widget.NewVBox()
	selectSubscription(config.CurrentProfileSubscription)
	mainWindow.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewRectangle(application.Settings().Theme().(*material.Theme).WindowBackground),
		fyne.NewContainerWithLayout(&mainWindowLayout{},
			drawer,
			material.NewToolbar(titleLabel, updateSubscriptionButton, editSubscriptionButton, removeSubscriptionButton,
				runButton, stopButton, addProfileButton, editProfileButton, settingsButton),
			widget.NewScrollContainer(profileList))))
	mainWindow.Show()
}

func refreshAll(resetSelection bool) {
	subscriptions := config.Subscriptions
	titles := make([]string, len(subscriptions)-1)
	for i, s := range subscriptions[1:] {
		titles[i] = s.Name
	}
	drawer.SubscriptionTitles = titles
	drawer.Refresh()
	if drawer.SelectedIndex < 0 || drawer.SelectedIndex >= len(config.Subscriptions) {
		resetSelection = true
	}
	if resetSelection {
		selectSubscription(0)
	} else {
		selectSubscription(drawer.SelectedIndex)
	}
}

func selectSubscription(index int) {
	drawer.SelectedIndex = index
	drawer.Refresh()
	titleLabel.SetText(config.Subscriptions[index].Name)
	profileList.Children = buildServerList(index, config.Subscriptions[index])
	profileList.Refresh()
}

// UI widgets that should be updated on data changed
var subscriptionsTabs *widget.TabContainer
var currentProfileWidget *ProfileWidget

func buildServerList(index int, subscription Subscription) []fyne.CanvasObject {
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
			currentProfileWidget = profileWidget
		}
		profileWidget.Refresh()
		list = append(list, profileWidget)
	}
	return list
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
	nameEntry := widget.NewEntry()
	urlEntry := widget.NewEntry()
	if edit {
		nameEntry.SetText(subscription.Name)
		urlEntry.SetText(subscription.Url)
	}
	form := widget.NewForm(
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Url", urlEntry))
	dialog.ShowCustomConfirm("Add subscription", "Add", "Cancel", form, func(confirmed bool) {
		if confirmed {
			if edit {
				// Update subscription
				subscription.Name = nameEntry.Text
				subscription.Url = urlEntry.Text
				SaveSubscription(index, subscription)
				refreshAll(false)
			} else {
				// Add new one
				AddSubscription(Subscription{
					Name: nameEntry.Text,
					Url:  urlEntry.Text,
				})
				refreshAll(true)
			}
		}
	}, mainWindow)

}

func onRefreshAction() {
	go func() {
		err := UpdateSubscription(drawer.SelectedIndex)
		if err != nil {
			customWidget.Toast(err.Error())
			return
		}
		refreshAll(false)
	}()
}

func onEditSubscriptionAction() {
	showEditSubscriptionDialog(config.CurrentProfileSubscription)
}

func onRemoveSubscriptionAction() {
	subscriptionIndex := drawer.SelectedIndex
	err := removeSubscription(subscriptionIndex)
	if err != nil {
		log.Println(err)
		return
	}
	refreshAll(false)
}

func onRunAction() {
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
		Socks:  fmt.Sprintf("%s:%d", config.LocalAddress, config.LocalPort),
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
				refreshAll(false)
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

func onSettingsAction() {
	localAddressEntry := widget.NewEntry()
	localAddressEntry.SetText(config.LocalAddress)
	localPortEntry := widget.NewEntry()
	localPortEntry.SetText(strconv.Itoa(config.LocalPort))
	form := widget.NewForm(
		widget.NewFormItem("Local address", localAddressEntry),
		widget.NewFormItem("Local port", localPortEntry))
	dialog.ShowCustomConfirm("Settings", "Save", "Cancel", form, func(confirmed bool) {
		if confirmed {
			// Update settings
			localAddress := localAddressEntry.Text
			localPort, err := strconv.Atoi(localPortEntry.Text)
			if err != nil {
				log.Println(err)
				return
			}
			config.LocalAddress = localAddress
			config.LocalPort = localPort
			if err = SaveConfig(); err != nil {
				log.Println(err)
			}
		}
	}, mainWindow)
}

type mainWindowLayout struct {
}

func (l *mainWindowLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	drawerWidth := objects[0].MinSize().Width
	// Drawer
	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(fyne.NewSize(drawerWidth, size.Height))
	// Toolbar background
	objects[1].Move(fyne.NewPos(drawerWidth, 0))
	objects[1].Resize(fyne.NewSize(size.Width-drawerWidth, toolbarHeight))
	// Content
	objects[2].Move(fyne.NewPos(drawerWidth, toolbarHeight))
	objects[2].Resize(fyne.NewSize(size.Width-drawerWidth, size.Height-toolbarHeight))
}

func (l *mainWindowLayout) MinSize(objects []fyne.CanvasObject) (size fyne.Size) {
	drawerMinSize := objects[0].MinSize()
	toolbarMinWidth := objects[1].MinSize().Width
	size.Width = drawerMinSize.Width + fyne.Max(toolbarMinWidth, 640)
	size.Height = fyne.Max(drawerMinSize.Height, toolbarHeight)
	return
}
