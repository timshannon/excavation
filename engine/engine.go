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
var pipeline horde3d.H3DRes
var running bool

func Init() error {

	//load settings from config file
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	if err = glfw.Init(); err != nil {
		return err
	}

	//TODO: Load from config
	if err := glfw.OpenWindow(cfg.Int("WindowWidth"),
		cfg.Int("WindowHeight"), 8, 8, 8, 8,
		cfg.Int("WindowDepth"), 8,
		cfg.Int("WindowMode")); err != nil {
		return err
	}

	//vsync from settings file
	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(windowTitle)

	if !horde3d.Init() {
		horde3d.DumpMessages()
		return errors.New("Error starting Horde3D.  Check Horde3D_log.html for more information")
	}

	//load pipeline
	//TODO: Resource manager
	pipeline = horde3d.AddResource(horde3d.ResTypes_Pipeline, "pipelines/hdr.pipeline.xml", 0)
	horde3d.LoadResourcesFromDisk("../content")

	//add camera
	Cam = horde3d.AddCameraNode(horde3d.RootNode, "Camera", pipeline)
	glfw.SetWindowSizeCallback(onResize)
	return nil

}

func StartMainLoop() {
	running = true

	for running {
		//TODO: taskmanager
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

	horde3d.SetNodeParamI(Cam, horde3d.Camera_ViewportXI, 0)
	horde3d.SetNodeParamI(Cam, horde3d.Camera_ViewportYI, 0)
	horde3d.SetNodeParamI(Cam, horde3d.Camera_ViewportWidthI, w)
	horde3d.SetNodeParamI(Cam, horde3d.Camera_ViewportHeightI, h)

	//TODO: Set clip distance? Config?
	horde3d.SetupCameraView(Cam, 45.0, float32(w)/float32(h), 0.1, 1000.0)
	horde3d.ResizePipelineBuffers(pipeline, w, h)

}
