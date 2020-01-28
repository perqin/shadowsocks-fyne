package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

type WrapperWidget struct {
	widget.BaseWidget
	Content fyne.CanvasObject
}

func (w *WrapperWidget) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)

	return &wrapperWidgetRenderer{content: w.Content}
}

type wrapperWidgetRenderer struct {
	content fyne.CanvasObject
}

func (w *wrapperWidgetRenderer) Layout(size fyne.Size) {
	w.content.Resize(size)
}

func (w *wrapperWidgetRenderer) MinSize() fyne.Size {
	return w.content.MinSize()
}

func (w *wrapperWidgetRenderer) Refresh() {
	w.content.Refresh()
}

func (w *wrapperWidgetRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (w *wrapperWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{w.content}
}

func (w *wrapperWidgetRenderer) Destroy() {
}
