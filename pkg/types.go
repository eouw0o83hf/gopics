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

func OverlayColor(colorA, alphaA, colorB, alphaB uint8) uint8 {
	aContrib := float64(alphaA) * float64(colorA) / 0xff
	bContrib := float64(alphaB) * float64(colorB) / 0xff

	total := uint8(aContrib + bContrib)
	if total > 0xff {
		return  0xff
	}
	return total
}

func (c ColorLayer) Overlay(l RenderedLayer, t float64) RenderedLayer {
	result := make([][]color.RGBA, len(c))

	for x := 0; x < len(c); x++ {
		result[x] = make([]color.RGBA, len(c[x]))
		for y := 0; y < len(c[x]); y++ {
			behind := l.GetRgba(x, y)
			current := c[x][y]

			result[x][y] = color.RGBA{
				R: OverlayColor(behind.R, behind.A, current.R, current.A),
				G: OverlayColor(behind.G, behind.A, current.G, current.A),
				B: OverlayColor(behind.B, behind.A, current.B, current.A),
				A: 0xff,
			}
		}
	}

	return ColorLayer(result)
}

// Layer defines a component of an image stack
// which may or may not contain a colorspace rendering
type Layer interface {
	Overlay(l RenderedLayer, t float64) RenderedLayer
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
