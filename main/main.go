package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	height  = 240
	width   = 256
	rows    = 240
	columns = 256
)

var (
	bus    *Bus
	frames = 0
	second = time.Tick(time.Second)
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

type pix struct {
	x, y, h, w int
	color      *color.RGBA
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

		var x, y, size float64 = 0, 0, 1
		// colors := len(colornames.Names)
		i := 0
		win.Clear(colornames.Black)
		// max := int((x * y) / 4)
		for y < rows {
			for x < columns {
				if rand.Float32() > 0.95 {
					// imd.Color = colornames.Map[colornames.Names[int(x+y)%colors]]
					imd.Color = colornames.Map[colornames.Names[int(x+y)%16]]
					imd.Push(pixel.V(x, y), pixel.V(x+size, y+size))
					imd.Rectangle(0)
					i++
				}
				x += size
			}
			x = 0
			y += size
		}
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
		// t2 := time.Now()
		// diff := t2.Sub(t1)
		// fmt.Println("Iterated over ", i, " cells and spent ", diff, " on that frame - ", 1.0/diff.Seconds(), " FPS")
		// fps := 1.0 / diff.Seconds()
		// fmt.Println(fps, "FPS")
	}
}
