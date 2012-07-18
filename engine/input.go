package engine

import (
	"github.com/jteeuwen/glfw"
	"strings"
)

const (
	StatePressed = iota
	StateReleased
)

const (
	DeviceKeyboard = iota
	DeviceMouse
	DeviceJoystick
)

const (
	MouseAxisX = iota
	MouseAxisY
)

type InputHandler func(input Input)

//Input is the current values of the given input
//  Device is the source of the input
//  State is if the button or key is pressed or released
//  Position is the mouse position on the x or y axis
//  AxisPosition is the 1.0 to -1.0 position on a joystick axis
//  Unused variables will be -1, such as state for a axis input
type Input struct {
	Device       Device
	State        int
	Position     int
	AxisPosition float32
}

//Device is the source of a type of input
// Type is either Keyboard, Mouse or Joystick
// DeviceIndex refers to which joystick 0 - 15
//	Mouse and Keyboard have indexes of -1	
// Button is which key or button the keyboard, mouse or joystick got pressed
//	-1 if the input source is an axial movement (mouse or joystick)
// Axis is the index of the axis that was the source of the input
//	Joystick axies are unlimited
//	Mouse 0 is x movement
//	Mouse 1 is y movement
type Device struct {
	DeviceType  int
	DeviceIndex int
	Button      int
	Axis        int
}

var controlInputs map[string]InputHandler

func init() {
	controlInputs = make(map[string]InputHandler)
}

//TODO: Setup callbacks on control cfg load

//A named control is bound to a Input Type
// Game code refers to the control name "strafeLeft" and gets the input type in the call
// to determine how to handle the event

func BindInput(controlName string, function InputHandler) {
	controlInputs[controlName] = function
}

//For printing to controls config file
func (d *Device) String() string {
	var prefix string
	var suffix string

	switch d.DeviceType {
	case DeviceKeyboard:
		prefix = "Key"
		if d.Button < 256 {
			suffix = string(d.Button)
		} else {
			//TODO: Lookup enumeration name? Reflection?
		}
	case DeviceMouse:
		prefix = "Mouse"
		if d.Button >= 0 {
			suffix = string(d.Button)
		} else if d.Axis == MouseAxisX {
			suffix = "AxisX"
		} else if d.Axis == MouseAxisY {
			suffix = "AxisY"
		}
	case DeviceJoystick:
		prefix = "Joy" + string(d.DeviceIndex)
		if d.Button >= 0 {
			suffix = string(d.Button)
		} else {
			suffix = "Axis" + string(d.Axis)
		}

	}

	return prefix + "_" + suffix
}

func getGlfwKeyName(key int) string {
	//TODO: Array of special key names
	return string(glfw.KeyBackspace)
}

//New Device creates a new Device from a control config string
func newDevice(cfgString string) *Device {
	dev := new(Device)
	var prefix string
	var suffix string
	str := strings.Split(deviceString, "_")
	prefix = str[0]
	suffix = str[1]

	switch prefix {
	case "Key":
		dev.DeviceType = DeviceKeyboard
		dev.Button = int(suffix)
	case "Mouse":
		dev.DeviceType = DeviceMouse
		if strings.HasPrefix(suffix, "AxisX") {
			dev.Axis = MouseAxisX
		} else if strings.HasPrefix(suffix, "AxisY") {
			dev.Axis = MouseAxisY
		} else {
			dev.Button = int(suffix)
		}

	case "Joy": //fix
		dev.DeviceType = DeviceJoystick
	}

	return dev
}
