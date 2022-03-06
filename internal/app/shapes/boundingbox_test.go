package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNewEmptyBoundingBox(t *testing.T) {
	box := NewEmptyBoundingBox()
	assert.Equal(t, geom.NewTupleOf(math.Inf(1), math.Inf(1), math.Inf(1), 1), box.Min)
	assert.Equal(t, geom.NewTupleOf(math.Inf(-1), math.Inf(-1), math.Inf(-1), 1), box.Max)
}

func TestNewBoundingBoxWithVolume(t *testing.T) {
	box := NewBoundingBoxF(-1, -2, -3, 3, 2, 1)
	assert.Equal(t, geom.NewPoint(-1, -2, -3), box.Min)
	assert.Equal(t, geom.NewPoint(3, 2, 1), box.Max)
}

func TestAddPointToBoundingBox(t *testing.T) {

	box := NewEmptyBoundingBox()
	p1 := geom.NewPoint(-5, 2, 0)
	p2 := geom.NewPoint(7, 0, -3)
	box.Add(p1)
	box.Add(p2)
	assert.Equal(t, geom.NewPoint(-5, 0, -3), box.Min)
	assert.Equal(t, geom.NewPoint(7, 2, 0), box.Max)
}

//
//func TestBoundsOfSphere(t *testing.T) {
//	s := NewSphere()
//	box := BoundsOf(s)
//	assert.Equal(t, geom.NewPoint(-1, -1, -1), box.Min)
//	assert.Equal(t, geom.NewPoint(1, 1, 1), box.Max)
//}
//
//func TestBoundsOfPlane(t *testing.T) {
//	p := NewPlane()
//	box := BoundsOf(p)
//	assert.Equal(t, geom.NewPoint(math.Inf(-1), 0, math.Inf(-1)), box.Min)
//	assert.Equal(t, geom.NewPoint(math.Inf(1), 0, math.Inf(1)), box.Max)
//}
//
//func TestBoundsOfCube(t *testing.T) {
//	c := NewCube()
//	box := BoundsOf(c)
//	assert.Equal(t, geom.NewPoint(-1, -1, -1), box.Min)
//	assert.Equal(t, geom.NewPoint(1, 1, 1), box.Max)
//}
//
//func TestBoundsOfInfiniteCylinder(t *testing.T) {
//	c := NewCylinder()
//	box := BoundsOf(c)
//	assert.Equal(t, geom.NewPoint(-1, math.Inf(-1), -1), box.Min)
//	assert.Equal(t, geom.NewPoint(1, math.Inf(1), 1), box.Max)
//}
//
//func TestBoundsOfFiniteCylinder(t *testing.T) {
//	c := NewCylinder()
//	c.MinY = -5
//	c.MaxY = 3
//	box := BoundsOf(c)
//	assert.Equal(t, geom.NewPoint(-1, -5, -1), box.Min)
//	assert.Equal(t, geom.NewPoint(1, 3, 1), box.Max)
//}

//
//func TestBoundsOfInfiniteCone(t *testing.T) {
//	c := NewCone()
//	box := BoundsOf(c)
//	assert.Equal(t, geom.NewPoint(math.Inf(-1), math.Inf(-1), math.Inf(-1)), box.Min)
//	assert.Equal(t, geom.NewPoint(math.Inf(1), math.Inf(1), math.Inf(1)), box.Max)
//}
//
//func TestBoundsOfFiniteCone(t *testing.T) {
//	c := NewCone()
//	c.MinY = -5
//	c.MaxY = 3
//	box := BoundsOf(c)
//	assert.Equal(t, NewPoint(-5, -5, -5), box.Min)
//	assert.Equal(t, NewPoint(5, 3, 5), box.Max)
//}

func TestBoundsOfTriangle(t *testing.T) {
	p1 := geom.NewPoint(-3, 7, 2)
	p2 := geom.NewPoint(6, 2, -4)
	p3 := geom.NewPoint(2, -1, -1)
	tri := &Triangle{P1: p1, P2: p2, P3: p3}
	box := BoundsOf(tri)
	assert.Equal(t, geom.NewPoint(-3, -1, -4), box.Min)
	assert.Equal(t, geom.NewPoint(6, 7, 2), box.Max)
}

