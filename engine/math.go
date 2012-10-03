package engine

import (
	"github.com/spate/vectormath"
)

func sliceToMatrix4(m4 *vectormath.Matrix4, slice []float32) {
	var col0, col1, col2, col3 *vectormath.Vector4
	vectormath.M4GetCol0(col0, m)
	vectormath.M4GetCol1(col1, m)
	vectormath.M4GetCol2(col2, m)
	vectormath.M4GetCol3(col3, m)

	sliceToVector4(col0, slice[:3])
	sliceToVector4(col1, slice[4:7])
	sliceToVector4(col2, slice[8:11])
	sliceToVector4(col3, slice[12:15])
	vectormath.M4SetCol0(m, col0)
	vectormath.M4SetCol1(m, col1)
	vectormath.M4SetCol2(m, col2)
	vectormath.M4SetCol3(m, col3)
}

func sliceToVector4(v4 *vectormath.Vector4, slice []float32) {
	vectormath.V4MakeFromElems(v4, slice[0], slice[1], slice[2], slice[3])
}

func matrix4ToSlice(slice []float32, m4 *vectormath.Matrix4) {
	var col0, col1, col2, col3 *vectormath.Vector4
	vectormath.M4GetCol0(col0, m)
	vectormath.M4GetCol1(col1, m)
	vectormath.M4GetCol2(col2, m)
	vectormath.M4GetCol3(col3, m)

	vector4ToSlice(slice[:3], col0)
	vector4ToSlice(slice[4:7], col1)
	vector4ToSlice(slice[8:11], col2)
	vector4ToSlice(slice[12:15], col3)
}

func vector4ToSlice(slice []float32, v4 *vectormath.Vector4) {
	slice[0] = vectormath.V4GetElem(v4, 0)
	slice[1] = vectormath.V4GetElem(v4, 1)
	slice[2] = vectormath.V4GetElem(v4, 2)
	slice[3] = vectormath.V4GetElem(v4, 3)
}
