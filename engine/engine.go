package engine

import (
	"excavation/engine/horde3d"
	"github.com/banthar/Go-SDL/sdl"
)

func Init() bool {
	var running bool = true

	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		panic(sdl.GetError())
	}
	sdl.WM_SetCaption("Excavation", "test")

	//set sdl video mode
	if sdl.SetVideoMode(800, 600, 32, sdl.OPENGL) == nil {
		panic(sdl.GetError())
	}

	horde3d.H3dInit()

	for running == true {
		switch event := sdl.PollEvent(); event.(type) {
		case *sdl.QuitEvent:
			running = false
			break
		}
		sdl.GL_SwapBuffers()
	}
	sdl.Quit()
	return true

}
