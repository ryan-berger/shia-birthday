package letter

import (
	"image/gif"
	"sync"
	"image"
)

type entry struct {
	letterGif *gif.GIF
	ready     chan struct{}
}

type Cache struct {
	cache map[rune]*entry
	mu *sync.Mutex
}

func NewCache(filler Filler) {
	letterCache := &Cache{
		cache:make(map[rune]*entry),
		mu:&sync.Mutex{},
	}

	for i, letter := range Alphabet {
		e := &entry{ready: make(chan struct{})}
		key := rune('a' + i)
		letterCache.cache[key] = e
		go func() {
			letterCache.cache[key].letterGif = letter.generateGif(filler)
			close(e.ready)
		}()
	}
}

func (letterCache *Cache) Get(letter rune) *gif.GIF {
	letterCache.mu.Lock()
	e := letterCache.cache[letter]
	letterCache.mu.Unlock()
	<- e.ready
	return e.letterGif
}

func (letterCache *Cache) GetLettersAt(text string, frame int) []*image.Paletted  {
	var frames []*image.Paletted

	for _, letter := range text {
		frames = append(frames, letterCache.Get(letter).Image[frame])
	}

	return frames
}