package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/perqin/shadowsocks-fyne/material"
	"image/color"
)

const drawerWidth = 256
const drawerMinHeight = 640

func newSubscriptionDrawer(subscriptions []Subscription) *subscriptionDrawer {
	titles := make([]string, len(subscriptions)-1)
	for i, s := range subscriptions[1:] {
		titles[i] = s.Name
	}
	d := &subscriptionDrawer{SubscriptionTitles: titles}
	d.ExtendBaseWidget(d)
	return d
}

type subscriptionDrawer struct {
	widget.BaseWidget
	// SubscriptionTitles does NOT include the Custom Profiles item
	SubscriptionTitles []string
	// SelectedIndex includes the Custom Profiles
	SelectedIndex int
}

func (d *subscriptionDrawer) CreateRenderer() fyne.WidgetRenderer {
	d.ExtendBaseWidget(d)

	customProfilesItem := material.NewDrawerItem("Custom Profiles", func() {
		selectSubscription(0)
	})
	subtitle := widget.NewLabel("Subscription")
	list := []fyne.CanvasObject{customProfilesItem, subtitle}
	list = append(list, d.buildSubscriptionItems()...)
	r := &subscriptionDrawerRenderer{drawer: d, objects: []fyne.CanvasObject{
		widget.NewScrollContainer(widget.NewVBox(list...)),
		material.NewDrawerItem("Add", func() {
			showEditSubscriptionDialog(-1)
		}),
	}}
	return r
}

func (d *subscriptionDrawer) buildSubscriptionItems() []fyne.CanvasObject {
	var subscriptionItems []fyne.CanvasObject
	for i, t := range d.SubscriptionTitles {
		// Offset for Custom Profiles
		subscriptionIndex := i + 1
		button := material.NewDrawerItem(t, func() {
			selectSubscription(subscriptionIndex)
		})
		button.Selected = d.SelectedIndex == subscriptionIndex
		button.Refresh()
		subscriptionItems = append(subscriptionItems, button)
	}
	return subscriptionItems
}

type subscriptionDrawerRenderer struct {
	drawer *subscriptionDrawer
	// objects contains the following items:
	// - scroller view containing titles and subtitle
	// - Add button
	objects []fyne.CanvasObject
}

func (r *subscriptionDrawerRenderer) Layout(size fyne.Size) {
	scrollContainer := r.objects[0]
	addButton := r.objects[1]
	addButtonHeight := addButton.MinSize().Height
	scrollContainer.Move(fyne.NewPos(0, 0))
	scrollContainer.Resize(fyne.NewSize(size.Width, size.Height-addButtonHeight))
	addButton.Move(fyne.NewPos(0, size.Height-addButtonHeight))
	addButton.Resize(fyne.NewSize(size.Width, addButtonHeight))
}

func (r *subscriptionDrawerRenderer) MinSize() (size fyne.Size) {
	size.Width = drawerWidth
	size.Height = drawerMinHeight
	return
}

func (r *subscriptionDrawerRenderer) Refresh() {
	box := r.objects[0].(*widget.ScrollContainer).Content.(*widget.Box)
	customProfilesItem := box.Children[0].(*material.DrawerItem)
	customProfilesItem.SetSelected(r.drawer.SelectedIndex == 0)
	subtitle := box.Children[1].(*widget.Label)
	subtitle.Refresh()
	subscriptionItems := box.Children[2:]
	if len(subscriptionItems) != len(r.drawer.SubscriptionTitles) {
		subscriptionItems = r.drawer.buildSubscriptionItems()
		box.Children = append(box.Children[:2], subscriptionItems...)
		box.Refresh()
	} else {
		for i, t := range r.drawer.SubscriptionTitles {
			button := subscriptionItems[i].(*material.DrawerItem)
			button.Text = t
			button.SetSelected(r.drawer.SelectedIndex == i+1)
		}
	}
	addButton := r.objects[1].(*material.DrawerItem)
	addButton.Refresh()
	canvas.Refresh(r.drawer)
}

func (r *subscriptionDrawerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *subscriptionDrawerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *subscriptionDrawerRenderer) Destroy() {
}
