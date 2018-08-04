package wordGif

import (
	"image/gif"
	"sync"
)

type entry struct {

}

type Cache struct {
	mu *sync.Mutex
	cache map[string]*gif.GIF
}

func (cache *Cache) Get(text string)  {
}

func NewCache()  {

}
