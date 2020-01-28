package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"image/color"
)

func NewColoredBox() *ColoredBox {
	b := &ColoredBox{}
	b.ExtendBaseWidget(b)
	return b
}

type ColoredBox struct {
	widget.BaseWidget
	backgroundColor color.Color
}

func (b *ColoredBox) SetBackgroundColor(c color.Color) {
	b.backgroundColor = c
	b.Refresh()
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (b *ColoredBox) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	return &coloredBoxRenderer{colorBox: b}
}

type coloredBoxRenderer struct {
	colorBox *ColoredBox
}

func (c *coloredBoxRenderer) Layout(size fyne.Size) {
	c.colorBox.Resize(size)
}

func (c *coloredBoxRenderer) MinSize() fyne.Size {
	return c.colorBox.Size()
}

func (c *coloredBoxRenderer) Refresh() {
	canvas.Refresh(c.colorBox)
}

func (c *coloredBoxRenderer) BackgroundColor() color.Color {
	return c.colorBox.backgroundColor
}

func (c *coloredBoxRenderer) Objects() []fyne.CanvasObject {
	return nil
}

func (c *coloredBoxRenderer) Destroy() {
}
