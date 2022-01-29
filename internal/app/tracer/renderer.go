package tracer

import (
	"math/rand"
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

	rays := make([]geom.Ray, 0)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// create ray
			cameraRay := geom.NewEmptyRay()
			ctx.rayForPixelPathTracer(x, y, &cameraRay)
			rays = append(rays, cameraRay)
		}
	}

	// with all rays complete, we should call the OpenCL trace with:
	// all rays, world scene to intersect with. How should we layout each ray?
	// a ray is a POINT 4xfloat64 and a DIRECTION 4xfloat64. So each ray can be treated as 8x*float64
	// where the first 4 elements is the POINT and the last 4 is the DIRECTION.
	rayData := ocl.BuildRayBufferCL(rays)
	objects, normalObjects, triangles := ocl.BuildSceneBufferCL(ctx.scene.Objects)

	// Render the scene
	result := ocl.Trace(rayData, objects, normalObjects, triangles, width, height, ctx.samples)

	// result now contains RGBA values for each pixel,
	j := 0
	for i := 0; i < len(result); i += 4 {
		x := j % width
		y := j / width
		ctx.canvas.WritePixelMutex(x, y, geom.NewColor(result[i], result[i+1], result[i+2]))
		j++
	}
}

func (ctx *Ctx) rayForPixelPathTracer(x, y int, out *geom.Ray) {
	pointInView := geom.NewPoint(0, 0, -1.0)
	subVec := geom.NewVector(0, 0, 0)
	pixel := geom.NewTuple()
	// We might move the random in-pixel offset into OpenCL
	//xOffset := ctx.camera.PixelSize * (float64(x) + ctx.rnd.Float64()) // 0.5
	//yOffset := ctx.camera.PixelSize * (float64(y) + ctx.rnd.Float64()) // 0.5
	xOffset := ctx.camera.PixelSize * (float64(x) + 0.5)
	yOffset := ctx.camera.PixelSize * (float64(y) + 0.5)

	// this feels a little hacky but actually works.
	worldX := ctx.camera.HalfWidth - xOffset
	worldY := ctx.camera.HalfHeight - yOffset

	pointInView[0] = worldX
	pointInView[1] = worldY

	geom.MultiplyByTuplePtr(&ctx.camera.Inverse, &pointInView, &pixel)
	geom.MultiplyByTuplePtr(&ctx.camera.Inverse, &originPoint, &out.Origin)

	geom.SubPtr(pixel, out.Origin, &subVec)
	geom.NormalizePtr(&subVec, &out.Direction)
}
