package graphics

// https://github.com/ozankasikci/go-image-merge/blob/master/go-image-merge.go

import (
	"image"
	"image/color"
	"image/draw"
)

// Specifies how the grid pixel size should be calculated
type gridSizeMode int

const (
	// The size in pixels is fixed for all the grids
	fixedGridSize gridSizeMode = iota
	// The size in pixels is set to the nth image size
	gridSizeFromImage
)

// Grid holds the data for each grid
type Grid struct {
	Image           image.Image
	BackgroundColor color.Color
	OffsetX         int
	OffsetY         int
	Grids           []*Grid
}

// MergedImage is responsible for merging the given images
type MergedImage struct {
	Grids           []*Grid
	ImageCountDX    int
	ImageCountDY    int
	BaseDir         string
	FixedGridSizeX  int
	FixedGridSizeY  int
	GridSizeMode    gridSizeMode
	GridSizeFromNth int
}

// NewMergedImage returns a new *MergedImage instance
func NewMergedImage(grids []*Grid, imageCountDX, imageCountDY int, opts ...func(*MergedImage)) *MergedImage {
	mi := &MergedImage{
		Grids:        grids,
		ImageCountDX: imageCountDX,
		ImageCountDY: imageCountDY,
	}

	for _, option := range opts {
		option(mi)
	}

	return mi
}

// OptBaseDir is an functional option to set the BaseDir field
func OptBaseDir(dir string) func(*MergedImage) {
	return func(mi *MergedImage) {
		mi.BaseDir = dir
	}
}

// OptGridSize is an functional option to set the GridSize X & Y
func OptGridSize(sizeX, sizeY int) func(*MergedImage) {
	return func(mi *MergedImage) {
		mi.GridSizeMode = fixedGridSize
		mi.FixedGridSizeX = sizeX
		mi.FixedGridSizeY = sizeY
	}
}

// OptGridSizeFromNthImageSize is an functional option to set the GridSize from the nth image
func OptGridSizeFromNthImageSize(n int) func(*MergedImage) {
	return func(mi *MergedImage) {
		mi.GridSizeMode = gridSizeFromImage
		mi.GridSizeFromNth = n
	}
}

func (m *MergedImage) mergeGrids(images []image.Image) (*image.RGBA, error) {
	var canvas *image.RGBA
	imageBoundX := 0
	imageBoundY := 0

	if m.GridSizeMode == fixedGridSize && m.FixedGridSizeX != 0 && m.FixedGridSizeY != 0 {
		imageBoundX = m.FixedGridSizeX
		imageBoundY = m.FixedGridSizeY
	} else if m.GridSizeMode == gridSizeFromImage {
		imageBoundX = images[m.GridSizeFromNth].Bounds().Dx()
		imageBoundY = images[m.GridSizeFromNth].Bounds().Dy()
	} else {
		imageBoundX = images[0].Bounds().Dx()
		imageBoundY = images[0].Bounds().Dy()
	}

	canvasBoundX := m.ImageCountDX * imageBoundX
	canvasBoundY := m.ImageCountDY * imageBoundY

	canvasMaxPoint := image.Point{canvasBoundX, canvasBoundY}
	canvasRect := image.Rectangle{image.Point{0, 0}, canvasMaxPoint}
	canvas = image.NewRGBA(canvasRect)

	// Draw grids one by one
	for i, grid := range m.Grids {
		img := images[i]
		x := i % m.ImageCountDX
		y := i / m.ImageCountDX
		minPoint := image.Point{x * imageBoundX, y * imageBoundY}
		maxPoint := minPoint.Add(image.Point{imageBoundX, imageBoundY})
		nextGridRect := image.Rectangle{minPoint, maxPoint}

		if grid.BackgroundColor != nil {
			draw.Draw(canvas, nextGridRect, &image.Uniform{grid.BackgroundColor}, image.Point{}, draw.Src)
			draw.Draw(canvas, nextGridRect, img, image.Point{}, draw.Over)
		} else {
			draw.Draw(canvas, nextGridRect, img, image.Point{}, draw.Src)
		}

		if grid.Grids == nil {
			continue
		}

		// Draw top layer grids
		for _, grid := range grid.Grids {
			gridRect := nextGridRect.Bounds().Add(image.Point{grid.OffsetX, grid.OffsetY})
			draw.Draw(canvas, gridRect, grid.Image, image.Point{}, draw.Over)
		}
	}

	return canvas, nil
}

// Merge merges the images according to given configuration
func (m *MergedImage) Merge() (*image.RGBA, error) {
	var images []image.Image
	for _, grid := range m.Grids {
		images = append(images, grid.Image)
	}
	return m.mergeGrids(images)
}
