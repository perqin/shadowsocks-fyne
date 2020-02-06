package material

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"github.com/icza/gox/imagex/colorx"
	"image/color"
	"log"
	"math"
)

// NewLightTheme returns a new light theme.
func NewLightTheme() *Theme {
	// <color name="ripple_material_dark">#33ffffff</color>
	// <color name="ripple_material_light">#1f000000</color>
	t := &Theme{
		WindowBackground: baseTheme.BackgroundColor(),
		hoverColor:       withAlphaByte(mustParseHexColor("#000"), 0x1f),
		primaryColor:     mustParseHexColor("#00BCD4"),
	}
	return t
}

func mustParseHexColor(hex string) color.Color {
	c, err := colorx.ParseHexColor(hex)
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

type Theme struct {
	WindowBackground         color.Color
	hoverColor, primaryColor color.Color
}

var baseTheme = theme.LightTheme()

func (t *Theme) BackgroundColor() color.Color {
	return color.Transparent
}

func (t *Theme) ButtonColor() color.Color {
	return baseTheme.ButtonColor()
}

func (t *Theme) DisabledButtonColor() color.Color {
	return baseTheme.DisabledButtonColor()
}

func (t *Theme) HyperlinkColor() color.Color {
	return baseTheme.HyperlinkColor()
}

func (t *Theme) TextColor() color.Color {
	return baseTheme.TextColor()
}

func (t *Theme) DisabledTextColor() color.Color {
	return baseTheme.DisabledTextColor()
}

func (t *Theme) IconColor() color.Color {
	return baseTheme.IconColor()
}

func (t *Theme) DisabledIconColor() color.Color {
	return baseTheme.DisabledIconColor()
}

func (t *Theme) PlaceHolderColor() color.Color {
	return baseTheme.PlaceHolderColor()
}

func (t *Theme) PrimaryColor() color.Color {
	return t.primaryColor
}

func (t *Theme) HoverColor() color.Color {
	return t.hoverColor
}

func (t *Theme) FocusColor() color.Color {
	return baseTheme.FocusColor()
}

func (t *Theme) ScrollBarColor() color.Color {
	return baseTheme.ScrollBarColor()
}

func (t *Theme) ShadowColor() color.Color {
	return baseTheme.ShadowColor()
}

func (t *Theme) TextSize() int {
	return baseTheme.TextSize()
}

func (t *Theme) TextFont() fyne.Resource {
	return baseTheme.TextFont()
}

func (t *Theme) TextBoldFont() fyne.Resource {
	return baseTheme.TextBoldFont()
}

func (t *Theme) TextItalicFont() fyne.Resource {
	return baseTheme.TextItalicFont()
}

func (t *Theme) TextBoldItalicFont() fyne.Resource {
	return baseTheme.TextBoldItalicFont()
}

func (t *Theme) TextMonospaceFont() fyne.Resource {
	return baseTheme.TextMonospaceFont()
}

func (t *Theme) Padding() int {
	return 0
}

func (t *Theme) IconInlineSize() int {
	return baseTheme.IconInlineSize()
}

func (t *Theme) ScrollBarSize() int {
	return baseTheme.ScrollBarSize()
}

func (t *Theme) ScrollBarSmallSize() int {
	return baseTheme.ScrollBarSmallSize()
}

func withAlpha(c color.Color, alpha float64) color.Color {
	return withAlphaByte(c, uint8(float32(math.Round(alpha*255))))
}

func withAlphaByte(c color.Color, alpha uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: alpha}
}
