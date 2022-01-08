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

func BuildSceneBufferCL(in []shapes.Shape) []CLObject {
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
			obj.MinY = 0.0
			obj.MaxY = 0.0
		case *shapes.Sphere:
			obj.Type = 1
			obj.MinY = 0.0
			obj.MaxY = 0.0
		case *shapes.Cylinder:
			obj.Type = 2
			obj.MinY = in[i].(*shapes.Cylinder).MinY
			obj.MaxY = in[i].(*shapes.Cylinder).MaxY
		default:
			obj.Type = 999
		}
		// finally, pad with 32 bytes
		obj.Padding1 = 0
		obj.Padding2 = 0
		obj.Padding3 = 0
		obj.Padding4 = 0

		objs = append(objs, obj)
	}
	return objs
}

// BuildSceneBuffer32 maps shapes to a CLObject32 slice:
// Transform:        4x4 float32, offset: 0
// Inverse:          4x4 float32, offset: 8
// InverseTranspose: 4x4 float32, offset: 16
// Color:            4x float32, offset: 24
// Emission:         4x float32, offset: 26
// RefractiveIndex:  1x float32, offset: 28
// Type:             1x Int64, offset: 29
//func BuildSceneBuffer32(in []shapes.Shape) []CLObject32 {
//	objs := make([]CLObject32, 0)
//	for i := range in {
//		obj := CLObject32{}
//		obj.Transform = asFloat324x4(in[i].GetTransform())
//		obj.Inverse = asFloat324x4(in[i].GetInverse())
//		obj.InverseTranspose = asFloat324x4(in[i].GetInverseTranspose())
//		obj.Color = asFloat321x4(in[i].GetMaterial().Color)
//		obj.Emission = asFloat321x4(in[i].GetMaterial().Emission)
//		obj.RefractiveIndex = float32(in[i].GetMaterial().RefractiveIndex)
//
//		switch in[i].(type) {
//		case *shapes.Plane:
//			obj.Type = 0
//		case *shapes.Sphere:
//			obj.Type = 1
//		default:
//			obj.Type = 999
//		}
//		// finally, pad with 24 bytes
//		obj.Padding = [6]int32{0, 0, 0, 0, 0, 0}
//
//		objs = append(objs, obj)
//	}
//	return objs
//}
//
//func asFloat324x4(mat geom.Mat4x4) [16]float32 {
//	out := [16]float32{}
//	for i := range mat {
//		out[i] = float32(mat[i])
//	}
//	return out
//}
//
//func asFloat321x4(vec geom.Tuple4) [4]float32 {
//	out := [4]float32{}
//	for i := range vec {
//		out[i] = float32(vec[i])
//	}
//	return out
//}
