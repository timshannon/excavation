package main

import (
	"excavation/engine"
	"excavation/engine/gui"
)

var mainMenu *engine.Gui

func loadMainMenu() {

	mainMenu = new(engine.Gui)
	mainMenu.UseMouse = true

	button := gui.MakeButton("newButton", "New Game", 0.1,
		engine.NewScreenArea(0, .5, .2, .5, engine.ScreenRelativeRight))
	//button.ShowBackground = false
	mainMenu.AddWidget(button)
	engine.LoadGui(mainMenu)
}
