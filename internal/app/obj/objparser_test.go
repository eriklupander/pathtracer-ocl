package obj

import (
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/eriklupander/pathtracer-ocl/internal/app/shapes"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"reflect"
	"testing"
)

func TestParseGibberish(t *testing.T) {
	gibberish := `There was a young lady named Bright
who traveled much faster than light.
She set out one day
in a relative way,
and came back the previous night.`
	result := ParseObj(gibberish)
	assert.Equal(t, 5, result.IgnoredLines)
}

func TestParseVerticies(t *testing.T) {

	data := `
v -1 1 0
v -1.0000 0.5000 0.0000
v 1 0 0
v 1 1 0
`
	res := ParseObj(data)
	assert.Equal(t, geom.NewPoint(-1, 1, 0), res.Verticies[1])
	assert.Equal(t, geom.NewPoint(-1, 0.5, 0), res.Verticies[2])
	assert.Equal(t, geom.NewPoint(1, 0, 0), res.Verticies[3])
	assert.Equal(t, geom.NewPoint(1, 1, 0), res.Verticies[4])
}

func TestParseTriangleFaces(t *testing.T) {
	data := `
v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
f 1 2 3
f 1 3 4
`
	parser := ParseObj(data)
	gr := parser.DefaultGroup()
	t1 := gr.Children[0].(*shapes.Triangle)
	t2 := gr.Children[1].(*shapes.Triangle)
	assert.Equal(t, t1.P1, parser.Verticies[1])
	assert.Equal(t, t1.P2, parser.Verticies[2])
	assert.Equal(t, t1.P3, parser.Verticies[3])
	assert.Equal(t, t2.P1, parser.Verticies[1])
	assert.Equal(t, t2.P2, parser.Verticies[3])
	assert.Equal(t, t2.P3, parser.Verticies[4])
}

func TestTriangulatePolygon(t *testing.T) {
	data := `
v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
v 0 2 0
f 1 2 3 4 5`
	parser := ParseObj(data)
	gr := parser.DefaultGroup()
	t1 := gr.Children[0].(*shapes.Triangle)
	t2 := gr.Children[1].(*shapes.Triangle)
	t3 := gr.Children[2].(*shapes.Triangle)

	assert.Equal(t, t1.P1, parser.Verticies[1])
	assert.Equal(t, t1.P2, parser.Verticies[2])
	assert.Equal(t, t1.P3, parser.Verticies[3])
	assert.Equal(t, t2.P1, parser.Verticies[1])
	assert.Equal(t, t2.P2, parser.Verticies[3])
	assert.Equal(t, t2.P3, parser.Verticies[4])
	assert.Equal(t, t3.P1, parser.Verticies[1])
	assert.Equal(t, t3.P2, parser.Verticies[4])
	assert.Equal(t, t3.P3, parser.Verticies[5])
}

func TestTrianglesInGroups(t *testing.T) {
	data := `
v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
g FirstGroup
f 1 2 3
g SecondGroup
f 1 3 4`

	parser := ParseObj(data)
	gr1 := parser.Groups["FirstGroup"]
	gr2 := parser.Groups["SecondGroup"]
	t1 := gr1.Children[0].(*shapes.Triangle)
	t2 := gr2.Children[0].(*shapes.Triangle)

	assert.Equal(t, t1.P1, parser.Verticies[1])
	assert.Equal(t, t1.P2, parser.Verticies[2])
	assert.Equal(t, t1.P3, parser.Verticies[3])
	assert.Equal(t, t2.P1, parser.Verticies[1])
	assert.Equal(t, t2.P2, parser.Verticies[3])
	assert.Equal(t, t2.P3, parser.Verticies[4])

}

func TestNormalData(t *testing.T) {
	data := `
vn 0 0 1
vn 0.707 0 -0.707
vn 1 2 3`

	parser := ParseObj(data)
	assert.Equal(t, parser.Normals[1], geom.NewVector(0, 0, 1))
	assert.Equal(t, parser.Normals[2], geom.NewVector(0.707, 0, -0.707))
	assert.Equal(t, parser.Normals[3], geom.NewVector(1, 2, 3))

}

func TestFacesWithNormals(t *testing.T) {
	data := `
v 0 1 0
v -1 0 0
v 1 0 0
vn -1 0 0
vn 1 0 0
vn 0 1 0
f 1//3 2//1 3//2
f 1/0/3 2/102/1 3/14/2`
	parser := ParseObj(data)

	g := parser.DefaultGroup()
	t1 := g.Children[0].(*shapes.Triangle)
	t2 := g.Children[1].(*shapes.Triangle)
	assert.Equal(t, t1.P1, parser.Verticies[1])
	assert.Equal(t, t1.P2, parser.Verticies[2])
	assert.Equal(t, t1.P3, parser.Verticies[3])
	assert.Equal(t, t1.N1, parser.Normals[3])
	assert.Equal(t, t1.N2, parser.Normals[1])
	assert.Equal(t, t1.N3, parser.Normals[2])
	dr1 := *t1
	dr2 := *t2
	assert.True(t, reflect.DeepEqual(dr1, dr2))
}

func TestParseGopherMaterials(t *testing.T) {
	data := `# Blender MTL File: 'gopher.blend'
# Material Count: 7

newmtl Body
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.000000 0.429367 0.640000
Ks 0.500000 0.500000 0.500000
Ni 1.000000
d 1.000000
illum 2

newmtl Eye-White
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.800000 0.800000 0.800000
Ks 1.000000 1.000000 1.000000
Ni 1.000000
d 1.000000
illum 2

newmtl Material
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.640000 0.640000 0.640000
Ks 0.500000 0.500000 0.500000
Ni 1.000000
d 1.000000
illum 2

newmtl Material.001
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.000000 0.000000 0.000000
Ks 0.000000 0.000000 0.000000
Ni 1.000000
d 1.000000
illum 2

newmtl NoseTop
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.000000 0.000000 0.000000
Ks 0.000000 0.000000 0.000000
Ni 1.000000
d 1.000000
illum 2

newmtl SkinColor
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.609017 0.353452 0.144174
Ks 0.500000 0.500000 0.500000
Ni 1.000000
d 1.000000
illum 2

newmtl Tooth
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.640000 0.640000 0.640000
Ks 0.500000 0.500000 0.500000
Ni 1.000000
d 1.000000
illum 2
`

	materials := ParseMtl(data)
	assert.Equal(t, 7, len(materials))
}

func TestProcessModel(t *testing.T) {
	// Model
	bytes, err := os.ReadFile("../../../assets/teapot.obj")
	assert.NoError(t, err)

	model := ParseObj(string(bytes)).ToGroup()
	model.SetTransform(geom.Translate(0, 1.2, 0))
	model.SetTransform(geom.RotateX(math.Pi / 2))
	model.SetTransform(geom.RotateY(-math.Pi / 2))
	model.SetTransform(geom.RotateX(-math.Pi / 8))
	shapes.Divide(model, 100)
	model.Bounds()
}

func TestProcessGlassModel(t *testing.T) {
	// Model
	bytes, err := os.ReadFile("../../../assets/glass.obj")
	assert.NoError(t, err)

	model := ParseObj(string(bytes)).ToGroup()
	scndGrp := model.Children[1]
	model.Children = []shapes.Shape{scndGrp}
	model.Bounds()
	model.SetTransform(geom.Translate(0, 1.2, 0))
	model.SetTransform(geom.RotateX(math.Pi / 2))
	model.SetTransform(geom.RotateY(-math.Pi / 2))
	model.SetTransform(geom.RotateX(-math.Pi / 8))
	shapes.Divide(model, 100)
	model.Bounds()
}
