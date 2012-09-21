// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package main

import (
	"excavation/engine"
	"flag"
	"fmt"
	"os"
	"strings"
)

//cmd line options
var (
	sceneFlag string
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

	//engine.AddTask("test", task, 0, 1)
	//starting the loop should be the last thing
	// after setting up the game
	engine.StartMainLoop()
}

func task(t *engine.Task) {
	children := engine.Root.Children()

	fmt.Println(len(children))
	for _, c := range children {
		fmt.Println(c.Name())
	}

	t.Stop()
}

func setCfgDefaults(cfg *engine.Config) {
	switch cfg.Name {
	case "excavation.cfg":
		cfg.SetValue("WindowWidth", 1024)
		cfg.SetValue("WindowHeight", 728)
		cfg.SetValue("WindowDepth", 24)
		cfg.SetValue("Fullscreen", false)
		cfg.SetValue("VSync", 1)
	case "controls.cfg":
		//cfg.SetValue("StrafeLeft", "Joy0_1")
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

	for _, v := range sceneNode.Children() {
		if v.Type() == engine.NodeTypeCamera {
			engine.SetMainCam(&engine.Camera{v})
		}
	}
}
