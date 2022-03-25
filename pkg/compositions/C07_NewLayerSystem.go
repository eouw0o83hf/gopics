package compositions

import (
	"gopics/pkg"
	"gopics/pkg/colors"
	"gopics/pkg/core"
	"gopics/pkg/layers"
)

func Run_C07() {
	layer := layers.GridLayer{
		Size:     10,
		OnColor:  colors.White,
		OffColor: colors.Black,
	}

	compiled := core.CompileLayers(layer)

	pkg.RenderAvi(compiled, 50, "07")
}
