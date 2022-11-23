package scenes

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/camera"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"image"
)

type Scene struct {
	Camera   camera.Camera
	Objects  []shapes.Shape
	Textures []image.Image
}
