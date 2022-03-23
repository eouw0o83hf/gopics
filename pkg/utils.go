package pkg

import (
	"image/color"
	"math"
)

// ScaleColor scales a given color by a scalar.
// Alpha is passed through untouched
func ScaleColor(c color.RGBA, scale float64) color.RGBA {
	apply := func(x uint8) uint8 {
		return uint8(float64(x) * scale)
	}

	return color.RGBA{
		R: apply(c.R),
		G: apply(c.G),
		B: apply(c.B),
		A: c.A,
	}
}

// Dist returns the cartesian distance betwen two points
func Dist(a, b, x, y float64) float64 {
	return math.Sqrt(
		math.Pow(a - x, 2) +
		math.Pow(b - y, 2))
}
