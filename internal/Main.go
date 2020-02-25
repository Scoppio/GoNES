package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/jroimartin/gocui"
	"golang.org/x/image/colornames"
)

const (
	height     = 480
	width      = 780
	rows       = 240
	columns    = 256
	swatchSize = 6
)

var (
	Nes          *Bus
	cpu          *CPU6502
	mapAsm       map[Word]string
	frames       = 0
	second       = time.Tick(time.Second)
	emulationRun = false
	residualTime = 0.0
	elapsedTime  = 0.0
	ROM_NAME     = "nestest"
)

func init() {
	Nes = CreateBus(CreateCPU(), CreatePPU())
	Nes.InsertCartridge(LoadCartridge("../test/roms/" + ROM_NAME + ".nes"))
	cpu = Nes.cpu
	Nes.Reset()
	mapAsm = cpu.Disassemble(0x0000, 0xFFFF)
	WriteDisassemble(mapAsm, "../output/"+ROM_NAME+".txt")
	Nes.Reset()
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	g.Mouse = true
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, tickEmulator); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, resetEmulator); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("views", 0, 0, 12, 7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "6502 Asm\nC code\nstack\nram\nc")
	}

	if v, err := g.SetView("ccode", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		s := []string{"void _init(void)",
			"{",
			"	ppu_off();			 // screen off",
			"	pal_bg(palette_bg);  //	load the BG palette",
			"	pal_spr(palette_sp); // load the Sprite palette",
			"	ppu_on_all();		 //	turn on screen",
			"	set_vram_buffer();	 // PPU pointed to VRAM Buffer",
			"	bank_spr(0);		 // set bank for sprite",
			"	reset_game();		 // reset game variables",
			"}",
			"",
			"void reset_game(void)",
			"{",
			"	cursor.card = NULL;",
			"	cursor.cell = 0;",
			"	(*table_ptr)[0] = NULL;",
			"	(*table_ptr)[1] = NULL;",
			"	(*table_ptr)[2] = NULL;",
			"	(*table_ptr)[3] = NULL;",
			"	red_size_pt = 0;",
			"	yellow_size_pt = 0;",
			"	green_size_pt = 0;",
			"	blue_size_pt = 0;",
			"	red_bc_count = 3;",
			"	yellow_bc_count = 3;",
			"	green_bc_count = 3;",
			"	blue_bc_count = 3;",
			"	round = 13;",
			"	round_score = 0;",
			"	pp = 0;",
			"	second_forever = 0;",
			"	map_registers = NULL;",
			"	challenge = 0;",
			"	kicked_client = FALSE;",
			"	",
			"	shuffle_decks();",
			"}"}

		t := strings.Join(s, "\n")

		fmt.Fprintln(v, t)
	}

	if v, err := g.SetView("assembly", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		s := []string{";",
			"; File generated by cc65 v 2.18 - Git 0f08ae2",
			";",
			"	.fopt		compiler,\"cc65 v 2.18 - Git 0f08ae2\"",
			"	.setcpu		\"6502\"",
			"	.smart		on",
			"	.autoimport	on",
			"	.case		on",
			"	.debuginfo	off",
			"	.importzp	sp, sreg, regsave, regbank",
			"	.importzp	tmp1, tmp2, tmp3, tmp4, ptr1, ptr2, ptr3, ptr4",
			"	.macpack	longbranch",
			"	.forceimport	__STARTUP__",
			"	.import		_pal_bg",
			"	.import		_pal_spr",
			"	.import		_ppu_wait_nmi",
			"	.import		_ppu_off",
			"	.import		_ppu_on_all",
			"	.import		_oam_clear",
			"	.import		_oam_spr",
			"	.import		_oam_meta_spr",
			"	.import		_pad_poll",
			"	.import		_bank_spr",
			"	.import		_rand8",
			"	.import		_set_rand",
			"	.import		_vram_adr",
			"	.import		_vram_put",
			"	.import		_set_vram_buffer",
			"	.import		_multi_vram_buffer_horz",
			"	.import		_clear_vram_buffer",
			"	.import		_get_pad_new",
			"	.import		_set_scroll_y",
			"	.import		_get_ppu_addr",
			"	.import		_get_at_addr",
			"	.import		_set_data_pointer",
			"	.import		_set_mt_pointer",
			"	.import		_buffer_4_mt",
			"	.import		_flush_vram_update_nmi",
			"	.export		_play_sound",
			"	.export		_i"}
		t := strings.Join(s, "\n")
		fmt.Fprintln(v, t)

	}

	if v, err := g.SetView("ram", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		s := []string{
			"00000000: 55 7A 6E 61 11 00 00 00 60 00 00 00 0B 00 00 00",
			"00000010: 6B 00 00 00 10 00 00 00 7B 00 00 00 EF D7 00 00",
			"00000020: 6A D8 00 00 78 08 00 00 E2 E0 00 00 72 12 00 00",
			"00000030: 54 F3 00 00 02 00 00 00 57 F3 00 00 2E 94 00 00",
			"00000040: 85 87 01 00 E6 0E 00 00 6B 96 01 00 01 00 00 00",
			"00000050: 56 F3 00 00 01 00 00 00 6C 96 01 00 01 00 00 00",
			"00000060: 03 02 02 40 8B 85 CB F2 05 03 03 02 01 8B C2 52",
			"00000070: 5E E9 BC 07 10 05 D8 28 5E C8 14 06 D0 8D 00 00",
			"00000080: 8A 03 00 9A 32 01 02 DB 25 00 01 20 01 CA 2B 0A",
			"00000090: 82 4F 01 CA 2B 00 01 4C 01 CB 1F 0A 82 4E 01 CB"}
		t := strings.Join(s, "\n")
		fmt.Fprintln(v, t)

	}

	if v, err := g.SetView("registers", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		s := []string{"a: 0x00",
			"tick-count: 0",
			"pc: 0x00",
			"x: 0x00",
			"y: 0x00",
			"stkptr: 0x00",
			"N: 0",
			"V: 0 ",
			"U: 0",
			"B: 0",
			"I: 0",
			"Z: 0",
			"C: 0"}
		t := strings.Join(s, "\n")
		fmt.Fprintln(v, t)
	}

	if v, err := g.SetView("stack", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		s := []string{"0x25",
			"0x11",
			"0x13",
			"0x01"}
		t := strings.Join(s, "\n")
		fmt.Fprintln(v, t)
	}

	return nil
}

