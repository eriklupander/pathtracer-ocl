package tracer

import (
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/raw"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	camera2 "github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	canvas2 "github.com/eriklupander/pathtracer-ocl/internal/app/canvas"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/scenes"
	"github.com/eriklupander/pathtracer-ocl/internal/ocl"
)

type Ctx struct {
	Id      int
	scene   *scenes.Scene
	canvas  *canvas2.Canvas
	camera  camera2.Camera
	samples int
	rnd     *rand.Rand
}

func NewCtx(id int, scene *scenes.Scene, canvas *canvas2.Canvas, samples int) *Ctx {
	return &Ctx{
		Id:      id,
		scene:   scene,
		canvas:  canvas,
		camera:  scene.Camera,
		samples: samples,
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// for the OCL pathtracer, call this in the main thread and pre-calculate all rays.
func (ctx *Ctx) renderPixelPathTracer(width, height int) {

	sceneObjects, triangles, groups := ocl.BuildSceneBufferCL(ctx.scene.Objects)

	clCamera := ocl.CLCamera{
		Width:       int32(ctx.camera.Width),
		Height:      int32(ctx.camera.Height),
		Fov:         ctx.camera.Fov,
		PixelSize:   ctx.camera.PixelSize,
		HalfWidth:   ctx.camera.HalfWidth,
		HalfHeight:  ctx.camera.HalfHeight,
		Aperture:    ctx.camera.Aperture,
		FocalLength: ctx.camera.FocalLength,
		//Transform:   ctx.camera.Transform,
		Inverse: ctx.camera.Inverse,
		Padding: [72]byte{},
	}

	// Render the scene
	result := ocl.Trace(sceneObjects, triangles, groups, cmd.Cfg.DeviceIndex, ctx.samples, clCamera, ctx.scene.Textures, ctx.scene.SphereTextures, ctx.scene.CubeTextures)

	// result now contains RGBA values for each pixel,
	// write .raw file
	rawData := raw.WriteRawImage(result, width, height)
	if err := ioutil.WriteFile("experiment.raw", rawData, os.FileMode(0755)); err != nil {
		logrus.WithError(err).Error("error writing .raw file to disk")
	}

	j := 0
	for i := 0; i < len(result); i += 4 {
		x := j % width
		y := j / width
		ctx.canvas.WritePixelMutex(x, y, geom.NewColor(result[i], result[i+1], result[i+2]))
		j++
	}
}
