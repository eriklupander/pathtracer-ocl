package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewSphere() *Sphere {

	return &Sphere{
		Basic: Basic{
			Id:               rand.Int63(),
			Transform:        geom.New4x4(),
			Inverse:          geom.New4x4(),
			InverseTranspose: geom.New4x4(),
			Material: material.Material{
				Color:           geom.Tuple4{1, .5, .5},
				Emission:        geom.Tuple4{0, 0, 0},
				RefractiveIndex: 1,
			},
		},
	}
}

type Sphere struct {
	Basic
	parent Shape
}

func (s *Sphere) IntersectLocal(ray geom.Ray) []Intersection {
	panic("implement me")
}

func (s *Sphere) NormalAtLocal(point geom.Tuple4, intersection *Intersection) geom.Tuple4 {
	panic("implement me")
}

func (s *Sphere) GetLocalRay() geom.Ray {
	panic("implement me")
}

func (s *Sphere) CastsShadow() bool {
	panic("implement me")
}

func (s *Sphere) Name() string {
	panic("implement me")
}

func (s *Sphere) ID() int64 {
	return s.Id
}

func (s *Sphere) GetParent() Shape {
	return s.parent
}

func (s *Sphere) GetTransform() geom.Mat4x4 {
	return s.Transform
}
func (s *Sphere) GetInverse() geom.Mat4x4 {
	return s.Inverse
}
func (s *Sphere) GetInverseTranspose() geom.Mat4x4 {
	return s.InverseTranspose
}
func (s *Sphere) GetMaterial() material.Material {
	return s.Material
}

// SetTransform passes a pointer to the Sphere on which to apply the translation matrix
func (s *Sphere) SetTransform(translation geom.Mat4x4) {
	s.Transform = geom.Multiply(s.Transform, translation)
	s.Inverse = geom.Inverse(s.Transform)
	s.InverseTranspose = geom.Transpose(s.Inverse)
}

// SetMaterial passes a pointer to the Sphere on which to set the material
func (s *Sphere) SetMaterial(m material.Material) {
	s.Material = m
}

func (s *Sphere) SetParent(shape Shape) {
	s.parent = shape
}
