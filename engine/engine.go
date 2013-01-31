// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
	"github.com/jteeuwen/glfw"
)

//Used to hold resize parms between cameras 
// the render cam will lookup near and far plane,
// but FOV will have to be set separately
type renderCam struct {
	camera      *Camera
	fallbackCam *Camera
	fov         float32
	nearPlane   float32
	farPlane    float32
}

var Root *Node
var mainCam *renderCam
var running bool
var frames int
var startTime float64
var controlCfg *Config
var standardCfg *Config
var paused bool

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
	initPhysics()

	//setup base camera
	pipeline, err := loadDefaultPipeline()
	if err != nil {
		panic(err)
	}

	mainCam = new(renderCam)
	mainCam.fallbackCam = AddCamera(Root, "FallbackCamera", pipeline)
	mainCam.camera = mainCam.fallbackCam
	mainCam.fov = 45

	resizeView(Cfg().Int("WindowWidth"), Cfg().Int("WindowHeight"))

	//Music and Audio
	initMusic()
	initAudio(cfg.String("AudioDevice"), cfg.Int("MaxAudioSources"), cfg.Int("MaxAudioBufferSize"))
	glfw.SetWindowSizeCallback(resizeView)

	return nil

}

func StartMainLoop() {
	running = true
	paused = false

	resetView()
	startTime = Time()
	for running {
		frames++
		joyUpdate()
		if !paused {
			runTasks()
			updateAudio()
			updatePhysics()
		}
		updateGui()
		horde3d.Render(mainCam.camera.H3DNode)
		horde3d.FinalizeFrame()
		horde3d.ClearOverlays()
		glfw.SwapBuffers()
	}

	ClearAll()
	phWorld.Destroy()
	horde3d.Release()
	glfw.Terminate()
	glfw.CloseWindow()
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

func SetMainCamera(camera *Camera) {
	mainCam.camera = camera
	mainCam.nearPlane = horde3d.GetNodeParamF(camera.H3DNode, horde3d.Camera_NearPlaneF, 0)
	mainCam.farPlane = horde3d.GetNodeParamF(camera.H3DNode, horde3d.Camera_FarPlaneF, 0)
	resetView()
}

func SetCameraFOV(fov float32) {
	mainCam.fov = fov
	resetView()
}

func MainCamera() *Camera {
	return mainCam.camera
}

func resetView() {
	w, h := glfw.WindowSize()
	resizeView(w, h)
}

func resizeView(w, h int) {
	if h == 0 {
		h = 1
	}

	mainCam.camera.SetViewport(0, 0, w, h)

	mainCam.camera.SetupView(mainCam.fov, float32(w)/float32(h), mainCam.nearPlane, mainCam.farPlane)
	mainCam.camera.Pipeline().ResizeBuffers(w, h)
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
	clearAllAudio()
	clearAllPhysics()
	//TODO: Close compressed data file if open
	//horde3d.Clear()

	children := Root.Children()
	for i := range children {
		children[i].Remove()
	}

	resList := ResourceList()
	for i := range resList {
		resList[i].Remove()
	}

	horde3d.ReleaseUnusedResources()

	SetMainCamera(mainCam.fallbackCam)
	//LoadAllResources()
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
