package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"math"
)

func SphericalMap(p geom.Tuple4) (float64, float64) {

	// compute the azimuthal angle
	// -π < theta <= π
	// angle increases clockwise as viewed from above,
	// which is opposite of what we want, but we'll fix it later.
	theta := math.Atan2(p[0], p[2])

	// vec is the vector pointing from the sphere's origin (the world origin)
	// to the point, which will also happen to be exactly equal to the sphere's
	// radius.
	vec := geom.NewVector(p[0], p[1], p[2])
	radius := geom.Magnitude(vec)

	// compute the polar angle
	// 0 <= phi <= π
	phi := math.Acos(p[1] / radius)

	// -0.5 < raw_u <= 0.5
	rawU := theta / (2 * math.Pi)

	// 0 <= u < 1
	// here's also where we fix the direction of u. Subtract it from 1,
	// so that it increases counterclockwise as viewed from above.
	u := 1 - (rawU + 0.5)

	// we want v to be 0 at the south pole of the sphere,
	// and 1 at the north pole, so we have to "flip it over"
	// by subtracting it from 1.
	v := 1 - phi/math.Pi

	return u, v
}
