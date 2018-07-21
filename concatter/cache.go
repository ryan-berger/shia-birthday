package concatter

import (
	"image"
	"os"
	"fmt"
	"image/gif"
)

type GifCache struct {
	cache [][]*image.Paletted
	filler *gif.GIF
}

func (gifCache *GifCache) GetFrames(text string, index int) []*image.Paletted {
	var frame []*image.Paletted

	for _, char := range text {
		frame = append(frame, gifCache.cache[rune(char - 'a')][index])
	}

	return frame
}

func (gifCache *GifCache) GetFiller(frame int) *image.Paletted {
	return gifCache.filler.Image[frame]
}

func NewGifCache(directory, filler string) *GifCache {
	gifCache := &GifCache{}

	f, _ := os.Open(filler)
	g, _ := gif.DecodeAll(f)
	f.Close()

	gifCache.filler = g

	for i := 0; i < 26; i++ {
		f, e := os.Open(fmt.Sprintf("%s/%d.gif", directory, i))

		if e != nil {
			panic(e)
		}

		g, e := gif.DecodeAll(f)

		if e != nil {
			panic(e)
		}

		gifCache.cache = append(gifCache.cache, g.Image)
		f.Close()
	}
	return gifCache
}
