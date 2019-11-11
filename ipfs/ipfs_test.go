package ipfs

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
	"strings"
	"testing"

	graphics "../graphics"
)

func TestEncodeDecode(t *testing.T) {
	red := `iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAApElEQVR42u3RAQ0AAAjDMO5fNCCDkC5z0HTVrisFCBABASIgQAQEiIAAAQJEQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAREQIAICBABASIgQAQECBAgAgJEQIAIyPcGFY7HnV2aPXoAAAAASUVORK5CYII=`
	drawing := graphics.Drawing{0, 0, "A", red}
	reader := encode(drawing)
	//fmt.Println(reader)
	decoded, err := decode(reader)
	if err != nil {
		t.Errorf("decoded == nil")
	}
	//fmt.Println(decoded)
	if decoded != drawing {
		t.Errorf("decoded != drawing")
	}
}

func TestMakeImage(t *testing.T) {
	data := `iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAASUlEQVQ4T2NkoBAwUqifgaAB/xkY/jMy4FaH1wCQZpALyTIAZjNZLkDWRLIB6BpINgA9VkYNACcE8hMSKEAH3gBCmY1gZiJkAAD82B4RwbuIHQAAAABJRU5ErkJggg==`
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	_, _, err := image.Decode(reader)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestReadAndMakeImage(t *testing.T) {
	hash := "QmcuDkAXgt4YNEVqKUcXskynaUALkQ6XdY9oFDfqgpCWg7"
	drawing, err := Read(hash)
	if err != nil {
		t.Errorf("drawing == nil")
	}
	// Must trim "data:image/png;base64," from dataURL
	data := strings.TrimPrefix(drawing.Data, "data:image/png;base64,")
	fmt.Println(data)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	_, _, decodeErr := image.Decode(reader)
	if decodeErr != nil {
		t.Errorf(decodeErr.Error())
	}
}
