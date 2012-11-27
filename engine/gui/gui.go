package gui

import (
	"code.google.com/p/gohorde/horde3d"
	"github.com/jteeuwen/glfw"
)

//Used for both menus and HUDs

const (
	RelativeAspect = iota
	RelativeLeft
	RelativeRight
)

var screenRatio float32
var activeGui *Gui
var tempArray [16]float32 //Rectangles only for now

//resolution independent
// relativePos can be used to set position relative to a side of the screen

//ActualPosition returns the actual position on the screen from the interpreted
// relative position
func ActualPosition(result []float32, position *Postion, size *Size, relative int) {
	//Y is not relative to aspect
	//vert1 
	result[1] = position.Y
	result[2] = 0
	result[3] = 1
	//vert2 
	result[5] = position.Y
	result[6] = 0
	result[7] = 0
	//vert3 
	result[9] = (position.Y + size.Height)
	result[10] = 1
	result[11] = 0
	//vert4 
	result[13] = (position.Y + size.Height)
	result[14] = 1
	result[15] = 1
	switch relative {
	case RelativeAspect:
		result[0] = position.X
		result[4] = (position.X + size.Width)
		result[8] = position.X
		result[12] = (position.X + size.Width)
	case RelativeLeft:
		result[0] = position.X / screenRatio
		result[4] = (position.X / screenRatio) + size.Width
		result[8] = position.X / screenRatio
		result[12] = (position.X / screenRatio) + size.Width
	case RelativeRight:
		result[0] = (position.X / screenRatio) + size.Width
		result[4] = position.X / screenRatio
		result[8] = (position.X / screenRatio) + size.Width
		result[12] = position.X / screenRatio
	}
}

func addOverlay(widget Widget) {
	ActualPosition(tempArray[:], widget.Position(), widget.Size(), widget.RelativeTo())
	horde3d.ShowOverlays(tempArray[:], 4, widget.Color().R, widget.Color().G,
		widget.Color().B, widget.Color().A, widget.Material(), 0)

}

type Postion struct {
	X, Y float32
}
type Size struct {
	Height, Width float32
}

type Color struct {
	R, G, B, A float32
}

type Widget interface {
	Update()
	Position() *Postion
	Size() *Size
	RelativeTo() int
	Color() *Color
	Material() horde3d.H3DRes
}

type Gui struct {
	Widgets []Widget
	//TODO: Hover, click, scroll, type

}

func LoadGui(gui *Gui) {
	horde3d.ClearOverlays()
	activeGui = gui
}

func (g *Gui) AddWidget(widget Widget) {
	g.Widgets = append(g.Widgets, widget)
}

func Update() {
	if activeGui != nil {
		horde3d.ClearOverlays()
		for g := range activeGui.Widgets {
			activeGui.Widgets[g].Update()
			addOverlay(activeGui.Widgets[g])
		}
	}
}

func UpdateScreenSize(w, h int) {
	screenRatio = float32(w) / float32(h)
}

func MousePos() (int, int) {
	return glfw.MousePos()
}
