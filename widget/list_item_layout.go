package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func NewTwoLineListItem(title, secondaryText string) fyne.CanvasObject {
	return fyne.NewContainerWithLayout(
		&twoLineListItemLayout{},
		widget.NewLabel(title),
		widget.NewLabel(secondaryText))
}

const twoLineListItemHeight = 64

type twoLineListItemLayout struct {
}

func (l *twoLineListItemLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	h1 := objects[0].MinSize().Height
	h2 := objects[1].MinSize().Height
	paddingH := (twoLineListItemHeight - h1 - h2 - 4) / 2
	width := size.Width - 32
	// Layout title
	objects[0].Move(fyne.Position{X: 16, Y: paddingH})
	objects[0].Resize(fyne.Size{Width: width, Height: h1})
	// Layout secondary text
	objects[1].Move(fyne.Position{X: 15, Y: paddingH + h1 + 4})
	objects[1].Resize(fyne.Size{Width: width, Height: h2})
}

func (l *twoLineListItemLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var size fyne.Size
	for _, o := range objects {
		size = size.Union(o.MinSize())
	}
	size.Height = twoLineListItemHeight
	return size
}
