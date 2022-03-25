package pkg

import (
	"bytes"
	"fmt"
	"github.com/icza/mjpeg"
	"gopics/pkg/core"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func LayerToImage(l RenderedLayer) *image.RGBA {
	width := l.GetWidth()
	height := l.GetHeight()

	img := image.NewRGBA(
		image.Rect(0, 0, width, height))

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Set(i, j, l.GetRgba(i, j))
		}
	}

	return img
}

func CopyToLatest(filename string) {
	temp, _ := ioutil.ReadFile(filename)

	ext := filepath.Ext(filename)
	target := fmt.Sprintf("latest%v", ext)

	ioutil.WriteFile(target, temp, 0644)
}

func getFilename(prefix, extension string) string {
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("renders/%v_%v.%v", prefix, timestamp, extension)
}

func RenderPng(l RenderedLayer, prefix string) {
	img := LayerToImage(l)

	filename := getFilename(prefix, "png")

	f, _ := os.Create(filename)
	png.Encode(f, img)

	CopyToLatest(filename)
}

func RenderAvi(layerGen func(t int) *core.Raster, frameCount int, prefix string) {
	frames := make([]bytes.Buffer, frameCount)

	completed := 0

	percentageChunk := frameCount / 20
	if percentageChunk == 0 {
		percentageChunk = 1
	}

	var width, height int

	for i := 0; i < frameCount; i++ {
		func(i int) {
			raster := layerGen(i)
			width = raster.Window.Dx()
			height = raster.Window.Dy()

			img := raster.ToImage()

			var buf bytes.Buffer
			jpeg.Encode(&buf, img, &jpeg.Options{
				Quality: 100,
			})
			frames[i] = buf

			completed++
			if completed%percentageChunk == 0 {
				fmt.Printf("%v%% rendered\n", int(completed*100/frameCount))
			}
		}(i)
	}

	//wg.Wait()

	filename := getFilename(prefix, "avi")

	movie, _ := mjpeg.New(filename, int32(width), int32(height), 30)
	for _, frame := range frames {
		movie.AddFrame(frame.Bytes())
	}
	movie.Close()

	CopyToLatest(filename)
}

func RenderAvi_Deprecated(layerGen func(float642 float64) RenderedLayer, frameCount, width, height int, prefix string) {
	frames := make([]bytes.Buffer, frameCount)

	//var wg sync.WaitGroup
	completed := 0

	percentageChunk := frameCount / 20
	if percentageChunk == 0 {
		percentageChunk = 1
	}

	for i := 0; i < frameCount; i++ {
		//wg.Add(1)
		func(i int) {
			img := LayerToImage(
				layerGen(float64(i)))

			var buf bytes.Buffer
			jpeg.Encode(&buf, img, &jpeg.Options{
				Quality: 100,
			})
			frames[i] = buf

			//wg.Done()

			completed++
			if completed%percentageChunk == 0 {
				fmt.Printf("%v%% rendered\n", int(completed*100/frameCount))
			}
		}(i)
	}

	//wg.Wait()

	filename := getFilename(prefix, "avi")

	movie, _ := mjpeg.New(filename, int32(width), int32(height), 30)
	for _, frame := range frames {
		movie.AddFrame(frame.Bytes())
	}
	movie.Close()

	CopyToLatest(filename)
}
