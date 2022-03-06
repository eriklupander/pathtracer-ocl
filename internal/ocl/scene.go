package ocl

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
)

// BuildSceneBuffer maps shapes to a float64 slice:
// Transform:        4x4 float64, offset: 0
// Inverse:          4x4 float64, offset: 16
// InverseTranspose: 4x4 float64, offset: 32
// Color:            4xfloat64, offset: 48
// Emission:         4xfloat64, offset: 52
// RefractiveIndex:  1xfloat64, offset: 56
// Type:             1xInt64, offset: 57
//func BuildSceneBuffer(in []shapes.Shape) []float64 {
//	objs := make([]float64, 0)
//	for i := range in {
//		transform := in[i].GetTransform()
//		objs = append(objs, transform[:]...)
//
//		inverse := in[i].GetInverse()
//		objs = append(objs, inverse[:]...)
//
//		inverseTranspose := in[i].GetInverseTranspose()
//		objs = append(objs, inverseTranspose[:]...)
//
//		color := in[i].GetMaterial().Color
//		objs = append(objs, color[:]...)
//
//		emission := in[i].GetMaterial().Emission
//		objs = append(objs, emission[:]...)
//
//		objs = append(objs, in[i].GetMaterial().RefractiveIndex)
//
//		switch in[i].(type) {
//		case *shapes.Plane:
//			objs = append(objs, 0.0)
//		case *shapes.Sphere:
//			objs = append(objs, 1.0)
//		case *shapes.Cylinder:
//			objs = append(objs, 2.0)
//		default:
//			objs = append(objs, 999)
//		}
//		// finally, pad
//		pad := [6]float64{}
//		objs = append(objs, pad[:]...)
//		objs = objs[:((i + 1) * 64)] // truncate to power of 64, just in case...
//	}
//	return objs
//}

func BuildSceneBufferCL(in []shapes.Shape) ([]CLObject, []CLTriangle) {
	triangles := make([]CLTriangle, 0)
	triOffset := 0
	objs := make([]CLObject, 0)
	for i := range in {
		obj := CLObject{}
		obj.Transform = in[i].GetTransform()
		obj.Inverse = in[i].GetInverse()
		obj.InverseTranspose = in[i].GetInverseTranspose()
		obj.Color = in[i].GetMaterial().Color
		obj.Emission = in[i].GetMaterial().Emission
		obj.RefractiveIndex = in[i].GetMaterial().RefractiveIndex

		switch in[i].(type) {
		case *shapes.Plane:
			obj.Type = 0
		case *shapes.Sphere:
			obj.Type = 1
		case *shapes.Cylinder:
			obj.Type = 2
			obj.MinY = in[i].(*shapes.Cylinder).MinY
			obj.MaxY = in[i].(*shapes.Cylinder).MaxY
		case *shapes.Cube:
			obj.Type = 3
		case *shapes.Group:
			obj.Type = 4
			obj.BBMin = in[i].(*shapes.Group).BoundingBox.Min
			obj.BBMax = in[i].(*shapes.Group).BoundingBox.Max
			// add any triangles to triangle list and keep track of offset and count
			for j := range in[i].(*shapes.Group).Children {
				o := in[i].(*shapes.Group).Children[j]
				if tri, ok := o.(*shapes.Triangle); ok {
					triangles = append(triangles, CLTriangle{
						P1:      tri.P1,
						P2:      tri.P2,
						P3:      tri.P3,
						E1:      tri.E1,
						E2:      tri.E2,
						N:       tri.N,
						N1:      tri.N1,
						N2:      tri.N2,
						N3:      tri.N3,
						Padding: [224]byte{},
					})
				}
			}
			obj.TriangleOffset = int32(triOffset)
			obj.TriangleCount = int32(len(in[i].(*shapes.Group).Children))
			triOffset += int(obj.TriangleCount)
		default:
			obj.Type = 999
		}

		obj.Reflectivity = in[i].GetMaterial().Reflectivity

		// finally, pad!
		obj.Padding3 = 0
		obj.Padding4 = 0
		obj.Padding5 = [448]byte{}

		objs = append(objs, obj)
	}
	return objs, triangles
}
