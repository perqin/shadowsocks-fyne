package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

func NewTappableWidget(content fyne.CanvasObject, onTapped func()) fyne.CanvasObject {
	w := &TappableWidget{BaseWidget: widget.BaseWidget{}, Content: content, OnTapped: onTapped}
	w.ExtendBaseWidget(w)
	return w
}

type TappableWidget struct {
	widget.BaseWidget
	Content  fyne.CanvasObject
	OnTapped func()
	hovered  bool
}

var _ fyne.Tappable = (*TappableWidget)(nil)
var _ desktop.Hoverable = (*TappableWidget)(nil)

func (w *TappableWidget) Tapped(*fyne.PointEvent) {
	if w.OnTapped != nil {
		w.OnTapped()
	}
}

func (w *TappableWidget) TappedSecondary(*fyne.PointEvent) {
}

func (w *TappableWidget) MouseIn(*desktop.MouseEvent) {
	w.hovered = true
	w.Refresh()
}

func (w *TappableWidget) MouseOut() {
	w.hovered = false
	w.Refresh()
}

func (w *TappableWidget) MouseMoved(*desktop.MouseEvent) {
}

func (w *TappableWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	return &tappableWidgetRenderer{widget: w, content: w.Content}
}

type tappableWidgetRenderer struct {
	widget  *TappableWidget
	content fyne.CanvasObject
}

func (w *tappableWidgetRenderer) Layout(size fyne.Size) {
	w.content.Resize(size)
}

func (w *tappableWidgetRenderer) MinSize() fyne.Size {
	return w.content.MinSize()
}

func (w *tappableWidgetRenderer) Refresh() {
	canvas.Refresh(w.widget)
	w.content.Refresh()
}

func (w *tappableWidgetRenderer) BackgroundColor() color.Color {
	if w.widget.hovered {
		return theme.HoverColor()
	} else {
		return color.Transparent
	}
}

func (w *tappableWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{w.content}
}

func (w *tappableWidgetRenderer) Destroy() {
}
