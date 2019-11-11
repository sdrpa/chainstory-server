package ipfs

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"

	graphics "../graphics"
	shell "github.com/ipfs/go-ipfs-api"
)

const (
	ipfsURL = "localhost:5001"
)

// Add - Store the drawing to IPFS
func Add(drawing graphics.Drawing) string {
	sh := shell.NewShell(ipfsURL)
	hash, err := sh.Add(encode(drawing))
	if err != nil {
		log.Panic(err)
	}
	return hash
}

// Read - Read the drawing out of IPFS
func Read(hash string) (graphics.Drawing, error) {
	sh := shell.NewShell(ipfsURL)
	reader, err := sh.Cat(hash)
	if err != nil {
		return graphics.MakeDrawing(), err
	}
	drawing, err := decode(reader)
	reader.Close()
	if err != nil {
		return graphics.MakeDrawing(), err
	}
	return drawing, nil
}

func encode(v interface{}) io.Reader {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(v)
	if err != nil {
		log.Panic(err)
	}
	return bytes.NewReader(buff.Bytes())
}

func decode(r io.Reader) (graphics.Drawing, error) {
	dec := gob.NewDecoder(r)
	var drawing graphics.Drawing
	err := dec.Decode(&drawing)
	if err != nil {
		return graphics.MakeDrawing(), err
	}
	return drawing, nil
}
