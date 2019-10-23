package main

import (
	"image/color"
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
	COARSE_X          = 5
)

var (
	DEFAULT_COLORS_2C02 [64]*color.RGBA
)

type loopyRegister struct {
	coarseX    byte
	coarseY    byte
	nametableX byte
	nametableY byte
	fineY      byte
}

func (l *loopyRegister) getAddress() Word {
	return Word(l.coarseX) | Word(l.coarseY)<<5 | Word(l.nametableX)<<10 | Word(l.nametableY)<<11 | Word(l.fineY)<<12
}

func (l *loopyRegister) increment() {
	l.add(1)
}

func (l *loopyRegister) add(value Word) {
	tbd := l.getAddress() + value
	l.coarseX = byte(tbd & 0x001F)
	l.coarseY = byte((tbd >> 5) & 0x001F)
	l.nametableX = byte(tbd >> 10 & 0x0001)
	l.nametableY = byte(tbd >> 11 & 0x0001)
	l.fineY = byte(tbd >> 12 & 0x0007)
}

func (l *loopyRegister) set(value Word) {
	tbd := value
	l.coarseX = byte(tbd & 0x001F)
	l.coarseY = byte((tbd >> 5) & 0x001F)
	l.nametableX = byte(tbd >> 10 & 0x0001)
	l.nametableY = byte(tbd >> 11 & 0x0001)
	l.fineY = byte(tbd >> 12 & 0x0007)
}

// CreateLoopyRegister : Create a cleared loopy register
func CreateLoopyRegister() *loopyRegister {
	return &loopyRegister{0, 0, 0, 0, 0}
}

// PPU2C02 : PPU
type PPU2C02 struct {
	bus                *Bus
	cart               *Cartridge
	nameTable          [2][1024]byte // VRAM
	paletteTable       [32]byte
	patternTable       [2][4096]byte // Pattern Memory
	paletteScreen      [64]*color.RGBA
	spriteScreen       *Sprite    // Sprite(256, 240)
	spriteNameTable    [2]*Sprite // Sprite(128, 128), Sprite(128, 128)
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
	addressLatch    byte
	ppuDataBuffer   byte

	NonMaskableInterrupt bool

	vRam  *loopyRegister
	tRam  *loopyRegister
	fineX byte

	bgNextTileId     byte
	bgNextTileAttrib byte
	bgNextTileLsb    byte
	bgNextTileMsb    byte

	bgShifterPatternLo Word
	bgShifterPatternHi Word
	bgShifterAttribLo  Word
	bgShifterAttribHi  Word
}

func init() {
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
}

// CreatePPU : CreatePPU cleared
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
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		false,

		CreateLoopyRegister(),
		CreateLoopyRegister(),
		0,

		0, 0, 0, 0,

		0, 0, 0, 0}
}

// GetFlag : GetFlag
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
		return (longRegister >> Word(flag) & Word(0x0001)) == 0
	}

	return (shortRegister >> byte(flag) & byte(0x01)) == 0
}

func (p *PPU2C02) GetFlagByte(flag Flag, at Register) byte {
	if p.GetFlag(flag, at) {
		return byte(0x01)
	}
	return byte(0x00)
}

// SetFlag : SetFlag
func (p *PPU2C02) SetFlag(flag Flag, at Register) {
	p.defineFlag(true, flag, at)
}

// ClearFlag : ClearFlag
func (p *PPU2C02) ClearFlag(flag Flag, at Register) {
	p.defineFlag(false, flag, at)
}

