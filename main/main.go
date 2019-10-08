package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	height          = 480
	width           = 680
	rows            = 240
	columns         = 256
	size    float64 = 2 // Pixel size modifier
)

var (
	bus      *Bus
	cpu      *CPU6502
	basicTxt *text.Text
	frames   = 0
	second   = time.Tick(time.Second)
	atlas    = text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))
)

func init() {
	cpu = &CPU6502{}
	memory := &Memory64k{}
	bus = &Bus{cpu, memory}
	cpu.ConnectBus(bus)
}

func main() {
	cpu.SetStatusRegisterFlag(V, true)
	fmt.Println(cpu)
	fmt.Println(rand.Float32())
	pixelgl.Run(run)
}

var imd *imdraw.IMDraw

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "GoNES",
		Bounds: pixel.R(0, 0, width, height),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd = imdraw.New(nil)
	basicTxt = text.New(pixel.V(0, 0), atlas)

	// last := time.Now()
	for !win.Closed() {
		win.Clear(colornames.Blue)

		// drawPixels()
		drawCpu()

		imd.Draw(win)
		basicTxt.Draw(win, pixel.IM)
		imd.Clear()
		basicTxt.Clear()

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}

	}
}

func drawPixels() {
	var x, y float64 = size, size * rows
	i := 0
	for y > 0 {
		for x < columns*size {
			if rand.Float32() > 0.5 {
				imd.Color = colornames.Map[colornames.Names[int(x+y)%16]]
				imd.Push(pixel.V(x, y), pixel.V(x+1, y+1))
				imd.Rectangle(0)
				i++
			}
			x += size
		}
		x = 0
		y -= size
	}
}

func drawCpu() {
	c := cpu
	DrawString(10, 400, "Status", colornames.White)
	DrawString(10, 420, "N", RedGreen(c.StatusRegister(N)))
	DrawString(10+32, 420, "V", RedGreen(c.StatusRegister(V)))
	DrawString(10+32*2, 420, "U", RedGreen(c.StatusRegister(U)))
	DrawString(10+32*3, 420, "B", RedGreen(c.StatusRegister(B)))
	DrawString(10+32*4, 420, "D", RedGreen(c.StatusRegister(D)))
	DrawString(10+32*5, 420, "I", RedGreen(c.StatusRegister(I)))
	DrawString(10+32*6, 420, "Z", RedGreen(c.StatusRegister(Z)))
	DrawString(10+32*7, 420, "C", RedGreen(c.StatusRegister(C)))
}

func RedGreen(b bool) color.RGBA {
	if b {
		return colornames.Green
	} else {
		return colornames.Red
	}
}

func DrawString(x, y float64, message string, color color.RGBA) {
	// basicTxt.Clear()
	basicTxt.Dot = pixel.V(x, height-y)
	basicTxt.Color = color
	fmt.Fprintln(basicTxt, message)
}
