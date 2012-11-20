// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
	"github.com/jteeuwen/glfw"
)

const (
	windowTitle = "Excavation"
)

//1 engine unit = 1 meter

var Root *Node
var MainCam *Camera
var running bool
var frames int
var startTime float64
var controlCfg *Config
var standardCfg *Config

func init() {
	Root = new(Node)
	Root.H3DNode = horde3d.RootNode
}

func Init(name string) error {
	appName = name
	//load settings from config file
	cfg, err := NewStandardCfg()

	if err != nil {
		return err
	}

	if err = cfg.Load(); err != nil {
		return err
	}

	standardCfg = cfg
	if err = glfw.Init(); err != nil {
		return err
	}

	var mode int
	if cfg.Bool("Fullscreen") {
		mode = glfw.Fullscreen
	} else {
		mode = glfw.Windowed
	}

	if err := glfw.OpenWindow(cfg.Int("WindowWidth"),
		cfg.Int("WindowHeight"), 8, 8, 8, 8,
		cfg.Int("WindowDepth"), 8,
		mode); err != nil {
		return err
	}

	glfw.SetSwapInterval(cfg.Int("VSync"))
	glfw.SetWindowTitle(windowTitle)
	glfw.Disable(glfw.MouseCursor)

	if !horde3d.Init() {
		horde3d.DumpMessages()
		return errors.New("Error starting Horde3D.  Check Horde3D_log.html for more information")
	}

	//setup input handling
	controls, err := NewControlCfg()

	controlCfg = controls
	if err != nil {
		return err
	}
	controlCfg.Load()
	loadBindingsFromCfg(controlCfg)

	glfw.SetKeyCallback(keyCallback)
	glfw.SetMouseButtonCallback(mouseButtonCallback)
	glfw.SetMousePosCallback(mousePosCallback)
	glfw.SetMouseWheelCallback(mouseWheelCallback)

	//load pipeline
	pipeline, err := LoadPipeline()
	if err != nil {
		return err
	}

	//add camera
	MainCam = AddCamera(Root, "MainCamera", pipeline)

	//Music and Audio
	initMusic()
	initAudio(cfg.String("AudioDevice"), cfg.Int("MaxAudioSources"), cfg.Int("MaxAudioBufferSize"))
	glfw.SetWindowSizeCallback(onResize)

	return nil

}

func StartMainLoop() {
	running = true

	//TODO: fix missed onresize call
	onResize(1024, 768)

	startTime = Time()
	for running {
		frames++
		joyUpdate()
		runTasks()
		updateAudio()
		horde3d.Render(MainCam.H3DNode)
		horde3d.FinalizeFrame()
		glfw.SwapBuffers()

		//TODO: handle with input and tasks
		running = glfw.Key(glfw.KeyEsc) == 0 &&
			glfw.WindowParam(glfw.Opened) == 1

	}

	horde3d.Release()
	glfw.Terminate()
	glfw.CloseWindow()
}

func StopMainLoop() {
	running = false
}

func Fps() float64 {
	return float64(frames) / (Time() - startTime)
}

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	MainCam.SetViewport(0, 0, w, h)

	//TODO: Set clip distance? Config?
	MainCam.SetupView(45.0, float32(w)/float32(h), 0.1, 1000.0)
	MainCam.Pipeline().ResizeBuffers(w, h)

}

func Time() float64 {
	return glfw.Time()
}

//Clear clears all rendering, physics, and sound resources, nodes, etc
func ClearAll() {
	horde3d.Clear()
	ClearAllAudio()
	//TODO: Clear Physics entities
	//TODO: Close compressed data file if open
}

func Cfg() *Config {
	return standardCfg
}

func ControlCfg() *Config {
	return controlCfg
}
