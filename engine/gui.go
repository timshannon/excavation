package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"github.com/jteeuwen/glfw"
)

//Used for both menus and HUDs

const (
	ScreenRelativeAspect = iota
	ScreenRelativeLeft
	ScreenRelativeRight
)

var screenRatio float32
var screenHeight int
var screenWidth int
var tempArray [16]float32 //Rectangles only for now
var activeGui *Gui

func initGui() {
	glfw.SetCharCallback(keyCollector)
}

func LoadGui(gui *Gui) {
	HaltInput()
	activeGui = gui
	activeGui.Load()
}

func UnloadGui() {
	activeGui.Unload()
	activeGui = nil
	ResumeInput()
}
func updateGui() {
	if activeGui != nil {
		activeGui.Update()
	}
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

//toVertex returns the actual position on the screen from the interpreted
// relative position
//TODO: Test for stretching between different screen ratios
// width should be constant, and not change between ratio changes
// everything else should be relative to screen size
func (s *ScreenArea) toVertex(result []float32) {
	//Y is not relative to aspect
	//vert1 
	result[1] = s.Position.Y
	result[2] = 0
	result[3] = 1
	//vert2 
	result[5] = s.Position.Y
	result[6] = 0
	result[7] = 0
	//vert3 
	result[9] = (s.Position.Y + s.Height)
	result[10] = 1
	result[11] = 0
	//vert4 
	result[13] = (s.Position.Y + s.Height)
	result[14] = 1
	result[15] = 1
	switch s.Position.RelativeTo {
	case ScreenRelativeAspect:
		result[0] = s.Position.X
		result[4] = (s.Position.X + s.Width)
		result[8] = s.Position.X
		result[12] = (s.Position.X + s.Width)
	case ScreenRelativeLeft:
		result[0] = s.Position.X * screenRatio
		result[4] = (s.Position.X * screenRatio) + s.Width
		result[8] = s.Position.X * screenRatio
		result[12] = (s.Position.X * screenRatio) + s.Width
	case ScreenRelativeRight:
		result[0] = (s.Position.X * screenRatio) + s.Width
		result[4] = s.Position.X * screenRatio
		result[8] = (s.Position.X * screenRatio) + s.Width
		result[12] = s.Position.X * screenRatio
	}
}

type ScreenPosition struct {
	X, Y       float32
	RelativeTo int
}

type ScreenArea struct {
	Position      *ScreenPosition
	Height, Width float32
}

type Color struct {
	R, G, B, A float32
}

type Overlay struct {
	Dimensions *ScreenArea
	Color      *Color
	Material   *Material
}

type Text struct {
	Text         string
	Position     *ScreenPosition
	Size         float32
	FontMaterial *Material
	Color        *Color
}

func (o *Overlay) Place() {
	o.Dimensions.toVertex(tempArray[:])
	horde3d.ShowOverlays(tempArray[:], 4, o.Color.R, o.Color.G,
		o.Color.B, o.Color.A, o.Material.H3DRes, 0)
}

func (t *Text) Place() {

}

//Widget is a collection of Overlays
type Widget interface {
	MouseArea() *ScreenArea
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
	var dimensions *ScreenArea
	//Loop through list backwards to check for topmost
	// widget that the mouse hits
	for i := len(g.Widgets) - 1; i >= 0; i-- {
		dimensions = g.Widgets[i].MouseArea()
		x, y = GuiMousePos(dimensions.Position.RelativeTo)
		if x >= dimensions.Position.X && x <= dimensions.Position.X+dimensions.Width &&
			y >= dimensions.Position.Y && y <= dimensions.Position.Y+dimensions.Height {
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
	case ScreenRelativeLeft:
		return (x / float32(screenWidth)), (y / float32(screenHeight))
	case ScreenRelativeRight:
		return (float32(screenWidth) - x) / float32(screenWidth),
			(float32(screenHeight) - y) / float32(screenHeight)
	case ScreenRelativeAspect:
		return (x / float32(screenWidth)) * screenRatio, (y / float32(screenHeight))
	}
	return -1, -1
}