package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"math/rand"
)

func NewCube() *Cube {
	m1 := geom.New4x4()  //NewMat4x4(make([]float64, 16))
	inv := geom.New4x4() //NewMat4x4(make([]float64, 16))
	invTranspose := geom.New4x4()

	return &Cube{
		Basic: Basic{
			Id:               rand.Int63(),
			Transform:        m1,
			Inverse:          inv,
			InverseTranspose: invTranspose,
			Material:         material.NewDefaultMaterial(),
		},
	}
}

type Cube struct {
	Basic
	parent Shape
}

func (c *Cube) ID() int64 {
	return c.Id
}
func (c *Cube) Lbl() string {
	return c.Label
}
func (c *Cube) GetTransform() geom.Mat4x4 {
	return c.Transform
}
func (c *Cube) GetInverse() geom.Mat4x4 {
	return c.Inverse
}
func (c *Cube) GetInverseTranspose() geom.Mat4x4 {
	return c.InverseTranspose
}

func (c *Cube) SetTransform(transform geom.Mat4x4) {
	c.Transform = geom.Multiply(c.Transform, transform)
	c.Inverse = geom.Inverse(c.Transform)
	c.InverseTranspose = geom.Transpose(c.Inverse)
}

func (c *Cube) GetMaterial() material.Material {
	return c.Material
}

func (c *Cube) SetMaterial(material material.Material) {
	c.Material = material
}

func (c *Cube) GetParent() Shape {
	return c.parent
}
func (c *Cube) SetParent(shape Shape) {
	c.parent = shape
}
