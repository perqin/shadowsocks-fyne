package widget

import "fyne.io/fyne"

func NewStackLayout() fyne.Layout {
	return &stackLayout{}
}

type stackLayout struct {
}

func (l *stackLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.Position{X: 0, Y: 0})
		s := size
		ms := o.MinSize()
		if ms.Width < s.Width {
			s.Width = ms.Width
		}
		if ms.Height < s.Height {
			s.Height = ms.Height
		}
		o.Resize(s)
	}
}

func (l *stackLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var size fyne.Size
	for _, o := range objects {
		size = size.Union(o.MinSize())
	}
	return size
}
