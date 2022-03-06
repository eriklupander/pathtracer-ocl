package geom

func NewRay(origin Tuple4, direction Tuple4) Ray {
	return Ray{Origin: origin, Direction: direction}
}
func NewEmptyRay() Ray {
	return Ray{Origin: NewTuple(), Direction: NewTuple()}
}

type Ray struct {
	Origin    Tuple4
	Direction Tuple4
}

func TransformRayPtr(r Ray, m1 Mat4x4, out *Ray) {
	MultiplyByTuplePtr(&m1, &r.Origin, &out.Origin)
	MultiplyByTuplePtr(&m1, &r.Direction, &out.Direction)
}
