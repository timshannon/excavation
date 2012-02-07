package horde3d

/*
#include "Horde3D.h"
*/
import "C"

func H3dInit() bool {
	return bool(C.h3dInit())
}
