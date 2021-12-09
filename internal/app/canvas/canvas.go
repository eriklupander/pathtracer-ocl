package canvas

import (
	"fmt"
	"github.com/eriklupander/pathtracer-ocl/internal/app/geom"
	"github.com/sirupsen/logrus"
	"sync"
)

type Canvas struct {
	W        int
	H        int
	MaxIndex int
	Pixels   []geom.Tuple4
}

func NewCanvas(w int, h int) *Canvas {
	pixels := make([]geom.Tuple4, w*h)
	for i, _ := range pixels {
		pixels[i] = geom.NewColor(0, 0, 0)
	}
	return &Canvas{W: w, H: h, Pixels: pixels, MaxIndex: w * h}
}

func (c *Canvas) WritePixel(col, row int, color geom.Tuple4) {
	if row < 0 || col < 0 || row >= c.H || col > c.W {
		fmt.Println("pixel was out of bounds")
		return
	}
	if row*col > c.MaxIndex {
		fmt.Println("pixel was out of max bounds index bounds")
		return
	}
	c.Pixels[c.toIdx(col, row)] = color
}

var mutex = sync.Mutex{}

func (c *Canvas) WritePixelMutex(col, row int, color geom.Tuple4) {
	if row < 0 || col < 0 || row >= c.H || col > c.W {
		logrus.Infof("pixel was out of bounds: %d > %d || %d > %d\n", row, c.H, col, c.W)
		return
	}
	if row*col > c.MaxIndex {
		fmt.Println("pixel was out of max bounds index bounds")
		return
	}
	mutex.Lock()
	c.Pixels[c.toIdx(col, row)] = color
	mutex.Unlock()
}

func (c *Canvas) WritePixelToIndex(idx int, color geom.Tuple4) {
	c.Pixels[idx] = color
}

func (c *Canvas) ColorAt(col, row int) geom.Tuple4 {
	return c.Pixels[c.toIdx(col, row)]
}

func (c *Canvas) toIdx(x, y int) int {
	return y*c.W + x
}
