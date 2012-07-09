package main

import (
	"excavation/engine"
	"flag"
)

func main() {
	//setup separate init
	// with in memory pipeline and hardcoded window values
	//if err := engine.Init(); err != nil {
	//panic("Error starting Excavation engine: " + err.Error())
	//}

	engine.StartMainLoop()

}
