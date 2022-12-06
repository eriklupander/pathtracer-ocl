package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"math"
	"math/rand"
)

func NewCylinder() *Cylinder {
	m1 := geom.New4x4()
	inv := geom.New4x4()
	invTranspose := geom.New4x4()

	return &Cylinder{
		Basic: Basic{
			Id:               rand.Int63(),
			Transform:        m1,
			Inverse:          inv,
			InverseTranspose: invTranspose,
			Material:         material.NewDefaultMaterial(),
		},
		MinY: math.Inf(-1),
		MaxY: math.Inf(1),
	}
}

func NewCylinderMM(min, max float64) *Cylinder {
	c := NewCylinder()
	c.MinY = min
	c.MaxY = max
	return c
}

func NewCylinderMMC(min, max float64, closed bool) *Cylinder {
	c := NewCylinder()
	c.MinY = min
	c.MaxY = max
	c.closed = closed
	return c
}

type Cylinder struct {
	Basic
	parent Shape
	MinY   float64
	MaxY   float64
	closed bool
}

func (c *Cylinder) ID() int64 {
	return c.Id
}
func (c *Cylinder) Lbl() string {
	return c.Label
}
func (c *Cylinder) GetTransform() geom.Mat4x4 {
	return c.Transform
}

func (c *Cylinder) GetInverse() geom.Mat4x4 {
	return c.Inverse
}
func (c *Cylinder) GetInverseTranspose() geom.Mat4x4 {
	return c.InverseTranspose
}

func (c *Cylinder) SetTransform(transform geom.Mat4x4) {
	c.Transform = geom.Multiply(c.Transform, transform)
	c.Inverse = geom.Inverse(c.Transform)
	c.InverseTranspose = geom.Transpose(c.Inverse)
}

func (c *Cylinder) GetMaterial() material.Material {
	return c.Material
}

func (c *Cylinder) SetMaterial(material material.Material) {
	c.Material = material
}

func (c *Cylinder) GetParent() Shape {
	return c.parent
}
func (c *Cylinder) SetParent(shape Shape) {
	c.parent = shape
}

func (c *Cylinder) Init() {}
