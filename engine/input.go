package engine

import (
	"github.com/jteeuwen/glfw"
)

type InputHandler func(input Input)

type Input struct {
	Dev	Device 
	State int
	Position     int
	AxisPosition float32
}


type Device struct {
	DeviceType int
	DeviceIndex int
	Button int
	Axis int

//A named control is bound to a Input Type
// Game code refers to "strafeLeft" and gets the input type in the call
// to determine how to handle the event

func BindInput(controlName string, function InputHandler) {

}
