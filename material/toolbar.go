package material

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const toolbarHeight = 56
const titleLabelHPadding = 16

// TODO: Replace with custom Widget
func NewToolbar(title *widget.Label, actions ...*ToolbarActionButton) fyne.CanvasObject {
	objects := []fyne.CanvasObject{canvas.NewRectangle(theme.PrimaryColor()), title}
	for _, a := range actions {
		objects = append(objects, a)
	}
	return fyne.NewContainerWithLayout(&toolbarLayout{}, objects...)
}

type toolbarLayout struct {
}

func (l *toolbarLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// Background
	objects[0].Move(fyne.NewPos(0, 0))
	objects[0].Resize(size)
	var actionsWidth int
	var actionWidthList []int
	for _, o := range objects[2:] {
		width := o.MinSize().Width
		actionsWidth += width
		actionWidthList = append(actionWidthList, width)
	}
	titleMinSize := objects[1].MinSize()
	vPadding := (size.Height - titleMinSize.Height) / 2
	objects[1].Move(fyne.NewPos(titleLabelHPadding, vPadding))
	objects[1].Resize(fyne.NewSize(size.Width-actionsWidth-titleLabelHPadding*2, titleMinSize.Height))
	left := size.Width - actionsWidth
	for i, o := range objects[2:] {
		o.Move(fyne.NewPos(left, 0))
		o.Resize(fyne.NewSize(actionWidthList[i], size.Height))
		left += actionWidthList[i]
	}
}

func (l *toolbarLayout) MinSize(objects []fyne.CanvasObject) (size fyne.Size) {
	titleMinSize := objects[1].MinSize()
	actionsWidth := 0
	for _, o := range objects[2:] {
		actionsWidth += o.MinSize().Width
	}
	size.Width = titleLabelHPadding*2 + titleMinSize.Width + actionsWidth
	size.Height = toolbarHeight
	return
}
