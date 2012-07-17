package engine

import (
	"github.com/jteeuwen/glfw"
)

var (
	pressEvents    map[int]KeyHandler
	releaseEvents  map[int]KeyHandler
	mButtonPress   map[int]MouseButtonHandler
	mButtonRelease map[int]MouseButtonHandler
)

//Retyped from glfw in case I ever swap it out with a different lib
// so I don't have references in my game code to external libraries, 
// the engine should handle that
type KeyHandler func(key int)
type MousePosHandler func(x, y int)
type MouseButtonHandler func(button int)
type MouseWheelHandler func(delta int)

func init() {
	pressEvents = make(map[int]KeyHandler)
	releaseEvents = make(map[int]KeyHandler)
}

func keyCallBack(key int, state int) {
	var function KeyHandler
	var ok bool

	if state == glfw.KeyPress {
		function, ok = pressEvents[key]
	} else if state == glfw.KeyRelease {
		function, ok = releaseEvents[key]
	}

	if ok {
		function(key)
	}
}

func mButtonCallBack(button int, state int) {
	var function MouseButtonHandler
	var ok bool

	if state == glfw.KeyPress {
		function, ok = mButtonPress[button]
	} else if state == glfw.KeyRelease {
		function, ok = mButtonRelease[button]
	}

	if ok {
		function(button)
	}
}

func BindKeyPress(key int, function KeyHandler) {
	pressEvents[key] = function
}

func BindKeyRelease(key int, function KeyHandler) {
	releaseEvents[key] = function
}

func BindMButtonPress(button int, function MouseButtonHandler) {
	mButtonPress[button] = function
}

func BindMButtonRelease(button int, function MouseButtonHandler) {
	mButtonRelease[button] = function
}

func MousePos() (int, int) {
	return glfw.MousePos()
}

func MouseWheelPos() int {
	return glfw.MouseWheel()
}

func BindMousePos(function MousePosHandler) {
	glfw.SetMousePosCallback(glfw.MousePosHandler(function))
}

func BindMouseWheel(function MouseWheelHandler) {
	glfw.SetMouseWheelCallback(glfw.MouseWheelHandler(function))
}

//TODO: Joystick Detection and handling
