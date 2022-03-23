package pkg

import "image/color"

// BoundingLinesFieldLayer defines a layer which is powered
// by a field, which in turn is defined by a top and bottom
// function (of x and t). y values between the functions
// are considered included in the field and receive the
// color given.
type BoundingLinesFieldLayer struct {
	TopFunc    func(x, t float64) float64
	BottomFunc func(x, t float64) float64
	Color      color.RGBA
	Width      int
	Height     int
}

func (f BoundingLinesFieldLayer) GetFieldStrength(x, y, t float64) float64 {
	lowVal := f.BottomFunc(x, t)
	highVal := f.TopFunc(x, t)

	if y >= lowVal && y <= highVal {
		return 1
	}
	if y <= lowVal && y >= highVal {
		return 1
	}
	return 0
}

func (f BoundingLinesFieldLayer) Rasterize(t float64) ColorLayer {
	result := make([][]color.RGBA, f.Width)

	for x := 0; x < f.Width; x++ {
		result[x] = make([]color.RGBA, f.Height)
		for y := 0; y < f.Height; y++ {
			field := f.GetFieldStrength(float64(x), float64(y), t)
			result[x][y] = ScaleColor(f.Color, field)
		}
	}

	return result
}

func (f BoundingLinesFieldLayer) Overlay(l RenderedLayer, t float64) RenderedLayer {
	return f.Rasterize(t).Overlay(l, t)
}
