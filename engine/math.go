package engine

import (
	"github.com/spate/vectormath"
	"math"
)

func sliceToMatrix4(m *vectormath.Matrix4, slice []float32) {
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			m.SetElem(c, r, slice[(r*4)+c])
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
			slice[(r*4)+c] = m.GetElem(c, r)
		}
	}
}

func vector3ToSlice(slice []float32, v3 *vectormath.Vector3) {
	slice[0] = v3.GetElem(0)
	slice[1] = v3.GetElem(1)
	slice[2] = v3.GetElem(2)
}

//TEST?
func M4MakeRotationOnly(matrix *vectormath.Matrix4) {
	matrix.SetElem(0, 3, 0)
	matrix.SetElem(1, 3, 0)
	matrix.SetElem(2, 3, 0)
	matrix.SetElem(3, 3, 1)
	matrix.SetElem(3, 0, 0)
	matrix.SetElem(3, 1, 0)
	matrix.SetElem(3, 2, 0)
}

//QuatFromEuler turns a euler vector into a quaternion
func QuatFromEuler(result *vectormath.Quat, euler *vectormath.Vector3) {
	// Assuming the angles are in radians.
	c1 := math.Cos(float64(euler.Y()))
	s1 := math.Sin(float64(euler.Y()))
	c2 := math.Cos(float64(euler.Z()))
	s2 := math.Sin(float64(euler.Z()))
	c3 := math.Cos(float64(euler.X()))
	s3 := math.Sin(float64(euler.X()))
	w := math.Sqrt(1.0+c1*c2+c1*c3-s1*s2*s3+c2*c3) / 2.0
	w4 := (4.0 * w)
	result.SetX(float32((c2*s3 + c1*s3 + s1*s2*c3) / w4))
	result.SetY(float32((s1*c2 + s1*c3 + c1*s2*s3) / w4))
	result.SetZ(float32((-s1*s3 + c1*s2*c3 + s2) / w4))
	result.SetW(float32(w))

}
