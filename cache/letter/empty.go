package letter

import (
	"image"
	"image/color"
	"image/draw"
)

type EmptyLetterFrame struct {
	*image.Paletted
}

func (empty *EmptyLetterFrame) SetBackground(rgba color.RGBA) {
	draw.Draw(empty, empty.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 0}}, image.ZP, draw.Over)
}

func (empty *EmptyLetterFrame) DrawImageAt(row, column int, img *image.Paletted) {
	bounds := getBoundsFromCoords(row, column, img.Bounds())
	draw.Draw(empty, bounds, img, image.ZP, draw.Over)
}

func getBoundsFromCoords(row, column int, rectangle image.Rectangle) image.Rectangle  {
	rectangle.Min.X = rectangle.Max.X * column
	rectangle.Max.X = rectangle.Min.X + rectangle.Max.X

	rectangle.Min.Y= rectangle.Max.Y * row
	rectangle.Max.Y = rectangle.Min.Y + rectangle.Max.Y

	return rectangle
}