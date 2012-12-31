// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"github.com/jteeuwen/glfw"
	"strconv"
	"strings"
)

//A named control is bound to a Input Type
// Game code refers to the control name "strafeLeft" and gets the input type in the call
// to determine how to handle the event
// Input types should be created and indexed on load, so as not to have to
// allocate a bunch of types over and over again with each input

const (
	StateReleased = iota
	StatePressed
)

const (
	DeviceKeyboard = iota
	DeviceMouse
	DeviceJoystick
)

const (
	MouseAxisPos = iota
	MouseAxisWheel
)

var (
	mouseAxisInputs map[int]*Input
	mouseBtnInputs  map[int]*Input
	keyInputs       map[int]*Input
	joyAxisInputs   map[int]*Input
	joyBtnInputs    map[int]*Input
	haltInput       bool
)

var inputHandlers map[string]InputHandler

func initInput() {
	resetInputs()

	glfw.SetKeyCallback(keyCallback)
	glfw.SetMouseButtonCallback(mouseButtonCallback)
	glfw.SetMousePosCallback(mousePosCallback)
	glfw.SetMouseWheelCallback(mouseWheelCallback)

	//Reload configs on write
	controlCfg.RegisterOnWriteHandler(reloadBindingsFromCfg)
	loadBindingsFromCfg(controlCfg)
}

func resetInputs() {
	mouseAxisInputs = make(map[int]*Input)
	mouseBtnInputs = make(map[int]*Input)
	keyInputs = make(map[int]*Input)
	joyAxisInputs = make(map[int]*Input)
	joyBtnInputs = make(map[int]*Input)
	inputHandlers = make(map[string]InputHandler)
}

type InputHandler func(input *Input)

//Input is the current values of the given input
//  Device is the source of the input
//  State is if the button or key is pressed or released
//  x,y is the mouse position on the x or y axis
//  x is also used for the mouse wheel position
//  AxisPosition is the 1.0 to -1.0 position on a joystick axis
type Input struct {
	controlName string
	Device      *Device
	State       int
	X           int
	Y           int
	AxisPos     float32
}

// JoyAxis returns the joystick's current axis value if the input is
// from a joystick axis, otherwise ok == false
func (input *Input) JoyAxis() (axisPos float32, ok bool) {
	if input.Device.Type == DeviceJoystick &&
		input.Device.Axis != -1 {
		axisPos = input.AxisPos
		ok = true
		return
	}

	return
}

//MousePos returns the mouse position if the input is
// from the mouse, otherwise ok == false
func (input *Input) MousePos() (x, y int, ok bool) {
	if input.Device.Type == DeviceMouse &&
		input.Device.Axis == MouseAxisPos {

		x = input.X
		y = input.Y
		ok = true
		return
	}

	return
}

//MouseWheelPos returns the position of the mouse wheel if
// the input is from a mouse wheel, otherwise ok == false
func (input *Input) MouseWheelPos() (delta int, ok bool) {
	if input.Device.Type == DeviceMouse &&
		input.Device.Axis == MouseAxisWheel {
		delta = input.X
		ok = true
		return
	}

	return
}

//ButtonState returns the state of the pressed button if the input
// is from a button press, otherwise ok == false
func (input *Input) ButtonState() (state int, ok bool) {
	if input.Device.Button != -1 {
		state = input.State
		ok = true
		return
	}

	return
}

func (input *Input) ControlName() string {
	return input.controlName
}

//Device is the source of a type of input
// Type is either Keyboard, Mouse or Joystick
// DeviceIndex refers to which joystick 0 - 15
//	Mouse and Keyboard have indexes of -1	
// Button is which key or button the keyboard, mouse or joystick got pressed
//	-1 if the input source is an axial movement (mouse or joystick)
// Axis is the index of the axis that was the source of the input
//	Joystick axies are unlimited
// 	mouse axis 1  is wheel
type Device struct {
	Type   int
	Index  int
	Button int
	Axis   int
}

//For printing to controls config file
func (d *Device) String() string {
	var prefix string
	var suffix string

	switch d.Type {
	case DeviceKeyboard:
		prefix = "Key"
		suffix = keyString(d.Button)
	case DeviceMouse:
		prefix = "Mouse"
		if d.Button >= 0 {
			suffix = strconv.Itoa(d.Button)
		} else {
			suffix = "Axis" + strconv.Itoa(d.Axis)
		}
	case DeviceJoystick:
		prefix = "Joy" + strconv.Itoa(d.Index)
		if d.Button >= 0 {
			suffix = strconv.Itoa(d.Button)
		} else {
			suffix = "Axis" + strconv.Itoa(d.Axis)
		}

	}

	return prefix + "_" + suffix
}

func keyString(key int) string {
	if key == 32 {
		return "Space"
	}
	if key > 256 && key <= 324 {
		return specialKeyString[key-257]
	}
	return string(key)
}

func keyInt(key string) int {
	if key == "Space" {
		return 32
	}
	if keyint, ok := specialKeyInt[key]; ok {
		return keyint
	}
	return int(key[0])

}

//joystick is used to store information about the
// current configured joystick
type joystick struct {
	index   int
	buttons []byte
	axes    []float32
}

var curJoystick *joystick

