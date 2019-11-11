package artisan

import (
	"image"
	"image/png"
	"os"
	"testing"

	graphics "../graphics"
)

func TestMakeImage(t *testing.T) {
	red := `iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAApElEQVR42u3RAQ0AAAjDMO5fNCCDkC5z0HTVrisFCBABASIgQAQEiIAAAQJEQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAQECBAgAgJEQIAIyPcGFY7HnV2aPXoAAAAASUVORK5CYII=`
	img := makeImage(graphics.Drawing{0, 0, "A", red})
	write(t, img, "./test-make-image.png")
}

func TestNewImageMegring(t *testing.T) {
	red := `iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAApElEQVR42u3RAQ0AAAjDMO5fNCCDkC5z0HTVrisFCBABASIgQAQEiIAAAQJEQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAQECBAgAgJEQIAIyPcGFY7HnV2aPXoAAAAASUVORK5CYII=`
	green := `iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAAo0lEQVR42u3RAQ0AAAjDMO5fNCCDkG4SmupdZwoQIAICRECACAgQAQECBIiAABEQIAICRECACIiAABEQIAICRECACIiAABEQIAICRECACIiAABEQIAICRECACIiAABEQIAICRECACIiAABEQIAICRECACIiAABEQIAICRECACIiAABEQIAICRECACIiAABEQIAICRECACAgQIEAEBIiAABGQ7w2x48edS3GF7AAAAABJRU5ErkJggg==`
	blue := `iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAAnElEQVR42u3RAQ0AAAgDoL9/aK3hHFSgyUw4o0KEIEQIQoQgRAhChAgRghAhCBGCECEIEYIQhAhBiBCECEGIEIQgRAhChCBECEKEIAQhQhAiBCFCECIEIQgRghAhCBGCECEIQYgQhAhBiBCECEEIQoQgRAhChCBECEIQIgQhQhAiBCFCEIIQIQgRghAhCBGCECFChCBECEKEIOS7BU5Hx50BmcQaAAAAAElFTkSuQmCC`
	magenta := `iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAAnElEQVR42u3RAQ0AAAgDIN8/9K3hHFQgnXY4I0KEIEQIQoQgRAhChAgRghAhCBGCECEIEYIQhAhBiBCECEGIEIQgRAhChCBECEKEIAQhQhAiBCFCECIEIQgRghAhCBGCECEIQYgQhAhBiBCECEEIQoQgRAhChCBECEIQIgQhQhAiBCFCEIIQIQgRghAhCBGCECFChCBECEKEIOS7BSwYK0jk0PkoAAAAAElFTkSuQmCC`
	xs := []graphics.Drawing{
		graphics.Drawing{0, 0, "A", red},
		graphics.Drawing{0, 1, "B", green},
		graphics.Drawing{0, 2, "C", blue},
		graphics.Drawing{0, 3, "D", magenta},
	}
	img := newImageMerging(xs, len(xs))
	write(t, img, "./new-image-merging.png")

	// magenta, green, red, blue
	xs1 := []graphics.Drawing{
		graphics.Drawing{0, 2, "A", red},
		graphics.Drawing{0, 1, "B", green},
		graphics.Drawing{0, 3, "C", blue},
		graphics.Drawing{0, 0, "D", magenta},
	}
	img1 := newImageMerging(xs1, len(xs1))
	write(t, img1, "./new-image-merging1.png")
}

func write(t *testing.T, img image.Image, path string) {
	out, err := os.Create(path)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = png.Encode(out, img)
	if err != nil {
		t.Errorf(err.Error())
	}
}
