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
	dimensions             *engine.ScreenArea
	BackgroundOverlay      *engine.Overlay
	BackgroundHoverOverlay *engine.Overlay
	BackgroundClickOverlay *engine.Overlay
	//text
	Text           *engine.Text
	TextHover      *engine.Text
	TextClick      *engine.Text
	hover          bool
	showBackground bool
	ClickEvent     func(sender string)
}

//MakeButton returns a button with the default background and colors
//  changes from default can be made by accessing exported variables
func MakeButton(name, text string, textSize float32, dimensions *engine.ScreenArea) *Button {
	defaultColor := engine.NewColor(118, 118, 118, 255)
	hoverColor := engine.NewColor(155, 155, 155, 50)
	textColor := engine.NewColor(255, 255, 255, 0)
	//TODO: Determine text position
	// Auto size button to hold text, center text vertically
	textPosition := dimensions.Position

	button := &Button{
		Name:                   name,
		BackgroundOverlay:      engine.NewOverlay(defaultBackground, defaultColor, dimensions),
		BackgroundHoverOverlay: engine.NewOverlay(defaultBackground, hoverColor, dimensions),
		BackgroundClickOverlay: engine.NewOverlay(defaultBackground, hoverColor, dimensions),
		showBackground:         true,
		Text:                   engine.NewText(text, textSize, defaultFont, textColor, textPosition),
		TextHover:              engine.NewText(text, textSize, defaultFont, textColor, textPosition),
		TextClick:              engine.NewText(text, textSize, defaultFont, textColor, textPosition),
	}
	button.dimensions = button.BackgroundOverlay.Dimensions
	return button

}

func (b *Button) ShowBackground(value bool) {
	if value {
		b.dimensions = b.BackgroundOverlay.Dimensions
	} else {
		//If no background, base mouse area on text
		b.dimensions = engine.NewScreenArea(b.Text.Position.X, b.Text.Position.Y,
			b.Text.Size, b.Text.Width(), b.Text.Position.RelativeTo)
	}
	b.showBackground = value
}

func (b *Button) MouseArea() *engine.ScreenArea {
	return b.dimensions
}

func (b *Button) Hover() {
	b.hover = true
}

func (b *Button) Update() {
	if b.hover {
		if b.showBackground {
			b.BackgroundHoverOverlay.Place()
		}
		if b.TextHover.Text != "" {
			b.TextHover.Place()
		}
		b.hover = false
	} else {
		if b.showBackground {
			b.BackgroundOverlay.Place()
		}
		if b.Text.Text != "" {
			b.Text.Place()
		}
	}
}

func (b *Button) Click(button int) {
	if button == 0 {
		if b.showBackground {
			b.BackgroundClickOverlay.Place()
		}
		if b.Text.Text != "" {
			b.TextClick.Place()
		}
		b.ClickEvent(b.Name)
	}
}

func (b *Button) Scroll(delta int) {
	//Nothing
	return
}
