package engine

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

func initGui() {
	glfw.SetCharCallback(keyCollector)
}

func keyCollector(key, state int) {
	if state == glfw.KeyPress {
		if gKeyCollector != nil {
			gKeyCollector(key)
		}
	}
}

type KeyCollector func(key int)

var gKeyCollector KeyCollector

//resolution independent
// relativePos can be used to set position relative to a side of the screen

//ActualPosition returns the actual position on the screen from the interpreted
// relative position
//TODO: Test for stretching between different screen ratios
// width should be constant, and not change between ratio changes
// everything else should be relative to screen size
func (d *Dimensions) ActualPosition(result []float32) {
	//Y is not relative to aspect
	//vert1 
	result[1] = d.Y
	result[2] = 0
	result[3] = 1
	//vert2 
	result[5] = d.Y
	result[6] = 0
	result[7] = 0
	//vert3 
	result[9] = (d.Y + d.Height)
	result[10] = 1
	result[11] = 0
	//vert4 
	result[13] = (d.Y + d.Height)
	result[14] = 1
	result[15] = 1
	switch d.RelativeTo {
	case RelativeAspect:
		result[0] = d.X
		result[4] = (d.X + d.Width)
		result[8] = d.X
		result[12] = (d.X + d.Width)
	case RelativeLeft:
		result[0] = d.X * screenRatio
		result[4] = (d.X * screenRatio) + d.Width
		result[8] = d.X * screenRatio
		result[12] = (d.X * screenRatio) + d.Width
	case RelativeRight:
		result[0] = (d.X * screenRatio) + d.Width
		result[4] = d.X * screenRatio
		result[8] = (d.X * screenRatio) + d.Width
		result[12] = d.X * screenRatio
	}
}

type Dimensions struct {
	X, Y, Height, Width float32
	RelativeTo          int
}

type Color struct {
	R, G, B, A float32
}

type Overlay struct {
	Dimensions *Dimensions
	Color      *Color
	Material   *Material
}

func (o *Overlay) Place() {
	o.Dimensions.ActualPosition(tempArray[:])
	horde3d.ShowOverlays(tempArray[:], 4, o.Color.R, o.Color.G,
		o.Color.B, o.Color.A, o.Material.H3DRes, 0)
}

//Widget is a collection of Overlays
type Widget interface {
	MouseArea() *Dimensions
	Update()
	Hover()
	Click()
	Scroll(int)
}

//Gui is a collection of Widgets
type Gui struct {
	Widgets      []Widget
	UseMouse     bool
	KeyCollect   KeyCollector
	prevTime     float64
	prevWheelPos int
}

func (g *Gui) ElapsedTime() float64 {
	return (glfw.Time() - g.prevTime)
}

//AddWidget adds a widget to the last / top location
// of the gui
func (g *Gui) AddWidget(widget Widget) {
	g.Widgets = append(g.Widgets, widget)
}

func (g *Gui) Load() {
	if g.UseMouse {
		glfw.Enable(glfw.MouseCursor)
	} else {
		glfw.Disable(glfw.MouseCursor)
	}
	gKeyCollector = g.KeyCollect

}

func (g *Gui) Unload() {
	glfw.Disable(glfw.MouseCursor)
	gKeyCollector = nil
}

func (g *Gui) Update() {
	horde3d.ClearOverlays()

	if widget, ok := g.WidgetUnderMouse(); ok {
		widget.Hover()
		if glfw.MouseButton(0) == glfw.KeyPress {
			widget.Click()
		}
		delta := glfw.MouseWheel()
		if delta != g.prevWheelPos {
			//TODO: Test delta
			widget.Scroll(g.prevWheelPos - delta)
		}
	}

	for i := range g.Widgets {
		g.Widgets[i].Update()
	}
	g.prevTime = glfw.Time()
}

func (g *Gui) WidgetUnderMouse() (Widget, bool) {
	var x, y float32
	var dimensions *Dimensions
	//Loop through list backwards to check for topmost
	// widget that the mouse hits
	for i := len(g.Widgets) - 1; i >= 0; i-- {
		dimensions = g.Widgets[i].MouseArea()
		x, y = GuiMousePos(dimensions.RelativeTo)
		if x >= dimensions.X && x <= dimensions.X+dimensions.Width &&
			y >= dimensions.Y && y <= dimensions.Y+dimensions.Height {
			return g.Widgets[i], true

		}
	}

	return nil, false
}

func updateScreenSize(w, h int) {
	screenHeight = h
	screenWidth = w
	screenRatio = float32(w) / float32(h)
}

func GuiMousePos(relative int) (x, y float32) {
	//Return position according to widget ratio positioning
	//  0.0 - 1.0
	gX, gY := glfw.MousePos()
	x = float32(gX)
	y = float32(gY)
	switch relative {
	case RelativeLeft:
		return (x / float32(screenWidth)), (y / float32(screenHeight))
	case RelativeRight:
		return (float32(screenWidth) - x) / float32(screenWidth),
			(float32(screenHeight) - y) / float32(screenHeight)
	case RelativeAspect:
		return (x / float32(screenWidth)) * screenRatio, (y / float32(screenHeight))
	}
	return -1, -1
}
