package common

import (
	"image/gif"
	"os"
)

func LoadGif(path string) (*gif.GIF, error)  {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return gif.DecodeAll(f)
}
