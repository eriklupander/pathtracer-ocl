package tracer

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
)

// Position multiplies direction of ray with the passed distance and adds the result onto the origin.
// Used for finding the position along a ray.
func Position(r geom.Ray, distance float64) geom.Tuple4 {
	add := geom.MultiplyByScalar(r.Direction, distance)
	pos := geom.Add(r.Origin, add)
	return pos
}

func TransformRay(r geom.Ray, m1 geom.Mat4x4) geom.Ray {
	origin := geom.MultiplyByTuple(m1, r.Origin)
	direction := geom.MultiplyByTuple(m1, r.Direction)
	return geom.NewRay(origin, direction)
}

func TransformRayPtr(r geom.Ray, m1 geom.Mat4x4, out *geom.Ray) {
	geom.MultiplyByTuplePtr(&m1, &r.Origin, &out.Origin)
	geom.MultiplyByTuplePtr(&m1, &r.Direction, &out.Direction)
}
