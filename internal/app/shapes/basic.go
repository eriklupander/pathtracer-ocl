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
	GetParent() Shape
	SetParent(shape Shape)
}

type Basic struct {
	Id               int64
	Transform        geom.Mat4x4
	Inverse          geom.Mat4x4
	InverseTranspose geom.Mat4x4
	Material         material.Material
}