func TestBoundingBox_MergeWith(t *testing.T) {
	b1 := NewBoundingBoxF(-5, -2, 0, 7, 4, 4)
	b2 := NewBoundingBoxF(8, -7, -2, 14, 2, 8)
	b1.MergeWith(b2)
	assert.Equal(t, geom.NewPoint(-5, -7, -2), b1.Min)
	assert.Equal(t, geom.NewPoint(14, 4, 8), b1.Max)
}

func TestBoundingBoxContainsPoint(t *testing.T) {

	BoundingBox := NewBoundingBoxF(5, -2, 0, 11, 4, 7)
	tests := []struct {
		point  geom.Tuple4
		result bool
	}{
		{geom.NewPoint(5, -2, 0), true},
		{geom.NewPoint(11, 4, 7), true},
		{geom.NewPoint(8, 1, 3), true},
		{geom.NewPoint(3, 0, 3), false},
		{geom.NewPoint(8, -4, 3), false},
		{geom.NewPoint(8, 1, -1), false},
		{geom.NewPoint(13, 1, 3), false},
		{geom.NewPoint(8, 5, 3), false},
		{geom.NewPoint(8, 1, 8), false},
	}

	for _, tc := range tests {
		res := BoundingBox.ContainsPoint(tc.point)
		assert.Equal(t, tc.result, res)
	}
}

func TestBoxContainsBox(t *testing.T) {

	BoundingBox := NewBoundingBoxF(5, -2, 0, 11, 4, 7)
	tests := []struct {
		min    geom.Tuple4
		max    geom.Tuple4
		result bool
	}{
		{geom.NewPoint(5, -2, 0), geom.NewPoint(11, 4, 7), true},
		{geom.NewPoint(6, -1, 1), geom.NewPoint(10, 3, 6), true},
		{geom.NewPoint(4, -3, -1), geom.NewPoint(10, 3, 6), false},
		{geom.NewPoint(6, -1, 1), geom.NewPoint(12, 5, 8), false},
	}

	for _, tc := range tests {
		res := BoundingBox.ContainsBox(NewBoundingBox(tc.min, tc.max))
		assert.Equal(t, tc.result, res)
	}
}

func TestTransformBoundingBox(t *testing.T) {
	box := NewBoundingBoxF(-1, -1, -1, 1, 1, 1)
	m1 := geom.Multiply(geom.RotateX(math.Pi/4), geom.RotateY(math.Pi/4))
	box2 := TransformBoundingBox(box, m1)

	assert.InEpsilon(t, -1.4142, box2.Min[0], geom.Epsilon)
	assert.InEpsilon(t, -1.7071, box2.Min[1], geom.Epsilon)
	assert.InEpsilon(t, -1.7071, box2.Min[2], geom.Epsilon)
	assert.InEpsilon(t, 1.4142, box2.Max[0], geom.Epsilon)
	assert.InEpsilon(t, 1.7071, box2.Max[1], geom.Epsilon)
	assert.InEpsilon(t, 1.7071, box2.Max[2], geom.Epsilon)
}

func TestQueryBBTransformInParentSpace(t *testing.T) {
	shape := NewSphere()
	shape.SetTransform(geom.Translate(1, -3, 5))
	shape.SetTransform(geom.Scale(0.5, 2, 4))
	box := ParentSpaceBounds(shape)
	assert.Equal(t, geom.NewPoint(0.5, -5, 1), box.Min)
	assert.Equal(t, geom.NewPoint(1.5, -1, 9), box.Max)
}

//func TestGroupBoundingBoxContainsAllItsChildren(t *testing.T) {
//
//	s := NewSphere()
//	s.SetTransform(geom.Translate(2, 5, -3))
//	s.SetTransform(geom.Scale(2, 2, 2))
//
//	c := NewCylinder()
//	c.MinY = -2
//	c.MaxY = 2
//	c.SetTransform(geom.Translate(-4, -1, 4))
//	c.SetTransform(geom.Scale(0.5, 1, 0.5))
//	g := NewGroup()
//	g.AddChild(s)
//	g.AddChild(c)
//	box := BoundsOf(g)
//	assert.Equal(t, geom.NewPoint(-4.5, -3, -5), box.Min)
//	assert.Equal(t, geom.NewPoint(4, 7, 4.5), box.Max)
//}

