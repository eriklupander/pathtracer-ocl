package material

import "github.com/eriklupander/pathtracer-ocl/internal/app/geom"

type Mtl struct {
	Ambient         geom.Tuple4
	Diffuse         geom.Tuple4
	Specular        geom.Tuple4
	Shininess       float64
	Reflectivity    float64
	Transparency    float64
	RefractiveIndex float64
	Name            string
}
