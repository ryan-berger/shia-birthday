package main

import (
	"strings"
	"image"
	"image/draw"
	"os"
	"image/gif"
	"image/png"
	"fmt"
	"unicode"
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

	generateLetter(&generatorInfo{
		letter: A,
		filler: shiaGif,
		content: shiaHeadGif,
	})
}


func generateLetter(info *generatorInfo) {
	trimmed := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, info.letter)
	newImage := image.NewRGBA(image.Rect(0, 0, 640, 640))

	fmt.Println(trimmed)

	fmt.Println(len(info.filler.Image))
	fmt.Println(len(info.content.Image))

	for i := 0; i < 60; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				adjustedIndex := (j * 5) + k
				var selectedGif *image.Paletted

				fmt.Println(string(rune(trimmed[adjustedIndex])))

				if trimmed[adjustedIndex] == 72 {
					selectedGif = info.filler.Image[i]
				} else {
					selectedGif = info.content.Image[i]
				}

				draw.Draw(newImage, image.Rect(k * 128, j * 128, (k* 128) + 128, (j * 128) + 128), selectedGif, image.ZP, draw.Over)
			}
		}
		f, e := os.Create(fmt.Sprintf("a/%d.png", i))
		if e != nil {
			fmt.Println(e)
		}
		er := png.Encode(f, newImage)
		if er != nil {
			fmt.Println(er)
		}
	}


}
