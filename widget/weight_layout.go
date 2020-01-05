package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var _ fyne.Layout = (*weightLayout)(nil)

type weightLayout struct {
	horizontal bool
}

func (g *weightLayout) weightOf(object fyne.CanvasObject) uint {
	if w, ok := object.(weighted); ok {
		return w.Weight()
	}
	return 0
}

// Layout is called to pack all child objects into a specified size.
// For a VWeightLayout this will pack objects into a single column where each item
// is full width but the height is weighted. A zero-weight item's height is the minimum required.
func (g *weightLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	weightCount := uint(0)
	total := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		weight := g.weightOf(child)
		if weight > 0 {
			weightCount += weight
			continue
		}
		if g.horizontal {
			total += child.MinSize().Width
		} else {
			total += child.MinSize().Height
		}
	}

	x, y := 0, 0
	var extra int
	if g.horizontal {
		extra = size.Width - total - (theme.Padding() * (len(objects) - 1))
	} else {
		extra = size.Height - total - (theme.Padding() * (len(objects) - 1))
	}
	extraCell := 0
	if weightCount > 0 {
		extraCell = int(float64(extra) / float64(weightCount))
	}

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		width := child.MinSize().Width
		height := child.MinSize().Height
		if weight := g.weightOf(child); weight > 0 {
			if g.horizontal {
				width = extraCell * int(weight)
			} else {
				height = extraCell * int(weight)
			}
		}

		child.Move(fyne.NewPos(x, y))
		if g.horizontal {
			x += theme.Padding() + width
			child.Resize(fyne.NewSize(width, size.Height))
		} else {
			y += theme.Padding() + height
			child.Resize(fyne.NewSize(size.Width, height))
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a WeightLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (g *weightLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	added := false
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if g.horizontal {
			minSize = minSize.Add(fyne.NewSize(child.MinSize().Width, 0))
			minSize.Height = fyne.Max(child.MinSize().Height, minSize.Height)
			if added {
				minSize.Width += theme.Padding()
			}
		} else {
			minSize = minSize.Add(fyne.NewSize(0, child.MinSize().Height))
			minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
			if added {
				minSize.Height += theme.Padding()
			}
		}
		added = true
	}

	return minSize
}

func NewHWeightLayout() fyne.Layout {
	return &weightLayout{true}
}

func NewVWeightLayout() fyne.Layout {
	return &weightLayout{false}
}

type weighted interface {
	Weight() uint
}

var _ weighted = (*weightedItem)(nil)

type weightedItem struct {
	WrapperWidget
	weight uint
}

func (w *weightedItem) Weight() uint {
	return w.weight
}

func NewWeightedItem(object fyne.CanvasObject, weight uint) fyne.CanvasObject {
	return &weightedItem{WrapperWidget{widget.BaseWidget{}, object}, weight}
}
