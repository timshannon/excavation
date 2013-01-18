// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
	"fmt"
	"github.com/jteeuwen/glfw"
)

var Root *Node
var MainCam *Camera
var running bool
var frames int
var startTime float64
var controlCfg *Config
var standardCfg *Config
var paused bool

//exported variable for changing the clip distance on the fly if need be
var ClipPlaneDistance float32 = 1000

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
	glfw.SetWindowTitle(name)
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

	initInput()

	initGui()
	setupRenderer()

	//Music and Audio
	initMusic()
	initAudio(cfg.String("AudioDevice"), cfg.Int("MaxAudioSources"), cfg.Int("MaxAudioBufferSize"))
	glfw.SetWindowSizeCallback(onResize)

	return nil

}

func StartMainLoop() {
	running = true
	paused = false

	onResize(Cfg().Int("WindowWidth"), Cfg().Int("WindowHeight"))
	startTime = Time()
	for running {
		frames++
		joyUpdate()
		if !paused {
			runTasks()
			updateAudio()
			//TODO: Physics
		}
		updateGui()
		horde3d.Render(MainCam.H3DNode)
		horde3d.FinalizeFrame()
		horde3d.ClearOverlays()
		glfw.SwapBuffers()
	}

	horde3d.Release()
	glfw.Terminate()
	glfw.CloseWindow()
}

func setupRenderer() {
	pipeline, err := loadDefaultPipeline()
	if err != nil {
		panic(err)
	}

	//add camera
	MainCam = AddCamera(Root, "MainCamera", pipeline)
	onResize(Cfg().Int("WindowWidth"), Cfg().Int("WindowHeight"))

}

func StopMainLoop() {
	running = false
}

func Fps() float64 {
	fps := float64(frames) / (Time() - startTime)
	frames = 0
	startTime = Time()
	return fps
}

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	MainCam.SetViewport(0, 0, w, h)

	MainCam.SetupView(45.0, float32(w)/float32(h), 0.1, ClipPlaneDistance)
	MainCam.Pipeline().ResizeBuffers(w, h)
	updateGuiScreenSize(w, h)

}

func Time() float64 {
	return glfw.Time()
}

var pauseStart, pausedTime float64

//Game time is the actual game time
// not including the time paused.  When the game is
// paused, the game time will not increment
func GameTime() float64 {
	return Time() - pausedTime
}

//Clear clears all rendering, physics, and sound resources, nodes, etc
func ClearAll() {
	removeAllTasks()
	UnloadAllGuis()
	ClearAllAudio()
	//TODO: Clear Physics entities
	//TODO: Close compressed data file if open
	//horde3d.Clear()
	children := Root.Children()
	for i := range children {
		//if children[i].Type() != NodeTypeCamera {
		children[i].Remove()
		//}
	}

	resList := ResourceList()
	for i := range resList {
		resList[i].Remove()
	}

	resList = ResourceList()
	for i := range resList {
		//wtf, remove doesn't remove
		fmt.Println("res: ", resList[i].Name())

	}

	fmt.Println("Clear")

	setupRenderer()
}

func Pause() {
	paused = true
	pauseStart = Time()
	pauseAllAudio()
	PauseMusic()
	//TODO: Pause physics?
}

func Resume() {
	paused = false
	pausedTime += Time() - pauseStart
	resumeAllAudio()
	ResumeMusic()
	//TODO: Resume Physics?

}

func Cfg() *Config {
	return standardCfg
}

func ControlCfg() *Config {
	return controlCfg
}
