package main

import (
	"image/color"
	"math/rand"
)

const (
	CONTROL_REGISTER = 0
	MASK_REGISTER    = 1
	STATUS_REGISTER  = 2
	OAM_ADDRESS      = 3
	OAM_DATA         = 4
	SCROLL_REGISTER  = 5
	ADDRESS_REGISTER = 6
	DATA_REGISTER    = 7

	VERTICAL_BLANK  = 7
	SPRITE_ZERO_HIT = 6
	SPRITE_OVERFLOW = 5

	GRAY_SCALE             = 0
	RENDER_BACKGROUND_LEFT = 1
	RENDER_SPRITES_LEFT    = 2
	RENDER_BACKGROUND      = 3
	RENDER_SPRITES         = 4
	ENHANCE_RED            = 5
	ENHANCE_GREEN          = 6
	ENHANCE_BLUE           = 7

	NAMETABLE_X       = 0
	NAMETABLE_Y       = 1
	INCREMENT_MODE    = 2
	PATTERN_SPRITE    = 3
	PATTERN_BACKGROUN = 4
	SPRITE_SIZE       = 5
	SLAVE_MODE        = 6
	ENABLE_NMI        = 7
)

var (
	DEFAULT_COLORS_2C02 = [0x40]*color.RGBA{
		&color.RGBA{84, 84, 84, 255},
		&color.RGBA{0, 30, 116, 255},
		&color.RGBA{8, 16, 144, 255},
		&color.RGBA{48, 0, 136, 255},
		&color.RGBA{68, 0, 100, 255},
		&color.RGBA{92, 0, 48, 255},
		&color.RGBA{84, 4, 0, 255},
		&color.RGBA{60, 24, 0, 255},
		&color.RGBA{32, 42, 0, 255},
		&color.RGBA{8, 58, 0, 255},
		&color.RGBA{0, 64, 0, 255},
		&color.RGBA{0, 60, 0, 255},
		&color.RGBA{0, 50, 60, 255},
		&color.RGBA{0, 0, 0, 255},
		&color.RGBA{0, 0, 0, 255},
		&color.RGBA{0, 0, 0, 255},

		&color.RGBA{152, 150, 152, 255},
		&color.RGBA{8, 76, 196, 255},
		&color.RGBA{48, 50, 236, 255},
		&color.RGBA{92, 30, 228, 255},
		&color.RGBA{136, 20, 176, 255},
		&color.RGBA{160, 20, 100, 255},
		&color.RGBA{152, 34, 32, 255},
		&color.RGBA{120, 60, 0, 255},
		&color.RGBA{84, 90, 0, 255},
		&color.RGBA{40, 114, 0, 255},
		&color.RGBA{8, 124, 0, 255},
		&color.RGBA{0, 118, 40, 255},
		&color.RGBA{0, 102, 120, 255},
		&color.RGBA{0, 0, 0, 255},
		&color.RGBA{0, 0, 0, 255},
		&color.RGBA{0, 0, 0, 255},

		&color.RGBA{236, 238, 236, 255},
		&color.RGBA{76, 154, 236, 255},
		&color.RGBA{120, 124, 236, 255},
		&color.RGBA{176, 98, 236, 255},
		&color.RGBA{228, 84, 236, 255},
		&color.RGBA{236, 88, 180, 255},
		&color.RGBA{236, 106, 100, 255},
		&color.RGBA{212, 136, 32, 255},
		&color.RGBA{160, 170, 0, 255},
		&color.RGBA{116, 196, 0, 255},
		&color.RGBA{76, 208, 32, 255},
		&color.RGBA{56, 204, 108, 255},
		&color.RGBA{56, 180, 204, 255},
		&color.RGBA{60, 60, 60, 255},
		&color.RGBA{0, 0, 0, 255},
		&color.RGBA{0, 0, 0, 255},

		&color.RGBA{236, 238, 236, 255},
		&color.RGBA{168, 204, 236, 255},
		&color.RGBA{188, 188, 236, 255},
		&color.RGBA{212, 178, 236, 255},
		&color.RGBA{236, 174, 236, 255},
		&color.RGBA{236, 174, 212, 255},
		&color.RGBA{236, 180, 176, 255},
		&color.RGBA{228, 196, 144, 255},
		&color.RGBA{204, 210, 120, 255},
		&color.RGBA{180, 222, 120, 255},
		&color.RGBA{168, 226, 144, 255},
		&color.RGBA{152, 226, 180, 255},
		&color.RGBA{160, 214, 228, 255},
		&color.RGBA{160, 162, 160, 255},
		&color.RGBA{0, 0, 0, 255},
		&color.RGBA{0, 0, 0, 255},
	}
)

