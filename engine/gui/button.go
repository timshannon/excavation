package gui

import (
	"excavation/engine"
)

const (
	defaultBackground = "textures/gui/defaultBackground.material.xml"
)

type Button struct {
	BackgroundOverlay      *engine.Overlay
	BackgroundHoverOverlay *engine.Overlay
	BackgroundClickOverlay *engine.Overlay
	ShowBackground         bool
	//text
	Text           string
	TextSize       float32
	FontMaterial   *engine.Material
	TextColor      *engine.Color
	TextHoverColor *engine.Color
	TextClickColor *engine.Color
	hover          bool
}

//MakeButton returns a button with the default background and colors
//  changes from default can be made by accessing exported variables
func MakeButton(text string, textSize float32, dimensions *engine.Dimensions) *Button {

	return

}

func (b *Button) MouseArea() *engine.Dimension {
	return b.Dimensions
}

func (b *Button) Hover() {
	b.hover = true
}

//Update()
//Click()
//Scroll(int)
