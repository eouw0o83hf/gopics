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
	width := 2000
	height := 1500

	// offsets: [x shift, poly0, poly1, poly2]
	makeLayer := func(offsets []float64, c color.RGBA, radius, power float64) GlowingCurveLayer {
		return GlowingCurveLayer{
			FieldFunc: func(x, t float64) float64 {
				x = x + offsets[0]
				coeffs := []float64{
					180 + offsets[1],
					200 + offsets[2],
					150 + offsets[3],
				}
				return coeffs[0]*math.Sin(x/55-t/100) +
					coeffs[1]*math.Cos(x/100-t/22) -
					coeffs[2]*math.Sin(math.Pi/3+x/70+t/1000)
			},
			BaseColor: c,
			Width:     width,
			Height:    height,
			Radius: radius,
			Power: power,
		}
	}

	layers := []GlowingCurveLayer{
		makeLayer(
			[]float64{0, 0, 0, 0},
			color.RGBA{0x10, 0xe0, 0xa0, 0xff},
			20, 2),
		makeLayer(
			[]float64{0, 0, 2, 2},
			color.RGBA{0xe0, 0xe0, 0x30, 0xc0},
			20, 2),
		makeLayer(
			[]float64{0, 5, -3, 10},
			color.RGBA{0xc0, 0x10, 0x60, 0xc0},
			20, 2),
		makeLayer(
			[]float64{-10, 0, 0, 0},
			color.RGBA{0xff, 0xff, 0xff, 0xff},
			10, 0.4),
		makeLayer(
			[]float64{10, 0, 0, 0},
			color.RGBA{0xff, 0xff, 0xff, 0xff},
			10, 0.4),
		makeLayer(
			[]float64{-3, 10, 10, 10},
			color.RGBA{0xff, 0x00, 0xff, 0xff},
			30, 1.5),
		makeLayer(
			[]float64{7, -14, 3, -11},
			color.RGBA{0xff, 0xff, 0x00, 0xff},
			5, 2),
	}
	RenderAvi(func(t float64) RenderedLayer {
		result := make([][]color.RGBA, width)

		for x := 0; x < width; x++ {
			result[x] = make([]color.RGBA, height)
			for y := 0; y < height; y++ {
				accum := NewRgbAa(nil)

				for i := -1; i <= 1; i++ {
					for j := -1; j <= 1; j++ {
						innerAccum := color.RGBA{0, 0, 0, 0}
						for _, l := range layers {
							dX := float64(i) / 2
							dY := float64(j) / 2

							innerAccum = l.GetPixel(innerAccum, float64(x) + dX, float64(y) + dY, t)
						}
						accum = accum.Add(innerAccum)
					}
				}

				result[x][y] = accum.ToColor()

				//accum := color.RGBA{0, 0, 0, 0}
				//for _, l := range layers {
				//	accum = l.GetPixel(accum, float64(x), float64(y), t)
				//}
				//result[x][y] = accum
			}
		}

		return ColorLayer(result)

		//accum := RenderedLayer(
		//	BackgroundLayer(width, height))
		//for _, l := range layers {
		//	accum = l.Overlay(accum, t)
		//}
		//return accum
	}, 500, width, height, "06")
	return

	//layer := GlowingCurveLayer{
	//	FieldFunc: func(x, t float64) float64 {
	//		//return 200 * math.Cos((x - t + math.Cos(t / 10) * 10) / 100)
	//		x = x - 5
	//		return 175*math.Sin(x/55-t/100) + 205*math.Cos(x/100-t/22) - 130*math.Sin(math.Pi/3+x/70+t/1000)
	//	},
	//	BaseColor: color.RGBA{
	//		R: 0x10,
	//		G: 0xe0,
	//		B: 0xa0,
	//		A: 0xff,
	//	},
	//	Width:  width,
	//	Height: height,
	//}
	//layer2 := GlowingCurveLayer{
	//	FieldFunc: func(x, t float64) float64 {
	//		return 180*math.Sin(x/55-t/100) + 200*math.Cos(x/100-t/22) - 150*math.Sin(math.Pi/3+x/70+t/1000)
	//	},
	//	BaseColor: color.RGBA{
	//		R: 0xe0,
	//		G: 0x50,
	//		B: 0x20,
	//		A: 0xff,
	//	},
	//	Width:  width,
	//	Height: height,
	//}
	//
	//RenderAvi(func(t float64) RenderedLayer {
	//	return layer2.Rasterize(t).Overlay(
	//		layer.Rasterize(t), t)
	//}, 100, width, height, "04")
	//return

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
		//{
		//	f: func(x float64) float64 {
		//		moded := float64(x) / 29
		//		return 100*math.Cos(moded+1) + 300*math.Sin(moded/8)
		//	},
		//	c: rgb(0x30, 0x60, 0xa0),
		//},
		//{
		//	f: func(x float64) float64 {
		//		moded := float64(x) / 100
		//		return 10*math.Sin(moded) - float64(x)
		//	},
		//	c: rgb(0x10, 0xc0, 0xa0),
		//},
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

	length := 50 //0

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
			if completed%10 == 0 {
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
