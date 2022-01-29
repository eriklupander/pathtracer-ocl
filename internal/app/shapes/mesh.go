package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"math/rand"
)

type Mesh struct {
	Basic
	Triangles []*Triangle
}

func (m *Mesh) IntersectLocal(ray geom.Ray) []Intersection {
	panic("implement me")
}

func (m *Mesh) NormalAtLocal(point geom.Tuple4, intersection *Intersection) geom.Tuple4 {
	panic("implement me")
}

func (m *Mesh) GetLocalRay() geom.Ray {
	panic("implement me")
}

func (m *Mesh) CastsShadow() bool {
	panic("implement me")
}

func (m *Mesh) Name() string {
	panic("implement me")
}

func NewMeshTri(triangles []*Triangle) *Mesh {
	mesh := NewMesh()
	mesh.Triangles = triangles
	return mesh
}
func NewMesh() *Mesh {
	m1 := geom.New4x4()
	inv := geom.New4x4()
	invTranspose := geom.New4x4()

	return &Mesh{
		Basic: Basic{
			Id:               rand.Int63(),
			Transform:        m1,
			Inverse:          inv,
			InverseTranspose: invTranspose,
			Material:         material.NewDefaultMaterial(),
		},
	}
}

func (m *Mesh) ID() int64 {
	return m.Id
}

func (m *Mesh) GetTransform() geom.Mat4x4 {
	return m.Transform
}

func (m *Mesh) GetInverse() geom.Mat4x4 {
	return m.Inverse
}

func (m *Mesh) GetInverseTranspose() geom.Mat4x4 {
	return m.InverseTranspose
}

func (m *Mesh) SetTransform(translation geom.Mat4x4) {
	m.Transform = geom.Multiply(m.Transform, translation)
	m.Inverse = geom.Inverse(m.Transform)
	m.InverseTranspose = geom.Transpose(m.Inverse)
}

func (m *Mesh) GetMaterial() material.Material {
	return m.Material
}

func (m *Mesh) SetMaterial(material material.Material) {
	m.Material = material
}

func (m *Mesh) GetParent() Shape {
	panic("impl me!")
}

func (m *Mesh) SetParent(shape Shape) {

}
