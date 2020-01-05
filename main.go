package main

import (
	customWidget "SSD-Go/widget"
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
			//subscriptionsTabs.Append(widget.NewTabItem(subscription.Name, widget.NewLabel(fmt.Sprintf("URL: %s", subscription.Url))))
			subscriptionsTabs.Append(widget.NewTabItem(subscription.Name, widget.NewVBox(
				widget.NewLabel(fmt.Sprintf("URL: %s", subscription.Url)),
				customWidget.NewSpacerContainer(false, widget.NewLabel("This is in a SpacerContainer")),
				widget.NewLabel("Bottom"))))
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
		}))
}

func main() {
	// TODO: Use Preference to r/w config
	if err := LoadConfig(); err != nil {
		log.Fatalln(err)
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
