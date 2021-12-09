package tracer

import (
	"fmt"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRay(t *testing.T) {
	r := geom.NewRay(geom.NewPoint(1, 2, 3), geom.NewVector(4, 5, 6))
	assert.True(t, geom.TupleEquals(r.Origin, geom.NewPoint(1, 2, 3)))
	assert.True(t, geom.TupleEquals(r.Direction, geom.NewVector(4, 5, 6)))
}

func TestDistanceFromPoint(t *testing.T) {
	r := geom.NewRay(geom.NewPoint(2, 3, 4), geom.NewVector(1, 0, 0))
	p1 := Position(r, 0)
	assert.Equal(t, geom.NewPoint(2, 3, 4), p1)
}

func TestTranslateRay(t *testing.T) {
	r := geom.NewRay(geom.NewPoint(1, 2, 3), geom.NewVector(0, 1, 0))
	m1 := geom.Translate(3, 4, 5)
	r2 := TransformRay(r, m1)
	assert.True(t, geom.TupleEquals(r2.Origin, geom.NewPoint(4, 6, 8)))
	assert.True(t, geom.TupleEquals(r2.Direction, geom.NewVector(0, 1, 0)))
}
func BenchmarkTransformRay(b *testing.B) {
	r := geom.NewRay(geom.NewPoint(1, 2, 3), geom.NewVector(0, 1, 0))
	m1 := geom.Translate(3, 4, 5)
	var r2 geom.Ray
	for i := 0; i < b.N; i++ {
		r2 = TransformRay(r, m1)
	}
	fmt.Printf("%v\n", r2)
}

func BenchmarkTransformRayPtr(b *testing.B) {
	r := geom.NewRay(geom.NewPoint(1, 2, 3), geom.NewVector(0, 1, 0))
	m1 := geom.Translate(3, 4, 5)
	var r2 geom.Ray
	for i := 0; i < b.N; i++ {
		TransformRayPtr(r, m1, &r2)
	}
	fmt.Printf("%v\n", r2)
}
