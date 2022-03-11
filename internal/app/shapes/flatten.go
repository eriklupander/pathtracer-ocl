package shapes

import "fmt"

// Flatten is meaningless since each level of the hierarchy seems to build upon its parent's transform.
func Flatten(g *Group, groups *[]*Group) {
	fmt.Println("handle group " + g.Label)
	triangles := make([]Shape, 0)
	for i := range g.Children {
		tri, ok := g.Children[i].(*Triangle)
		if ok {
			triangles = append(triangles, tri)
		}
	}
	newGroup := &Group{
		Id:               g.Id,
		Transform:        g.Transform,
		Inverse:          g.Inverse,
		InverseTranspose: g.InverseTranspose,
		Material:         g.Material,
		Mtl:              g.Mtl,
		Label:            g.Label,
		parent:           g.parent,
		Children:         triangles,
	}
	newGroup.Bounds()
	*groups = append(*groups, newGroup)

	// now, handle subgroups
	for i := range g.Children {
		subgroup, ok := g.Children[i].(*Group)
		if ok {
			Flatten(subgroup, groups)
		}
	}
}
