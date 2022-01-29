package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
)

//type Shape interface {
//	ID() int64
//	GetTransform() geom.Mat4x4
//	GetInverse() geom.Mat4x4
//	GetInverseTranspose() geom.Mat4x4
//	SetTransform(transform geom.Mat4x4)
//	GetMaterial() material.Material
//	SetMaterial(material material.Material)
//	GetParent() Shape
//	SetParent(shape Shape)
//}

type Basic struct {
	Id               int64
	Transform        geom.Mat4x4
	Inverse          geom.Mat4x4
	InverseTranspose geom.Mat4x4
	Material         material.Material
}

func max(values ...float64) float64 {
	c := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > c {
			c = values[i]
		}
	}
	return c
}

func min(values ...float64) float64 {
	c := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < c {
			c = values[i]
		}
	}
	return c
}
