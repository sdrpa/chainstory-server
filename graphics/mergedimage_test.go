package graphics

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"testing"
)

func randImage(w, h int) image.Image {
	randColor := func() color.RGBA {
		cs := []color.RGBA{
			color.RGBA{255, 0, 0, 255},
			color.RGBA{255, 128, 0, 255},
			color.RGBA{255, 255, 0, 255},
			color.RGBA{128, 255, 0, 255},
			color.RGBA{0, 255, 0, 255},
			color.RGBA{0, 255, 128, 255},
			color.RGBA{0, 255, 255, 255},
			color.RGBA{0, 128, 255, 255},
			color.RGBA{0, 0, 255, 255},
			color.RGBA{128, 0, 255, 255},
			color.RGBA{255, 0, 255, 255},
			color.RGBA{255, 0, 128, 255},
		}
		min := 0
		max := len(cs) - 1
		idx := rand.Intn(max-min+1) + min
		//fmt.Println(idx)
		return cs[idx]
	}

	createSolid := func(w, h int, color color.RGBA) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				img.Set(x, y, color)
			}
		}
		return img
	}

	createRand := func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				img.Set(x, y, randColor())
			}
		}
		return img
	}

	colored := createRand(w, h)
	transparent := createSolid(w, h, color.RGBA{255, 255, 255, 0})
	flag := rand.Float32() < 0.5
	if flag {
		return colored
	}
	return transparent
}

func TestMergedImage(t *testing.T) {
	var grids []*Grid
	for i := 0; i < 4096; i++ {
		g := Grid{Image: randImage(16, 16)}
		grids = append(grids, &g)
	}
	rgba, err := NewMergedImage(grids, 64, 64).Merge()
	if err != nil {
		log.Panic(err)
	}
	file, err := os.Create("./test.png")
	if err != nil {
		log.Panic(err)
	}
	encodeErr := png.Encode(file, rgba)
	if encodeErr != nil {
		log.Panic(encodeErr)
	}
}
