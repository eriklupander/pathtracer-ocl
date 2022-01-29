package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
)

type Shape interface {
	ID() int64
	GetTransform() geom.Mat4x4
	GetInverse() geom.Mat4x4
	GetInverseTranspose() geom.Mat4x4
	SetTransform(transform geom.Mat4x4)
	GetMaterial() material.Material
	SetMaterial(material material.Material)
	IntersectLocal(ray geom.Ray) []Intersection
	NormalAtLocal(point geom.Tuple4, intersection *Intersection) geom.Tuple4
	GetLocalRay() geom.Ray
	GetParent() Shape
	SetParent(shape Shape)
	CastsShadow() bool
	Name() string
}

func WorldToObject(shape Shape, point geom.Tuple4) geom.Tuple4 {
	if shape.GetParent() != nil {
		point = WorldToObject(shape.GetParent(), point)
	}
	return geom.MultiplyByTuple(shape.GetInverse(), point)
}

func WorldToObjectPtr(shape Shape, point geom.Tuple4, out *geom.Tuple4) {
	if shape.GetParent() != nil {
		WorldToObjectPtr(shape.GetParent(), point, &point)
	}
	i := shape.GetInverse()
	geom.MultiplyByTuplePtr(&i, &point, out)
}

func NormalToWorld(shape Shape, normal geom.Tuple4) geom.Tuple4 {
	normal = geom.MultiplyByTuple(shape.GetInverseTranspose(), normal)
	normal[3] = 0.0 // set w to 0
	normal = geom.Normalize(normal)

	if shape.GetParent() != nil {
		normal = NormalToWorld(shape.GetParent(), normal)
	}
	return normal
}

func NormalToWorldPtr(shape Shape, normal *geom.Tuple4) {
	it := shape.GetInverseTranspose()
	geom.MultiplyByTuplePtr(&it, normal, normal)
	normal[3] = 0.0 // set w to 0
	geom.NormalizePtr(normal, normal)

	if shape.GetParent() != nil {
		NormalToWorldPtr(shape.GetParent(), normal)
	}
}
