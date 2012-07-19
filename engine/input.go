package engine

import (
	"strconv"
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
//  x,y is the mouse position on the x or y axis
//  AxisPosition is the 1.0 to -1.0 position on a joystick axis
//  Unused variables will be -1, such as state for a axis input
type Input struct {
	Device       Device
	State        int
	X            int
	Y            int
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

func init() {
}

//For printing to controls config file
func (d *Device) String() string {
	var prefix string
	var suffix string

	switch d.DeviceType {
	case DeviceKeyboard:
		prefix = "Key"
		suffix = strconv.Itoa(d.Button)
	case DeviceMouse:
		prefix = "Mouse"
		if d.Button >= 0 {
			suffix = strconv.Itoa(d.Button)
		} else {
			suffix = "Axis"
		}
	case DeviceJoystick:
		prefix = "Joy" + strconv.Itoa(d.DeviceIndex)
		if d.Button >= 0 {
			suffix = strconv.Itoa(d.Button)
		} else {
			suffix = "Axis" + strconv.Itoa(d.Axis)
		}

	}

	return prefix + "_" + suffix
}

//New Device creates a new Device from a control config string
func newDevice(cfgString string) *Device {
	dev := new(Device)
	var prefix string
	var suffix string
	str := strings.Split(cfgString, "_")
	prefix = str[0]
	suffix = str[1]

	switch {
	case prefix == "Key":
		dev.DeviceType = DeviceKeyboard
		dev.Button, _ = strconv.Atoi(suffix)
	case prefix == "Mouse":
		dev.DeviceType = DeviceMouse
		if !strings.HasPrefix(suffix, "Axis") {
			dev.Button, _ = strconv.Atoi(suffix)
		}

	case strings.Contains(prefix, "Joy"):
		dev.DeviceType = DeviceJoystick
		dev.DeviceIndex, _ = strconv.Atoi(strings.TrimLeft(prefix, "Joy"))
		if strings.HasPrefix(suffix, "Axis") {
			dev.Axis, _ = strconv.Atoi(strings.TrimLeft(suffix, "Axis"))
		} else {
			//Button
			dev.Button, _ = strconv.Atoi(suffix)
		}
	}

	return dev
}

//TODO: Setup callbacks on control cfg load

//A named control is bound to a Input Type
// Game code refers to the control name "strafeLeft" and gets the input type in the call
// to determine how to handle the event

func BindInput(controlName string, function InputHandler) {
} //loadBindingFromCfg loads and binds the inputs, creating indexed device
// entries for use with input callbacks
func loadBindingsFromCfg(cfg *Config) {

}

func keyCallBack(key, state int) {

}

func mouseButtonCallback(button, state int) {

}

func mousePosCallback(x, y int) {

}
