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
	bus      *Bus
	cpu      *CPU6502
	basicTxt *text.Text
	mapAsm   map[Word]string
	frames   = 0
	second   = time.Tick(time.Second)
	atlas    = text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))
)

func init() {
	cpu = &CPU6502{}
	memory := &Memory64k{}
	bus = &Bus{cpu, memory}
	cpu.ConnectBus(bus)

	s := "A20A8E0000A2038E0100AC0000A900186D010088D0FA8D0200EAEAEA"
	offset := Word(0x8000)
	for i := 0; i < len(s); i += 2 {
		bus.Write(offset, ByteToHex(s[i])<<4|ByteToHex(s[i+1]))
		offset++
	}
	memory.mem[0xFFFC] = 0x00
	memory.mem[0xFFFD] = 0x80
	cpu.Reset()
}

func main() {
	mapAsm = cpu.Disassemble(0x0000, 0xFFFF)
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
			mapAsm = cpu.Disassemble(0x0000, 0xFFFF)
		}
		win.Clear(colornames.Darkblue)

		// drawPixels()
		drawRam(2, 12, 0x0000, 16, 16)
		drawRam(2, 196, 0x8000, 16, 16)
		drawCpu(408, 12)
		drawCode(408, 88, 26)

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

func drawCpu(x, y float64) {
	c := cpu
	redGreen := func(b bool) color.RGBA {
		if b {
			return colornames.White
		} else {
			return colornames.Red
		}
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

func drawRam(x, y float64, addr Word, rows, columns int) {
	nRamX := x
	nRamY := y
	for row := 0; row < rows; row++ {
		var sOffset bytes.Buffer
		sOffset.WriteByte('$')
		sOffset.WriteString(Hex(uint32(addr), 4))
		sOffset.WriteByte(':')
		for col := 0; col < columns; col++ {
			v, e := bus.Read(addr, true)
			if e != nil {
				//
			}
			if cpu.pc == addr {
				sOffset.WriteByte('>')
			} else {
				sOffset.WriteByte(' ')
			}
			sOffset.WriteString(Hex(uint32(v), 2))
			addr += 1
		}
		drawString(nRamX, nRamY, sOffset.String(), colornames.White)
		nRamY += 11
	}
}

func drawCode(x, y float64, lines int) {
	////
	it_a := cpu.pc
	nLineY := (lines>>1)*11 + int(y)
	if val, ok := mapAsm[it_a]; ok {
		drawString(x, float64(nLineY), val, colornames.Cyan)
		for nLineY < ((lines * 10) + int(y)) {
			it_a++
			// if it_a != mapAsm.end() {
			if val2, ok2 := mapAsm[it_a]; ok2 {
				if len(val2) > 0 {
					nLineY += 11
					drawString(x, float64(nLineY), val2, colornames.White)

				}
			}
		}
	}

	it_a = cpu.pc

	nLineY = (lines>>1)*11 + int(y)
	if _, ok3 := mapAsm[it_a]; ok3 {
		for float64(nLineY) > y {
			it_a--
			if val4, ok4 := mapAsm[it_a]; ok4 {
				if len(val4) > 0 {
					nLineY -= 11
					drawString(x, float64(nLineY), val4, colornames.White)
				}
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
