package main

import (
	"image/gif"
	"os"
	"fmt"
	"io"
)

func main() {
	shiaFile, _ := os.Open("images/shia.gif")
	shiaHeadFile, _ := os.Open("images/parrot.gif")

	defer shiaFile.Close()
	defer shiaHeadFile.Close()
	//splitAndJoin(shiaFile, "shia", 10)
	splitAndJoin(shiaHeadFile, "parrot", 30)
}

func splitAndJoin(reader io.Reader, name string, numberFrames int) {
	g, e := gif.DecodeAll(reader)

	newGif := &gif.GIF{}
	*newGif = *g

	for i := 0; i < numberFrames/len(newGif.Image) + 1; i++ {
		newGif.Image = append(newGif.Image, g.Image...)
		newGif.Delay = append(newGif.Delay, g.Delay...)
		newGif.Disposal = append(newGif.Disposal, g.Disposal...)
	}

	f, _ := os.Create(fmt.Sprintf("%s.gif", name))
	e = gif.EncodeAll(f, newGif)
	f.Close()
	if e != nil {
		fmt.Println(e)
	}
}
