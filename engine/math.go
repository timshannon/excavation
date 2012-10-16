package engine

import (
	"github.com/spate/vectormath"
)

func sliceToMatrix4(m *vectormath.Matrix4, slice []float32) {
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			m.SetElem(c, r, slice[(c*4)+r])
		}
	}
}

func sliceToVector3(v3 *vectormath.Vector3, slice []float32) {
	v3.SetElem(0, slice[0])
	v3.SetElem(1, slice[1])
	v3.SetElem(2, slice[2])
}

func matrix4ToSlice(slice []float32, m *vectormath.Matrix4) {
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			slice[(c*4)+r] = m.GetElem(c, r)
		}
	}
}

func vector3ToSlice(slice []float32, v3 *vectormath.Vector3) {
	slice[0] = v3.GetElem(0)
	slice[1] = v3.GetElem(1)
	slice[2] = v3.GetElem(2)
}