func tickEmulator() {
	Nes.ExecuteOperation()
}

func resetEmulator() {
	Nes.Reset()
	for !Nes.cpu.Complete() {
		Nes.Clock()
	}
}

func memoryPointer(g *gocui.Gui) error {

	drawCPU(516, 15)
	drawCode(516, 112, 16)

	// drawRAM(2, 12, 0x0000, 16, 16)
	// draw palette selected
	drawRect(float64(int(516)+int(selectedPalette)*(swatchSize*5)-1), 132, swatchSize*4+2, swatchSize+2, &colornames.White)

	for p := 0; p < 8; p++ {
		for s := 0; s < 4; s++ {
			drawRect(float64(516+p*(swatchSize*5)+s*swatchSize), 133, swatchSize, swatchSize, Nes.ppu.GetColorFromPaletteRam(byte(p), byte(s)))
		}
	}

	elapsedTime = -lastUpdate.Sub(time.Now()).Seconds()
	lastUpdate = time.Now()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func run() {
	lastUpdate := time.Now()
	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}
		win.Clear(colornames.Darkblue)
		// win.Clear(color.Black)

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
	drawString(x, y, fmt.Sprintln("Clock: ", c.cycles), colornames.White)
	drawString(x, y, fmt.Sprintln("GlobalClock: ", clock_count), colornames.White)
	drawString(x, y, fmt.Sprintln("ADD ABS: ", fmt.Sprintf("0x%X", c.address_abs)), colornames.White)
	drawString(x, y, fmt.Sprintln("ADD REL: ", fmt.Sprintf("0x%X", c.address_rel)), colornames.White)
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
	// basicTxt.Dot = pixel.V(x, height-y)
	// basicTxt.Color = color
	fmt.Fprintln(basicTxt, message)
}