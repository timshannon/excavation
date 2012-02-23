package engine

import (
	"excavation/engine/horde3d"
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
)

func Init() bool {
	test := sdl.BUTTON_LEFT
	test2 := horde3d.H3DEmitter_DelayF
	fmt.Println("tests: ", test, test2)
	return true
}
