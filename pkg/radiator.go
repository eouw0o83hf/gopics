package pkg

import (
	"image"
	"image/color"
	"math"
)

func GetRadiator(t float64, width, height int) [][]*color.RGBA {
	// initialize field
	field := make([][]float64, width)
	for x := 0; x < width; x++ {
		field[x] = make([]float64, height)
	}

	// define the radiator
	radiator := func(x float64) (y, heat float64) {
		y = 180*math.Sin(x/55-t/100) + 200*math.Cos(x/100-t/22) - 150*math.Sin(math.Pi /3 + x/70 + t/1000)
		// 0-1
		heat = math.Cos(x/300 - t/40)
		heat = heat * heat
		if heat < 0 {
			heat = 0
		}
		return y, heat
	}

	// i,j are used to navigate the array field
	// x,y are the locations on the graph
	for i := 0; i < width; i++ {
		x := i - width/2

		jVal, heat := radiator(float64(x))
		yVal := jVal + float64(height)/2

		radius := 30
		for r := -radius; r <= radius; r++ {
			for s := -radius; s <= radius; s++ {
				y0 := int(yVal) + s

				i0 := i + r
				j0 := y0
				if i0 < 0 || j0 < 0 || i0 >= width || j0 >= width {
					continue
				}

				distance := math.Sqrt(
					math.Pow(float64(r), 2) +
						math.Pow(float64(y0)-yVal, 2))
				// square it
				distance = distance * distance / 4
				fieldStrength := heat / distance
				if fieldStrength >= 0.01 {
					field[i0][j0] = field[i0][j0] + fieldStrength
					if field[i0][j0] > 1 {
						field[i0][j0] = 1
					}
				}
			}
		}
	}

	result := make([][]*color.RGBA, width)
	for i := 0; i < width; i++ {
		result[i] = make([]*color.RGBA, height)
		for j := 0; j < height; j++ {

			rCoef := (math.Sin(float64(i)/80-t/40) + 1) / 2
			gCoef := (math.Cos(math.Pi/3+float64(i)/70-t/45) + 1) / 2
			bCoef := (math.Sin(math.Pi/2+float64(i)/59-t/50) + 1) / 2

			r := uint8(rCoef * field[i][j] * 0xff)
			g := uint8(gCoef * field[i][j] * 0xff)
			b := uint8(bCoef * field[i][j] * 0xff)

			result[i][j] = &color.RGBA{
				R: r,
				G: g,
				B: b,
				A: 0xff,
			}
		}
	}

	return result
}

func GetRadiatorImage(t float64, width, height int) *image.RGBA {
	radiator := GetRadiator(t, width, height)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, radiator[x][y])
		}
	}

	return img
}
