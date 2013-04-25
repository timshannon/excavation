package main

import (
	"excavation/engine"
	"excavation/engine/gui"
)

var mainMenu *engine.Gui

func loadMainMenu() {

	mainMenu = engine.NewGui()
	mainMenu.UseMouse = true
	mainMenu.HaltInput = true

	mainMenu.Bind(closeMenu, "Key_Esc")

	//background Image
	img := gui.MakeImage("jupiterBackground", "overlays/gui/mainMenu/jupiter.material.xml",
		engine.NewScreenArea(0, 0, 1, 1.8, engine.ScreenRelativeLeft))
	mainMenu.AddWidget(img)

	//New
	btnNew := gui.MakeButton("new", "New Game", 0.05,
		engine.NewScreenArea(0.1, .7, .03, .5, engine.ScreenRelativeLeft))
	btnNew.ShowBackground(false)
	btnNew.Text.SetColor(engine.NewColor(75, 75, 75, 255))
	btnNew.TextHover.SetColor(engine.NewColor(100, 100, 100, 255))
	btnNew.TextClick.SetColor(engine.NewColor(255, 255, 255, 255))

	btnNew.ClickEvent = mainMenuButtons

	//Quit
	btnQuit := gui.MakeButton("quit", "Quit", 0.05,
		engine.NewScreenArea(0.1, .75, .03, .5, engine.ScreenRelativeLeft))
	//btnQuit.ShowBackground(false)
	btnQuit.Text.SetColor(engine.NewColor(75, 75, 75, 255))
	btnQuit.TextHover.SetColor(engine.NewColor(100, 100, 100, 255))
	btnQuit.TextClick.SetColor(engine.NewColor(100, 100, 100, 255))

	btnQuit.ClickEvent = mainMenuButtons

	mainMenu.AddWidget(btnNew)
	mainMenu.AddWidget(btnQuit)

	btnTest := gui.MakeButton("test", "test a lot of text ", 72,
		engine.NewScreenArea(0.3, .1, .25, .25, engine.ScreenRelativeLeft))
	btnTest.Text.SetColor(engine.NewColor(100, 100, 100, 255))
	//btnTest.ShowBackground(false)
	mainMenu.AddWidget(btnTest)
	engine.Println("Width: ", btnTest.MouseArea().PixelWidth())
	engine.Println("Height: ", btnTest.MouseArea().PixelHeight())

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
