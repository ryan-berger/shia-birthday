package concatter

import (
	"image"
	"sync"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"github.com/pborman/uuid"
)

var cache = NewGifCache("letter-generator/letters", "images/shia.gif")

type frame struct {
	letters []*image.Paletted
	filler  *image.Paletted
	index   int
}

type frameInformation struct {
	frame    *frame
	finished chan *workResult
}

type workResult struct {
	frame *image.Paletted
	index int
}

type gifRequest struct {
	id   uuid.UUID
	text string
}

type gifResult struct {
	requestId uuid.UUID
	text      string
	result    *gif.GIF
}

type Worker struct {
	textChan  chan *gifRequest
	workQueue chan *frameInformation
	gifResult chan *gifResult
	group     *sync.WaitGroup
}

func (w *Worker) work(info *frameInformation) {
	defer w.group.Done()
	work := info.frame
	newImage := image.NewPaletted(image.Rect(0, 0, len(work.letters)*80+20, 140), palette.Plan9)
	drawBorder(newImage, work.index, len(work.letters)*5)
	drawLines(newImage, work.index, len(work.letters)+1)
	for index, img := range work.letters {
		draw.Draw(newImage, image.Rect((index*60)+(index*20)+20, 20, (index*60)+(index*20)+80, 120), img, image.ZP, draw.Over)
	}

	info.finished <- &workResult{newImage, work.index}
}

func (w *Worker) listenForWork() {
	for {
		select {
		case work := <-w.workQueue:
			w.work(work)
		}
	}
}

func (w *Worker) spinUpWorkers() {
	for i := 0; i < 60; i++ {
		go w.listenForWork()
	}
}

func (w *Worker) makeFrames(text string) chan *workResult {
	result := make(chan *workResult, 60)

	for i := 0; i < 60; i++ {
		w.group.Add(1)
		r := &frame{
			letters: cache.GetFrames(text, i),
			index:   i,
		}

		w.workQueue <- &frameInformation{r, result}
	}

	return result
}

func processFrames(results chan *workResult) *gif.GIF {
	close(results)
	images := [60]*image.Paletted{}

	for res := range results {
		images[res.index] = res.frame
	}

	newGif := &gif.GIF{}
	for _, img := range images {
		newGif.Image = append(newGif.Image, img)
		newGif.Delay = append(newGif.Delay, 0)
	}
	return newGif
}

func (w *Worker) listenForRequests() {
	for {
		select {
		case request := <-w.textChan:
			results := w.makeFrames(request.text)
			w.group.Wait()
			w.gifResult <- &gifResult{requestId: request.id, result: processFrames(results), text: request.text}
		}
	}
}

func NewWorker(textChan chan *gifRequest, gifResult chan *gifResult) *Worker {
	worker := &Worker{
		textChan:  textChan,
		gifResult: gifResult,
		workQueue: make(chan *frameInformation),
	}
	worker.group = &sync.WaitGroup{}
	worker.spinUpWorkers()
	go worker.listenForRequests()
	return worker
}

func drawBorder(img *image.Paletted, frame, length int) {
	filler := cache.GetFiller(frame)
	for i := 0; i < length; i++ {
		draw.Draw(img, image.Rect(i*20, 0, (i*20)+20, 20), filler, image.ZP, draw.Over)
		draw.Draw(img, image.Rect(i*20, 120, (i*20)+20, 140), filler, image.ZP, draw.Over)
	}
}

func drawLines(img *image.Paletted, frame, length int) {
	filler := cache.GetFiller(frame)
	for i := 0; i < length; i++ {
		for j := 0; j < 5; j++ {
			draw.Draw(img, image.Rect(i*80, (j*20)+20, (i*80)+20, (j*20)+40), filler, image.ZP, draw.Over)
		}
	}
}
