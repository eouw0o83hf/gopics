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
			x := i + r.Window.Min.X
			y := j + r.Window.Min.Y

			f(i, j, x, y)
		}
	}
}

func (r Raster) GetRgba(i, j int) color.RGBA {
	return r.Bitmap[i][j]
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
