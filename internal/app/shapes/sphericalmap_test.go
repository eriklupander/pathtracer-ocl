package shapes

import (
	"fmt"
	geom "github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_SphericalMap(t *testing.T) {
	tcases := []struct {
		point geom.Tuple4
		u, v  float64
	}{
		{point: geom.NewTupleOf(0, 0, -1, 1), u: 0.0, v: 0.5},
		{point: geom.NewTupleOf(1, 0, 0, 1), u: 0.25, v: 0.5},
		{point: geom.NewTupleOf(0, 0, 1, 1), u: 0.5, v: 0.5},
		{point: geom.NewTupleOf(-1, 0, 0, 1), u: 0.75, v: 0.5},
		{point: geom.NewTupleOf(0, 1, 0, 1), u: 0.5, v: 1.0},
		{point: geom.NewTupleOf(0, -1, 0, 1), u: 0.5, v: 0},
		{point: geom.NewTupleOf(0.957443, -0.280411, 0.068360, 1), u: 0.261344, v: 0.261344},
		{point: geom.NewTupleOf(math.Sqrt(2.0)/2.0, math.Sqrt(2.0)/2.0, 0, 1), u: 0.25, v: 0.75},
	}
	for _, tc := range tcases {
		t.Run(fmt.Sprintf("%v", tc.point), func(t *testing.T) {
			u, v := SphericalMap(tc.point)
			assert.Equal(t, tc.u, u)
			assert.Equal(t, tc.v, v)
		})
	}
}
