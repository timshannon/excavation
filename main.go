package main

import (
	"excavation/engine"
)

func main() {

	engine.SetDefaultConfigHandler(setCfgDefaults)

	if err := engine.Init("excavation"); err != nil {
		panic("Error starting Excavation: " + err.Error())
	}

	engine.StartMainLoop()

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
		cfg.SetValue("StrafeLeft", 'A')
	}

}