//
//func TestCSGBoundingBoxContainsAllItsChildren(t *testing.T) {
//
//	left := NewSphere()
//	right := NewSphere()
//	right.SetTransform(geom.Translate(2, 3, 4))
//	csg := NewCSG("difference", left, right)
//	box := BoundsOf(csg)
//	assert.Equal(t, geom.NewPoint(-1, -1, -1), box.Min)
//	assert.Equal(t, geom.NewPoint(3, 4, 5), box.Max)
//}

func TestIntersectBoundingBoxWithRayAtOrigin(t *testing.T) {

	box := NewBoundingBoxF(-1, -1, -1, 1, 1, 1)

	testcases := []struct {
		origin    geom.Tuple4
		direction geom.Tuple4
		result    bool
	}{
		{geom.NewPoint(5, 0.5, 0), geom.NewVector(-1, 0, 0), true},
		{geom.NewPoint(-5, 0.5, 0), geom.NewVector(1, 0, 0), true},
		{geom.NewPoint(0.5, 5, 0), geom.NewVector(0, -1, 0), true},
		{geom.NewPoint(0.5, -5, 0), geom.NewVector(0, 1, 0), true},
		{geom.NewPoint(0.5, 0, 5), geom.NewVector(0, 0, -1), true},
		{geom.NewPoint(0.5, 0, -5), geom.NewVector(0, 0, 1), true},
		{geom.NewPoint(0, 0.5, 0), geom.NewVector(0, 0, 1), true},
		{geom.NewPoint(-2, 0, 0), geom.NewVector(2, 4, 6), false},
		{geom.NewPoint(0, -2, 0), geom.NewVector(6, 2, 4), false},
		{geom.NewPoint(0, 0, -2), geom.NewVector(4, 6, 2), false},
		{geom.NewPoint(2, 0, 2), geom.NewVector(0, 0, -1), false},
		{geom.NewPoint(0, 2, 2), geom.NewVector(0, -1, 0), false},
		{geom.NewPoint(2, 2, 0), geom.NewVector(-1, 0, 0), false},
	}

	for _, tc := range testcases {
		direction := geom.Normalize(tc.direction)
		r := geom.NewRay(tc.origin, direction)
		assert.Equal(t, tc.result, IntersectRayWithBox(r, box))
	}
}

func TestIntersectNonCubicBoundingBoxWithRay(t *testing.T) {

	box := NewBoundingBoxF(5, -2, 0, 11, 4, 7)

	testcases := []struct {
		origin    geom.Tuple4
		direction geom.Tuple4
		result    bool
	}{
		{geom.NewPoint(15, 1, 2), geom.NewVector(-1, 0, 0), true},
		{geom.NewPoint(-5, -1, 4), geom.NewVector(1, 0, 0), true},
		{geom.NewPoint(7, 6, 5), geom.NewVector(0, -1, 0), true},
		{geom.NewPoint(9, -5, 6), geom.NewVector(0, 1, 0), true},
		{geom.NewPoint(8, 2, 12), geom.NewVector(0, 0, -1), true},
		{geom.NewPoint(6, 0, -5), geom.NewVector(0, 0, 1), true},
		{geom.NewPoint(8, 1, 3.5), geom.NewVector(0, 0, 1), true},
		{geom.NewPoint(9, -1, -8), geom.NewVector(2, 4, 6), false},
		{geom.NewPoint(8, 3, -4), geom.NewVector(6, 2, 4), false},
		{geom.NewPoint(9, -1, -2), geom.NewVector(4, 6, 2), false},
		{geom.NewPoint(4, 0, 9), geom.NewVector(0, 0, -1), false},
		{geom.NewPoint(8, 6, -1), geom.NewVector(0, -1, 0), false},
		{geom.NewPoint(12, 5, 4), geom.NewVector(-1, 0, 0), false},
	}

	for _, tc := range testcases {
		direction := geom.Normalize(tc.direction)
		r := geom.NewRay(tc.origin, direction)
		assert.Equal(t, tc.result, IntersectRayWithBox(r, box))
	}
}

