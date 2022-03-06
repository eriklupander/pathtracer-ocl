package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"math/rand"
)

type Group struct {
	Id               int64
	Transform        geom.Mat4x4
	Inverse          geom.Mat4x4
	InverseTranspose geom.Mat4x4
	Material         material.Material
	Mtl              material.Mtl
	Label            string
	parent           Shape
	Children         []Shape
	savedRay         geom.Ray

	InnerRays   []geom.Ray
	XsCache     Intersections
	BoundingBox *BoundingBox

	CastShadow bool
}

func NewGroup() *Group {
	m1 := geom.New4x4()
	inv := geom.New4x4()

	cachedXs := make([]Intersection, 16)
	innerRays := make([]geom.Ray, 0)

	return &Group{
		Id:               rand.Int63(),
		Transform:        m1,
		Inverse:          inv,
		InverseTranspose: geom.New4x4(),
		Children:         make([]Shape, 0),
		savedRay:         geom.NewRay(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 0)),

		XsCache:   cachedXs,
		InnerRays: innerRays,

		BoundingBox: NewEmptyBoundingBox(),
		CastShadow:  true,
	}
}

func (g *Group) ID() int64 {
	return g.Id
}

func (g *Group) GetTransform() geom.Mat4x4 {
	return g.Transform
}

func (g *Group) GetInverse() geom.Mat4x4 {
	return g.Inverse
}
func (g *Group) GetInverseTranspose() geom.Mat4x4 {
	return g.InverseTranspose
}

func (g *Group) SetTransform(transform geom.Mat4x4) {
	g.Transform = geom.Multiply(g.Transform, transform)
	g.Inverse = geom.Inverse(g.Transform)
	g.InverseTranspose = geom.Transpose(g.Inverse)
}

func (g *Group) GetMaterial() material.Material {
	return g.Material
}

func (g *Group) SetMaterial(material material.Material) {
	g.Material = material
	//for _, c := range g.Children {
	//	c.SetMaterial(material)
	//}
}

//func (g *Group) IntersectLocal(ray geom.Ray) []Intersection {
//
//	if g.BoundingBox != nil && !IntersectRayWithBox(ray, g.BoundingBox) {
//		//calcstats.Incr()
//		return nil
//	}
//
//	g.XsCache = g.XsCache[:0]
//	for idx := range g.Children {
//		geom.TransformRayPtr(ray, g.Children[idx].GetInverse(), &g.InnerRays[idx])
//		lxs := g.Children[idx].IntersectLocal(g.InnerRays[idx])
//		if len(lxs) > 0 {
//			g.XsCache = append(g.XsCache, lxs...)
//		}
//	}
//
//	if len(g.XsCache) > 1 {
//		sort.Sort(g.XsCache)
//	}
//	return g.XsCache
//}

func (g *Group) NormalAtLocal(point geom.Tuple4, intersection *Intersection) geom.Tuple4 {
	panic("not applicable to a group")
}

func (g *Group) GetLocalRay() geom.Ray {
	panic("not applicable to a group")
}

func (g *Group) AddChildren(shapes ...Shape) {
	for i := 0; i < len(shapes); i++ {
		g.AddChild(shapes[i])
	}
}

func (g *Group) AddChild(s Shape) {
	g.Children = append(g.Children, s)
	s.SetParent(g)

	// allocate memory for inner rays each time a child is added.
	g.InnerRays = append(g.InnerRays, geom.NewRay(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 0)))

	// recalculate bounds
	g.BoundingBox.MergeWith(BoundsOf(s))
}

func (g *Group) Bounds() {
	g.BoundingBox = BoundsOf(g)
}

func (g *Group) CastsShadow() bool {
	return g.CastShadow
}

func (g *Group) GetParent() Shape {
	return g.parent
}
func (g *Group) SetParent(shape Shape) {
	g.parent = shape
}
func (g *Group) BoundsToCube() *Cube {
	TransformBoundingBox(g.BoundingBox, g.Transform)
	xscale := (g.BoundingBox.Max[0] - g.BoundingBox.Min[0]) / 2
	yscale := (g.BoundingBox.Max[1] - g.BoundingBox.Min[1]) / 2
	zscale := (g.BoundingBox.Max[2] - g.BoundingBox.Min[2]) / 2
	x := g.BoundingBox.Min[0] + xscale
	y := g.BoundingBox.Min[1] + yscale
	z := g.BoundingBox.Min[2] + zscale

	c := NewCube()
	c.SetTransform(g.Transform)
	c.SetTransform(geom.Translate(x, y, z))
	c.SetTransform(geom.Scale(xscale, yscale, zscale))

	m := material.NewDefaultMaterial()
	//m.Transparency = 0.95
	m.Color = geom.NewColor(0.8, 0.7, 0.9)
	c.SetMaterial(m)
	return c
}
func (g *Group) Name() string {
	return g.Label
}
