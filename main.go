package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"log"
)

// UI widgets that should be updated on data changed
var window fyne.Window
var tabs *widget.TabContainer

func updateTabs() {
	if len(config.Subscriptions) == 0 {
		window.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("No subscriptions.")))
	} else {
		for len(tabs.Items) != 0 {
			tabs.RemoveIndex(0)
		}
		for _, subscription := range config.Subscriptions {
			tabs.Append(widget.NewTabItem(subscription.Name, widget.NewLabel(fmt.Sprintf("URL: %s", subscription.Url))))
		}
		window.SetContent(tabs)
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

func main() {
	// TODO: Use Preference to r/w config
	if err := LoadConfig(); err != nil {
		log.Fatalln(err)
	}

	application := app.New()

	window = application.NewWindow("SSD Go")
	window.SetMainMenu(newMainMenu(window))
	// At least 1 tab is required, or index out of range is thrown
	tabs = widget.NewTabContainer(widget.NewTabItem("", widget.NewVBox()))
	tabs.SetTabLocation(widget.TabLocationLeading)
	window.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})
	window.CenterOnScreen()

	// Update UI
	updateTabs()

	window.ShowAndRun()
}
