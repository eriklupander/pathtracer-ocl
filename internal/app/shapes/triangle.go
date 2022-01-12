package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
)

type Triangle struct {
	P1       geom.Tuple4 // 32 bytes
	P2       geom.Tuple4 // 32 bytes
	P3       geom.Tuple4 // 32 bytes
	E1       geom.Tuple4 // 32 bytes
	E2       geom.Tuple4 // 32 bytes
	N        geom.Tuple4 // 32 bytes
}

func (t *Triangle) ID() int64 {
	return -1
}

func (t *Triangle) GetTransform() geom.Mat4x4 {
	return geom.IdentityMatrix
}
func (t *Triangle) GetInverse() geom.Mat4x4 {
	return geom.IdentityMatrix
}
func (t *Triangle) GetInverseTranspose() geom.Mat4x4 {
	return geom.IdentityMatrix
}

func (t *Triangle) SetTransform(transform geom.Mat4x4) {
	panic("implement me")
}

func NewTriangle(p1 geom.Tuple4, p2 geom.Tuple4, p3 geom.Tuple4) *Triangle {

	e1 := geom.Sub(p2, p1)
	e2 := geom.Sub(p3, p1)
	n := geom.Normalize(geom.Cross(e2, e1))
	return &Triangle{P1: p1, P2: p2, P3: p3, E1: e1, E2: e2, N: n}
}
