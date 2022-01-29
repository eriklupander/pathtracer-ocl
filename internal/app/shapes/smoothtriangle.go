package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"math"
)

const TriThreshold = 0.00000000001

func DefaultSmoothTriangle() *SmoothTriangle {
	return NewSmoothTriangle(
		geom.NewPoint(0, 1, 0),
		geom.NewPoint(-1, 0, 0),
		geom.NewPoint(1, 0, 0),
		geom.NewVector(0, 1, 0),
		geom.NewVector(-1, 0, 0),
		geom.NewVector(1, 0, 0))
}

func NewSmoothTriangle(p1 geom.Tuple4, p2 geom.Tuple4, p3 geom.Tuple4, n1 geom.Tuple4, n2 geom.Tuple4, n3 geom.Tuple4) *SmoothTriangle {

	e1 := geom.Sub(p2, p1)
	e2 := geom.Sub(p3, p1)
	n := geom.Normalize(geom.Cross(e2, e1))

	// for barycentric
	//d00 := Dot(e1, e1)
	//d01 := Dot(e1, e2)
	//d11 := Dot(e2, e2)
	//denom := d00*d11 - d01*d01

	return &SmoothTriangle{P1: p1, P2: p2, P3: p3, E1: e1, E2: e2, N: n, N1: n1, N2: n2, N3: n3,
		Material:   material.NewDefaultMaterial(),
		CastShadow: true,
		Label:      "SmoothTriangle",
		xs:         make([]Intersection, 1, 1),
		//p1ToOrigin:NewVector(0,0,0),
		//D00:d00,
		//D01:d01,
		//D11:d11,
		//Denom: denom,
	}
}

type SmoothTriangle struct {
	P1         geom.Tuple4
	P2         geom.Tuple4
	P3         geom.Tuple4
	E1         geom.Tuple4
	E2         geom.Tuple4
	N          geom.Tuple4
	N1         geom.Tuple4
	N2         geom.Tuple4
	N3         geom.Tuple4
	Material   material.Material
	CastShadow bool
	Label      string

	D00   float64
	D01   float64
	D11   float64
	Denom float64

	parent Shape

	//x Intersection
	p1ToOrigin    geom.Tuple4
	originCrossE1 geom.Tuple4
	dirCrossE2    geom.Tuple4

	xs Intersections
}

// Barycentric computes barycentric coordinates (u, v, w) for point p with respect to triangle defined by pre-computed
// vectors E1 and E2, which was derived into points d00, d01, d11 and denominator in constructor func.
func (s *SmoothTriangle) Barycentric(p geom.Tuple4, u *float64, v *float64, w *float64) {

	v2 := geom.NewTuple()
	geom.SubPtr(p, s.P1, &v2)

	d20 := geom.DotPtr(&v2, &s.E1)
	d21 := geom.DotPtr(&v2, &s.E2)

	*v = (s.D11*d20 - s.D01*d21) / s.Denom
	*w = (s.D00*d21 - s.D01*d20) / s.Denom
	*u = 1.0 - *v - *w
}

func (s *SmoothTriangle) CastsShadow() bool {
	return s.CastShadow
}

func (s *SmoothTriangle) ID() int64 {
	return -1
}

func (s *SmoothTriangle) GetTransform() geom.Mat4x4 {
	return geom.IdentityMatrix
}

func (s *SmoothTriangle) GetInverse() geom.Mat4x4 {
	return geom.IdentityMatrix
}
func (s *SmoothTriangle) GetInverseTranspose() geom.Mat4x4 {
	return geom.IdentityMatrix
}

func (s *SmoothTriangle) SetTransform(transform geom.Mat4x4) {
	panic("implement me")
}

func (s *SmoothTriangle) GetMaterial() material.Material {
	return s.Material
}

func (s *SmoothTriangle) SetMaterial(material material.Material) {
	s.Material = material
}

func (s *SmoothTriangle) IntersectLocal(ray geom.Ray) []Intersection {
	geom.CrossProduct(&ray.Direction, &s.E2, &s.dirCrossE2)
	determinant := geom.DotPtr(&s.E1, &s.dirCrossE2)
	if math.Abs(determinant) < TriThreshold {
		return nil
	}

	// Triangle misses over P1-P3 edge
	f := 1.0 / determinant
	for i := 0; i < 4; i++ {
		s.p1ToOrigin[i] = ray.Origin[i] - s.P1[i]
	}
	//p1ToOrigin := Sub(ray.Origin, s.P1)
	u := f * geom.DotPtr(&s.p1ToOrigin, &s.dirCrossE2)
	if u < 0 || u > 1 {
		return nil
	}

	geom.CrossProduct(&s.p1ToOrigin, &s.E1, &s.originCrossE1)
	v := f * geom.DotPtr(&ray.Direction, &s.originCrossE1)
	if v < 0 || (u+v) > 1 {
		return nil
	}
	tdist := f * geom.DotPtr(&s.E2, &s.originCrossE1)

	s.xs[0] = NewIntersectionUV(tdist, s, u, v)
	return s.xs
}

func (s *SmoothTriangle) NormalAtLocal(point geom.Tuple4, intersection *Intersection) geom.Tuple4 {
	return geom.Add(geom.Add(geom.MultiplyByScalar(s.N2, intersection.U),
		geom.MultiplyByScalar(s.N3, intersection.V)),
		geom.MultiplyByScalar(s.N1, 1-intersection.U-intersection.V))
}

func (s *SmoothTriangle) GetLocalRay() geom.Ray {
	panic("implement me")
}

func (s *SmoothTriangle) GetParent() Shape {
	return s.parent
}

func (s *SmoothTriangle) SetParent(shape Shape) {
	s.parent = shape
}
func (s *SmoothTriangle) Name() string {
	return s.Label
}
