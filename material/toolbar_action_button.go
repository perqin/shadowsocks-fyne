package material

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

const toolbarActionButtonSize = 48
const toolbarActionButtonIconSize = 24

func NewToolbarActionButton(icon fyne.Resource, onTapped func()) *ToolbarActionButton {
	b := &ToolbarActionButton{icon: icon, onTapped: onTapped}
	b.ExtendBaseWidget(b)
	return b
}

type ToolbarActionButton struct {
	widget.BaseWidget
	icon     fyne.Resource
	hovered  bool
	onTapped func()
}

func (b *ToolbarActionButton) Tapped(*fyne.PointEvent) {
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *ToolbarActionButton) TappedSecondary(*fyne.PointEvent) {
}

func (b *ToolbarActionButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

func (b *ToolbarActionButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

func (b *ToolbarActionButton) MouseMoved(*desktop.MouseEvent) {
}

func (b *ToolbarActionButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	bg := canvas.NewCircle(theme.HoverColor())
	bg.Hide()
	i := canvas.NewImageFromResource(b.icon)
	return &toolbarActionButtonRenderer{button: b, background: bg, icon: i, objects: []fyne.CanvasObject{bg, i}}
}

type toolbarActionButtonRenderer struct {
	button     *ToolbarActionButton
	background *canvas.Circle
	icon       *canvas.Image
	objects    []fyne.CanvasObject
}

func (r *toolbarActionButtonRenderer) Layout(size fyne.Size) {
	r.background.Move(fyne.NewPos((size.Width-toolbarActionButtonSize)/2, (size.Height-toolbarActionButtonSize)/2))
	r.background.Resize(fyne.NewSize(toolbarActionButtonSize, toolbarActionButtonSize))
	r.icon.Move(fyne.NewPos((size.Width-toolbarActionButtonIconSize)/2, (size.Height-toolbarActionButtonIconSize)/2))
	r.icon.Resize(fyne.NewSize(toolbarActionButtonIconSize, toolbarActionButtonIconSize))
}

func (r *toolbarActionButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(toolbarActionButtonSize, toolbarActionButtonSize)
}

func (r *toolbarActionButtonRenderer) Refresh() {
	if r.button.hovered {
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.icon.Refresh()
	canvas.Refresh(r.button)
}

func (r *toolbarActionButtonRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *toolbarActionButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *toolbarActionButtonRenderer) Destroy() {
}
