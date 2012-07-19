package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
	"github.com/jteeuwen/glfw"
)

const (
	windowTitle = "Excavation"
)

//globally accessible camera and pipeline
var Cam horde3d.H3DNode
var pipeline *Resource
var running bool

func Init(name string) error {
	setAppName(name)
	//load settings from config file
	cfg, err := NewStandardCfg()
	if err != nil {
		return err
	}
	if err = cfg.Load(); err != nil {
		return err
	}

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

	if !horde3d.Init() {
		horde3d.DumpMessages()
		return errors.New("Error starting Horde3D.  Check Horde3D_log.html for more information")
	}

	//setup input handling
	controlCfg, _ := NewControlCfg()
	loadBindingsFromCfg(controlCfg)

	glfw.SetKeyCallback(keyCallback)
	glfw.SetMouseButtonCallback(mouseButtonCallback)
	glfw.SetMousePosCallback(mousePosCallback)
	glfw.SetMouseWheelCallback(mouseWheelCallback)

	//load pipeline
	pipeline, err = LoadPipeline()
	if err != nil {
		return err
	}

	//add camera
	//TODO: camera and node types
	Cam = horde3d.AddCameraNode(horde3d.RootNode, "Camera", pipeline.H3DRes)
	glfw.SetWindowSizeCallback(onResize)
	return nil

}

func StartMainLoop() {
	running = true

	for running {
		joyUpdate()
		runTasks()
		horde3d.Render(Cam)
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

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	//TODO:camera Type
	horde3d.SetNodeParamI(Cam, horde3d.Camera_ViewportXI, 0)
	horde3d.SetNodeParamI(Cam, horde3d.Camera_ViewportYI, 0)
	horde3d.SetNodeParamI(Cam, horde3d.Camera_ViewportWidthI, w)
	horde3d.SetNodeParamI(Cam, horde3d.Camera_ViewportHeightI, h)

	//TODO: Set clip distance? Config?
	horde3d.SetupCameraView(Cam, 45.0, float32(w)/float32(h), 0.1, 1000.0)
	horde3d.ResizePipelineBuffers(pipeline.H3DRes, w, h)

}

func Time() float64 {
	return glfw.Time()
}
