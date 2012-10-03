package engine

import (
	"github.com/spate/vectormath"
)

func sliceToMatrix4(m *vectormath.Matrix4, slice []float32) {
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			vectormath.M4SetElem(m, c, r, slice[(r*4)+c])
		}
	}
}

func sliceToVector3(v3 *vectormath.Vector3, slice []float32) {
	vectormath.V3MakeFromElems(v3, slice[0], slice[1], slice[2])
}

func matrix4ToSlice(slice []float32, m *vectormath.Matrix4) {
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			slice[(r*4)+c] = vectormath.M4GetElem(m, c, r)
		}
	}
}

func vector3ToSlice(slice []float32, v3 *vectormath.Vector3) {
	slice[0] = vectormath.V3GetElem(v3, 0)
	slice[1] = vectormath.V3GetElem(v3, 1)
	slice[2] = vectormath.V3GetElem(v3, 2)
}
