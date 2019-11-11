package graphics

// Drawing represents pixels data
type Drawing struct {
	Version int
	Index   int
	Message string
	Data    string
}

// MakeDrawing - Creates a transparent 1x1 drawing
func MakeDrawing() Drawing {
	// Transparent image dataURL
	const transparent = `iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=`
	return Drawing{0, 0, "transparent", transparent}
}
