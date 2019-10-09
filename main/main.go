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

	offset := Word(0x8000)
	s := "A20A8E0000A2038E0100AC0000A900186D010088D0FA8D0200EAEAEA"
	for i := 0; i < len(s); i++ {
		bus.Write(offset+Word(i), s[i])
	}
	memory.mem[0xFFFC] = 0x00
	memory.mem[0xFFFD] = 0x80
	cpu.Reset()
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
		if win.JustPressed(pixelgl.KeyR) {
			// Reset
			cpu.Reset()
		}
		if win.JustPressed(pixelgl.KeyI) {
			// I
			cpu.SetStatusRegisterFlag(I, !cpu.StatusRegister(I))
		}
		if win.JustPressed(pixelgl.KeyN) {
			// N
			cpu.SetStatusRegisterFlag(N, !cpu.StatusRegister(N))
		}
		if win.JustPressed(pixelgl.KeyV) {
			// V
			cpu.SetStatusRegisterFlag(V, !cpu.StatusRegister(V))
		}
		if win.JustPressed(pixelgl.KeyC) {
			// C
			cpu.SetStatusRegisterFlag(C, !cpu.StatusRegister(C))
		}
		if win.JustPressed(pixelgl.KeyZ) {
			// Z
			cpu.SetStatusRegisterFlag(Z, !cpu.StatusRegister(Z))
		}
		if win.JustPressed(pixelgl.KeyB) {
			// B
			cpu.SetStatusRegisterFlag(B, !cpu.StatusRegister(B))
		}
		if win.JustPressed(pixelgl.KeyU) {
			// U
			cpu.SetStatusRegisterFlag(U, !cpu.StatusRegister(U))
		}
		if win.JustPressed(pixelgl.KeySpace) {
			// SPACE
			cpu.Clock()
			for !cpu.Complete() {
				cpu.Clock()
			}
		}
		win.Clear(colornames.Darkblue)

		// drawPixels()
		drawCpu()
		drawRam(2, 3, 0x0000, 28)
		drawRam(322, 3, 0x8000, 28)

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
	redGreen := func(b bool) color.RGBA {
		if b {
			return colornames.White
		} else {
			return colornames.Red
		}
	}
	drawString(10, 414, "CPU State", colornames.White)
	drawString(10, 428, fmt.Sprintln("A: ", fmt.Sprintf("0x%X", c.a)), colornames.White)
	drawString(10+32*3, 428, fmt.Sprintln("X: ", fmt.Sprintf("0x%X", c.x)), colornames.White)
	drawString(10+32*6, 428, fmt.Sprintln("Y: ", fmt.Sprintf("0x%X", c.y)), colornames.White)
	drawString(10+32*9, 428, fmt.Sprintln("ADD ABS: ", fmt.Sprintf("0x%X", c.address_abs)), colornames.White)
	drawString(10+32*15, 428, fmt.Sprintln("ADD REL: ", fmt.Sprintf("0x%X", c.address_rel)), colornames.White)
	drawString(10, 442, "N", redGreen(c.StatusRegister(N)))
	drawString(10+32, 442, "V", redGreen(c.StatusRegister(V)))
	drawString(10+32*2, 442, "U", redGreen(c.StatusRegister(U)))
	drawString(10+32*3, 442, "B", redGreen(c.StatusRegister(B)))
	drawString(10+32*4, 442, "D", redGreen(c.StatusRegister(D)))
	drawString(10+32*5, 442, "I", redGreen(c.StatusRegister(I)))
	drawString(10+32*6, 442, "Z", redGreen(c.StatusRegister(Z)))
	drawString(10+32*7, 442, "C", redGreen(c.StatusRegister(C)))
	drawString(10, 456, fmt.Sprintln("Clock: ", c.cycles), colornames.White)
	drawString(10+32*3, 456, fmt.Sprintln("PC: ", fmt.Sprintf("0x%X", c.pc)), colornames.White)
	drawString(10, 472, fmt.Sprintln("GlobalClock: ", clock_count), colornames.White)
}

func drawRam(x, y float64, memOffset Word, rows int) {
	c := cpu
	row := rows
	text := c.Disassemble(memOffset, memOffset+Word(rows*2))
	i := 0
	for row > 0 {
		l := len(text[memOffset+Word(i)])
		line := 0
		if l > 0 {
			line++
			y = y + float64(14*line)
			color := colornames.White
			if c.pc == memOffset+Word(i) {
				color = colornames.Greenyellow
			}
			drawString(x, y, text[memOffset+Word(i)], color)
			row--
		}
		i++
	}
}

func drawString(x, y float64, message string, color color.RGBA) {
	basicTxt.Dot = pixel.V(x, height-y)
	basicTxt.Color = color
	fmt.Fprintln(basicTxt, message)
}

func drawStringNL(message string, color color.RGBA) {
	basicTxt.Color = color
	fmt.Fprintln(basicTxt, message)
}
