package engine

import (
	"github.com/jteeuwen/glfw"
)

var pressEvents map[int]func()
var releaseEvents map[int]func()

func init() {
	pressEvents = make(map[int]func())
	releaseEvents = make(map[int]func())
}

func keyHandler(key int, state int) {
	var function func()
	var ok bool

	if state == glfw.KeyPress {
		function, ok = pressEvents[key]
	} else if state == glfw.KeyRelease {
		function, ok = releaseEvents[key]
	}

	if ok {
		function()
	}
}

func BindKeyPress(key int, function func()) {
	pressEvents[key] = function
}

func BindKeyRelease(key int, function func()) {
	releaseEvents[key] = function
}

func MousePos() (int, int) {
	return glfw.MousePos()
}
