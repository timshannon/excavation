package gui

import (
	"excavation/engine"
)

type Button struct {
	Dimensions     *engine.Dimensions
	Background     *engine.Material
	Color          *engine.Color
	HoverColor     *engine.Color
	ShowBackground bool
	//text
	Text           string
	TextColor      *engine.Color
	TextHoverColor *engine.Color
	hover          bool
}

//func AddButton(dimensions *Dimensions, background *horde3d.H3DRes, color *Color, hoverColor *Color,
//text string, textColor *Color, textHoverColor *Color) *Button {
//return

//}
