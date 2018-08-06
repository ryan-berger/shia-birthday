package wordGif

import (
	"image/gif"
	"sync"
)

type entry struct {
	value *gif.GIF
	ready chan struct{}
}

type Cache struct {
	mu *sync.Mutex
	cache map[string]*entry
	workerPool *WorkerPool
}

func (cache *Cache) Get(text string) *gif.GIF {
	cache.mu.Lock()
	e := cache.cache[text]
	if e == nil {
		newE := &entry{
			ready: make(chan struct{}),
		}
		cache.cache[text] = newE
		cache.mu.Unlock()

		res := make(chan *gif.GIF)
		cache.workerPool.MakeRequest(text, res)
		newE.value = <- res
		close(e.ready)
	} else {
		cache.mu.Unlock()
		<- e.ready
	}

	return e.value
}

func NewCache() *Cache {
	return &Cache{
		mu: &sync.Mutex{},
		cache: make(map[string]*entry),
		workerPool: NewWorkerPool(),
	}
}
