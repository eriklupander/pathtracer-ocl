package tracer

import (
	"fmt"
	"github.com/eriklupander/pathtracer-ocl/cmd"
	"github.com/eriklupander/pathtracer-ocl/internal/app/canvas"
	"github.com/eriklupander/pathtracer-ocl/internal/app/scenes"
	"testing"
)

func TestPathTracer_Render(t *testing.T) {
	cmd.FromConfig()
	cmd.Cfg.Width = 1
	cmd.Cfg.Height = 1
	canvas := canvas.NewCanvas(1, 1)
	testee := NewCtx(1, scenes.OCLScene()(), canvas, 1)

	testee.renderPixelPathTracer(1, 1)
}

func Test_ConvertToHex(t *testing.T) {
	numbers := [][]float64{
		{0.635774, 0.565133, 0.494491},
		{0.900000, 0.800000, 0.700000},
		{0.435455, 0.129024, 0.112896},
		{0.750000, 0.250000, 0.250000},
		{0.074067, 0.021946, 0.057607},
		{0.250000, 0.250000, 0.750000},
		{0.666599, 0.175565, 0.345644},
	}

	for _, row := range numbers {
		r := int(row[0] * 255)
		g := int(row[1] * 255)
		b := int(row[2] * 255)
		fmt.Printf("%X%X%X\n", r, g, b)
	}
}

type BB struct {
	Inside bool // replace Max, Min
	Left int
	Right int
	TriangleCount int
	Offset int
}

func Test_FlatRecursion(t *testing.T) {
	// build "tree"
	tree := make([]BB, 0)
	tree = append(tree, BB{Inside: true, Left: 1, Right: 2, TriangleCount: 0, Offset: 0}) // 0 == root
	tree = append(tree, BB{Inside: true, Left: 5, Right: 6, TriangleCount: 0, Offset: 0}) // 1 == left, leaf
	tree = append(tree, BB{Inside: true, Left: 3, Right: 4, TriangleCount: 0, Offset: 0}) // 2 == right
	tree = append(tree, BB{Inside: false, Left: -1, Right: -1, TriangleCount: 3, Offset: 5}) // 3 == right-left
	tree = append(tree, BB{Inside: true, Left: -1, Right: -1, TriangleCount: 7, Offset: 8}) // 4 == right-right
	tree = append(tree, BB{Inside: true, Left: -1, Right: -1, TriangleCount: 10, Offset: 15}) // 5 == right-right
	tree = append(tree, BB{Inside: true, Left: -1, Right: 7, TriangleCount: 11, Offset: 26}) // 6 == left-right
	tree = append(tree, BB{Inside: true, Left: -1, Right: -1, TriangleCount: 13, Offset: 37}) // 7 == left-right-right

	visited := make([]int, len(tree))
	vIndex := 0
	//previousIndex := -1
	index := 0
	depth := 0
	stack := make([]int, 8)
	iterations := 0
	stack[depth] = -1
	depth++
	for index > -1 {
		index = stack[depth]
		iterations++
		// check if intersected. Only if we intersected we are interested in the contents
		if tree[index].Inside {
			// if there are triangles here, render them!
			if tree[index].TriangleCount > 0 {
				fmt.Printf("Test intersection of triangles %d to %d\n", tree[index].Offset, tree[index].Offset + tree[index].TriangleCount)
			}
			// then continue traverse, but only go deeper if node hasn't been visited already.
			if tree[index].Left > -1 && !contains(visited, tree[index].Left) {
				depth++
				stack[depth] = tree[index].Left

			} else if tree[index].Right > -1 && !contains(visited, tree[index].Right) {
				depth++
				stack[depth] = tree[index].Right

			} else {
				visited[vIndex] = index
				vIndex++
				depth--
				index = stack[depth]
			}
		} else {
			visited[vIndex] = index
			vIndex++
			depth--
			index = stack[depth]
		}

	}
	fmt.Printf("finished after %d iterations\n", iterations)
}

func Test_TraverseWithStack(t *testing.T) {

}

func contains(sl []int, num int) bool {
	for _, v := range sl {
		if  v == num {
			return true
		}
	}
	return false
}