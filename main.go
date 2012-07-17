package main

import (
	"excavation/engine"
)

func main() {
	if err := engine.Init("excavation"); err != nil {
		panic("Error starting Excavation: " + err.Error())
	}

	engine.StartMainLoop()

}
