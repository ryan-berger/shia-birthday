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
)

type generatorInfo struct {
	letter  string
	filler  *gif.GIF
	content *gif.GIF
	index   int
}


func main() {
	shiaFile, _ := os.Open("images/shia.gif")
	shiaHeadFile, _ := os.Open("images/shiaHead.gif")

	defer shiaFile.Close()
	defer shiaFile.Close()

	shiaGif, _ := gif.DecodeAll(shiaFile)
	shiaHeadGif, _ := gif.DecodeAll(shiaHeadFile)

	workerChannel := make(chan *generatorInfo)
	waitGroup := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
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

	waitGroup.Wait()
}

func generateLetter(work chan *generatorInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	for info := range work {
		trimmed := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, info.letter)

		letterGif := &gif.GIF{}

		for i := 0; i < 60; i++ {
			newImage := image.NewPaletted(image.Rect(0, 0, 640, 640), palette.Plan9)

			for j := 0; j < 5; j++ {
				for k := 0; k < 5; k++ {
					adjustedIndex := (j * 5) + k
					var selectedGif *image.Paletted
					if rune(trimmed[adjustedIndex]) == rune('-') {
						selectedGif = info.filler.Image[i]
					} else {
						selectedGif = info.content.Image[i]
					}

					draw.Draw(newImage, image.Rect(k*128, j*128, (k*128)+128, (j*128)+128), selectedGif, image.ZP, draw.Over)
				}
			}

			fmt.Println(fmt.Sprintf("finished cycle #%d", i))

			letterGif.Image = append(letterGif.Image, newImage)
			letterGif.Delay = append(letterGif.Delay, 0)
		}
		f, _ := os.Create(fmt.Sprintf("letter-generator/letters/%d.gif", info.index))
		gif.EncodeAll(f, letterGif)
	}
}
