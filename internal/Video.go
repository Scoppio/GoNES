package main

import (
	"bytes"
	"fmt"
	"image/color"
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
	height     = 480
	width      = 780
	rows       = 240
	columns    = 256
	swatchSize = 6
)

var (
	Nes             *Bus
	cpu             *CPU6502
	basicTxt        *text.Text
	mapAsm          map[Word]string
	frames               = 0
	second               = time.Tick(time.Second)
	atlas                = text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))
	emulationRun         = false
	residualTime         = 0.0
	elapsedTime          = 0.0
	selectedPalette byte = 0
	imd             *imdraw.IMDraw
	win             *pixelgl.Window
)

func init() {
	Nes = CreateBus(CreateCPU(), CreatePPU())
	Nes.InsertCartridge(LoadCartridge("../test/roms/nestest.nes"))
	cpu = Nes.cpu
	cpu.ConnectBus(Nes)
	cpu.Reset()
}

func main() {
	mapAsm = cpu.Disassemble(0x0000, 0xFFFF)

	Nes.Reset()
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:       "GoNES",
		Bounds:      pixel.R(0, 0, width, height),
		Undecorated: false,
	}
	win, _ = pixelgl.NewWindow(cfg)

	imd = imdraw.New(nil)
	basicTxt = text.New(pixel.V(0, 0), atlas)

	lastUpdate := time.Now()
	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}
		win.Clear(colornames.Darkblue)
		// win.Clear(color.Black)

		if emulationRun {
			if residualTime > 0.0 {
				residualTime -= elapsedTime
			} else {
				residualTime += 1.0/60.0 - elapsedTime
				Nes.Clock()
				for !Nes.ppu.Complete() {
					Nes.Clock()
				}
				Nes.ppu.frameComplete = false

			}
		} else {
			if win.JustPressed(pixelgl.KeyC) {
				// One microcode clock
				Nes.Clock()
				for !Nes.cpu.Complete() {
					Nes.Clock()
				}
				if Nes.ppu.Complete() {
					Nes.ppu.frameComplete = false
				}
			}
			if win.JustPressed(pixelgl.KeyF) {
				// One full frame
				Nes.Clock()
				for !Nes.ppu.Complete() {
					Nes.Clock()
				}

				for !Nes.cpu.Complete() {
					Nes.Clock()
				}
				Nes.ppu.frameComplete = false
			}
		}

		if win.JustPressed(pixelgl.KeyR) {
			// Reset
			Nes.Reset()
			for !Nes.cpu.Complete() {
				Nes.Clock()
			}
		}
		if win.Pressed(pixelgl.KeyP) {
			// One microcode clock
			selectedPalette++
			selectedPalette &= 0x07
		}

		if win.JustPressed(pixelgl.KeySpace) {
			// SPACE
			emulationRun = !emulationRun
		}

		drawCPU(516, 15)
		drawCode(516, 112, 16)

		// drawRAM(2, 12, 0x0000, 16, 16)

		drawRect(float64(int(516)+int(selectedPalette)*(swatchSize*5)-1), 132, swatchSize*4+2, swatchSize+2, &colornames.White)

		for p := 0; p < 8; p++ {
			for s := 0; s < 4; s++ {
				drawRect(float64(516+p*(swatchSize*5)+s*swatchSize), 133, swatchSize, swatchSize, Nes.ppu.GetColorFromPaletteRam(byte(p), byte(s)))
			}
		}

		drawSprite(516, 2, Nes.ppu.GetPatternTable(0, selectedPalette), 1)
		drawSprite(648, 2, Nes.ppu.GetPatternTable(1, selectedPalette), 1)
		//
		drawSprite(0, 0, Nes.ppu.GetScreen(), 2)

		basicTxt.Draw(win, pixel.IM) // .Moved(c).Scaled(c, scale))
		basicTxt.Clear()

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
		elapsedTime = -lastUpdate.Sub(time.Now()).Seconds()
		lastUpdate = time.Now()
	}
}

func drawRect(x, y float64, w, h int, c *color.RGBA) {
	m := CreateSprite(w, h)
	for row := 0; row < w; row++ {
		for col := 0; col < h; col++ {
			m.SetPixel(row, col, c)
		}
	}
	drawSprite(x, y, m, 1)
}

func drawSprite(x, y float64, sprite *Sprite, scale float64) {
	p := pixel.PictureDataFromImage(Frame(sprite))
	x += (float64(sprite.w) * scale) / 2
	y += (float64(sprite.h) * scale) / 2
	pixel.NewSprite(p, p.Bounds()).
		Draw(win, pixel.IM.Moved(pixel.V(x, y)).
			Scaled(pixel.V(float64(sprite.w), float64(sprite.h)), scale))
}

func drawCPU(x, y float64) {
	c := cpu
	redGreen := func(b bool) color.RGBA {
		if b {
			return colornames.White
		}
		return colornames.Red
	}
	drawString(x, y, "CPU", colornames.White)
	drawString(x-30+64, y, "N", redGreen(c.StatusRegister(N)))
	drawString(x-30+80, y, "V", redGreen(c.StatusRegister(V)))
	drawString(x-30+96, y, "U", redGreen(c.StatusRegister(U)))
	drawString(x-30+112, y, "B", redGreen(c.StatusRegister(B)))
	drawString(x-30+128, y, "D", redGreen(c.StatusRegister(D)))
	drawString(x-30+144, y, "I", redGreen(c.StatusRegister(I)))
	drawString(x-30+160, y, "Z", redGreen(c.StatusRegister(Z)))
	drawString(x-30+178, y, "C", redGreen(c.StatusRegister(C)))
	drawString(x, y+12, fmt.Sprintln("PC: ", fmt.Sprintf("$%s [%d]", Hex(uint32(c.pc), 4), c.pc)), colornames.White)
	drawString(x, y+24, fmt.Sprintln("A : ", fmt.Sprintf("$%s   [%d]", Hex(uint32(c.a), 2), c.a)), colornames.White)
	drawString(x, y+36, fmt.Sprintln("X : ", fmt.Sprintf("$%s   [%d]", Hex(uint32(c.x), 2), c.x)), colornames.White)
	drawString(x, y+48, fmt.Sprintln("Y : ", fmt.Sprintf("$%s   [%d]", Hex(uint32(c.y), 2), c.y)), colornames.White)
	drawString(x, y+60, fmt.Sprintln("Stack P: ", fmt.Sprintf("$%s", Hex(uint32(c.stkp), 4))), colornames.White)
	drawString(x, y+72, fmt.Sprintln("Clock Count: ", ClockCount), colornames.White)
	drawString(x, y+84, fmt.Sprintln("Operation Count: ", OperationCount), colornames.White)
	// drawString(x, y, fmt.Sprintln("Clock: ", c.cycles), colornames.White)
	// drawString(x, y, fmt.Sprintln("GlobalClock: ", clock_count), colornames.White)
	// drawString(x, y, fmt.Sprintln("ADD ABS: ", fmt.Sprintf("0x%X", c.address_abs)), colornames.White)
	// drawString(x, y, fmt.Sprintln("ADD REL: ", fmt.Sprintf("0x%X", c.address_rel)), colornames.White)
}

func drawRAM(x, y float64, addr Word, rows, columns int) {
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
