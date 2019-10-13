package main

import (
	"image/color"
)

type Sprite struct {
	x, y, h, w  int
	pixelMatrix [][]*color.RGBA
}

func CreateSprite(h, w int) *Sprite {

	pixelMatrix := make([][]*color.RGBA, h)
	for i := range pixelMatrix {
		pixelMatrix[i] = make([]*color.RGBA, w)
		for k := range pixelMatrix[i] {
			pixelMatrix[i][k] = nil
		}
	}

	sprite := &Sprite{0, 0, h, w, pixelMatrix}
	return sprite
}

func (s *Sprite) SetPixel(x, y int, color *color.RGBA) {
	xx := x % len(s.pixelMatrix)
	yy := y % len(s.pixelMatrix[1])
	s.pixelMatrix[xx][yy] = color
}
