package material

import (
	"fyne.io/fyne/theme"
	"image/color"
	"math"
)

func primaryColorOf(alpha float32) color.Color {
	r, g, b, _ := theme.PrimaryColor().RGBA()
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(float32(math.Round(float64(alpha) * 255)))}
}

func hoverColor() color.Color {
	return color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(float32(math.Round(0.12 * 255)))}
}
