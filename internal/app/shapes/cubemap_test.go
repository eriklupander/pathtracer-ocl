package shapes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFaceFromPoint1(t *testing.T) {
	type args struct {
		point geom.Tuple4
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"left", args{point: geom.NewPoint(-1, 0.5, -0.25)}, 1},
		{"right", args{point: geom.NewPoint(1.1, -0.75, 0.8)}, 0},
		{"front", args{point: geom.NewPoint(0.1, 0.6, 0.9)}, 4},
		{"back", args{point: geom.NewPoint(-0.7, 0, -2)}, 5},
		{"up", args{point: geom.NewPoint(0.5, 1, 0.9)}, 2},
		{"down", args{point: geom.NewPoint(-0.2, -1.3, 1.1)}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FaceFromPoint(tt.args.point), "FaceFromPoint(%v)", tt.args.point)
		})
	}
}

func Test_uv_pattern_at(t *testing.T) {
	alignCheck := AlignCheck{
		main: geom.Tuple4{1, 1, 1},
		ul:   geom.Tuple4{1, 0, 0},
		ur:   geom.Tuple4{1, 1, 0},
		bl:   geom.Tuple4{0, 1, 0},
		br:   geom.Tuple4{0, 1, 1},
	}
	type args struct {
		u float64
		v float64
	}
	tests := []struct {
		args args
		want geom.Tuple4
	}{
		{args: args{u: 0.5, v: 0.5}, want: alignCheck.main},
		{args: args{u: 0.1, v: 0.9}, want: alignCheck.ul},
		{args: args{u: 0.9, v: 0.9}, want: alignCheck.ur},
		{args: args{u: 0.1, v: 0.1}, want: alignCheck.bl},
		{args: args{u: 0.9, v: 0.1}, want: alignCheck.br},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equalf(t, tt.want, uv_pattern_at(alignCheck, tt.args.u, tt.args.v), "uv_pattern_at(%v, %v, %v)", alignCheck, tt.args.u, tt.args.v)
		})
	}
}

func Test_CubeMapUp(t *testing.T) {

	p1 := geom.NewPoint(-0.5, 1, -0.5)
	p2 := geom.NewPoint(0.5, 1, 0.5)

	u, v := cubeUVTop(p1)
	assert.Equal(t, 0.25, u)
	assert.Equal(t, 0.75, v)
	u, v = cubeUVTop(p2)
	assert.Equal(t, 0.75, u)
	assert.Equal(t, 0.25, v)
}

func Test_CubeMapDown(t *testing.T) {

	p1 := geom.NewPoint(-0.5, -1, 0.5)
	p2 := geom.NewPoint(0.5, -1, -0.5)

	u, v := cubeUVBottom(p1)
	assert.Equal(t, 0.25, u)
	assert.Equal(t, 0.75, v)
	u, v = cubeUVBottom(p2)
	assert.Equal(t, 0.75, u)
	assert.Equal(t, 0.25, v)
}

func Test_UVAlignCheck(t *testing.T) {
	main := geom.NewColor(1, 1, 1)
	ul := geom.NewColor(1, 0, 0)
	ur := geom.NewColor(1, 1, 0)
	bl := geom.NewColor(0, 1, 0)
	br := geom.NewColor(0, 1, 1)
	pattern := uv_align_check(main, ul, ur, bl, br)

	assert.Equal(t, uv_pattern_at(pattern, 0.5, 0.5), main)
	assert.Equal(t, uv_pattern_at(pattern, 0.1, 0.9), ul)
	assert.Equal(t, uv_pattern_at(pattern, 0.9, 0.9), ur)
	assert.Equal(t, uv_pattern_at(pattern, 0.1, 0.1), bl)
	assert.Equal(t, uv_pattern_at(pattern, 0.9, 0.1), br)
}

func Test_CubeMapping(t *testing.T) {
	red := geom.NewColor(1, 0, 0)
	yellow := geom.NewColor(1, 1, 0)
	brown := geom.NewColor(1, 0.5, 0)
	green := geom.NewColor(0, 1, 0)
	cyan := geom.NewColor(0, 1, 1)
	blue := geom.NewColor(0, 0, 1)
	purple := geom.NewColor(1, 0, 1)
	white := geom.NewColor(1, 1, 1)
	left := uv_align_check(yellow, cyan, red, blue, brown)
	front := uv_align_check(cyan, red, yellow, brown, green)
	right := uv_align_check(red, yellow, purple, green, white)
	back := uv_align_check(green, purple, cyan, white, blue)
	up := uv_align_check(brown, cyan, purple, red, yellow)
	down := uv_align_check(purple, brown, green, blue, white)
	pattern := map[int]AlignCheck{
		1: left,
		0: right,
		4: front,
		5: back,
		2: up,
		3: down,
	}

	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-1, 0, 0)), yellow)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-1, 0.9, -0.9)), cyan)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-1, 0.9, 0.9)), red)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-1, -0.9, -0.9)), blue)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-1, -0.9, 0.9)), brown)

	// FRONT
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0, 0, 1)), cyan)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-0.9, 0.9, 1)), red)   // upper left
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0.9, 0.9, 1)), yellow) // upper right
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-0.9, -0.9, 1)), brown)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0.9, -0.9, 1)), green)

	assert.Equal(t, pattern_at(pattern, geom.NewPoint(1, 0, 0)), red)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(1, 0.9, 0.9)), yellow)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(1, 0.9, -0.9)), purple)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(1, -0.9, 0.9)), green)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(1, -0.9, -0.9)), white)

	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0, 0, -1)), green)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0.9, 0.9, -1)), purple)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-0.9, 0.9, -1)), cyan)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0.9, -0.9, -1)), white)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-0.9, -0.9, -1)), blue)

	// UP
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0, 1, 0)), brown)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-0.9, 1, -0.9)), cyan)  // upper right
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0.9, 1, -0.9)), purple) // upper left
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-0.9, 1, 0.9)), red)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0.9, 1, 0.9)), yellow)

	// DOWN
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0, -1, 0)), purple)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-0.9, -1, 0.9)), brown)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0.9, -1, 0.9)), green)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(-0.9, -1, -0.9)), blue)
	assert.Equal(t, pattern_at(pattern, geom.NewPoint(0.9, -1, -0.9)), white)

}
