package gui

import (
	"excavation/engine"
)

const (
	defaultBackground = "overlays/gui/default/background.material.xml"
	defaultFont       = "overlays/gui/default/font.material.xml"
)

type Button struct {
	Name                   string
	BackgroundOverlay      *engine.Overlay
	BackgroundHoverOverlay *engine.Overlay
	BackgroundClickOverlay *engine.Overlay
	ShowBackground         bool
	//text
	Text      *engine.Text
	TextHover *engine.Text
	TextClick *engine.Text
	hover     bool
}

//MakeButton returns a button with the default background and colors
//  changes from default can be made by accessing exported variables
func MakeButton(name, text string, textSize float32, dimensions *engine.ScreenArea) *Button {
	defaultColor := &engine.Color{118, 118, 118, 1}
	hoverColor := &engine.Color{155, 155, 155, 50}
	textColor := &engine.Color{255, 255, 255, 255}
	//TODO: Determine text position
	// Auto size button to hold text, center text vertically
	textPosition := dimensions.Position

	button := &Button{
		Name:                   name,
		BackgroundOverlay:      engine.NewOverlay(defaultBackground, defaultColor, dimensions),
		BackgroundHoverOverlay: engine.NewOverlay(defaultBackground, hoverColor, dimensions),
		BackgroundClickOverlay: engine.NewOverlay(defaultBackground, hoverColor, dimensions),
		ShowBackground:         true,
		Text:                   engine.NewText(text, textSize, defaultFont, textColor, textPosition),
		TextHover:              engine.NewText(text, textSize, defaultFont, textColor, textPosition),
		TextClick:              engine.NewText(text, textSize, defaultFont, textColor, textPosition),
	}
	return button

}

//textPosition returns the position of the text based on the size
// and dimensions of the button
//func textPosition(dimensions *engine.Dimensions, textSize) *engine.ScreenPosition {

//}

func (b *Button) MouseArea() *engine.ScreenArea {
	return b.BackgroundOverlay.Dimensions
}

func (b *Button) Hover() {
	b.hover = true
}

func (b *Button) Update() {
	if b.hover {
		b.BackgroundHoverOverlay.Place()
		b.TextHover.Place()
	} else {
		b.BackgroundOverlay.Place()
		b.Text.Place()
	}
}

func (b *Button) Click() {
	b.BackgroundClickOverlay.Place()
	b.TextClick.Place()
}

func (b *Button) Scroll(delta int) {
	//Nothing
	return
}
