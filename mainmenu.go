package main

import (
	"excavation/engine"
	"excavation/engine/gui"
)

var mainMenu *engine.Gui

func loadMainMenu() {

	mainMenu = new(engine.Gui)
	mainMenu.UseMouse = true
	mainMenu.HaltInput = true

	//New
	btnNew := gui.MakeButton("new", "New Game", 0.05,
		engine.NewScreenArea(0.1, .7, .03, .5, engine.ScreenRelativeLeft))
	btnNew.ShowBackground(false)
	btnNew.Text.Color = engine.NewColor(175, 175, 175, 255)
	btnNew.TextHover.Color = engine.NewColor(255, 255, 255, 255)
	btnNew.TextClick.Color = engine.NewColor(255, 255, 255, 255)

	btnNew.ClickEvent = mainMenuButtons

	//Quit
	btnQuit := gui.MakeButton("quit", "Quit", 0.05,
		engine.NewScreenArea(0.1, .75, .03, .5, engine.ScreenRelativeLeft))
	btnQuit.ShowBackground(false)
	btnQuit.Text.Color = engine.NewColor(175, 175, 175, 255)
	btnQuit.TextHover.Color = engine.NewColor(255, 255, 255, 255)

	btnQuit.ClickEvent = mainMenuButtons

	mainMenu.AddWidget(btnNew)
	mainMenu.AddWidget(btnQuit)
	engine.LoadGui(mainMenu)
}

func mainMenuButtons(sender string) {
	switch sender {
	case "quit":
		engine.StopMainLoop()
	case "new":
		//TODO: Fix
		loadScene("test")
	}
}
