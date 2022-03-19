package pkg

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"
)

func MakeAPng() {
	width := 1000
	height := 1000

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	min := func(a, b int) uint8 {
		if a < b {
			return uint8(a)
		}
		return uint8(b)
	}
	rgb := func(r, g, b int) color.RGBA {
		return color.RGBA{
			min(r, 0xff),
			min(g, 0xff),
			min(b, 0xff),
			0xff}
	}

	wavesToColors := []struct {
		f func(int) float64
		c color.RGBA
	}{
		{
			f: func(x int) float64 {
				moded := float64(x) / 100
				return math.Sin(moded - 2) * 300
			},
			c: rgb(0xff, 0x60, 0x10),
		},
		{
			f: func(x int) float64 {
				moded := float64(x) / 29
				return 100*math.Cos(moded + 1) + 300*math.Sin(moded/8)
			},
			c: rgb(0x30, 0x60, 0xa0),
		},
		{
			f: func(x int) float64 {
				moded := float64(x) / 100
				return 10 * math.Sin(moded) - float64(x)
			},
			c: rgb(0x10, 0xc0, 0xa0),
		},
	}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			viewportX := x - (width / 2)
			viewportY := (height / 2) - y

			accum := rgb(0, 0, 0)

			for _, wave := range wavesToColors {
				target := wave.f(viewportX)
				diff := math.Abs(target - float64(viewportY))

				maxDiff := float64(180)
				if diff < maxDiff {
					interpolation := math.Pow(diff / maxDiff, 2)
					glow := rgb(
						int(float64(wave.c.R)*interpolation),
						int(float64(wave.c.G)*interpolation),
						int(float64(wave.c.B)*interpolation))

					accum = rgb(
						int(accum.R+glow.R),
						int(accum.G+glow.G),
						int(accum.B+glow.B))
				}
			}

			img.Set(x, y, accum)
		}
	}

	// Encode as PNG.
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("renders/00_%v.png", timestamp)
	f, _ := os.Create(filename)
	png.Encode(f, img)

	f, _ = os.Create("latest.png")
	png.Encode(f, img)
}