func (p *PPU2C02) defineFlag(val bool, flag Flag, at Register) {

	switch at {
	case CONTROL_REGISTER:
		if val {
			p.controlRegister |= byte(1 << uint(flag))
		} else {
			p.controlRegister &= ^byte(1 << uint(flag))
		}
		break
	case MASK_REGISTER:
		if val {
			p.maskRegister |= byte(1 << uint(flag))
		} else {
			p.maskRegister &= ^byte(1 << uint(flag))
		}
		break
	case STATUS_REGISTER:
		if val {
			p.statusRegister |= byte(1 << uint(flag))
		} else {
			p.statusRegister &= ^byte(1 << uint(flag))
		}
		break
	case SCROLL_REGISTER:
		if val {
			p.scrollRegister |= byte(1 << uint(flag))
		} else {
			p.scrollRegister &= ^byte(1 << uint(flag))
		}
		break
	case DATA_REGISTER:
		if val {
			p.dataRegister |= byte(1 << uint(flag))
		} else {
			p.dataRegister &= ^byte(1 << uint(flag))
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

// Complete : Complete
func (p *PPU2C02) Complete() bool {
	return p.frameComplete
}

// GetScreen : GetScreen
func (p *PPU2C02) GetScreen() *Sprite {
	return p.spriteScreen
}

// GetNameTable : GetNameTable
func (p *PPU2C02) GetNameTable(i int) *Sprite {
	return p.spriteNameTable[i]
}

// GetColorFromPaletteRam : GetColorFromPaletteRam
func (p *PPU2C02) GetColorFromPaletteRam(palette, pixelValue byte) *color.RGBA {
	idx, _ := p.PPURead(0x3F00+Word(palette)<<2+Word(pixelValue), false)
	idx &= 0x3F
	c := p.paletteScreen[idx]
	return c
}

// GetPatternTable : GetPatternTable
func (p *PPU2C02) GetPatternTable(i, palette byte) *Sprite {
	var tileY, tileX, row, col, offset int
	for tileY = 0; tileY < 16; tileY++ {
		for tileX = 0; tileX < 16; tileX++ {
			offset = tileY*256 + tileX*16
			for row = 0; row < 8; row++ {
				pos := Word(i)*Word(0x1000) + Word(offset+row)
				tileLSB, _ := p.PPURead(pos+0, true)
				tileMSB, _ := p.PPURead(pos+8, true)
				for col = 0; col < 8; col++ {
					pixel := (tileMSB & 0x01) + (tileLSB & 0x01)
					tileLSB = tileLSB >> 1
					tileMSB = tileMSB >> 1

					p.spritePatternTable[i].SetPixel(
						tileX*8+(7-col),
						tileY*8+row,
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

// PPURead : PPURead
func (p *PPU2C02) PPURead(address Word, readOnly bool) (byte, error) {
	var data byte = 0
	address &= 0x3FFF

	if d, ok := p.cart.PPURead(address); ok {
		data = d
	} else if address >= 0x0000 && address <= 0x1FFF {
		data = p.patternTable[(address&0x1000)>>12][address&0x0FFF]
	} else if address >= 0x2000 && address <= 0x3EFF {
		if p.cart.Mirror == VERTICAL {
			if address >= 0x0000 && address <= 0x03FF {
				p.nameTable[0][address&0x03FF] = data
			} else if address >= 0x0400 && address <= 0x03FF {
				p.nameTable[1][address&0x03FF] = data
			} else if address >= 0x0800 && address <= 0x0BFF {
				p.nameTable[0][address&0x03FF] = data
			} else if address >= 0x0C00 && address <= 0x0FFF {
				p.nameTable[1][address&0x03FF] = data
			}
		} else if p.cart.Mirror == HORIZONTAL {
			if address >= 0x0000 && address <= 0x03FF {
				p.nameTable[0][address&0x03FF] = data
			} else if address >= 0x0400 && address <= 0x03FF {
				p.nameTable[0][address&0x03FF] = data
			} else if address >= 0x0800 && address <= 0x0BFF {
				p.nameTable[1][address&0x03FF] = data
			} else if address >= 0x0C00 && address <= 0x0FFF {
				p.nameTable[1][address&0x03FF] = data
			}
		}
	} else if address >= 0x3F00 && address <= 0x3FFF {
		address &= 0x001F

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
		//
	} else if address >= 0x0000 && address <= 0x1FFF {
		p.patternTable[(address&0x1000)>>12][address&0x0FFF] = data
	} else if address >= 0x2000 && address <= 0x3EFF {
		if p.cart.Mirror == VERTICAL {
			if address >= 0x0000 && address <= 0x03FF {
				data = p.nameTable[0][address&0x03FF]
			} else if address >= 0x0400 && address <= 0x03FF {
				data = p.nameTable[1][address&0x03FF]
			} else if address >= 0x0800 && address <= 0x0BFF {
				data = p.nameTable[0][address&0x03FF]
			} else if address >= 0x0C00 && address <= 0x0FFF {
				data = p.nameTable[1][address&0x03FF]
			}
		} else if p.cart.Mirror == HORIZONTAL {
			if address >= 0x0000 && address <= 0x03FF {
				data = p.nameTable[0][address&0x03FF]
			} else if address >= 0x0400 && address <= 0x03FF {
				data = p.nameTable[0][address&0x03FF]
			} else if address >= 0x0800 && address <= 0x0BFF {
				data = p.nameTable[1][address&0x03FF]
			} else if address >= 0x0C00 && address <= 0x0FFF {
				data = p.nameTable[1][address&0x03FF]
			}
		}
	} else if address >= 0x3F00 && address <= 0x3FFF {
		address &= 0x001F
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

// CPURead : CPURead
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
		p.ppuDataBuffer, _ = p.PPURead(p.vRam.getAddress(), false)
		if p.vRam.getAddress() >= 0x3F00 {
			data = p.ppuDataBuffer
		}
		p.vRam.increment()
		break
	}
	return data, nil
}

func (p *PPU2C02) CPUWrite(address Word, data byte) error {

	switch address {
	case CONTROL_REGISTER:
		p.controlRegister = data
		p.tRam.nametableX = p.GetFlagByte(NAMETABLE_X, CONTROL_REGISTER)
		p.tRam.nametableY = p.GetFlagByte(NAMETABLE_Y, CONTROL_REGISTER)
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
		if p.addressLatch == 0 {
			p.fineX = data & 0x07
			p.tRam.coarseX = data >> 3

			p.addressLatch = 1
		} else {
			p.tRam.fineY = data & 0x07
			p.tRam.coarseY = data >> 3
			p.addressLatch = 0
		}
		// p.scrollRegister = data
		break
	case ADDRESS_REGISTER:
		if p.addressLatch == 0 {
			p.tRam.set((p.tRam.getAddress() & 0x00FF) | (Word(data) << 8))
			p.addressLatch = 1
		} else {
			p.tRam.set((p.tRam.getAddress() & 0xFF00) | Word(data))
			p.vRam.set(p.tRam.getAddress())
			p.addressLatch = 0
		}
		break
	case DATA_REGISTER:
		p.PPUWrite(p.vRam.getAddress(), data)
		increment := Word(1)
		if p.GetFlag(INCREMENT_MODE, CONTROL_REGISTER) {
			increment = Word(32)
		}
		p.vRam.add(increment)
		break
	}

	return nil
}

// Clock : Bus clock implementation pulses the clock to all things attached to it
func (p *PPU2C02) Clock() {

	incrementScrollX := func(p *PPU2C02) {
		if p.GetFlag(RENDER_BACKGROUND, MASK_REGISTER) || p.GetFlag(RENDER_SPRITES, MASK_REGISTER) {
			if p.vRam.coarseX == byte(31) {
				p.vRam.coarseX = 0
				p.vRam.nametableX = ^p.vRam.nametableX
			} else {
				p.vRam.coarseX++
			}
		}
	}

	incrementScrollY := func(p *PPU2C02) {
		if p.GetFlag(RENDER_BACKGROUND, MASK_REGISTER) || p.GetFlag(RENDER_SPRITES, MASK_REGISTER) {
			if p.vRam.fineY < 7 {
				p.vRam.fineY++
			} else {
				p.vRam.fineY = 0

				if p.vRam.coarseY == 29 {
					p.vRam.coarseY = 0
					p.vRam.nametableY = ^p.vRam.nametableY
				} else if p.vRam.coarseY == 31 {
					p.vRam.coarseY = 0
				} else {
					p.vRam.coarseY++
				}
			}

			if p.vRam.coarseX == byte(31) {
				p.vRam.coarseX = 0
				p.vRam.nametableX = ^p.vRam.nametableX
			} else {
				p.vRam.coarseX++
			}
		}
	}

	transferAddressX := func(p *PPU2C02) {
		if p.GetFlag(RENDER_BACKGROUND, MASK_REGISTER) || p.GetFlag(RENDER_SPRITES, MASK_REGISTER) {
			p.vRam.nametableX = p.tRam.nametableX
			p.vRam.coarseX = p.tRam.coarseX
		}
	}

	transferAddressY := func(p *PPU2C02) {
		if p.GetFlag(RENDER_BACKGROUND, MASK_REGISTER) || p.GetFlag(RENDER_SPRITES, MASK_REGISTER) {
			p.vRam.fineY = p.tRam.fineY
			p.vRam.nametableY = p.tRam.nametableY
			p.vRam.coarseY = p.tRam.coarseY
		}
	}

	loadBackgroundShifters := func(p *PPU2C02) {

		p.bgShifterPatternLo = p.bgShifterPatternLo&0xFF00 | Word(p.bgNextTileLsb)
		p.bgShifterPatternHi = p.bgShifterPatternHi&0xFF00 | Word(p.bgNextTileMsb)

		lo := Word(0x00)
		if p.bgNextTileAttrib*0x01 != 0 {
			lo = Word(0x00FF)
		}
		p.bgShifterAttribLo = p.bgShifterAttribLo&0xFF00 | lo

		hi := Word(0x00)
		if p.bgNextTileAttrib*0x02 != 0 {
			hi = Word(0x00FF)
		}
		p.bgShifterAttribHi = p.bgShifterAttribHi&0xFF00 | hi
	}

	updateShifters := func(p *PPU2C02) {
		if p.GetFlag(RENDER_BACKGROUND, MASK_REGISTER) {
			p.bgShifterPatternLo <<= 1
			p.bgShifterPatternHi <<= 1
			p.bgShifterAttribLo <<= 1
			p.bgShifterAttribHi <<= 1
		}
	}

	if p.scanLine >= -1 && p.scanLine < 240 {

		if p.scanLine == -1 && p.cycle == 1 {
			p.ClearFlag(VERTICAL_BLANK, STATUS_REGISTER)
		}

		if (p.cycle >= 2 && p.cycle < 258) || (p.cycle >= 321 && p.cycle < 338) {
			updateShifters(p)
			switch (p.cycle - 1) % 8 {
			case 0:
				loadBackgroundShifters(p)
				p.bgNextTileId, _ = p.PPURead(Word(0x0200)|(p.vRam.getAddress()&0x0FFF), false)
				break
			case 2:
				p.bgNextTileAttrib, _ = p.PPURead(
					Word(0x23C0)|(Word(p.vRam.nametableY)<<11)|(Word(p.vRam.nametableX)<<10)|((Word(p.vRam.coarseY)>>2)<<3)|(Word(p.vRam.coarseX)>>2),
					false)
				if p.vRam.coarseY&0x02 != 0x00 {
					p.bgNextTileAttrib >>= 4
				}
				if p.vRam.coarseX&0x02 != 0x00 {
					p.bgNextTileAttrib >>= 2
				}
				p.bgNextTileAttrib &= 0x03
				break
			case 4:
				p.bgNextTileLsb, _ = p.PPURead((Word(p.GetFlagByte(PATTERN_BACKGROUN, CONTROL_REGISTER))<<12)+Word(p.bgNextTileId)<<4+Word(p.vRam.fineY)+0, false)
				break
			case 6:
				p.bgNextTileLsb, _ = p.PPURead(Word(p.GetFlagByte(PATTERN_BACKGROUN, CONTROL_REGISTER))<<12+Word(p.bgNextTileId)<<4+Word(p.vRam.fineY)+8, false)
				break
			case 7:
				incrementScrollX(p)
				break
			}

		}

		if p.cycle == 256 {
			incrementScrollY(p)
		}
		if p.cycle == 257 {
			transferAddressX(p)
		}
		if p.scanLine == -1 && p.cycle >= 280 && p.cycle < 305 {
			transferAddressY(p)
		}
	}

	// if p.scanline == 240 // post render scanline

	if p.scanLine >= 241 && p.scanLine < 261 {
		if p.scanLine == 241 && p.cycle == 1 {
			// Effectively end of frame, so set vertical blank flag
			p.SetFlag(VERTICAL_BLANK, STATUS_REGISTER)

			// If the control register tells us to emit a NMI when
			// entering vertical blanking period, do it! The CPU
			// will be informed that rendering is complete so it can
			// perform operations with the PPU knowing it wont
			// produce visible artefacts

			if p.GetFlag(ENABLE_NMI, CONTROL_REGISTER) {
				p.NonMaskableInterrupt = true
			}
		}
	}

	bgPixel := byte(0x00)
	bgPalette := byte(0x00)

	if p.GetFlag(RENDER_BACKGROUND, MASK_REGISTER) {
		bitMux := Word(0x8000) >> p.fineX
		pixelLo := byte(0)
		if p.bgShifterAttribLo&bitMux > 0 {
			pixelLo = 1
		}
		pixelHi := byte(0)
		if p.bgShifterAttribHi&bitMux > 0 {
			pixelHi = 1
		}

		bgPixel = pixelHi<<1 | pixelLo

		palLo := byte(0)
		if p.bgShifterAttribLo&bitMux > 0 {
			palLo = 1
		}
		palHi := byte(0)
		if p.bgShifterAttribHi&bitMux > 0 {
			palHi = 1
		}

		bgPalette = palHi<<1 | palLo
	}

	// Paint pixel
	p.GetScreen().SetPixel(int(p.cycle)-1, int(p.scanLine), p.GetColorFromPaletteRam(bgPalette, bgPixel))

	p.cycle++
	if p.cycle >= 341 {
		p.cycle = 0
		p.scanLine++
		if p.scanLine >= 261 {
			p.scanLine = -1
			p.frameComplete = true
		}
	}
}

// InsertCartridge : sets the pointer to the cartridge in the PPU
func (p *PPU2C02) InsertCartridge(c *Cartridge) {
	p.cart = c
}
