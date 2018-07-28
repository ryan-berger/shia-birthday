package main

import (
	"strings"
	"image"
	"image/draw"
	"os"
	"image/gif"
	"fmt"
	"unicode"
	"image/color/palette"
	"github.com/ryan-berger/shia-birthday/letter-generator"
	"sync"
	"image/color"
)

type generatorInfo struct {
	letter  string
	filler  *gif.GIF
	content *gif.GIF
	index   int
}


func main() {
	shiaFile, e := os.Open("images/gentlemanparrot.gif")

	if e != nil {
		panic(e)
	}

	shiaHeadFile, e := os.Open("images/jediparrot.gif")

	if e != nil {
		panic(e)
	}

	defer shiaFile.Close()
	defer shiaHeadFile.Close()

	shiaGif, _ := gif.DecodeAll(shiaFile)
	shiaHeadGif, _ := gif.DecodeAll(shiaHeadFile)

	workerChannel := make(chan *generatorInfo)
	waitGroup := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		waitGroup.Add(1)
		go generateLetter(workerChannel, waitGroup)
	}

	for i := 0; i < 26; i++ {
		workerChannel <- &generatorInfo{
			letter: letter_generator.Alphabet[i],
			filler: shiaGif,
			content: shiaHeadGif,
			index: i,
		}
	}

	close(workerChannel)

	waitGroup.Wait()
}

func generateLetter(work chan *generatorInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	for info := range work {
		letterGif := &gif.GIF{}

		for i := 0; i < 10; i++ {
			letterGif.Image = append(letterGif.Image, image.NewPaletted(image.Rect(0, 0, 60, 100), palette.Plan9))
			letterGif.Delay = append(letterGif.Delay, 0)
		}

		trimmed := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, info.letter)

		for i := 0; i < 10; i++ {
			draw.Draw(letterGif.Image[i], letterGif.Image[i].Bounds(), &image.Uniform{color.RGBA{255, 255, 255 , 0}}, image.ZP, draw.Over)
			for j := 0; j < 5; j++ {
				for k := 1; k < 4; k++ {
					adjustedIndex := (j * 5) + k
					var selectedGif *image.Paletted
					if rune(trimmed[adjustedIndex]) == rune('-') {
						selectedGif = info.filler.Image[i]
					} else if rune(trimmed[adjustedIndex]) == rune('0') {
						selectedGif = info.content.Image[i]
					}
					draw.Draw(letterGif.Image[i], image.Rect((k - 1)*20, j*20, ((k - 1) *20)+20, (j*20)+20), selectedGif, image.ZP, draw.Over)
				}
			}
			fmt.Println(fmt.Sprintf("letter: %d cycle: %d", info.index, i))
		}
		f, _ := os.Create(fmt.Sprintf("letter-generator/letters/%d.gif", info.index))
		gif.EncodeAll(f, letterGif)
	}
}
