package main

import (
	"excavation/engine"
	"excavation/engine/gui"
	"fmt"
)

var mainMenu *engine.Gui

func loadMainMenu() {

	mainMenu = new(engine.Gui)
	mainMenu.UseMouse = true

	btnNew := gui.MakeButton("new", "New Game", 0.1,
		engine.NewScreenArea(0.1, .2, .1, .5, engine.ScreenRelativeLeft))
	btnQuit := gui.MakeButton("quit", "Quit", 0.1,
		engine.NewScreenArea(0.1, .35, .1, .5, engine.ScreenRelativeAspect))

	btnQuit.ClickEvent = mainMenuButtons
	btnNew.ClickEvent = mainMenuButtons

	mainMenu.AddWidget(btnNew)
	mainMenu.AddWidget(btnQuit)
	engine.LoadGui(mainMenu)
}

func mainMenuButtons(sender string) {
	switch sender {
	case "quit":
		fmt.Println("Quit")
		engine.StopMainLoop()
	case "new":
		loadScene("test")
	}
}
