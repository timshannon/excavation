package main

import (
	"excavation/engine"
	"fmt"
)

func main() {

	engine.SetDefaultConfigHandler(setCfgDefaults)

	if err := engine.Init("excavation"); err != nil {
		panic("Error starting Excavation: " + err.Error())
	}

	engine.BindInput("StrafeLeft", inputTest)
	engine.StartMainLoop()

}

func inputTest(input *engine.Input) {
	fmt.Println("StrafeLeft")
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
		cfg.SetValue("StrafeLeft", "Joy0_1")
	}

}
