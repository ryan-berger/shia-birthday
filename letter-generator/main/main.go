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
)

type generatorInfo struct {
	letter   string
	filler   *gif.GIF
	content  *gif.GIF
}

const A = `
SHHHS
SHSHS
SHHHS
SHSHS
SHSHS
`

func main()  {
	shiaFile, _ := os.Open("images/shia.gif")
	shiaHeadFile, _ := os.Open("images/shiaHead.gif")

	defer shiaFile.Close()
	defer shiaFile.Close()

	shiaGif, _ := gif.DecodeAll(shiaFile)
	shiaHeadGif, _ := gif.DecodeAll(shiaHeadFile)

	finishedChan := make(chan interface{}, 26)

	for i := 0; i < 26 ; i++ {
		generateLetter(&generatorInfo{
			letter: letter_generator.Alphabet[i],
			filler: shiaGif,
			content: shiaHeadGif,
		}, finishedChan)

	}

	for i := 0; i < 26; i++ {
		select {
			case <- finishedChan:
				fmt.Print("finished")
		}
	}
}


func generateLetter(info *generatorInfo, finished chan interface{}) {
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

				if trimmed[adjustedIndex] == 72 {
					selectedGif = info.filler.Image[i]
				} else {
					selectedGif = info.content.Image[i]
				}

				draw.Draw(newImage, image.Rect(k * 128, j * 128, (k* 128) + 128, (j * 128) + 128), selectedGif, image.ZP, draw.Over)
			}
		}

		letterGif.Image = append(letterGif.Image, newImage)
	}


}
