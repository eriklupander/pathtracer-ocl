package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"math/rand"
)

var up = geom.NewVector(0, 1, 0)

func NewPlane() *Plane {
	m1 := geom.New4x4()
	inv := geom.New4x4()
	invTranspose := geom.New4x4()
	return &Plane{
		Basic: Basic{
			Id:               rand.Int63(),
			Transform:        m1,
			Inverse:          inv,
			InverseTranspose: invTranspose,
			Material: material.Material{
				Color:           geom.Tuple4{0, .5, 1},
				Emission:        geom.Tuple4{0, 0, 0},
				RefractiveIndex: 0,
			},
		},
	}
}

type Plane struct {
	Basic
	parent Shape
}

func (p *Plane) ID() int64 {
	return p.Id
}
func (p *Plane) Lbl() string {
	return p.Label
}
func (p *Plane) GetTransform() geom.Mat4x4 {
	return p.Transform
}
func (p *Plane) GetInverse() geom.Mat4x4 {
	return p.Inverse
}
func (p *Plane) GetInverseTranspose() geom.Mat4x4 {
	return p.InverseTranspose
}

func (p *Plane) GetMaterial() material.Material {
	return p.Material
}

// SetTransform passes a pointer to the Plane on which to apply the translation matrix
func (p *Plane) SetTransform(translation geom.Mat4x4) {
	p.Transform = geom.Multiply(p.Transform, translation)
	p.Inverse = geom.Inverse(p.Transform)
	p.InverseTranspose = geom.Transpose(p.Inverse)
}

// SetMaterial passes a pointer to the Plane on which to set the material
func (p *Plane) SetMaterial(m material.Material) {
	p.Material = m
}

func (p *Plane) GetParent() Shape {
	return p.parent
}
func (p *Plane) SetParent(shape Shape) {
	p.parent = shape
}
