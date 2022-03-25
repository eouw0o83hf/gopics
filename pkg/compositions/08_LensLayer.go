package compositions

import (
	"gopics/pkg"
	"gopics/pkg/colors"
	"gopics/pkg/core"
	"gopics/pkg/layers"
)

func Run_C08() {
	grid := layers.GridLayer{
		Size:     10,
		OnColor:  colors.White,
		OffColor: colors.Black,
	}

	lens := layers.SineLensLayer{
		Width: 50,
		Radius:     80,
		Depth:      1,
	}

	compiled := core.CompileLayers(grid, lens)

	pkg.RenderAvi(compiled, 50, "08")
}
