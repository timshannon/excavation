package engine

import (
	"excavation/engine/horde3d"
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
)

func Init() bool {
	var running bool = true
	var cam horde3d.H3DNode = 0
	var testRes horde3d.H3DRes = 1234

	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		panic(sdl.GetError())
	}
	sdl.WM_SetCaption("Excavation", "test")

	//set sdl video mode
	if sdl.SetVideoMode(800, 600, 32, sdl.OPENGL) == nil {
		panic(sdl.GetError())
	}

	horde3d.H3dInit()
	fmt.Println("Version: ", horde3d.H3dGetVersionString())

	for running == true {
		switch event := sdl.PollEvent(); event.(type) {
		case *sdl.QuitEvent:
			running = false
			break
		}
		horde3d.H3dRender(cam)
		testRes = horde3d.H3dFindResource(horde3d.H3DResTypes_Material, "test")
		horde3d.H3dIsResLoaded(testRes)
		sdl.GL_SwapBuffers()
	}
	horde3d.H3dRelease()
	sdl.Quit()
	return true

}
