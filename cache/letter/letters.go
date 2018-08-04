package letter

import (
	"strings"
	"image/gif"
	"image"
	"image/color/palette"
	"image/color"
)

const (
	FILLER  = '-'
	CONTENT = '0'
)

const (
	COLUMNS = 3
	ROWS    = 5
)


type Filler interface {
	GetContent(frame int) *image.Paletted
	GetFiller(frame int) *image.Paletted
}


type Letter string

func (letter Letter) Split() []string  {
	return strings.Split(string(letter), "\n")
}

func (letter Letter) generateGif(filler Filler) *gif.GIF  {
	letterGif := &gif.GIF{}

	for i := 0; i < 10; i++ {
		letterGif.Image = append(letterGif.Image, image.NewPaletted(image.Rect(0, 0, 60, 100), palette.Plan9))
		letterGif.Delay = append(letterGif.Delay, 0)
	}

	trimmed := letter.Split()


	for i := 0; i < 10; i++ {
		emptyImage := &EmptyLetterFrame{Paletted: letterGif.Image[i]}
		emptyImage.SetBackground(color.RGBA{255, 255, 255, 0})
		for j := 0; j < (COLUMNS * ROWS); j++ {
			row := j % ROWS
			column := j % COLUMNS

			var selectedGif *image.Paletted

			if trimmed[row][column + 1] == FILLER{
				selectedGif = filler.GetFiller(i)
			} else if trimmed[row][column + 1] == CONTENT {
				selectedGif = filler.GetFiller(i)
			}


			emptyImage.DrawImageAt(row, column, selectedGif)
		}
	}

	return letterGif
}


var Alphabet = [26]Letter {
`--0--
-0-0-
-000-
-0-0-
-0-0-`,

`-00--
-0-0-
-00--
-0-0-
-00--`,

`--00-
-0---
-0---
-0---
--00-`,

`-00--
-0-0-
-0-0-
-0-0-
-00--`,

`-000-
-0---
-00--
-0---
-000-`,

`-000-
-0---
-00--
-0---
-0---`,

`--00-
-0---
-0-0-
-0-0-
--00-`,

`-0-0-
-0-0-
-000-
-0-0-
-0-0-`,

`-000-
--0--
--0--
--0--
-000-`,

`--00-
---0-
---0-
-0-0-
--0--`,

`-0-0-
-0-0-
-00--
-0-0-
-0-0-`,

`-0---
-0---
-0---
-0---
-000-`,

`-0-0-
-000-
-000-
-0-0-
-0-0-`,

`-000-
-0-0-
-0-0-
-0-0-
-0-0-`,

`--0--
-0-0-
-0-0-
-0-0-
--0--`,

`-00--
-0-0-
-00--
-0---
-0---`,

`--0--
-0-0-
-0-0-
--00-
---0-`,

`-00--
-0-0-
-00--
-0-0-
-0-0-`,

`--00-
-0---
--0--
---0-
-00--`,

`-000-
--0--
--0--
--0--
--0--`,

`-0-0-
-0-0-
-0-0-
-0-0-
-000-`,

`-0-0-
-0-0-
-0-0-
-0-0-
--0--`,

`-0-0-
-0-0-
-000-
-000-
-0-0-`,

`-0-0-
-0-0-
--0--
-0-0-
-0-0-`,

`-0-0-
-0-0-
--0--
--0--
--0--`,

`-000-
---0-
--0--
-0---
-000-`,
}


