package obj

import (
	"fmt"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/material"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"os"
	"strconv"
	"strings"
)

func ParseObj(data string) *Obj {
	out := &Obj{
		Verticies: make([]geom.Tuple4, 0),
		Groups:    make(map[string]*shapes.Group),
	}

	var mats map[string]*material.Mtl

	// fill index 0 with placeholder
	out.Verticies = append(out.Verticies, geom.NewPoint(0, 0, 0))
	out.Normals = append(out.Normals, geom.NewVector(0, 0, 0))
	rows := strings.Split(data, "\n")
	var currentGroup = "DefaultGroup"
	var currentMaterial = material.NewDefaultMaterial()
	out.Groups[currentGroup] = shapes.NewGroup()
	out.Groups[currentGroup].Label = currentGroup

	for _, row := range rows {
		if strings.TrimSpace(row) != "" {
			parts := strings.Fields(strings.TrimSpace(row))
			switch parts[0] {
			case "mtllib":
				fileName := parts[1]
				matData, err := os.ReadFile(fileName) // "./assets/models/" +
				if err != nil {
					panic(err.Error())
				}
				mats = ParseMtl(string(matData))
			case "usemtl":
				currentMaterial = toMaterial(mats[parts[1]])
				out.Groups[currentGroup].SetMaterial(currentMaterial)
				fmt.Printf("Set material '%v' on object '%v'\n", mats[parts[1]].Name, currentGroup)
			case "v":

				x, _ := strconv.ParseFloat(parts[1], 64)
				y, _ := strconv.ParseFloat(parts[2], 64)
				z, _ := strconv.ParseFloat(parts[3], 64)
				out.Verticies = append(out.Verticies, geom.NewPoint(x, y, z))

			case "vn":
				x, _ := strconv.ParseFloat(parts[1], 64)
				y, _ := strconv.ParseFloat(parts[2], 64)
				z, _ := strconv.ParseFloat(parts[3], 64)
				out.Normals = append(out.Normals, geom.NewVector(x, y, z))

			case "f":
				// 1/1/1 == vertex/texture/normal

				if !strings.Contains(row, "/") { // ONLY verticies
					for i := 2; i < len(parts)-1; i++ {
						idx1, _ := strconv.Atoi(parts[1])
						idx2, _ := strconv.Atoi(parts[i])
						idx3, _ := strconv.Atoi(parts[i+1])
						tri := shapes.NewTriangle3P(
							out.Verticies[idx1],
							out.Verticies[idx2],
							out.Verticies[idx3])
						out.Groups[currentGroup].AddChild(tri)
					}
				} else {

					for i := 2; i < len(parts)-1; i++ {
						subparts1 := strings.Split(parts[1], "/")
						subparts2 := strings.Split(parts[i], "/")
						subparts3 := strings.Split(parts[i+1], "/")

						// Vertices
						idx1, _ := strconv.Atoi(subparts1[0])
						idx2, _ := strconv.Atoi(subparts2[0])
						idx3, _ := strconv.Atoi(subparts3[0])

						// Future texture coordinates
						_, _ = strconv.Atoi(subparts1[1])
						_, _ = strconv.Atoi(subparts2[1])
						_, _ = strconv.Atoi(subparts3[1])

						// Normal
						var normIdx1, normIdx2, normIdx3 int
						if len(subparts1) == 3 {
							normIdx1, _ = strconv.Atoi(subparts1[2])
							normIdx2, _ = strconv.Atoi(subparts2[2])
							normIdx3, _ = strconv.Atoi(subparts3[2])
						}

						tri := shapes.NewTriangle(
							out.Verticies[idx1],
							out.Verticies[idx2],
							out.Verticies[idx3],
							out.Normals[normIdx1],
							out.Normals[normIdx2],
							out.Normals[normIdx3])
						tri.Material = currentMaterial
						out.Groups[currentGroup].AddChild(tri)
					}
				}
			case "g":
				fallthrough
			case "o":
				currentGroup = strings.Fields(strings.TrimSpace(row))[1]
				if _, exists := out.Groups[currentGroup]; !exists {
					out.Groups[currentGroup] = shapes.NewGroup()
					if len(parts) > 1 {
						out.Groups[currentGroup].Label = parts[1]
					}
				}
			default:
				out.IgnoredLines++
			}
		} else {
			out.IgnoredLines++
		}
	}
	tris := 0
	for i := range out.Groups {
		tris += len(out.Groups[i].Children)
	}
	fmt.Println("Loaded object:")
	fmt.Printf("Groups:    %d\n", len(out.Groups))
	fmt.Printf("Triangles: %d\n", tris)
	fmt.Printf("Verticies: %d\n", len(out.Verticies))
	fmt.Printf("Normals:   %d\n", len(out.Normals))
	return out
}

