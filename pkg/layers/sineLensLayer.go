package layers

import (
	"gopics/pkg/core"
	"math"
)

type SineLensLayer struct {
	// width of lens
	Width float64
	// radius of lens (should be wider than width)
	Radius 		float64
	Depth      float64
}

func (s SineLensLayer) fieldStrength(x int, t int) float64 {
	xf := float64(x)
	tf := float64(t)

	return math.Sin(xf/100+tf/50) * 150
}

func (s SineLensLayer) Render(r *core.Raster, t int) *core.Raster {
	for i := 0; i < r.Window.Dx(); i++ {
		x := r.ItoX(i)
		yBase := s.fieldStrength(x, t)

		for z := -s.Width; z <= s.Width; z++ {
			rho := z * (math.Pi / 2) / s.Radius
			lensDepth := math.Cos(rho)

			lensSlope := -math.Sin(rho)       // derivative of depth
			lensAngle := math.Atan(lensSlope) // slope to angle

			reflectedAngle := math.Pi - lensAngle * 2

			d := s.Depth + lensDepth
			dy := d * math.Tan(reflectedAngle)
			y := yBase + z + dy

			// todo AA
			underlying := r.GetRgba(i, r.YToJ(int(y)))
			r.SetRgba(i, r.YToJ(int(yBase + z)), underlying)
		}
	}

	return r
}
