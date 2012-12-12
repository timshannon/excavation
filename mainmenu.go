package main

import (
	"excavation/engine"
	"excavation/engine/gui"
)

var mainMenu *engine.Gui

func loadMainMenu() {

	mainMenu = new(engine.Gui)
	mainMenu.UseMouse = true
	button := gui.MakeButton("newButton", "New Game", 5,
		engine.NewScreenArea(.5, .5, .2, .5, engine.ScreenRelativeAspect))
	mainMenu.AddWidget(button)

	engine.LoadGui(mainMenu)
}
