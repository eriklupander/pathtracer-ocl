package material

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
)

type Material struct {
	Color           geom.Tuple4
	Emission        geom.Tuple4
	RefractiveIndex float64
	Reflectivity    float64
	Textured        bool
	TextureID       uint8
	TextureScaleX   float64
	TextureScaleY   float64
	TexturedNM      bool
	TextureIDNM     uint8
	TextureScaleXNM float64
	TextureScaleYNM float64
	IsEnvMap        bool
}

func NewDefaultMaterial() Material {
	return Material{
		Color:           geom.Tuple4{1, 1, 1},
		Emission:        geom.Tuple4{0, 0, 0},
		RefractiveIndex: 1.0,
	}
}

func NewDiffuse(r, g, b float64) Material {
	return Material{
		Color:           geom.Tuple4{r, g, b},
		Emission:        geom.Tuple4{0, 0, 0},
		RefractiveIndex: 1.0,
	}
}
func NewGlass() Material {
	return Material{
		Color:           geom.Tuple4{1, 1, 1},
		Emission:        geom.Tuple4{0, 0, 0},
		RefractiveIndex: 1.52,
		Reflectivity:    0.05,
	}
}
func NewMirror() Material {
	return Material{
		Color:           geom.Tuple4{1, 1, 1},
		Emission:        geom.Tuple4{0, 0, 0},
		RefractiveIndex: 1.0,
		Reflectivity:    1.0,
	}
}
func NewLightBulb() Material {
	return Material{
		Color:           geom.Tuple4{1, 1, 1},
		Emission:        geom.Tuple4{8, 8, 8},
		RefractiveIndex: 1.0,
	}
}
