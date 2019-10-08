package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	height  = 480
	width   = 680
	rows    = 240
	columns = 256
)

var (
	bus      *Bus
	basicTxt text
	frames   = 0
	second   = time.Tick(time.Second)
	atlas    = text.NewAtlas(basicfont.Face7x13, text.ASCII)
)

func init() {
	cpu := &CPU6502{}
	memory := &Memory64k{}
	bus = &Bus{cpu, memory}
	cpu.ConnectBus(bus)
}

func main() {
	cpu := bus.cpu
	cpu.SetStatusRegisterFlag(V, true)
	fmt.Println(cpu)
	fmt.Println(rand.Float32())
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "GoNES",
		Bounds: pixel.R(0, 0, width, height),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	// last := time.Now()
	for !win.Closed() {
		var x, y, size float64 = 0, 0, 2
		i := 0
		// if frames%2 == 0 {
		win.Clear(colornames.Blue)
		// }
		for y < rows*size {
			for x < columns*size {
				if rand.Float32() > 0.9 {
					// imd.Color = colornames.Map[colornames.Names[int(x+y)%16]]
					// imd.Push(pixel.V(x, y), pixel.V(x+1, y+1))
					// imd.Rectangle(0)
					i++
				}
				x += size
			}
			x = 0
			y += size
		}

		// basicTxt.Draw(win, pixel.IM)
		imd.Draw(win)
		imd.Clear()
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

func drawCpu(c *CPU6502) {
	DrawString(0.0, 0.0, "Status")
	DrawString(0.0+64, 0.0, fmt.Sprintln("N: ", c.StatusRegisterAsWord(N)))
	DrawString(0.0+64*2, 0.0, fmt.Sprintln("V: ", c.StatusRegisterAsWord(V)))
	DrawString(0.0+64*4, 0.0, fmt.Sprintln("U: ", c.StatusRegisterAsWord(U)))
	DrawString(0.0+64*5, 0.0, fmt.Sprintln("B: ", c.StatusRegisterAsWord(B)))
	DrawString(0.0+64*6, 0.0, fmt.Sprintln("D: ", c.StatusRegisterAsWord(D)))
	DrawString(0.0+64*7, 0.0, fmt.Sprintln("I: ", c.StatusRegisterAsWord(I)))
	DrawString(0.0+64*8, 0.0, fmt.Sprintln("Z: ", c.StatusRegisterAsWord(Z)))
	DrawString(0.0+64*9, 0.0, fmt.Sprintln("C: ", c.StatusRegisterAsWord(C)))
}

func DrawString(x, y float64, message string) {
	basicTxt = text.New(pixel.V(x, y), atlas)
	// if color == nil {
	// 	color = colornames.White
	// }
	basicTxt.Color = colornames.White
	fmt.Fprintln(basicTxt, message)
}
