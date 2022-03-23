package pkg

import (
	"image/color"
	"math"
)

type GlowingCurveLayer struct {
	FieldFunc func(x, t float64) float64
	BaseColor color.RGBA
	Width     int
	Height    int
	Radius float64
	Power float64
}

func (f GlowingCurveLayer) GetFieldStrength(x, y, t float64) float64 {
	x = x - float64(f.Width)/2
	y = y - float64(f.Height)/2

	yActual := f.FieldFunc(x, t)
	dist := math.Abs(y - yActual)
	if dist <= f.Radius {
		strength := f.Power / dist
		if strength >= 1 {
			return  1
		}
		return strength
	}
	return 0
}

func (f GlowingCurveLayer) Rasterize(t float64) ColorLayer {
	result := make([][]color.RGBA, f.Width)

	for x := 0; x < f.Width; x++ {
		result[x] = make([]color.RGBA, f.Height)
		for y := 0; y < f.Height; y++ {
			field := f.GetFieldStrength(float64(x), float64(y), t)
			result[x][y] = ScaleColor(f.BaseColor, field)
		}
	}

	return result
}

func (f GlowingCurveLayer) Overlay(l RenderedLayer, t float64) RenderedLayer {
	return f.Rasterize(t).Overlay(l, t)
}

func (f GlowingCurveLayer) GetPixel(c color.RGBA, x, y, t float64) color.RGBA {
	field := f.GetFieldStrength(x, y, t)
	scaled := ScaleColor(f.BaseColor, field)
	return OverlayColor(c, scaled)
}