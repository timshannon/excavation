package horde3d

/*
#cgo LDFLAGS: -lHorde3D
#include "goHorde3D.h"
*/
import "C"

func H3dInit() int {
	return int(C.h3dInit())
}

func H3dGetOption(
