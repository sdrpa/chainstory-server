// Package artisan is constantly checking the blockchain and build the image
package artisan

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
	"time"

	blockchain "../blockchain"
	graphics "../graphics"
	ipfs "../ipfs"
)

// Private cache holding fetched objects
type object struct {
	Drawing     graphics.Drawing
	Transaction blockchain.Transaction
}

var cache map[int]object

// Creates an image from a drawing
func makeImage(drawing graphics.Drawing) image.Image {
	// Must trim leading "data:image/png;base64," from dataURL if exists
	data := strings.TrimPrefix(drawing.Data, "data:image/png;base64,")
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Panic(err)
	}
	return img
}

// Merge drawings to an image
// Note: drawings count must be power of 2
func newImageMerging(xs []graphics.Drawing, count int) *image.RGBA {
	// Map to [index]Drawing
	m := make(map[int]graphics.Drawing, count)
	for _, d := range xs {
		m[d.Index] = d
	}
	// Insert transparent drawing where necessary
	var ys []graphics.Drawing
	for i := 0; i < count; i++ {
		value, exists := m[i]
		if exists {
			ys = append(ys, value)
		} else {
			ys = append(ys, graphics.MakeDrawing())
		}
	}
	var grids []*graphics.Grid
	for _, drawing := range ys {
		grid := graphics.Grid{Image: makeImage(drawing)}
		grids = append(grids, &grid)
	}
	// Determine number of columns and rows
	sqrt := int(math.Sqrt(float64(count)))
	rgba, err := graphics.NewMergedImage(grids, sqrt, sqrt).Merge()
	if err != nil {
		log.Panic(err)
	}
	return rgba
}

// Publish - Publish image to a file at path
func publish(img *image.RGBA, path string) {
	out, err := os.Create(path)
	if err != nil {
		log.Panic(err)
	}
	err = png.Encode(out, img)
	if err != nil {
		log.Panic(err)
	}
	out.Close()
}

func paint(client blockchain.Client, offset int, c chan int) {
	padded := fmt.Sprintf("%4v", offset)
	fmt.Println("Fetching offset: " + padded + ". " + time.Now().Format(time.RFC850))

	transactions := client.GetValidTransactions(offset)
	fmt.Println("            got: " + fmt.Sprintf("%4v", len(transactions)) + " transactions")
	// Read drawing from IPFS and cache the drawing
	for _, transaction := range transactions {
		drawing, err := ipfs.Read(transaction.Asset.Data)
		if err == nil {
			cache[drawing.Index] = object{drawing, transaction}
		}
	}
	// Loop throught the cache and build drawings array
	var ds []graphics.Drawing
	count := 4096
	for i := 0; i < count; i++ {
		// If there is a drawing for index append
		if object, ok := cache[i]; ok {
			ds = append(ds, object.Drawing)
		} else {
			transparent16x16 := `iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAQAAAC1+jfqAAAAEUlEQVR42mNkIAAYRxWMJAUAE5gAEdz4t9QAAAAASUVORK5CYII=`
			ds = append(ds, graphics.Drawing{0, i, "", transparent16x16})
		}
	}
	// Build the image from drawings array
	img := newImageMerging(ds, count)
	publish(img, "./public/image.png")

	// Advance and wrap offset if necessary
	offset = offset + 100
	if offset > 4000 || len(transactions) == 0 {
		offset = 0
	}

	c <- offset
}

// MARK: - Public

// Run - runs a scheduler
func Run(client blockchain.Client) {
	cache = make(map[int]object, 4096)
	c := make(chan int)
	var offset = 0

	go paint(client, offset, c)

	for offset := range c {
		go func(offset int) {
			time.Sleep(time.Second * 5) // 5 sec timeout
			paint(client, offset, c)
		}(offset)
	}
}

// DrawingBase64 - Returns base64 encoded string of image at path
func DrawingBase64(path string) ([]int, string) {
	img, err := os.Open(path)
	defer img.Close()
	if err != nil {
		log.Panic(err)
	}
	// Create a new buffer base on file size
	info, _ := img.Stat()
	var size int64 = info.Size()
	buf := make([]byte, size)
	// Read file content into buffer
	reader := bufio.NewReader(img)
	reader.Read(buf)
	// Convert the buffer bytes to base64 string
	base64String := base64.StdEncoding.EncodeToString(buf)
	// Build colored indexes array
	var indexes []int
	for i := 0; i < 4096; i++ {
		if _, ok := cache[i]; ok {
			indexes = append(indexes, i)
		}
	}
	return indexes, base64String
}

// DrawingAt - Returns drawing at index from cache or error if not available
func DrawingAt(index int) (graphics.Drawing, blockchain.Transaction, error) {
	if object, ok := cache[index]; ok {
		return object.Drawing, object.Transaction, nil
	}
	return graphics.MakeDrawing(), blockchain.MakeTransaction(), errors.New("drawing at index == nil")
}
