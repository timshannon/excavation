// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package engine

import (
	"bitbucket.org/tshannon/gohorde/horde3d"
	"github.com/jteeuwen/glfw"
)

//Used for both menus and HUDs

const (
	//position is based on the aspect ratio of the screen
	ScreenRelativeAspect = iota
	//position is independant of the aspect ratio and position from the left
	ScreenRelativeLeft
	//position is independant of the aspect ratio and position from the right (0 is right 1 is left)
	ScreenRelativeRight
)

var screenRatio float32
var screenHeight int
var screenWidth int
var tempArray [16]float32 //Rectangles only for now
var activeGuis []*Gui

func initGui() {
	glfw.SetCharCallback(charCollector)
	activeGuis = make([]*Gui, 0, 5)
}

//LoadGui pushes a gui onto a stack of guis, only top most in the stack
// is used for interaction
func LoadGui(gui *Gui) {
	activeGuis = append([]*Gui{gui}, activeGuis...)
	gui.load()
}

//UnloadGui pops a gui off the stack of guis
func UnloadGui() {
	if len(activeGuis) == 0 {
		return
	}

	gui := activeGuis[0]
	gui.unload()

	activeGuis = activeGuis[1:]
	if len(activeGuis) != 0 {
		//Reset input and mouse
		activeGuis[0].load()
	}
}

//UnloadAllGuis unload all the guis on the stack and resets
// the engine and inputs back to normal operation
func UnloadAllGuis() {
	for i := range activeGuis {
		activeGuis[i].unload()
	}
	activeGuis = activeGuis[0:0]
}

func updateGui() {
	horde3d.ClearOverlays()
	for i := range activeGuis {
		if i == 0 {
			activeGuis[i].handleInput()
		}
		if len(activeGuis) > 0 {
			activeGuis[i].update()
		}
	}

	updateDebugPrint()
}

func charCollector(key, state int) {
	if state == glfw.KeyPress {
		if gCharCollector != nil {
			gCharCollector(key)
		}
	}
}

type CharCollector func(key int)

var gCharCollector CharCollector

//toVertex returns the actual position on the screen from the interpreted
// relative position
func (s *ScreenArea) toVertex(result []float32) {
	//Y is not relative to aspect
	//verts are added counter clockwise
	//vert1
	result[0] = s.X()
	result[1] = s.Position.Y
	result[2] = 0
	result[3] = 1
	//vert2
	result[4] = s.X()
	result[5] = (s.Position.Y + s.Height)
	result[6] = 0
	result[7] = 0
	//vert3
	result[8] = s.X2()
	result[9] = (s.Position.Y + s.Height)
	result[10] = 1
	result[11] = 0
	//vert4
	result[12] = s.X2()
	result[13] = s.Position.Y
	result[14] = 1
	result[15] = 1
}

type ScreenPosition struct {
	X, Y       float32
	RelativeTo int
}

type ScreenArea struct {
	Position      *ScreenPosition
	Width, Height float32
}

//X1 Returns the first x vertex position
// based on the relative positioning of the screen area
func (s *ScreenArea) X() float32 {
	switch s.Position.RelativeTo {
	case ScreenRelativeAspect:
		return s.Position.X
	case ScreenRelativeLeft:
		return s.Position.X * screenRatio
	case ScreenRelativeRight:
		return (screenRatio - (s.Position.X * screenRatio)) - s.Width
	}
	return 0
}

//X2 Returns the second x (width) position
// based on the relative positioning of the screen area
func (s *ScreenArea) X2() float32 {
	switch s.Position.RelativeTo {
	case ScreenRelativeAspect:
		return (s.Position.X + s.Width)
	case ScreenRelativeLeft:
		return (s.Position.X * screenRatio) + s.Width
	case ScreenRelativeRight:
		return (screenRatio - (s.Position.X * screenRatio))
	}

	return 0
}

//PixelHeight is the height in actual pixels as relating to the current
// screen resolution
func (s *ScreenArea) PixelHeight() int {
	return int(float32(screenHeight) * s.Height)
}

