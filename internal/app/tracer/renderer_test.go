package tracer

import (
	"fmt"
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/canvas"
	"github.com/eriklupander/pathtracer-ocl/internal/app/scenes"
	"testing"
)

func TestPathTracer_Render(t *testing.T) {
	cmd.FromConfig()
	cmd.Cfg.Width = 1
	cmd.Cfg.Height = 1
	canvas := canvas.NewCanvas(1, 1)
	testee := NewCtx(1, scenes.OCLScene()(), canvas, 1)

	testee.renderPixelPathTracer(1, 1)
}

func Test_ConvertToHex(t *testing.T) {
	numbers := [][]float64{
		{0.635774, 0.565133, 0.494491},
		{0.900000, 0.800000, 0.700000},
		{0.435455, 0.129024, 0.112896},
		{0.750000, 0.250000, 0.250000},
		{0.074067, 0.021946, 0.057607},
		{0.250000, 0.250000, 0.750000},
		{0.666599, 0.175565, 0.345644},
	}

	for _, row := range numbers {
		r := int(row[0] * 255)
		g := int(row[1] * 255)
		b := int(row[2] * 255)
		fmt.Printf("%X%X%X\n", r,g,b)
	}
}