type PPU2C02 struct {
	bus                *Bus
	cart               *Cartridge
	nameTable          [2][1024]byte
	paletteTable       [32]byte
	patternTable       [2][4096]byte
	paletteScreen      [64]*color.RGBA
	spriteScreen       *Sprite    // Sprite(256, 240)
	spriteNameTable    [2]*Sprite // Sprite (128, 128), Sprite(128, 128)
	spritePatternTable [2]*Sprite // Sprite(256, 240), Sprite(256, 240)
	frameComplete      bool
	scanLine           int16
	cycle              int16

	// Registers
	controlRegister byte
	maskRegister    byte
	statusRegister  byte
	// --
	// --
	scrollRegister  byte
	addressRegister Word // 14 bits
	dataRegister    byte

	addressLatch  byte
	ppuDataBuffer byte
	ppuAddress    Word
}

func (p *PPU2C02) GetFlag(flag Flag, at Register) bool {
	var shortRegister byte
	var longRegister Word

	switch at {
	case CONTROL_REGISTER:
		shortRegister = p.controlRegister
		break
	case MASK_REGISTER:
		shortRegister = p.maskRegister
		break
	case STATUS_REGISTER:
		shortRegister = p.statusRegister
		break
	case SCROLL_REGISTER:
		shortRegister = p.scrollRegister
		break
	case DATA_REGISTER:
		shortRegister = p.dataRegister
		break
	case ADDRESS_REGISTER:
		longRegister = p.addressRegister
		return (longRegister >> flag & 0x0001) == 0
	}

	return (shortRegister >> flag & 0x01) == 0
}

func (p *PPU2C02) SetFlag(flag Flag, at Register) {
	p.defineFlag(true, flag, at)
}

func (p *PPU2C02) ClearFlag(flag Flag, at Register) {
	p.defineFlag(false, flag, at)
}

func (p *PPU2C02) defineFlag(val bool, flag Flag, at Register) {

	switch at {
	case CONTROL_REGISTER:
		if val {
			p.controlRegister |= byte(1 << uint(at))
		} else {
			p.controlRegister &= ^byte(1 << uint(at))
		}
		break
	case MASK_REGISTER:
		if val {
			p.maskRegister |= byte(1 << uint(at))
		} else {
			p.maskRegister &= ^byte(1 << uint(at))
		}
		break
	case STATUS_REGISTER:
		if val {
			p.statusRegister |= byte(1 << uint(at))
		} else {
			p.statusRegister &= ^byte(1 << uint(at))
		}
		break
	case SCROLL_REGISTER:
		if val {
			p.scrollRegister |= byte(1 << uint(at))
		} else {
			p.scrollRegister &= ^byte(1 << uint(at))
		}
		break
	case DATA_REGISTER:
		if val {
			p.dataRegister |= byte(1 << uint(at))
		} else {
			p.dataRegister &= ^byte(1 << uint(at))
		}
		break
	case ADDRESS_REGISTER:
		if val {
			p.addressRegister |= Word(1 << uint(flag))
		} else {
			p.addressRegister &= ^Word(1 << uint(flag))
		}
	}
}

func (p *PPU2C02) Complete() bool {
	return p.frameComplete
}

func (p *PPU2C02) GetScreen() *Sprite {
	return p.spriteScreen
}

func (p *PPU2C02) GetNameTable(i int) *Sprite {

	return p.spriteNameTable[i]
}

func (p *PPU2C02) GetColorFromPaletteRam(palette, pixelValue byte) *color.RGBA {
	idx, _ := p.PPURead(0x3F00+Word(palette<<2+pixelValue), true)
	return p.paletteScreen[idx]
}

func (p *PPU2C02) GetPatternTable(i, palette byte) *Sprite {

	for x := 0; x < p.spritePatternTable[i].w; x++ {
		for y := 0; y < p.spritePatternTable[i].h; y++ {
			offset := y*256 + x*16
			for row := 0; row < 8; row++ {
				tileLSB, _ := p.PPURead(Word(int(i)*0x1000+offset+row+0), true)
				tileMSB, _ := p.PPURead(Word(int(i)*0x1000+offset+row+8), true)
				for col := 0; col < 8; col++ {
					pixel := (tileMSB & 0x01) + (tileLSB & 0x01)
					tileLSB = tileLSB >> 1
					tileMSB = tileMSB >> 1

					p.spritePatternTable[i].SetPixel(x*8+(7-col),
						y*8+row,
						p.GetColorFromPaletteRam(palette, pixel))

				}
			}
		}
	}

	return p.spritePatternTable[i]
}

// ConnectBus : connects the CPU to the Bus
func (p *PPU2C02) ConnectBus(bus *Bus) {
	p.bus = bus
}

