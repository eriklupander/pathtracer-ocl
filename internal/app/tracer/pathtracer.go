package tracer

import (
	"github.com/eriklupander/pathtracer-ocl/cmd"
	canvas2 "github.com/eriklupander/pathtracer-ocl/internal/app/canvas"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/scenes"
	"github.com/sirupsen/logrus"
	"image"
	"image/png"
	"math"
	"os"
	"time"
)

var originPoint = geom.NewPoint(0, 0, 0)

func Render(sceneFactory func() *scenes.Scene) {

	st := time.Now()
	canvas := canvas2.NewCanvas(cmd.Cfg.Width, cmd.Cfg.Height)

	// Create the render contexts, one per worker
	renderContext := NewCtx(0, sceneFactory(), canvas, cmd.Cfg.Samples)
	renderContext.renderPixelPathTracer(cmd.Cfg.Width, cmd.Cfg.Height)

	logrus.Infof("Finished in %v\n", time.Now().Sub(st))
	writeImagePNG(canvas, "out.png")
}

func writeImagePNG(canvas *canvas2.Canvas, filename string) {
	logrus.Infof("writing output to file %v\n", filename)
	myImage := image.NewRGBA(image.Rect(0, 0, canvas.W, canvas.H))
	writeDataToPNG(canvas, myImage)
	outputFile, _ := os.Create(filename)
	defer outputFile.Close()
	_ = png.Encode(outputFile, myImage)
}

func writeDataToPNG(canvas *canvas2.Canvas, myImage *image.RGBA) {
	for i := 0; i < len(canvas.Pixels); i++ {
		myImage.Pix[i*4] = clamp(canvas.Pixels[i][0])
		myImage.Pix[i*4+1] = clamp(canvas.Pixels[i][1])
		myImage.Pix[i*4+2] = clamp(canvas.Pixels[i][2])
		myImage.Pix[i*4+3] = 255
	}
}

func clamp(clr float64) uint8 {
	c := clr * 255.0
	rounded := math.Round(c)
	if rounded > 255.0 {
		rounded = 255.0
	} else if rounded < 0.0 {
		rounded = 0.0
	}
	return uint8(rounded)
}
