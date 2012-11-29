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
var screenHeight int
var screenWidth int
var tempArray [16]float32 //Rectangles only for now

func Init() {
	glfw.SetCharCallback(keyCollector)
	glfw.SetMouseButtonCallback(click)
}

func keyCollector(key, state int) {
	if state == glfw.KeyPress {
		gKeyCollector(key)
	}
}

func click(button, state int) {
	if state == glfw.KeyPress {
		getWidgetFromMouse().Click()
	}
}

func getWidgetFromMouse() Widget {
}

type KeyCollector func(key int)

var gKeyCollector KeyCollector

//resolution independent
// relativePos can be used to set position relative to a side of the screen

//ActualPosition returns the actual position on the screen from the interpreted
// relative position
func ActualPosition(result []float32, position *Position, size *Size) {
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
	switch position.RelativeTo {
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

func AddOverlay(widget Widget) {
	ActualPosition(tempArray[:], widget.Position(), widget.Size())
	horde3d.ShowOverlays(tempArray[:], 4, widget.Color().R, widget.Color().G,
		widget.Color().B, widget.Color().A, widget.Material(), 0)

}

type Position struct {
	X, Y       float32
	RelativeTo int
}
type Size struct {
	Height, Width float32
}

type Color struct {
	R, G, B, A float32
}

type Widget interface {
	Update()
	Position() *Position
	Size() *Size
	Color() *Color
	Material() horde3d.H3DRes
	Hover()
	Click()
}

type Gui struct {
	Widgets    []Widget
	UseMouse   bool
	KeyCollect KeyCollector
}

func (g *Gui) AddWidget(widget Widget) {
	g.Widgets = append(g.Widgets, widget)
}

func (g *Gui) Load() {
	if g.UseMouse {
		glfw.Enable(glfw.MouseCursor)
	} else {
		glfw.Disable(glfw.MouseCursor)
	}
	//set GLFW callbacks
	gKeyCollector = g.KeyCollect

}

func (g *Gui) Unload() {
	glfw.Disable(glfw.MouseCursor)
	gKeyCollector = nil
}

func (g *Gui) Update() {
	horde3d.ClearOverlays()
	for i := range g.Widgets {
		g.Widgets[i].Update()
	}
}

func UpdateScreenSize(w, h int) {
	screenHeight = h
	screenWidth = w
	screenRatio = float32(w) / float32(h)
}

func MousePos() *Position {
	//Return position according to widget ratio positioning
	//  0.0 - 1.0
	x, y := glfw.MousePos()

	return &Position{float32(x/screenWidth) * screenRatio, float32(y / screenHeight), RelativeAspect}
}
