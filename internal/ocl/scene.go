package ocl

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
)

// All below should go into a struct to avoid package-scoped state.
var globalTriangleOffset = int32(0)
var globalGroupOffset = int32(-1)

var triangles = make([]CLTriangle, 0) // global list of ALL triangles
var groups = make([]CLGroup, 0)       // global list of ALL groups

func BuildSceneBufferCL(in []shapes.Shape) ([]CLObject, []CLTriangle, []CLGroup) {

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

			// when we encounter a group...
			obj.GroupOffset = globalGroupOffset + 1
			BuildCLGroup(in[i].(*shapes.Group))
		default:
			obj.Type = 999
		}

		obj.Reflectivity = in[i].GetMaterial().Reflectivity

		// finally, pad!
		obj.Padding3 = 0
		obj.Padding4 = 0
		obj.Padding5 = [452]byte{}

		objs = append(objs, obj)
	}
	return objs, triangles, groups
}

func BuildCLGroup(group *shapes.Group) int32 {
	groups = append(groups, CLGroup{Children: [16]int32{}, Padding: [52]byte{}})
	globalGroupOffset++
	localGroupID := globalGroupOffset
	groups[localGroupID].Color = group.GetMaterial().Color
	groups[localGroupID].Emission = group.GetMaterial().Emission
	groups[localGroupID].BBMin = group.BoundingBox.Min
	groups[localGroupID].BBMax = group.BoundingBox.Max

	// for troubleshooting, pass label to padding
	for i, b := range group.Label {
		groups[localGroupID].Padding[i] = byte(b)
	}

	var localTrianglesAdded = int32(0)

	// first add triangles belonging to THIS group. But first, record offset when starting.
	groups[localGroupID].TriOffset = globalTriangleOffset
	for _, child := range group.Children {
		tri, ok := child.(*shapes.Triangle)
		if ok {
			clTriangle := CLTriangle{
				P1: tri.P1,
				P2: tri.P2,
				P3: tri.P3,
				E1: tri.E1,
				E2: tri.E2,
				N:  tri.N,
				N1: tri.N1,
				N2: tri.N2,
				N3: tri.N3,
			}
			triangles = append(triangles, clTriangle)
			localTrianglesAdded++
		}
	}
	groups[localGroupID].TriCount = localTrianglesAdded
	globalTriangleOffset += localTrianglesAdded

	// Once we're done with the triangles, start iterating over any subgroups
	// Start by recording CURRENT offset
	numSumGroups := 0
	for _, child := range group.Children {
		grChild, ok := child.(*shapes.Group)
		if ok {
			// if group, recurse
			groups[localGroupID].Children[numSumGroups] = BuildCLGroup(grChild)
			numSumGroups++
		}
	}
	if numSumGroups > 0 {
		groups[localGroupID].ChildGroupCount = int32(numSumGroups)
	} else {
		// mark as having no subgroups
		groups[localGroupID].ChildGroupCount = int32(-1)
	}

	return localGroupID
}
