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


func (s *Sprite) GetPixel(x, y int) *color.RGBA, bool {
	if y > -1 && x > -1 {
		y = y % s.h
		x = x % s.w
		return s.pixelMatrix[x][y], true
	}

	return nil, false
}


func loadTextures() *image.RGBA {
	p, err := png.Decode(bytes.NewReader(textureData))
	if err != nil {
		panic(err)
	}

	m := image.NewRGBA(p.Bounds())

	draw.Draw(m, m.Bounds(), p, image.ZP, draw.Src)

	return m
}

var textureData = []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13}

func Frame(X1POS, Y1POS float64, sprite *Sprite) *image.RGBA {
	m := image.NewRGBA(image.Rect(X1POS, Y1POS, sprite.w, sprite.h))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if c, ok := sprite.GetPixel(x, y); ok {
				m.Set(x, y, c)
			}
		}
	}

	return m
}