package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

const twoLineListItemHeight = 64
const textHPadding = 16
const indicatorWidth = 8

func NewProfileWidget() *ProfileWidget {
	w := &ProfileWidget{}
	w.ExtendBaseWidget(w)
	return w
}

type ProfileWidget struct {
	widget.BaseWidget
	Profile  Profile
	Selected bool
	OnTapped func()
	hovered  bool
}

func (w *ProfileWidget) Tapped(*fyne.PointEvent) {
	w.Selected = true
	w.Refresh()
	if w.OnTapped != nil {
		w.OnTapped()
	}
}

func (w *ProfileWidget) TappedSecondary(*fyne.PointEvent) {
}

func (w *ProfileWidget) MouseIn(*desktop.MouseEvent) {
	w.hovered = true
	w.Refresh()
}

func (w *ProfileWidget) MouseOut() {
	w.hovered = false
	w.Refresh()
}

func (w *ProfileWidget) MouseMoved(*desktop.MouseEvent) {
}

func (w *ProfileWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	indicator := canvas.NewRectangle(theme.PrimaryColor())
	title := widget.NewLabel(w.getTitle())
	secondaryText := widget.NewLabel(w.getSecondaryText())
	return &profileWidgetRenderer{profileWidget: w, objects: []fyne.CanvasObject{indicator, title, secondaryText}}
}

func (w *ProfileWidget) getTitle() string {
	return w.Profile.Name
}

func (w *ProfileWidget) getSecondaryText() string {
	return fmt.Sprintf("%s:%d", w.Profile.Server, w.Profile.ServerPort)
}

type profileWidgetRenderer struct {
	profileWidget *ProfileWidget
	// [0]: Indicator
	// [1]: Title
	// [2]: Secondary text
	objects []fyne.CanvasObject
}

func (r *profileWidgetRenderer) Layout(size fyne.Size) {
	indicator := r.objects[0]
	indicator.Move(fyne.NewPos(0, 0))
	indicator.Resize(fyne.NewSize(indicatorWidth, size.Height))

	title := r.objects[1]
	titleMinSize := title.MinSize()
	secondaryText := r.objects[2]
	secondaryTextMinSize := secondaryText.MinSize()
	vPadding := (size.Height - titleMinSize.Height - secondaryTextMinSize.Height) / 2
	textWidth := size.Width - textHPadding*2
	title.Move(fyne.NewPos(textHPadding, vPadding))
	title.Resize(fyne.NewSize(textWidth, titleMinSize.Height))
	secondaryText.Move(fyne.NewPos(textHPadding, vPadding+titleMinSize.Height))
	secondaryText.Resize(fyne.NewSize(textWidth, secondaryTextMinSize.Height))
}

func (r *profileWidgetRenderer) MinSize() fyne.Size {
	var size fyne.Size
	size.Height = twoLineListItemHeight
	titleMinSize := r.objects[1].MinSize()
	secondaryTextMinSize := r.objects[2].MinSize()
	textWidth := titleMinSize.Width
	if secondaryTextMinSize.Width > textWidth {
		textWidth = secondaryTextMinSize.Width
	}
	size.Width = textHPadding*2 + textWidth
	return size
}

func (r *profileWidgetRenderer) Refresh() {
	indicator := r.objects[0]
	if r.profileWidget.Selected {
		indicator.Show()
	} else {
		indicator.Hide()
	}

	title := r.objects[1].(*widget.Label)
	title.Text = r.profileWidget.getTitle()
	title.Refresh()

	secondaryText := r.objects[2].(*widget.Label)
	secondaryText.Text = r.profileWidget.getSecondaryText()
	secondaryText.Refresh()

	canvas.Refresh(r.profileWidget)
}

func (r *profileWidgetRenderer) BackgroundColor() color.Color {
	if r.profileWidget.hovered {
		return theme.HoverColor()
	} else {
		return theme.BackgroundColor()
	}
}

func (r *profileWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *profileWidgetRenderer) Destroy() {
}
