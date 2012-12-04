package gui

import (
	"code.google.com/p/gohorde/horde3d"
)

type Button struct {
	Dimensions     *Dimensions
	Background     *horde3d.H3DRes
	Color          *Color
	HoverColor     *Color
	ShowBackground bool
	//text
	Text           string
	TextColor      *Color
	TextHoverColor *Color
	hover          bool
}

func AddButton(dimensions *Dimensions, background *horde3d.H3DRes, color *Color, hoverColor *Color,
	text string, textColor *Color, textHoverColor *Color) *Button {
	return

}