//PixelWidth is the width in actual pixels as relating to the current
// screen resolution
func (s *ScreenArea) PixelWidth() int {
	switch s.Position.RelativeTo {
	case ScreenRelativeAspect:
		return int(s.Width * float32(screenWidth))
	case ScreenRelativeLeft, ScreenRelativeRight:
		return int((float32(screenWidth) / screenRatio) * s.Width)
	}
	return 0
}

func NewScreenArea(x, y, width, height float32, relativeTo int) *ScreenArea {
	return &ScreenArea{&ScreenPosition{x, y, relativeTo}, width, height}
}

//255 based color
// translates to horde float based color
type Color struct {
	r, g, b, a int
}

func NewColor(r, g, b, a int) *Color {
	return &Color{r, g, b, a}
}

func (c *Color) R() float32 {
	return c.toHordeColor(c.r)
}

func (c *Color) SetR(r int) { c.r = r }

func (c *Color) G() float32 {
	return c.toHordeColor(c.g)
}

func (c *Color) SetG(g int) { c.g = g }

func (c *Color) B() float32 {
	return c.toHordeColor(c.b)
}

func (c *Color) SetB(b int) { c.b = b }

func (c *Color) A() float32 {
	return c.toHordeColor(c.a)
}

func (c *Color) SetA(a int) { c.a = a }

func (c *Color) toHordeColor(color int) float32 {
	return float32(color) / 255.0
}

type Overlay struct {
	Dimensions *ScreenArea
	Color      *Color
	Material   *Material
}

func NewOverlay(materialLocation string, color *Color, dimensions *ScreenArea) *Overlay {
	material, err := NewMaterial(materialLocation)
	err = material.Load()
	if err != nil {
		RaiseError(err)
	}

	return &Overlay{dimensions, color, material}
}

func NewScreenPosition(X, Y float32, relativeTo int) *ScreenPosition {
	return &ScreenPosition{X, Y, relativeTo}
}

//BitmapText is text drawn on the screen using a bitmap based fonts
//  Can be faster for small sets of text vs freetype raster
//  Doesn't require regerating and reloading the resource
//  Limited by horde's overlay limit.
//  Recommended for use in small quick text changes (debug messages)
//  For large font sizes, or large amounts of text use the
//  Text type which uses a freetype rasterizer to genrate the text
//  and only ever uses one resource, but any changes made to that
//  text need to be re-rasterized
type BitmapText struct {
	Text         string
	Position     *ScreenPosition
	Size         float32
	FontMaterial *Material
	Color        *Color
}

func NewBitmapText(text string, size float32, materialLocation string,
	color *Color, position *ScreenPosition) *BitmapText {
	material, err := NewMaterial(materialLocation)
	err = material.Load()
	if err != nil {
		RaiseError(err)
	}
	return &BitmapText{text, position, size, material, color}
}

func (t *BitmapText) Place() {
	var newX float32
	switch t.Position.RelativeTo {
	case ScreenRelativeAspect:
		newX = t.Position.X
	case ScreenRelativeLeft:
		newX = t.Position.X * screenRatio
	case ScreenRelativeRight:
		newX = (screenRatio - (t.Position.X * screenRatio)) - t.Width()
	}
	horde3d.ShowText(t.Text, newX, t.Position.Y, t.Size, t.Color.R(),
		t.Color.G(), t.Color.B(), t.FontMaterial.H3DRes)

}

func (t *BitmapText) Width() float32 {
	return (float32(len(t.Text)) * t.Size * .5) + (t.Size * .25)
}

func (o *Overlay) Place() {
	o.Dimensions.toVertex(tempArray[:])
	horde3d.ShowOverlays(tempArray[:], 4, o.Color.R(), o.Color.G(),
		o.Color.B(), o.Color.A(), o.Material.H3DRes, 0)
}

