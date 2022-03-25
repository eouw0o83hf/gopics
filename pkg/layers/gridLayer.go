package layers

import (
	"gopics/pkg/core"
	"image/color"
)

type GridLayer struct {
	Size     int
	OnColor  color.RGBA
	OffColor color.RGBA
}

func (g GridLayer) Render(r *core.Raster, t int) *core.Raster {
	r.Traverse(func(i, j, x, y int) {
		xMatch := i/g.Size%2 == 0
		yMatch := j/g.Size%2 == 0

		var c color.RGBA
		if xMatch != yMatch {
			c = g.OnColor
		} else {
			c = g.OffColor
		}

		r.Bitmap[i][j] = c
	})

	return r
}
