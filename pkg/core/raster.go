package core

import (
	"image"
	"image/color"
)

type Raster struct {
	Bitmap [][]color.RGBA
	Window *image.Rectangle
}

func NewRaster(r *image.Rectangle) *Raster {
	width := r.Dx()
	height := r.Dy()

	bitmap := make([][]color.RGBA, width)
	for i := 0; i < width; i++ {
		bitmap[i] = make([]color.RGBA, height)
	}

	return &Raster{
		Bitmap: bitmap,
		Window: r,
	}
}

// Traverse calls a given function for every point
// in the Raster, providing the bitmap indexes and
// true coordinate values. This may be called in
// parallel.
func (r Raster) Traverse(f func(i, j, x, y int)) {
	width := r.Window.Dx()
	height := r.Window.Dy()

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			f(i, j, r.ItoX(i), r.JToY(j))
		}
	}
}

// ItoX converts i (raster index) to x (field coord)
func (r Raster) ItoX(i int) int {
	return i + r.Window.Min.X
}

func (r Raster) XtoI(x int) int {
	return x - r.Window.Min.X
}

// JToY converts j (raster index) to y (field coord)
func (r Raster) JToY(j int) int {
	return  j + r.Window.Min.Y
}

func (r Raster) YToJ(y int) int {
	return y - r.Window.Min.Y
}

func (r Raster) GetRgba(i, j int) color.RGBA {
	return r.Bitmap[i][j]
}

func (r Raster) SetRgba(i, j int, c color.RGBA) {
	r.Bitmap[i][j] = c
}

func (r Raster) ToImage() *image.RGBA {
	width := r.Window.Dx()
	height := r.Window.Dy()

	img := image.NewRGBA(
		image.Rect(0, 0, width, height))

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Set(i, j, r.GetRgba(i, j))
		}
	}

	return img
}
