package pkg

import (
	"image/color"
)

// default implementation of a layer that's just colors
type ColorLayer [][]color.RGBA

func BackgroundLayer(width, height int) ColorLayer {
	black := color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0xff,
	}

	layer := make([][]color.RGBA, width)
	for i := 0; i < width; i++ {
		layer[i] = make([]color.RGBA, height)
		for j := 0; j < height; j++ {
			layer[i][j] = black
		}
	}

	return layer
}

func (c ColorLayer) GetRgba(x, y int) *color.RGBA {
	return &c[x][y]
}

func (c ColorLayer) GetWidth() int {
	if c == nil {
		return 0
	}
	return len(c)
}

func (c ColorLayer) GetHeight() int {
	if c == nil || len(c) == 0 {
		return 0
	}
	return len(c[0])
}

func (c ColorLayer) Overlay(l RenderedLayer, t float64) RenderedLayer {
	result := make([][]color.RGBA, len(c))

	for x := 0; x < len(c); x++ {
		result[x] = make([]color.RGBA, len(c[x]))
		for y := 0; y < len(c[x]); y++ {
			behind := l.GetRgba(x, y)
			current := c[x][y]

			result[x][y] = OverlayColor(*behind, current)
		}
	}

	return ColorLayer(result)
}

func (a ColorLayer) GetPixel(c color.RGBA, x, y, t float64) color.RGBA {
	current := a[int(x)][int(y)]
	return OverlayColor(c, current)
}

// Layer defines a component of an image stack
// which may or may not contain a colorspace rendering
type Layer interface {
	// l is the underlying layer
	Overlay(l RenderedLayer, t float64) RenderedLayer
	// c is the underlying color
	GetPixel(c color.RGBA , x, y, t float64) color.RGBA
}

// RenderedLayer denotes a Layer which has been
// rendered, thus containing a colorspace matrix
type RenderedLayer interface {
	Layer
	GetRgba(x, y int) *color.RGBA
	GetWidth() int
	GetHeight() int
}

// Field defines a component which exists as
// the backing field to a particular Layer
type Field interface {
	GetFieldStrength(x, y, t float64) float64
}

// accumulator for an anti-aliasing color
type RgbAa struct {
	R float64
	G float64
	B float64
	Count int
}

func NewRgbAa(c *color.RGBA) RgbAa {
	if c == nil {
		return RgbAa{
			R:     0,
			G:     0,
			B:     0,
			Count: 0,
		}
	}

	return RgbAa{
		R:     float64(c.R),
		G:     float64(c.G),
		B:     float64(c.B),
		Count: 1,
	}
}

func (a RgbAa) Add(c color.RGBA) RgbAa  {
	return RgbAa{
		R:     a.R + float64(c.R),
		G:     a.G + float64(c.G),
		B:     a.B + float64(c.B),
		Count: a.Count + 1,
	}
}

func (a RgbAa) ToColor() color.RGBA {
	return color.RGBA{
		R: uint8(a.R / float64(a.Count)),
		G: uint8(a.G / float64(a.Count)),
		B: uint8(a.B / float64(a.Count)),
		A: 0xff,
	}
}
