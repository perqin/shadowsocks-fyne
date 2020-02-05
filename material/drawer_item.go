package material

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

const drawerItemHeight = 48
const drawerItemHPadding = 8
const drawerItemVPadding = 4
const drawerItemTextHPadding = 16

func NewDrawerItem(text string, onTapped func()) *DrawerItem {
	i := &DrawerItem{Text: text, onTapped: onTapped}
	i.ExtendBaseWidget(i)
	return i
}

type DrawerItem struct {
	widget.BaseWidget
	Text     string
	Selected bool
	hovered  bool
	onTapped func()
}

func (i *DrawerItem) Tapped(*fyne.PointEvent) {
	if i.onTapped != nil {
		i.onTapped()
	}
}

func (i *DrawerItem) TappedSecondary(*fyne.PointEvent) {
}

func (i *DrawerItem) MouseIn(*desktop.MouseEvent) {
	i.hovered = true
	i.Refresh()
}

func (i *DrawerItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

func (i *DrawerItem) MouseMoved(*desktop.MouseEvent) {
}

func (i *DrawerItem) SetSelected(selected bool) {
	i.Selected = selected
	i.Refresh()
}

func (i *DrawerItem) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)
	bg := canvas.NewRectangle(color.Transparent)
	t := canvas.NewText(i.Text, theme.TextColor())
	r := &drawerItemRenderer{item: i, background: bg, text: t, objects: []fyne.CanvasObject{bg, t}}
	return r
}

type drawerItemRenderer struct {
	item       *DrawerItem
	background *canvas.Rectangle
	text       *canvas.Text
	objects    []fyne.CanvasObject
}

func (r *drawerItemRenderer) Layout(size fyne.Size) {
	r.background.Move(fyne.NewPos(drawerItemHPadding, drawerItemVPadding))
	r.background.Resize(fyne.NewSize(size.Width-drawerItemHPadding*2, size.Height-drawerItemVPadding*2))
	textHeight := r.text.MinSize().Height
	textVPadding := (size.Height - textHeight) / 2
	r.text.Move(fyne.NewPos(drawerItemTextHPadding, textVPadding))
	r.text.Resize(fyne.NewSize(size.Width-drawerItemTextHPadding*2, size.Height-textVPadding*2))
}

func (r *drawerItemRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.text.MinSize().Width+drawerItemHPadding*2, drawerItemHeight)
}

func (r *drawerItemRenderer) Refresh() {
	var backgroundColor color.Color
	switch {
	case r.item.Selected:
		backgroundColor = primaryColorOf(0.12)
	case r.item.hovered:
		backgroundColor = hoverColor()
	default:
		backgroundColor = color.Transparent
	}
	r.background.FillColor = backgroundColor
	// TODO: Needed? r.background.Refresh()
	r.text.Text = r.item.Text
	canvas.Refresh(r.item)
}

func (r *drawerItemRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *drawerItemRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *drawerItemRenderer) Destroy() {
}
