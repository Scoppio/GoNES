package main

import (
	"image/color"
)

type Sprite struct {
	x, y, h, w  int
	pixelMatrix [][]*color.RGBA
}

func CreateSprite(w, h int) *Sprite {
	pixelMatrix := make([][]*color.RGBA, w)
	for i := range pixelMatrix {
		pixelMatrix[i] = make([]*color.RGBA, h)
		for k := range pixelMatrix[i] {
			pixelMatrix[i][k] = nil
		}
	}

	sprite := &Sprite{0, 0, h, w, pixelMatrix}
	return sprite
}

func (s *Sprite) SetPixel(x, y int, color *color.RGBA) {
	if y > -1 && x > -1 {
		y = y % s.h
		x = x % s.w
		s.pixelMatrix[x][y] = color
	}
}