//New Device creates a new Device from a control config string
func newInput(cfgString string) *Input {
	dev := &Device{-1, -1, -1, -1}
	input := new(Input)

	var prefix string
	var suffix string
	str := strings.Split(cfgString, "_")
	prefix = str[0]
	suffix = str[1]

	switch {
	case prefix == "Key":
		dev.Type = DeviceKeyboard
		dev.Button = keyInt(suffix)
	case prefix == "Mouse":
		dev.Type = DeviceMouse
		if strings.HasPrefix(suffix, "Axis") {
			dev.Axis, _ = strconv.Atoi(strings.TrimLeft(suffix, "Axis"))
		} else {
			dev.Button, _ = strconv.Atoi(suffix)
		}

	case strings.Contains(prefix, "Joy"):
		dev.Type = DeviceJoystick
		dev.Index, _ = strconv.Atoi(strings.TrimLeft(prefix, "Joy"))
		if strings.HasPrefix(suffix, "Axis") {
			dev.Axis, _ = strconv.Atoi(strings.TrimLeft(suffix, "Axis"))
		} else {
			//Button
			dev.Button, _ = strconv.Atoi(suffix)
		}
	}

	input.Device = dev
	return input
}

//BindInput binds a config input entry to a input handler function
// Note, you can bind multiple controlsNames to the same function in one call
// ex: BindInput(function, "Strafe_Left", "Strafe_Right")
func BindInput(function InputHandler, controlName ...string) {
	for i := range controlName {
		inputHandlers[controlName[i]] = function
	}
}

//BindDirectInput binds an input to a function directly
// without requiring a config entry for it
// ex: BindDirectInput(function, "Key_A", "Key_Space")
func BindDirectInput(function InputHandler, input ...string) {
	for i := range input {
		addBinding(input[i], input[i])
		inputHandlers[input[i]] = function
	}
}

//loadBindingFromCfg loads and binds the inputs, creating indexed input
// entries for use with input callbacks
func loadBindingsFromCfg(cfg *Config) {
	for k, v := range cfg.values {
		addBinding(k, v.(string))
	}
}

func addBinding(controlName, cfgString string) {
	input := newInput(cfgString)
	input.controlName = controlName
	device := input.Device

	switch device.Type {
	case DeviceKeyboard:
		keyInputs[device.Button] = input
	case DeviceMouse:
		if device.Button != -1 {
			mouseBtnInputs[device.Button] = input
		} else {
			mouseAxisInputs[device.Axis] = input
		}
	case DeviceJoystick:
		if device.Button != -1 {
			joyBtnInputs[device.Button] = input
		} else {
			joyAxisInputs[device.Axis] = input
		}
		if curJoystick == nil {
			//currently not supporting multiple binds
			// from multiple joysticks
			// the current joystick is set to the first
			// bound joystick
			curJoystick = new(joystick)
			curJoystick.index = device.Index
			numAxes := glfw.JoystickParam(device.Index, glfw.Axes)
			numButtons := glfw.JoystickParam(device.Index, glfw.Buttons)
			curJoystick.axes = make([]float32, numAxes)
			curJoystick.buttons = make([]byte, numButtons)
		}

	}
}

//keyCallBack handles the glfw callback and executes the configured
// inputhandler for the given input
func keyCallback(key, state int) {
	input, ok := keyInputs[key]
	if ok && !haltInput {
		input.State = state
		if function, ok := inputHandlers[input.controlName]; ok {
			function(input)
		}
	}
}

//mouseButtonCallBack handles the glfw callback and executes the configured
// inputhandler for the given input
func mouseButtonCallback(button, state int) {
	input, ok := mouseBtnInputs[button]
	if ok && !haltInput {
		input.State = state
		if function, ok := inputHandlers[input.controlName]; ok {
			function(input)
		}
	}
}

//mousePosCallBack handles the glfw callback and executes the configured
// inputhandler for the given input
func mousePosCallback(x, y int) {
	input, ok := mouseAxisInputs[MouseAxisPos]
	if ok && !haltInput {
		input.X = x
		input.Y = y
		if function, ok := inputHandlers[input.controlName]; ok {
			function(input)
		}
	}
}

func mouseWheelCallback(delta int) {
	input, ok := mouseAxisInputs[MouseAxisWheel]
	if ok && !haltInput {
		input.X = delta
		if function, ok := inputHandlers[input.controlName]; ok {
			function(input)
		}
	}
}

//joyUpdate updates the joystick input values and executes the configured
// input handler for the given input
func joyUpdate() {
	if curJoystick != nil {
		var results int
		var input *Input
		var ok bool
		results = glfw.JoystickButtons(curJoystick.index, curJoystick.buttons)

		for i := 0; i < results; i++ {
			input, ok = joyBtnInputs[i]
			if ok && !haltInput {
				input.State = int(curJoystick.buttons[i])
				if function, ok := inputHandlers[input.controlName]; ok {
					function(input)
				}
			}
		}

		results = glfw.JoystickPos(curJoystick.index, curJoystick.axes)
		for i := 0; i < results; i++ {
			input, ok = joyAxisInputs[i]
			if ok {
				input.AxisPos = curJoystick.axes[i]
				if function, ok := inputHandlers[input.controlName]; ok {
					function(input)
				}
			}
		}
	}
}

func MousePos() (int, int) {
	return glfw.MousePos()
}

func SetMousePos(x, y int) {
	glfw.SetMousePos(x, y)
}

func HaltInput() {
	//TODO: Fix, input groups?
	// halt some inputs while allowing others
	haltInput = true
}

func ResumeInput() {
	haltInput = false
}

func reloadBindingsFromCfg(cfg *Config) {
	resetInputs()
	loadBindingsFromCfg(cfg)
}