func CreatePPU() *PPU2C02 {
	return &PPU2C02{
		nil,             // bus
		nil,             // cart
		[2][1024]byte{}, // nameTable
		[32]byte{},      //paletteTable
		[2][4096]byte{}, //patternTable
		DEFAULT_COLORS_2C02,
		CreateSprite(256, 240),
		[2]*Sprite{
			CreateSprite(256, 240),
			CreateSprite(256, 240)},
		[2]*Sprite{
			CreateSprite(128, 128),
			CreateSprite(128, 128)},
		false,
		0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0,
	}
}

// PPURead : PPURead
func (p *PPU2C02) PPURead(address Word, readOnly bool) (byte, error) {
	var data byte = 0
	address &= 0x3FFF

	if d, ok := p.cart.PPURead(address); ok {
		data = d
	} else if address >= 0x0000 && address <= 0x1FFF {
		data = p.patternTable[(address*0x1000)>>12][address&0x0FFF]
	} else if address >= 0x2000 && address <= 0x3EFF {

	} else if address >= 0x3F00 && address <= 0x3FFF {
		address = address & 0x001F
		if address == 0x0010 {
			address = 0x0000
		} else if address == 0x0014 {
			address = 0x0004
		} else if address == 0x0018 {
			address = 0x0008
		} else if address == 0x001C {
			address = 0x000C
		}
		data = p.paletteTable[address]
	}

	return data, nil
}

// PPUWrite : PPUWrite
func (p *PPU2C02) PPUWrite(address Word, data byte) error {
	address &= 0x3FFF

	if p.cart.CPUWrite(address, data) {

	} else if address >= 0x0000 && address <= 0x1FFF {
		p.patternTable[(address*0x1000)>>12][address&0x0FFF] = data
	} else if address >= 0x2000 && address <= 0x3EFF {

	} else if address >= 0x3F00 && address <= 0x3FFF {
		address = address & 0x001F
		if address == 0x0010 {
			address = 0x0000
		} else if address == 0x0014 {
			address = 0x0004
		} else if address == 0x0018 {
			address = 0x0008
		} else if address == 0x001C {
			address = 0x000C
		}
		p.paletteTable[address] = data
	}
	return nil
}

func (p *PPU2C02) CPURead(address Word, readOnly bool) (byte, error) {
	var data byte = 0

	switch address {
	case CONTROL_REGISTER:
		break
	case MASK_REGISTER:
		break
	case STATUS_REGISTER:
		data = p.statusRegister&0xE0 | p.ppuDataBuffer&0x1F
		p.ClearFlag(VERTICAL_BLANK, STATUS_REGISTER)
		p.addressLatch = 0
		break
	case OAM_ADDRESS:
		break
	case OAM_DATA:
		break
	case SCROLL_REGISTER:
		break
	case ADDRESS_REGISTER:
		break
	case DATA_REGISTER:
		data = p.ppuDataBuffer
		p.ppuDataBuffer, _ = p.PPURead(p.ppuAddress, false)
		if p.ppuAddress > 0x3F00 {
			data = p.ppuDataBuffer
		}
		break
	}
	return data, nil
}

func (p *PPU2C02) CPUWrite(address Word, data byte) error {

	switch address {
	case CONTROL_REGISTER:
		p.controlRegister = data
		break
	case MASK_REGISTER:
		p.maskRegister = data
		break
	case STATUS_REGISTER:
		p.statusRegister = data
		break
	case OAM_ADDRESS:
		break
	case OAM_DATA:
		break
	case SCROLL_REGISTER:
		p.scrollRegister = data
		break
	case ADDRESS_REGISTER:
		if p.addressLatch == 0 {
			p.ppuAddress = p.ppuAddress&0x00FF | Word(data)<<8
			p.addressRegister = 1
		} else {
			p.ppuAddress = p.ppuAddress&0xFF00 | Word(data)
			p.addressRegister = 0
		}
		break
	case DATA_REGISTER:
		p.PPUWrite(p.ppuAddress, data)
		break
	}

	return nil
}

// Clock : Bus clock implementation pulses the clock to all things attached to it
func (p *PPU2C02) Clock() {
	// c := byte(p.bus.cpu.opcode) % byte(len(p.paletteScreen))
	c := byte(rand.Intn(len(p.paletteScreen)))
	p.GetScreen().SetPixel(int(p.cycle)-1, int(p.scanLine), p.paletteScreen[c])
	p.cycle++
	if p.cycle >= 341 {
		p.cycle = 0
		p.scanLine++
		if p.scanLine >= 240 {
			p.SetFlag(VERTICAL_BLANK, STATUS_REGISTER)
		}
		if p.scanLine >= 261 {
			p.scanLine = -1
			p.frameComplete = true
			p.ClearFlag(VERTICAL_BLANK, STATUS_REGISTER)
		}
	}
}

// InsertCartridge : sets the pointer to the cartridge in the PPU
func (p *PPU2C02) InsertCartridge(c *Cartridge) {
	p.cart = c
}
