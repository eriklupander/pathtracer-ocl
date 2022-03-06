package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetup(t *testing.T) {
	st := DefaultTriangle()
	assert.Equal(t, geom.NewPoint(0, 1, 0), st.P1)
}

//func TestSmoothTriWithUV(t *testing.T) {
//	st := DefaultTriangle()
//	i := NewIntersectionUV(3.5, st, 0.2, 0.4)
//	assert.Equal(t, 0.2, i.U)
//	assert.Equal(t, 0.4, i.V)
//}
//func TestIntersectWithTriStoresUV(t *testing.T) {
//	tri := DefaultTriangle()
//	r := geom.NewRay(geom.NewPoint(-0.2, 0.3, -2), geom.NewVector(0, 0, 1))
//	xs := tri.IntersectLocal(r)
//	assert.InEpsilon(t, 0.45, xs[0].U, geom.Epsilon)
//	assert.InEpsilon(t, 0.25, xs[0].V, geom.Epsilon)
//}
