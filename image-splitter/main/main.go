package main

import (
	"image/gif"
	"os"
	"image"
	"image/draw"
	"fmt"
	"image/png"
	"io"
	"image/color/palette"
)

func main() {
	shiaFile, _ := os.Open("images/shia.gif")
	shiaHeadFile, _ := os.Open("images/shiaHead.gif")

	defer shiaFile.Close()
	defer shiaHeadFile.Close()
	//fmt.Println(SplitAnimatedGIF(shiaHeadFile))

	joinGif("shia", 15)
	joinGif("shiaHead", 12)
}

func joinGif(directory string, numberFrames int)  {
	outGif := &gif.GIF{}
	fmt.Println()
	for j := 0; j < (60 / numberFrames); j++ {
		for i := 0; i < numberFrames; i++ {
			f, _ := os.Open(fmt.Sprintf("%s/%d.png", directory, i))

			img, e := png.Decode(f)
			palleted := image.NewPaletted(img.Bounds(), palette.Plan9)
			draw.Draw(palleted, palleted.Bounds(), img, img.Bounds().Min, draw.Over)

			if e != nil {
				fmt.Println(e)
			}
			f.Close()
			outGif.Image = append(outGif.Image, palleted)
			outGif.Delay = append(outGif.Delay, 0)
		}
	}

	newFile, _ := os.Create(directory + ".gif")
	defer newFile.Close()

	fmt.Println(len(outGif.Image))

	gif.EncodeAll(newFile, outGif)
}

func SplitAnimatedGIF(reader io.Reader) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error while decoding: %s", r)
		}
	}()

	gif, err := gif.DecodeAll(reader)

	if err != nil {
		return err
	}

	imgWidth, imgHeight := getGifDimensions(gif)

	overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), gif.Image[0], image.ZP, draw.Src)

	for i, srcImg := range gif.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.ZP, draw.Over)

		// save current frame "stack". This will overwrite an existing file with that name
		file, err := os.Create(fmt.Sprintf("%s%d%s", "shiaHead/", i, ".png"))
		if err != nil {
			return err
		}

		err = png.Encode(file, overpaintImage)
		if err != nil {
			return err
		}

		file.Close()
	}

	return nil
}

func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}