//Widget is a collection of Overlays
type Widget interface {
	Name() string
	MouseArea() *ScreenArea
	Update()
	Hover()
	Click(int)
	Scroll(int)
	Unload()
}

//Gui is a collection of Widgets
type Gui struct {
	Widgets       []Widget
	UseMouse      bool
	HaltInput     bool
	CharCollect   CharCollector
	prevTime      float64
	prevWheelPos  int
	mousePress    [8]bool
	inputs        *inputGroup
	prevMousePosX int
	prevMousePosY int
}

func NewGui() *Gui {
	gui := new(Gui)
	gui.inputs = newInputGroup()
	return gui
}

func (g *Gui) Bind(function InputHandler, input string) {
	g.inputs.bind(function, input, input)
}

func (g *Gui) ElapsedTime() float64 {
	return (glfw.Time() - g.prevTime)
}

//AddWidget adds a widget to the last / top location
// of the gui
func (g *Gui) AddWidget(widget Widget) {
	g.Widgets = append(g.Widgets, widget)
}

//Removes a widget from the gui.
func (g *Gui) RemoveWidget(name string) {
	for i := 0; i < len(g.Widgets); i++ {
		if g.Widgets[i].Name() == name {
			if len(g.Widgets) > 1 {
				g.Widgets = append(g.Widgets[:i], g.Widgets[i+1:]...)
			} else {
				g.Widgets = g.Widgets[0:0]
			}
		}
	}
}

func (g *Gui) load() {
	//TODO; Might be overkill
	err := LoadAllResources()
	if err != nil {
		RaiseError(err)
	}

	if g.HaltInput {
		loadInputGroup(g.inputs)
	} else {
		unloadInputGroup()
	}

	if g.UseMouse {
		g.prevMousePosX, g.prevMousePosY = glfw.MousePos()

		glfw.Enable(glfw.MouseCursor)
	} else {
		glfw.Disable(glfw.MouseCursor)
	}
	gCharCollector = g.CharCollect
}

func (g *Gui) unload() {
	glfw.Disable(glfw.MouseCursor)
	glfw.SetMousePos(g.prevMousePosX, g.prevMousePosY)
	glfw.PollEvents()
	unloadInputGroup()
	gCharCollector = nil
	for i := range g.Widgets {
		g.Widgets[i].Unload()
	}
}

func (g *Gui) mouseClick(button int) bool {
	if glfw.MouseButton(button) == glfw.KeyPress {
		g.mousePress[button] = true
		return false
	} else if glfw.MouseButton(button) == glfw.KeyRelease {
		if g.mousePress[button] {
			g.mousePress[button] = false
			return true
		}
	}
	return false
}

//handleInput only processes input for the topmost
// gui on the stack, basically creating modal guis
// as well has menus on top of game huds or other game guis
func (g *Gui) handleInput() {
	if g.UseMouse {
		if widget, ok := g.WidgetUnderMouse(); ok {
			widget.Hover()
			for i := range g.mousePress {
				if g.mouseClick(i) {
					widget.Click(i)
				}
			}
			delta := glfw.MouseWheel()
			if delta != g.prevWheelPos {
				//TODO: Test delta
				widget.Scroll(g.prevWheelPos - delta)
			}
		}
	}

}

func (g *Gui) update() {

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
		x, y = g.MousePos(ScreenRelativeAspect)
		if x >= dimensions.X() && x <= dimensions.X2() &&
			y >= dimensions.Position.Y && y <= dimensions.Position.Y+dimensions.Height {
			return g.Widgets[i], true

		}
	}

	return nil, false
}

func updateGuiScreenSize(w, h int) {
	screenHeight = h
	screenWidth = w
	screenRatio = float32(w) / float32(h)
}

func ScreenRatio() float32 {
	return screenRatio
}

func (g *Gui) MousePos(relative int) (x, y float32) {
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
			(y / float32(screenHeight))
	case ScreenRelativeAspect:
		return (x / float32(screenWidth)) * screenRatio, (y / float32(screenHeight))
	}
	return -1, -1
}