//
//func TestIntersectRayGroupWithMiss(t *testing.T) {
//	s := NewSphere()
//	g := NewGroup()
//	g.AddChild(s)
//	g.Bounds()
//	r := geom.NewRay(geom.NewPoint(0, 0, -5), geom.NewVector(0, 1, 0))
//	in := geom.NewRay(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 0)) // Pass this as pointer for intermediate calc
//	IntersectRayWithShapePtr(g, r, &in)
//
//	// savedRay should have default values if the sphere's intersect was not called
//	assert.Equal(t, 0.0, s.savedRay.Origin[0])
//	assert.Equal(t, 0.0, s.savedRay.Origin[1])
//	assert.Equal(t, 0.0, s.savedRay.Origin[2])
//	assert.Equal(t, 0.0, s.savedRay.Direction[0])
//	assert.Equal(t, 0.0, s.savedRay.Direction[1])
//	assert.Equal(t, 0.0, s.savedRay.Direction[2])
//
//}
//
//
//func IntersectRayWithShapePtr(s Shape, r2 geom.Ray, in *geom.Ray) []Intersection {
//	//calcstats.Incr()
//	// transform ray with inverse of shape transformation matrix to be able to intersect a translated/rotated/skewed shape
//	geom.TransformRayPtr(r2, s.GetInverse(), in)
//
//	// Call the intersect function provided by the shape implementation (e.g. Sphere, Plane osv)
//	return s.IntersectLocal(*in)
//}
//
//func TestIntersectRayGroupWithHit(t *testing.T) {
//	s := NewSphere()
//	g := NewGroup()
//	g.AddChild(s)
//	g.Bounds()
//	r := geom.NewRay(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
//	in := geom.NewRay(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 0)) // Pass this as pointer for intermediate calc
//	IntersectRayWithShapePtr(g, r, &in)
//
//	// savedRay should have val form ray if the sphere's intersect was called
//	assert.Equal(t, 0.0, s.savedRay.Direction[0])
//	assert.Equal(t, 0.0, s.savedRay.Direction[1])
//	assert.Equal(t, 1.0, s.savedRay.Direction[2])
//
//}
//
//func TestIntersectRayWithCSGMissesBox(t *testing.T) {
//	left := NewSphere()
//	right := NewSphere()
//	csg := NewCSG("difference", left, right)
//	csg.Bounds()
//	r := NewRay(geom.NewPoint(0, 0, -5), geom.NewVector(0, 1, 0))
//	in := NewRay(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 0)) // Pass this as pointer for intermediate calc
//	IntersectRayWithShapePtr(csg, r, &in)
//	assert.Equal(t, 0.0, left.savedRay.Direction[0])
//	assert.Equal(t, 0.0, right.savedRay.Direction[0])
//	assert.Equal(t, 0.0, left.savedRay.Direction[1])
//	assert.Equal(t, 0.0, right.savedRay.Direction[1])
//	assert.Equal(t, 0.0, left.savedRay.Direction[2])
//	assert.Equal(t, 0.0, right.savedRay.Direction[2])
//}
//
//func TestIntersectRayWithCSGHitsBox(t *testing.T) {
//	left := NewSphere()
//	right := NewSphere()
//	csg := NewCSG("difference", left, right)
//	csg.Bounds()
//	r := NewRay(geom.NewPoint(0, 0, -5), geom.NewVector(0, 0, 1))
//	in := NewRay(geom.NewPoint(0, 0, 0), geom.NewVector(0, 0, 0)) // Pass this as pointer for intermediate calc
//	IntersectRayWithShapePtr(csg, r, &in)
//	assert.Equal(t, 1.0, left.savedRay.Direction[2])
//	assert.Equal(t, 1.0, right.savedRay.Direction[2])
//}
