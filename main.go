// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package main

import (
	"excavation/engine"
	"excavation/entity"
	"flag"
	"os"
	"strings"
)

//cmd line options
var (
	sceneFlag string
	camera    *engine.Node
)

func init() {
	flag.StringVar(&sceneFlag, "scene", "", "Load a specific scene directly, instead of the main menu.")

	flag.Parse()
}

func main() {

	engine.SetDefaultConfigHandler(setCfgDefaults)

	if err := engine.Init("excavation"); err != nil {
		panic("Error starting Excavation: " + err.Error())
	}

	if sceneFlag != "" {
		loadScene(sceneFlag)
	} else {
		//TODO: Load Main Menu
	}

	//starting the loop should be the last thing
	// after setting up the game
	engine.StartMainLoop()
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

//loadScene loads a scene and all associated resources in the given
//scenefile and loads the entities and properties
func loadScene(scene string) {
	if !strings.HasSuffix(scene, ".scene.xml") {
		scene = scene + ".scene.xml"
	}
	//Clear any old scene data and resources
	engine.ClearAll()
	sceneRes, err := engine.NewScene(scene)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(sceneRes.FullPath()); err != nil {
		if os.IsNotExist(err) {
			//TODO: Load Main Menu instead
			panic("Scene file " + scene + " doesn't exist")

		}
	}
	err = sceneRes.Load()

	err = engine.LoadAllResources()
	sceneNode, err := engine.AddNodes(engine.Root, sceneRes)

	if err != nil {
		panic(err)
	}

	children := sceneNode.Children()
	for c := range children {
		//load entities
		if children[c].Attachment() != "" {
			err = entity.LoadEntity(children[c], children[c].Attachment())
		}
		if err != nil {
			panic(err)
		}

	}
}
