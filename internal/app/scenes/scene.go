package scenes

import (
	"bytes"
	"github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"github.com/sirupsen/logrus"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type Scene struct {
	Camera  camera.Camera
	Objects []shapes.Shape

	// Standard textures are typically 1:1, 2048x2048
	Textures []image.Image

	// Sphere map textures are typically 2:1 format, for example 3840x1920
	SphereTextures []image.Image

	// Cube textures use a 4:3 format with 6 sides forming a cross. Example is 4096x3072
	CubeTextures []image.Image
}

func LoadImage(path string) image.Image {
	t0, err := os.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	var img0 image.Image
	if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
		img0, err = jpeg.Decode(bytes.NewBuffer(t0))
	} else if strings.HasSuffix(path, ".png") {
		img0, err = png.Decode(bytes.NewBuffer(t0))
	} else {
		logrus.Fatalf("unsupported texture image format: %s", path)
	}

	if err != nil {
		panic(err.Error())
	}

	switch img0.(type) {
	case *image.NRGBA:
		return img0
	default:
		convertedImg := image.NewNRGBA(image.Rect(0, 0, img0.Bounds().Dx(), img0.Bounds().Dy()))
		draw.Draw(convertedImg, convertedImg.Bounds(), img0, img0.Bounds().Min, draw.Src)
		return convertedImg
	}
}
