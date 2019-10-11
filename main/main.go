package main

import (
	"bytes"
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
	Nes      *Bus
	cpu      *CPU6502
	basicTxt *text.Text
	mapAsm   map[rune]string
	frames   = 0
	second   = time.Tick(time.Second)
	atlas    = text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))
)

func init() {
	Nes = CreateBus(CreateCPU(), CreatePPU())
	Nes.InsertCartridge(LoadCartridge("../testroms/nestest.nes"))
	cpu = Nes.cpu
	cpu.ConnectBus(Nes)
	cpu.Reset()
}

func main() {
	mapAsm = cpu.Disassemble(0x0000, 0x07FF)
	Nes.Reset()
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
			for !cpu.Complete() {
				cpu.Clock()
			}
		}
		if win.JustPressed(pixelgl.KeyC) {
			// I
			cpu.Clock()
			for !cpu.Complete() {
				cpu.Clock()
			}
		}

		win.Clear(colornames.Darkblue)

		drawCPU(516, 2)
		drawCode(516, 72, 26)

		drawRAM(2, 12, 0x0000, 16, 16)

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

func drawCPU(x, y float64) {
	c := cpu
	redGreen := func(b bool) color.RGBA {
		if b {
			return colornames.White
		}
		return colornames.Red
	}
	drawString(x, y, "CPU State", colornames.White)
	drawString(x+15+64, y, "N", redGreen(c.StatusRegister(N)))
	drawString(x+15+80, y, "V", redGreen(c.StatusRegister(V)))
	drawString(x+15+96, y, "U", redGreen(c.StatusRegister(U)))
	drawString(x+15+112, y, "B", redGreen(c.StatusRegister(B)))
	drawString(x+15+128, y, "D", redGreen(c.StatusRegister(D)))
	drawString(x+15+144, y, "I", redGreen(c.StatusRegister(I)))
	drawString(x+15+160, y, "Z", redGreen(c.StatusRegister(Z)))
	drawString(x+15+178, y, "C", redGreen(c.StatusRegister(C)))
	drawString(x, y+12, fmt.Sprintln("PC: ", fmt.Sprintf("$%s [%d]", Hex(uint32(c.pc), 4), c.pc)), colornames.White)
	drawString(x, y+24, fmt.Sprintln("A : ", fmt.Sprintf("$%s   [%d]", Hex(uint32(c.a), 2), c.a)), colornames.White)
	drawString(x, y+36, fmt.Sprintln("X : ", fmt.Sprintf("$%s   [%d]", Hex(uint32(c.x), 2), c.x)), colornames.White)
	drawString(x, y+48, fmt.Sprintln("Y : ", fmt.Sprintf("$%s   [%d]", Hex(uint32(c.y), 2), c.y)), colornames.White)
	drawString(x, y+60, fmt.Sprintln("Stack P: ", fmt.Sprintf("$%s", Hex(uint32(c.stkp), 4))), colornames.White)
	// drawString(x, y, fmt.Sprintln("Clock: ", c.cycles), colornames.White)
	// drawString(x, y, fmt.Sprintln("GlobalClock: ", clock_count), colornames.White)
	// drawString(x, y, fmt.Sprintln("ADD ABS: ", fmt.Sprintf("0x%X", c.address_abs)), colornames.White)
	// drawString(x, y, fmt.Sprintln("ADD REL: ", fmt.Sprintf("0x%X", c.address_rel)), colornames.White)
}

func drawRAM(x, y float64, addr rune, rows, columns int) {
	RAMX := x
	RAMY := y
	for row := 0; row < rows; row++ {
		var sOffset bytes.Buffer
		sOffset.WriteByte('$')
		sOffset.WriteString(Hex(uint32(addr), 4))
		sOffset.WriteByte(':')
		for col := 0; col < columns; col++ {
			v, e := Nes.CPURead(addr, true)
			if e != nil {
				//
			}
			if cpu.pc == addr {
				sOffset.WriteByte('>')
			} else {
				sOffset.WriteByte(' ')
			}
			sOffset.WriteString(Hex(uint32(v), 2))
			addr++
		}
		drawString(RAMX, RAMY, sOffset.String(), colornames.White)
		RAMY += 11
	}
}

func drawCode(x, y float64, lines int) {
	//
	pc := cpu.pc
	yPos := float64(lines>>1)*11 + y
	if ida, ok := mapAsm[pc]; ok {
		drawString(x, yPos, ida, colornames.Cyan)
		for yPos < float64(lines)*10+y {
			pc++
			if ida, ok = mapAsm[pc]; ok {
				if len(ida) > 0 {
					yPos += 11
					drawString(x, yPos, ida, colornames.White)
				}
			}
		}
	}
	pc = cpu.pc
	yPos = float64(lines>>1)*11 + y
	if _, ok := mapAsm[pc]; ok {
		for yPos > y {
			pc--
			if adi, ok := mapAsm[pc]; ok {
				yPos -= 11
				drawString(x, yPos, adi, colornames.White)
			}
		}
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
