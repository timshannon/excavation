package main

import (
	"excavation/engine"
	"excavation/engine/gui"
)

var mainMenu *engine.Gui

//TODO: Main Menu vs Game Menu
func loadMainMenu() {

	engine.Pause()
	mainMenu = engine.NewGui()
	mainMenu.UseMouse = true
	mainMenu.HaltInput = true

	mainMenu.Bind(closeMenu, "Key_Esc")

	//background Image
	//img := gui.MakeImage("jupiterBackground", "overlays/gui/mainMenu/jupiter.material.xml",
	//engine.NewScreenArea(0, 0, 1.8, 1, engine.ScreenRelativeLeft))
	//mainMenu.AddWidget(img)

	//New
	btnNew := gui.MakeButton("new", "New Game", .04,
		engine.NewScreenArea(0.1, .7, .21, .05, engine.ScreenRelativeLeft))
	btnNew.ShowBackground(false)

	btnNew.Text.SetColor(engine.NewColor(75, 75, 75, 255))
	btnNew.TextHover.SetColor(engine.NewColor(100, 100, 100, 255))
	btnNew.TextClick.SetColor(engine.NewColor(255, 255, 255, 255))

	btnNew.ClickEvent = mainMenuButtons
	mainMenu.AddWidget(btnNew)

	//Quit
	btnQuit := gui.MakeButton("quit", "Quit", .04,
		engine.NewScreenArea(0.1, .75, .1, .05, engine.ScreenRelativeLeft))
	btnQuit.ShowBackground(false)

	btnQuit.Text.SetColor(engine.NewColor(75, 75, 75, 255))
	btnQuit.TextHover.SetColor(engine.NewColor(100, 100, 100, 255))
	btnQuit.TextClick.SetColor(engine.NewColor(100, 100, 100, 255))

	btnQuit.ClickEvent = mainMenuButtons

	mainMenu.AddWidget(btnQuit)

	engine.LoadGui(mainMenu)

}

func mainMenuButtons(sender string) {
	switch sender {
	case "quit":
		engine.StopMainLoop()
	case "new":
		loadScene("test")
		engine.Resume()
	}
}

func closeMenu(input *engine.Input) {
	if state, ok := input.ButtonState(); ok && state == engine.StateReleased {
		engine.UnloadGui()
		engine.Resume()
	}
}
