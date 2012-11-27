package gui

import (
	"code.google.com/p/gohorde/horde3d"
	"github.com/jteeuwen/glfw"
)

//Used for both menus and HUDs

const (
	RelativeNormal = iota
	RelativeLeft
	RelativeRight
)

var screenWidth, screenHeight int
var activeWidgets []Widget
var tempArray [4]float32

//resolution independent
// scale based on a 640x480 grid
// relativePos can be used to set position relative to a side of the screen
// each widget has a hardcoded material

type Position struct {
	X, Y, U, V float32
	Relative   int
}

//ActualPosition returns the actual position on the screen from the interpreted
// and scaled WPosition
func (p *Position) ActualPosition() []float32 {
	//TODO: add actual translation to position
	tempArray[0] = p.X
	tempArray[1] = p.Y
	tempArray[2] = p.U
	tempArray[3] = p.V
	return tempArray[:]
}

type Widget interface {
	Add()
	Update()
}

func LoadGui(widgets []Widget) {
	horde3d.ClearOverlays()
	activeWidgets = widgets
	for g := range widgets {
		widgets[g].Add()
	}
}

func Update() {
	if activeWidgets != nil {
		for g := range activeWidgets {
			activeWidgets[g].Update()
		}
	}
}

func UpdateScreenSize(w, h int) {
	screenHeight = h
	screenWidth = w
}

func MousePos() (int, int) {
	return glfw.MousePos()
}

//TODO: Hover, click, scroll, type