func ComputeVertexNormals(tris []*shapes.Triangle) {
	// brute force approach.
	// for every triangle we already have the surface normal,
	//    get x,y,z coord of P1
	//        iterate over all other triangles and find all faces where x, y or z == P1.x
	//        for each found face (tri), add its surface normal to N1
	//        when done, normalize or do some magnitude shit to get a unit length vertex normal that's
	//        the "average" of adjacent faces and self.
	//    repeat for P2 and P3. Do not touch N
	numNormals := 0
	numFaces := 0
	for i := range tris {
		t := tris[i]
		n1 := t.N
		n2 := t.N
		n3 := t.N
		// check vertex P1
		for j := range tris {
			if i == j {
				continue
			}
			// check P1
			if geom.TupleEquals(t.P1, tris[j].P1) || geom.TupleEquals(t.P1, tris[j].P2) || geom.TupleEquals(t.P1, tris[j].P3) {
				n1 = geom.Add(n1, tris[j].N)
			}
			// check P2
			if geom.TupleEquals(t.P2, tris[j].P1) || geom.TupleEquals(t.P2, tris[j].P2) || geom.TupleEquals(t.P2, tris[j].P3) {
				n2 = geom.Add(n2, tris[j].N)
			}
			// check P1
			if geom.TupleEquals(t.P3, tris[j].P1) || geom.TupleEquals(t.P3, tris[j].P2) || geom.TupleEquals(t.P3, tris[j].P3) {
				n3 = geom.Add(n3, tris[j].N)
			}

		}
		tris[i].N1 = geom.Normalize(n1)
		tris[i].N2 = geom.Normalize(n2)
		tris[i].N3 = geom.Normalize(n3)
		numNormals++
	}
	fmt.Printf("computed %d vertex normals from %d faces\n", numNormals, numFaces)
}

// toMaterial is a temp fix to convert our legacy materials as MTL materials
func toMaterial(mtl *material.Mtl) material.Material {
	m := material.Material{}
	//m.Name = mtl.Name

	r := mtl.Ambient[0] + mtl.Diffuse[0] + mtl.Specular[0]
	g := mtl.Ambient[1] + mtl.Diffuse[1] + mtl.Specular[1]
	b := mtl.Ambient[2] + mtl.Diffuse[2] + mtl.Specular[2]
	m.Color = geom.NewColor(r, g, b)
	//m.Ambient = avg(mtl.Ambient)
	//m.Diffuse = avg(mtl.Diffuse)
	//m.Specular = avg(mtl.Specular)
	//m.Transparency = mtl.Transparency
	m.RefractiveIndex = mtl.RefractiveIndex
	//m.Shininess = mtl.Shininess
	return m
}
func avg(t geom.Tuple4) float64 {
	return (t[1] + t[2]) / 2
}

type Obj struct {
	Verticies    []geom.Tuple4
	Normals      []geom.Tuple4
	Groups       map[string]*shapes.Group
	IgnoredLines int
}

func (o *Obj) ToGroup() *shapes.Group {
	g := shapes.NewGroup()
	g.Label = "ROOT"
	for _, v := range o.Groups {
		g.AddChild(v)
	}
	return g
}

func (o *Obj) DefaultGroup() *shapes.Group {
	return o.Groups["DefaultGroup"]
}

/*
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.000000 0.429367 0.640000
Ks 0.500000 0.500000 0.500000
Ni 1.000000
d 1.000000
illum 2
*/
func ParseMtl(data string) map[string]*material.Mtl {
	rows := strings.Split(data, "\n")
	out := make(map[string]*material.Mtl)

	var current string
	for _, row := range rows {
		if strings.TrimSpace(row) != "" {
			parts := strings.Fields(strings.TrimSpace(row))
			switch parts[0] {
			case "newmtl":
				name := parts[1]
				m := &material.Mtl{}
				m.Name = name
				out[name] = m
				current = name
			case "Ns":
				out[current].Shininess, _ = strconv.ParseFloat(parts[1], 64)
			case "Ka":
				r, _ := strconv.ParseFloat(parts[1], 64)
				g, _ := strconv.ParseFloat(parts[2], 64)
				b, _ := strconv.ParseFloat(parts[3], 64)
				out[current].Ambient = geom.NewColor(r, g, b)
			case "Kd":
				r, _ := strconv.ParseFloat(parts[1], 64)
				g, _ := strconv.ParseFloat(parts[2], 64)
				b, _ := strconv.ParseFloat(parts[3], 64)
				out[current].Diffuse = geom.NewColor(r, g, b)
			case "Ks":
				r, _ := strconv.ParseFloat(parts[1], 64)
				g, _ := strconv.ParseFloat(parts[2], 64)
				b, _ := strconv.ParseFloat(parts[3], 64)
				out[current].Specular = geom.NewColor(r, g, b)
			case "Ni":
				out[current].RefractiveIndex, _ = strconv.ParseFloat(parts[1], 64)
			case "d":
				n, _ := strconv.ParseFloat(parts[1], 64)
				out[current].Transparency = 1 - n
			default:
				// ignore..
			}
		}
	}
	return out
}
