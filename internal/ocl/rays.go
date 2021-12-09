package ocl

import "github.com/eriklupander/pathtracer-ocl/internal/app/geom"

// BuildRayBuffer passes each ray to a float64 slice, using 8 floats per ray. origin + direction.
func BuildRayBuffer(rays []geom.Ray) []float64 {
	rayData := make([]float64, 0)
	for i := range rays {
		rayData = append(rayData, rays[i].Origin[:]...)
		rayData = append(rayData, rays[i].Direction[:]...)
	}
	return rayData
}

func BuildRayBufferCL(rays []geom.Ray) []CLRay {
	rayData := make([]CLRay, 0)
	for i := range rays {
		rayData = append(rayData, CLRay{Origin: rays[i].Origin, Direction: rays[i].Direction})
	}
	return rayData
}

// BuildRayBufferCL32 builds a slice of CLRay32 from the passed geom.Ray slice, i.e. using float32 instead of float64.
func BuildRayBufferCL32(rays []geom.Ray) []CLRay32 {
	rayData := make([]CLRay32, 0)
	for i := range rays {
		origin := rays[i].Origin
		direction := rays[i].Direction
		rayData = append(rayData, CLRay32{
			Origin:    [4]float32{float32(origin.Get(0)), float32(origin.Get(1)), float32(origin.Get(2)), float32(origin.Get(3))},
			Direction: [4]float32{float32(direction.Get(0)), float32(direction.Get(1)), float32(direction.Get(2)), float32(direction.Get(3))},
		})
	}
	return rayData
}
