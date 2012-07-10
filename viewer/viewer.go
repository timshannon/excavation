package main

import (
	"excavation/engine"
	"flag"
)

func main() {
	engine.SetCfgFileName("viewer.cfg")

	if err := engine.Init(); err != nil {
		panic("Error starting Excavation engine: " + err.Error())
	}

	engine.StartMainLoop()

}
