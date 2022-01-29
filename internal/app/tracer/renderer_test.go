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
	Inside        bool // replace Max, Min
	Children      [2]int
	TriangleCount int
	Offset        int
}

func Test_FlatRecursion(t *testing.T) {
	// build "tree"
	tree := make([]BB, 0)
	tree = append(tree, BB{Inside: true, Children: [2]int{1, 2}, TriangleCount: 0, Offset: 0})     // 0 == root
	tree = append(tree, BB{Inside: true, Children: [2]int{5, 6}, TriangleCount: 5, Offset: 0})     // 1 == left
	tree = append(tree, BB{Inside: true, Children: [2]int{3, 4}, TriangleCount: 0, Offset: 0})     // 2 == right
	tree = append(tree, BB{Inside: false, Children: [2]int{-1, -1}, TriangleCount: 3, Offset: 5})  // 3 == right-left
	tree = append(tree, BB{Inside: true, Children: [2]int{-1, -1}, TriangleCount: 7, Offset: 8})   // 4 == right-right
	tree = append(tree, BB{Inside: true, Children: [2]int{-1, -1}, TriangleCount: 10, Offset: 15}) // 5 == right-right
	tree = append(tree, BB{Inside: true, Children: [2]int{-1, 7}, TriangleCount: 11, Offset: 26})  // 6 == left-right
	tree = append(tree, BB{Inside: true, Children: [2]int{-1, -1}, TriangleCount: 13, Offset: 37}) // 7 == left-right-right

	visited := make([]int, len(tree))
	vIndex := 0
	//previousIndex := -1
	index := 0
	depth := 0

	for index > -1 {
		// check if intersected. Only if we intersected we are interested in the contents
		if tree[index].Inside {
			// if there are triangles here, render them!

			for j := range tree[index].Children {
				if tree[index].Children[j] > -1 && !contains(visited, tree[index].Children[j]) {
					depth++
					index = tree[index].Children[j]
				}
			}

			// once children have been checked, see if there is something to render here.
			if tree[index].TriangleCount > 0 {
				fmt.Printf("Test intersection of triangles %d to %d\n", tree[index].Offset, tree[index].Offset+tree[index].TriangleCount)
			}

			visited[vIndex] = index
			vIndex++
			depth--
			//index =

		} else {
			visited[vIndex] = index
			vIndex++
			depth--
			// index =
		}

	}
}

func Test_TraverseWithStack(t *testing.T) {
	tree := make([]*NODE, 0)
	tree = append(tree, &NODE{Inside: true, Children: [2]int{1, 2}, TriangleCount: 0, Offset: 0})     // 0 == root
	tree = append(tree, &NODE{Inside: true, Children: [2]int{5, 6}, TriangleCount: 5, Offset: 0})     // 1 == left
	tree = append(tree, &NODE{Inside: true, Children: [2]int{3, 4}, TriangleCount: 0, Offset: 0})     // 2 == right
	tree = append(tree, &NODE{Inside: false, Children: [2]int{-1, -1}, TriangleCount: 3, Offset: 5})  // 3 == right-left
	tree = append(tree, &NODE{Inside: true, Children: [2]int{-1, -1}, TriangleCount: 7, Offset: 8})   // 4 == right-right
	tree = append(tree, &NODE{Inside: true, Children: [2]int{-1, -1}, TriangleCount: 10, Offset: 15}) // 5 == left-left
	tree = append(tree, &NODE{Inside: true, Children: [2]int{7, 8}, TriangleCount: 11, Offset: 26})   // 6 == left-right
	tree = append(tree, &NODE{Inside: true, Children: [2]int{-1, -1}, TriangleCount: 13, Offset: 37}) // 7 == left-right-left
	tree = append(tree, &NODE{Inside: true, Children: [2]int{-1, -1}, TriangleCount: 15, Offset: 50}) // 8 == left-right-right
	traverseIndex(tree)
}

var MAX_RECURSION_DEPTH = 200

type NODE struct {
	Inside        bool // replace Max, Min
	Children      [2]int
	TriangleCount int
	Offset        int
}

func (n *NODE) ChildCount() int {
	if n.Children[0] > -1 && n.Children[1] > -1 {
		return 2
	}
	if n.Children[0] == -1 || n.Children[1] == -1 {
		return 1
	}
	return 0
}

func traverse(tree []*NODE) {
	// 1) Create an empty stack S.
	var S = make([]*NODE, MAX_RECURSION_DEPTH)
	currentSIndex := 0

	// 2) Initialize current node as root
	current := tree[0]

	for current != nil || currentSIndex > -1 {

		for current != nil && current.Inside {
			// 3) Push the current node to S
			S[currentSIndex] = current
			currentSIndex++

			if current.Children[0] > -1 {
				current = tree[current.Children[0]]
			} else {
				current = nil
			}
		}
		/* Current must be NULL at this povar */
		currentSIndex--
		if currentSIndex == -1 {
			return
		}
		current = S[currentSIndex]

		fmt.Printf("Check %d triangles from node\n", current.TriangleCount)
		/*
		 * we have visited the node and its left subtree. Now, it's right subtree's turn
		 */
		if current.Children[1] != -1 {
			current = tree[current.Children[1]]
		} else {
			current = nil
		}
	}
}

// WORKS! Should be portable to OpenCL given fixed size stack.
func traverseIndex(tree []*NODE) {
	// 1) Create an empty stack.
	var stack = make([]int, MAX_RECURSION_DEPTH)

	// Stack index, i.e. current "depth" of stack
	currentSIndex := 0

	// Tree index, i.e. which "node index" we're currently processing
	currentNodeIndex := 0

	// Initialize current node as root
	current := tree[currentNodeIndex]

	for current != nil || currentSIndex > -1 {

		for current != nil && current.Inside {
			// Push the current node index to the Stack, i.e. add at current index and then increment the stack depth.
			stack[currentSIndex] = currentNodeIndex
			currentSIndex++

			// if the left child is populated (i.e. > -1), update currentNodeIndex with left child and
			// update the pointer to the current node
			if current.Children[0] > -1 {
				currentNodeIndex = current.Children[0]
				current = tree[current.Children[0]]
			} else {
				// If no left child, mark current as nil, so we can exit the inner for.
				current = nil
			}
		}

		// We pop our stack by decrementing (remember, the last iteration above resulting an increment, but no push. (Fix?)
		currentSIndex--
		if currentSIndex == -1 {
			return
		}

		// get the popped item by fetching the node index from the current stack index.
		current = tree[stack[currentSIndex]]

		// print contents. In our tracer, we'll iterate over all triangles and record triangle/ray intersections...
		fmt.Printf("Check %d triangles from node\n", current.TriangleCount)

		// we're done with the left subtree, check if there's a right-hand node.
		if current.Children[1] != -1 {
			// if there's a right-hand node, update the node index and the current node.
			currentNodeIndex = current.Children[1]
			current = tree[current.Children[1]]
		} else {
			// if no right-hand side, set current to nil.
			current = nil
		}
	}
}

func contains(sl []int, num int) bool {
	for _, v := range sl {
		if v == num {
			return true
		}
	}
	return false
}
