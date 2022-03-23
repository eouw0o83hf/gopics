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

func overlayColorComponent(colorA, alphaA, colorB, alphaB uint8) uint8 {
	aContrib := float64(alphaA) * float64(colorA) / 0xff
	bContrib := float64(alphaB) * float64(colorB) / 0xff

	total := uint8(aContrib + bContrib)
	if total > 0xff {
		return  0xff
	}
	return total
}

func OverlayColor(under, over color.RGBA) color.RGBA {
	return color.RGBA{
		R: overlayColorComponent(under.R, under.A, over.R, over.A),
		G: overlayColorComponent(under.G, under.A, over.G, over.A),
		B: overlayColorComponent(under.B, under.A, over.B, over.A),
		A: 0xff,
	}
}