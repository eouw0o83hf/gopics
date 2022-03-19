package pkg

import (
	"bytes"
	"fmt"
	"github.com/icza/mjpeg"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math"
	"sync"
	"time"
)

func MakeAPng() {
	width := 1000
	height := 1000

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

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
		f func(float642 float64) float64
		c color.RGBA
	}{
		{
			f: func(x float64) float64 {
				moded := float64(x) / 100
				return math.Sin(moded-2) * 300
			},
			c: rgb(0xff, 0x60, 0x10),
		},
		{
			f: func(x float64) float64 {
				moded := float64(x) / 29
				return 100*math.Cos(moded+1) + 300*math.Sin(moded/8)
			},
			c: rgb(0x30, 0x60, 0xa0),
		},
		{
			f: func(x float64) float64 {
				moded := float64(x) / 100
				return 10*math.Sin(moded) - float64(x)
			},
			c: rgb(0x10, 0xc0, 0xa0),
		},
	}

	// Set color for each pixel.
	// initial image implementation with some waves
	_ = func(xOffset, yOffset float64) *image.RGBA {
		img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {

				xEffective := float64(x) + xOffset
				yEffective := float64(y) + yOffset

				viewportX := xEffective - (float64(width) / 2)
				viewportY := (float64(height) / 2) - yEffective

				accum := rgb(0, 0, 0)

				for _, wave := range wavesToColors {
					target := wave.f(viewportX)
					diff := math.Abs(target - float64(viewportY))

					maxDiff := float64(180)
					if diff < maxDiff {
						interpolation := math.Pow(diff/maxDiff, 2)
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

		return img
	}

	//renderPng := func() {
	//	img := getImage(0, 0)
	//
	//	// Encode as PNG.
	//	timestamp := time.Now().Format("20060102_150405")
	//	filename := fmt.Sprintf("renders/00_%v.png", timestamp)
	//	f, _ := os.Create(filename)
	//	png.Encode(f, img)
	//
	//	f, _ = os.Create("latest.png")
	//	png.Encode(f, img)
	//}

	//outGif := &gif.GIF{}

	length := 500

	var wg sync.WaitGroup

	frames := make([]bytes.Buffer, length)
	completed := 0

	for i := 0; i < length; i++ {
		localI := i
		wg.Add(1)
		go func() {
			//x := float64(localI) * 2
			//y := float64(localI) * 0.3
			//
			//img := getImage(x, y)
			img := GetRadiatorImage(float64(localI), width, height)

			var buf bytes.Buffer
			jpeg.Encode(&buf, img, &jpeg.Options{
				Quality: 100,
			})
			frames[localI] = buf

			wg.Done()
			completed++
			if completed % 10 == 0 {
				fmt.Printf("finished frame %v / %v\n", completed, length)
			}
		}()
	}

	wg.Wait()

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("renders/01_%v.avi", timestamp)

	movie, _ := mjpeg.New(filename, int32(width), int32(height), 30)
	for _, frame := range frames {
		movie.AddFrame(frame.Bytes())
	}
	movie.Close()

	temp, _ := ioutil.ReadFile(filename)
	ioutil.WriteFile("latest.avi", temp, 0644)
}
