package tracer

import (
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
	testee := NewCtx(1, scenes.OCLScene()(), canvas)

	testee.renderPixelPathTracer(1, 1)
}
