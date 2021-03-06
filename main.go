// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"excavation/engine"
	"excavation/entity"
	"flag"
	"fmt"
	"github.com/jteeuwen/glfw"
	"runtime"
	"strings"
)

const (
	name = "excavation"
)

//cmd line options
var (
	sceneFlag string
)

func init() {
	flag.StringVar(&sceneFlag, "scene", "", "Load a specific scene directly, instead of the main menu.")
}

func main() {
	flag.Parse()

	engine.SetDefaultConfigHandler(setCfgDefaults)
	engine.SetErrorHandler(errHandler)

	if err := engine.Init(name); err != nil {
		panic("Error starting Excavation: " + err.Error())
	}

	if sceneFlag != "" {
		loadScene(sceneFlag)
	} else {
		loadMainMenu()
	}

	//Bind Esc to menu
	engine.BindInput(loadMenu, "Key_Esc")

	//todo: temp for testing frame independence
	engine.BindInput(ToggleVSync, "Key_F3")

	//starting the loop should be the last thing
	// after setting up the game
	engine.StartMainLoop()
}

func loadMenu(input *engine.Input) {
	if state, ok := input.ButtonState(); ok && state == engine.StateReleased {
		loadMainMenu()
	}
}

//ResetEngine Reloads all config from disk and reopens a new glfw window
// used for after video settings are changed
func ResetEngine() {
	engine.ClearAll()
	engine.StopMainLoop()
	//FIXME
	main()
}

func errHandler(err error) {
	fmt.Println(err)
}

func showFPS(t *engine.Task) {
	engine.Print("FPS: ", engine.Fps())
	t.Wait(0.25)
}

var vsync int

func ToggleVSync(input *engine.Input) {
	if state, ok := input.ButtonState(); ok {
		if state == engine.StatePressed {
			if vsync == 0 {
				vsync = 1
			} else {
				vsync = 0
			}
			glfw.SetSwapInterval(vsync)
		}
	}
}

//loadScene loads a scene and all associated resources in the given
//scenefile and loads the entities and properties
func loadScene(scene string) {
	if !strings.HasSuffix(scene, ".scene.xml") {
		scene = scene + ".scene.xml"
	}
	//Clear any old scene data and resources
	engine.ClearAll()
	runtime.GC()
	//TODO:  Loading screen, and camera management

	sceneRes, err := engine.NewScene(scene)
	if err != nil {
		panic(err)
	}

	err = sceneRes.Load()
	if err != nil {
		//TODO: Load Main Menu instead
		panic("Scene file " + scene + " doesn't exist")
	}

	err = engine.LoadAllResources()
	sceneNode, err := engine.Root.AddScene(sceneRes)

	if err != nil {
		panic(err)
	}

	children := sceneNode.Children()
	for c := range children {
		//load entities
		//TODO: Recurse Tree? Only first level entities get loaded now
		if children[c].Attachment() != "" {
			err = entity.LoadEntity(children[c], children[c].Attachment())
		}
		if err != nil {
			panic(err)
		}

	}
	//TODO: gui
	engine.AddTask("FPS", showFPS, nil, 0, 0.25)
}

func setCfgDefaults(cfg *engine.Config) {
	switch cfg.Name {
	case "excavation.cfg":
		cfg.SetValue("WindowWidth", 1024)
		cfg.SetValue("WindowHeight", 728)
		cfg.SetValue("WindowDepth", 24)
		cfg.SetValue("Fullscreen", false)
		cfg.SetValue("VSync", 0)
		cfg.SetValue("InvertMouse", true)
		cfg.SetValue("MouseSensitivity", 0.3)
		cfg.SetValue("AudioDevice", "")
		cfg.SetValue("MaxAudioSources", 16)
		cfg.SetValue("MaxAudioBufferSize", 5242880)
	case "controls.cfg":
		cfg.SetValue("Forward", "Key_W")
		cfg.SetValue("Backward", "Key_S")
		cfg.SetValue("StrafeLeft", "Key_A")
		cfg.SetValue("StrafeRight", "Key_D")
		cfg.SetValue("MoveUp", "Key_E")
		cfg.SetValue("MoveDown", "Key_Space")
		cfg.SetValue("PitchYaw", "Mouse_Axis0")
		cfg.SetValue("PitchUp", "Key_Up")
		cfg.SetValue("PitchDown", "Key_Down")
		cfg.SetValue("YawLeft", "Key_Left")
		cfg.SetValue("YawRight", "Key_Right")
	}

}
