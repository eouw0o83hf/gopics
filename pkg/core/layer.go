package core

import "gopics/pkg/windows"

type Layer interface {
	Render(r *Raster, t int) *Raster
}

// Compiles layers into a func which renders them
func CompileLayers(s ...Layer) func(t int) *Raster {
	return func(t int) *Raster {
		acc := NewRaster(&windows.Testing)
		for _, i := range s {
			acc = i.Render(acc, t)
		}
		return acc
	}
}