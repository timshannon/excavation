// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package main

import (
	"excavation/engine"
	"flag"
	"fmt"
	"github.com/spate/vectormath"
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

	engine.AddTask("test", task, 0, 1)
	//starting the loop should be the last thing
	// after setting up the game
	engine.StartMainLoop()
}

func task(t *engine.Task) {
	fmt.Println("Test Task")
	newVec := new(vectormath.Vector3)
	newV4 := new(vectormath.Vector4)

	relMat := new(vectormath.Matrix4)
	camera.RelativeTransMat(relMat)

	fmt.Println("Before: ", relMat)
	vectormath.V3SetZ(newVec, -0.1)
	vectormath.M4MulV3(newV4, relMat, newVec)

	vectormath.V3SetX(newVec, vectormath.V4GetX(newV4))
	vectormath.V3SetY(newVec, vectormath.V4GetY(newV4))
	vectormath.V3SetZ(newVec, vectormath.V4GetZ(newV4))

	vectormath.M4SetTranslation(relMat, newVec)
	camera.SetRelativeTransMat(relMat)
	fmt.Println("After: ", relMat)

	t.Wait(.5)
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

	children := sceneNode.Children()
	for c := range children {
		if children[c].Type() == engine.NodeTypeCamera {
			engine.SetMainCam(&engine.Camera{children[c]})
			camera = children[c]
		}
	}
}
