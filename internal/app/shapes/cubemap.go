package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"math"
)

func FaceFromPoint(point geom.Tuple4) int {

	absX := math.Abs(point[0])
	absY := math.Abs(point[1])
	absZ := math.Abs(point[2])
	coord := max(absX, absY, absZ)

	if coord == point[0] {
		return 0 // right
	}
	if coord == -point[0] {
		return 1 //"left"
	}
	if coord == point[1] {
		return 2 //"up"
	}
	if coord == -point[1] {
		return 3 //"down"
	}
	if coord == point[2] {
		return 4 // "front"
	}
	return 5 //"back"
}

func cubeUVFront(point geom.Tuple4) (float64, float64) {
	u := math.Mod(point[0]+1.0, 2) / 2.0
	v := math.Mod(point[1]+1.0, 2) / 2.0
	return u, v
}
func cubeUVBack(point geom.Tuple4) (float64, float64) {
	u := math.Mod(1.0-point[0], 2) / 2.0
	v := math.Mod(point[1]+1.0, 2) / 2.0
	return u, v
}
func cubeUVLeft(point geom.Tuple4) (float64, float64) {
	u := math.Mod(point[2]+1.0, 2) / 2.0
	v := math.Mod(point[1]+1.0, 2) / 2.0
	return u, v
}
func cubeUVRight(point geom.Tuple4) (float64, float64) {
	u := math.Mod(1.0-point[2], 2) / 2.0
	v := math.Mod(point[1]+1.0, 2) / 2.0
	return u, v
}
func cubeUVTop(point geom.Tuple4) (float64, float64) {
	u := math.Mod(point[0]+1.0, 2) / 2.0
	v := math.Mod(1.0-point[2], 2) / 2.0
	return u, v
}
func cubeUVBottom(point geom.Tuple4) (float64, float64) {
	u := math.Mod(point[0]+1.0, 2) / 2.0
	v := math.Mod(point[2]+1.0, 2) / 2.0
	return u, v
}
func pattern_at(cube_map map[int]AlignCheck, point geom.Tuple4) geom.Tuple4 {
	face := FaceFromPoint(point)

	var u, v float64

	if face == 1 {
		u, v = cubeUVLeft(point)
	} else if face == 0 {
		u, v = cubeUVRight(point)
	} else if face == 4 {
		u, v = cubeUVFront(point)
	} else if face == 5 {
		u, v = cubeUVBack(point)
	} else if face == 2 {
		u, v = cubeUVTop(point)
	} else { // down
		u, v = cubeUVBottom(point)
	}

	return uv_pattern_at(cube_map[face], u, v)
}

type AlignCheck struct {
	main, ul, ur, bl, br geom.Tuple4
}

func uv_align_check(main, ul, ur, bl, br geom.Tuple4) AlignCheck {
	return AlignCheck{
		main: main,
		ul:   ul,
		ur:   ur,
		bl:   bl,
		br:   br,
	}
}

func uv_pattern_at(align_check AlignCheck, u, v float64) geom.Tuple4 {
	// remember: v=0 at the bottom, v=1 at the top
	if v > 0.8 {
		if u < 0.2 {
			return align_check.ul
		}
		if u > 0.8 {
			return align_check.ur
		}
	} else if v < 0.2 {
		if u < 0.2 {
			return align_check.bl
		}
		if u > 0.8 {
			return align_check.br
		}
	}
	return align_check.main
}
