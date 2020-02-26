package main

import (
	"image/color"
)

const (
	controlRegister = 0
	maskRegister    = 1
	statusRegister  = 2
	oamAddress      = 3
	oamData         = 4
	scrollRegister  = 5
	addressRegister = 6
	dataRegister    = 7

	// status register
	// unused 0-4
	spriteOverflow = 5
	spriteZeroHit  = 6
	verticalBlank  = 7

	grayScale            = 0
	renderBackgroundLeft = 1
	renderSpritesLeft    = 2
	renderBackground     = 3
	renderSprites        = 4
	enhanceRed           = 5
	enhanceGreen         = 6
	enhanceBlue          = 7

	nametableX        = 0
	nametableY        = 1
	incrementMode     = 2
	patternSprite     = 3
	patternBackground = 4
	spriteSize        = 5
	enableNMI         = 7
)

var (
	defaultColors [64]*color.RGBA
)

// LoopyRegister : LoopyRegister
type LoopyRegister struct {
	coarseX    byte
	coarseY    byte
	nametableX byte
	nametableY byte
	fineY      byte
}

func (l *LoopyRegister) getAddress() Word {
	return Word(l.coarseX) | Word(l.coarseY)<<5 | Word(l.nametableX)<<10 | Word(l.nametableY)<<11 | Word(l.fineY)<<12
}

func (l *LoopyRegister) increment() {
	l.add(1)
}

func (l *LoopyRegister) add(value Word) {
	tbd := l.getAddress() + value
	l.coarseX = byte(tbd & 0x001F)
	l.coarseY = byte((tbd >> 5) & 0x001F)
	l.nametableX = byte(tbd >> 10 & 0x0001)
	l.nametableY = byte(tbd >> 11 & 0x0001)
	l.fineY = byte(tbd >> 12 & 0x0007)
}

func (l *LoopyRegister) set(value Word) {
	tbd := value
	l.coarseX = byte(tbd & 0x001F)
	l.coarseY = byte((tbd >> 5) & 0x001F)
	l.nametableX = byte(tbd >> 10 & 0x0001)
	l.nametableY = byte(tbd >> 11 & 0x0001)
	l.fineY = byte(tbd >> 12 & 0x0007)
}

// CreateLoopyRegister : Create a cleared loopy register
func CreateLoopyRegister() *LoopyRegister {
	return &LoopyRegister{0, 0, 0, 0, 0}
}

// PPU2C02 : PPU
type PPU2C02 struct {
	bus           *Bus
	cart          *Cartridge
	nameTable     [2][1024]byte // VRAM
	paletteTable  [32]byte
	patternTable  [2][4096]byte // Pattern Memory
	paletteScreen [64]*color.RGBA
	// spriteScreen       *Sprite    // Sprite(256, 240)
	// spriteNameTable    [2]*Sprite // Sprite(128, 128), Sprite(128, 128)
	// spritePatternTable [2]*Sprite // Sprite(256, 240), Sprite(256, 240)
	frameComplete bool
	scanLine      int16
	cycle         int16

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

	vRAM  *LoopyRegister
	tRAM  *LoopyRegister
	fineX byte

	bgNextTileID     byte
	bgNextTileAttrib byte
	bgNextTileLsb    byte
	bgNextTileMsb    byte

	bgShifterPatternLo Word
	bgShifterPatternHi Word
	bgShifterAttribLo  Word
	bgShifterAttribHi  Word
}

func init() {
	defaultColors = [0x40]*color.RGBA{
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
		defaultColors,
		// CreateSprite(256, 240),
		// [2]*Sprite{
		// 	CreateSprite(256, 240),
		// 	CreateSprite(256, 240)},
		// [2]*Sprite{
		// 	CreateSprite(128, 128),
		// 	CreateSprite(128, 128)},
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

// InsertCartridge : sets the pointer to the cartridge in the PPU
func (p *PPU2C02) InsertCartridge(c *Cartridge) {
	p.cart = c
}

// GetFlag : GetFlag
func (p *PPU2C02) GetFlag(flag Flag, at Register) bool {
	var shortRegister byte
	var longRegister Word

	switch at {
	case controlRegister:
		shortRegister = p.controlRegister
		break
	case maskRegister:
		shortRegister = p.maskRegister
		break
	case statusRegister:
		shortRegister = p.statusRegister
		break
	case scrollRegister:
		shortRegister = p.scrollRegister
		break
	case dataRegister:
		shortRegister = p.dataRegister
		break
	case addressRegister:
		longRegister = p.addressRegister
		return (longRegister >> Word(flag) & Word(0x0001)) == 0
	}

	return (shortRegister >> byte(flag) & byte(0x01)) == 0
}

// GetFlagByte : GetFlagByte
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
	case controlRegister:
		if val {
			p.controlRegister |= byte(1 << uint(flag))
		} else {
			p.controlRegister &= ^byte(1 << uint(flag))
		}
		break
	case maskRegister:
		if val {
			p.maskRegister |= byte(1 << uint(flag))
		} else {
			p.maskRegister &= ^byte(1 << uint(flag))
		}
		break
	case statusRegister:
		if val {
			p.statusRegister |= byte(1 << uint(flag))
		} else {
			p.statusRegister &= ^byte(1 << uint(flag))
		}
		break
	case scrollRegister:
		if val {
			p.scrollRegister |= byte(1 << uint(flag))
		} else {
			p.scrollRegister &= ^byte(1 << uint(flag))
		}
		break
	case dataRegister:
		if val {
			p.dataRegister |= byte(1 << uint(flag))
		} else {
			p.dataRegister &= ^byte(1 << uint(flag))
		}
		break
	case addressRegister:
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

// // GetScreen : GetScreen
// func (p *PPU2C02) GetScreen() *Sprite {
// 	return p.spriteScreen
// }

// // GetNameTable : GetNameTable
// func (p *PPU2C02) GetNameTable(i int) *Sprite {
// 	return p.spriteNameTable[i]
// }

// GetColorFromPaletteRAM : GetColorFromPaletteRAM
func (p *PPU2C02) GetColorFromPaletteRAM(palette, pixelValue byte) *color.RGBA {
	idx, _ := p.PPURead(0x3F00+Word(palette)<<2+Word(pixelValue), false)
	idx &= 0x3F
	c := p.paletteScreen[idx]
	return c
}

// GetPatternTable : GetPatternTable
// func (p *PPU2C02) GetPatternTable(i, palette byte) *Sprite {
// 	var tileY, tileX, row, col, offset int

// 	for tileY = 0; tileY < 16; tileY++ {

// 		for tileX = 0; tileX < 16; tileX++ {

// 			offset = tileY*256 + tileX*16

// 			for row = 0; row < 8; row++ {

// 				pos := Word(i)*Word(0x1000) + Word(offset) + Word(row)
// 				tileLSB, _ := p.PPURead(pos+0, true)
// 				tileMSB, _ := p.PPURead(pos+8, true)

// 				for col = 0; col < 8; col++ {

// 					pixel := (tileMSB & 0x01) + (tileLSB & 0x01)
// 					tileLSB = tileLSB >> 1
// 					tileMSB = tileMSB >> 1

// 					// p.spritePatternTable[i].SetPixel(
// 					// 	tileX*8+(7-col),
// 					// 	tileY*8+row,
// 					// 	p.GetColorFromPaletteRAM(palette, pixel))

// 				}
// 			}
// 		}
// 	}

// 	return p.spritePatternTable[i]
// }

// ConnectBus : connects the CPU to the Bus
func (p *PPU2C02) ConnectBus(bus *Bus) {
	p.bus = bus
}

// PPURead : PPURead
func (p *PPU2C02) PPURead(address Word, readOnly bool) (byte, error) {
	var data byte = 0
	address &= 0x3FFF

	if _, ok := p.cart.PPURead(address); ok {
		// data = d
	} else if address >= 0x0000 && address <= 0x1FFF {
		data = p.patternTable[(address&0x1000)>>12][address&0x0FFF]
	} else if address >= 0x2000 && address <= 0x3EFF {

		address &= 0x0FFF

		if p.cart.Mirror == Vertical {
			if address >= 0x0000 && address <= 0x03FF {
				p.nameTable[0][address&0x03FF] = data
			} else if address >= 0x0400 && address <= 0x07FF {
				p.nameTable[1][address&0x03FF] = data
			} else if address >= 0x0800 && address <= 0x0BFF {
				p.nameTable[0][address&0x03FF] = data
			} else if address >= 0x0C00 && address <= 0x0FFF {
				p.nameTable[1][address&0x03FF] = data
			}
		} else if p.cart.Mirror == Horizontal {
			if address >= 0x0000 && address <= 0x03FF {
				p.nameTable[0][address&0x03FF] = data
			} else if address >= 0x0400 && address <= 0x07FF {
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
		if p.GetFlag(grayScale, maskRegister) {
			data = p.paletteTable[address] & 0x30
		} else {
			data = p.paletteTable[address] & 0x3F
		}
	}

	return data, nil
}

// PPUWrite : PPUWrite
func (p *PPU2C02) PPUWrite(address Word, data byte) error {
	address &= 0x3FFF

	if p.cart.PPUWrite(address, data) {
		// left empty
	} else if address >= 0x0000 && address <= 0x1FFF {
		p.patternTable[(address&0x1000)>>12][address&0x0FFF] = data
	} else if address >= 0x2000 && address <= 0x3EFF {

		address &= 0x0FFF

		if p.cart.Mirror == Vertical {
			if address >= 0x0000 && address <= 0x03FF {
				data = p.nameTable[0][address&0x03FF]
			} else if address >= 0x0400 && address <= 0x07FF {
				data = p.nameTable[1][address&0x03FF]
			} else if address >= 0x0800 && address <= 0x0BFF {
				data = p.nameTable[0][address&0x03FF]
			} else if address >= 0x0C00 && address <= 0x0FFF {
				data = p.nameTable[1][address&0x03FF]
			}
		} else if p.cart.Mirror == Horizontal {
			if address >= 0x0000 && address <= 0x03FF {
				data = p.nameTable[0][address&0x03FF]
			} else if address >= 0x0400 && address <= 0x07FF {
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
	if readOnly {
		switch address {
		case controlRegister:
			data = p.controlRegister
			break
		case maskRegister:
			data = p.maskRegister
			break
		case statusRegister:
			data = p.statusRegister
			break
		case oamAddress:
			break
		case oamData:
			break
		case scrollRegister:
			break
		case addressRegister:
			break
		case dataRegister:
			break
		}
	} else {
		switch address {
		case controlRegister:
			break
		case maskRegister:
			break
		case statusRegister:
			data = (p.statusRegister & 0xE0) | (p.ppuDataBuffer & 0x1F)
			p.ClearFlag(verticalBlank, statusRegister)
			p.addressLatch = 0
			break
		case oamAddress:
			break
		case oamData:
			break
		case scrollRegister:
			break
		case addressRegister:
			break
		case dataRegister:
			data = p.ppuDataBuffer
			p.ppuDataBuffer, _ = p.PPURead(p.vRAM.getAddress(), false)
			if p.vRAM.getAddress() >= 0x3F00 {
				data = p.ppuDataBuffer
			}
			if p.GetFlag(incrementMode, controlRegister) {
				p.vRAM.add(32)
			} else {
				p.vRAM.increment()
			}
			break
		}
	}

	return data, nil
}

// CPUWrite : write data to PPU
func (p *PPU2C02) CPUWrite(address Word, data byte) error {

	switch address {
	case controlRegister:
		p.controlRegister = data
		p.tRAM.nametableX = p.GetFlagByte(nametableX, controlRegister)
		p.tRAM.nametableY = p.GetFlagByte(nametableY, controlRegister)
		break
	case maskRegister:
		p.maskRegister = data
		break
	case statusRegister:
		break
	case oamAddress:
		break
	case oamData:
		break
	case scrollRegister:
		if p.addressLatch == 0 {
			p.fineX = data & 0x07
			p.tRAM.coarseX = data >> 3
			p.addressLatch = 1
		} else {
			p.tRAM.fineY = data & 0x07
			p.tRAM.coarseY = data >> 3
			p.addressLatch = 0
		}
		// p.scrollRegister = data
		break
	case addressRegister:
		if p.addressLatch == 0 {
			p.tRAM.set((p.tRAM.getAddress() & 0x00FF) | (Word(data&0x3F) << 8))
			p.addressLatch = 1
		} else {
			p.tRAM.set((p.tRAM.getAddress() & 0xFF00) | Word(data))
			p.vRAM.set(p.tRAM.getAddress())
			p.addressLatch = 0
		}
		break
	case dataRegister:
		p.PPUWrite(p.vRAM.getAddress(), data)

		if p.GetFlag(incrementMode, controlRegister) {
			p.vRAM.add(32)
		} else {
			p.vRAM.increment()
		}

		break
	}

	return nil
}

// Reset : reset the PPU
func (p *PPU2C02) Reset() {
	p.fineX = 0x00
	p.addressLatch = 0x00
	p.ppuDataBuffer = 0x00
	p.scanLine = 0
	p.cycle = 0
	p.bgNextTileID = 0x00
	p.bgNextTileAttrib = 0x00
	p.bgNextTileLsb = 0x00
	p.bgNextTileMsb = 0x00
	p.bgShifterPatternLo = 0x0000
	p.bgShifterPatternHi = 0x0000
	p.bgShifterAttribLo = 0x0000
	p.bgShifterAttribHi = 0x0000
	p.statusRegister = 0x00
	p.maskRegister = 0x00
	p.controlRegister = 0x00
	p.vRAM.set(0x0000)
	p.tRAM.set(0x0000)
}

// Clock : Bus clock implementation pulses the clock to all things attached to it
func (p *PPU2C02) Clock() {

	incrementScrollX := func(p *PPU2C02) {
		if p.GetFlag(renderBackground, maskRegister) || p.GetFlag(renderSprites, maskRegister) {
			if p.vRAM.coarseX == byte(31) {
				p.vRAM.coarseX = 0
				p.vRAM.nametableX = ^p.vRAM.nametableX
			} else {
				p.vRAM.coarseX++
			}
		}
	}

	incrementScrollY := func(p *PPU2C02) {
		if p.GetFlag(renderBackground, maskRegister) || p.GetFlag(renderSprites, maskRegister) {
			if p.vRAM.fineY < 7 {
				p.vRAM.fineY++
			} else {
				p.vRAM.fineY = 0

				if p.vRAM.coarseY == byte(29) {
					p.vRAM.coarseY = 0
					p.vRAM.nametableY = ^p.vRAM.nametableY
				} else if p.vRAM.coarseY == byte(31) {
					p.vRAM.coarseY = 0
				} else {
					p.vRAM.coarseY++
				}
			}
		}
	}

	transferAddressX := func(p *PPU2C02) {
		if p.GetFlag(renderBackground, maskRegister) || p.GetFlag(renderSprites, maskRegister) {
			p.vRAM.nametableX = p.tRAM.nametableX
			p.vRAM.coarseX = p.tRAM.coarseX
		}
	}

	transferAddressY := func(p *PPU2C02) {
		if p.GetFlag(renderBackground, maskRegister) || p.GetFlag(renderSprites, maskRegister) {
			p.vRAM.fineY = p.tRAM.fineY
			p.vRAM.nametableY = p.tRAM.nametableY
			p.vRAM.coarseY = p.tRAM.coarseY
		}
	}

	loadBackgroundShifters := func(p *PPU2C02) {

		p.bgShifterPatternLo = p.bgShifterPatternLo&0xFF00 | Word(p.bgNextTileLsb)
		p.bgShifterPatternHi = p.bgShifterPatternHi&0xFF00 | Word(p.bgNextTileMsb)

		lo := byte(0x00)
		if p.bgNextTileAttrib&0x01 != 0 {
			lo = 0xFF
		}
		p.bgShifterAttribLo = p.bgShifterAttribLo&0xFF00 | Word(lo)

		hi := byte(0x00)
		if p.bgNextTileAttrib&0x02 != 0 {
			hi = 0xFF
		}
		p.bgShifterAttribHi = p.bgShifterAttribHi&0xFF00 | Word(hi)
	}

	updateShifters := func(p *PPU2C02) {
		if p.GetFlag(renderBackground, maskRegister) {
			p.bgShifterPatternLo = p.bgShifterPatternLo << 1
			p.bgShifterPatternHi = p.bgShifterPatternHi << 1

			p.bgShifterAttribLo = p.bgShifterAttribLo << 1
			p.bgShifterAttribHi = p.bgShifterAttribHi << 1
		}
	}

	if p.scanLine >= -1 && p.scanLine < 240 {

		if p.scanLine == 0 && p.cycle == 0 {
			// Odd frame, cycle skip
			p.cycle = 1
		}

		if p.scanLine == -1 && p.cycle == 1 {
			p.ClearFlag(verticalBlank, statusRegister)
		}

		if (p.cycle >= 2 && p.cycle < 258) || (p.cycle >= 321 && p.cycle < 338) {
			updateShifters(p)
			switch (p.cycle - 1) % 8 {
			case 0:
				loadBackgroundShifters(p)
				p.bgNextTileID, _ = p.PPURead(Word(0x0200)|(p.vRAM.getAddress()&0x0FFF), false)
				break
			case 2:
				p.bgNextTileAttrib, _ = p.PPURead(
					Word(0x23C0)|(Word(p.vRAM.nametableY)<<11)|(Word(p.vRAM.nametableX)<<10)|((Word(p.vRAM.coarseY)>>2)<<3)|(Word(p.vRAM.coarseX)>>2),
					false)
				if p.vRAM.coarseY&0x02 > 0 {
					p.bgNextTileAttrib = p.bgNextTileAttrib >> 4
				}
				if p.vRAM.coarseX&0x02 > 0 {
					p.bgNextTileAttrib = p.bgNextTileAttrib >> 2
				}
				p.bgNextTileAttrib &= 0x03
				break
			case 4:
				p.bgNextTileLsb, _ = p.PPURead((Word(p.GetFlagByte(patternBackground, controlRegister))<<12)+Word(p.bgNextTileID)<<4+Word(p.vRAM.fineY)+0, false)
				break
			case 6:
				p.bgNextTileMsb, _ = p.PPURead(Word(p.GetFlagByte(patternBackground, controlRegister))<<12+Word(p.bgNextTileID)<<4+Word(p.vRAM.fineY)+8, false)
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
			loadBackgroundShifters(p)
			transferAddressX(p)
		}

		// Superfluous reads of tile id at end of scanline
		if p.cycle == 338 || p.cycle == 340 {
			p.bgNextTileID, _ = p.PPURead(0x2000|(p.vRAM.getAddress()&0x0FFF), false)
		}

		if p.scanLine == -1 && p.cycle >= 280 && p.cycle < 305 {
			transferAddressY(p)
		}
	}

	if p.scanLine == 240 {
		// DO NOTHING - POST RENDER SCANLINE
	}

	if p.scanLine >= 241 && p.scanLine < 261 {
		if p.scanLine == 241 && p.cycle == 1 {
			// Effectively end of frame, so set vertical blank flag
			p.SetFlag(verticalBlank, statusRegister)

			// If the control register tells us to emit a NMI when
			// entering vertical blanking period, do it! The CPU
			// will be informed that rendering is complete so it can
			// perform operations with the PPU knowing it wont
			// produce visible artefacts

			if p.GetFlag(enableNMI, controlRegister) {
				p.NonMaskableInterrupt = true
			}
		}
	}

	// bgPixel := byte(0x00)
	// bgPalette := byte(0x00)

	// if p.GetFlag(renderBackground, maskRegister) {
	// 	bitMux := Word(0x8000) >> p.fineX
	// 	p0Pixel := byte(0)
	// 	if p.bgShifterAttribLo&bitMux > 0 {
	// 		p0Pixel = 1
	// 	}
	// 	p1Pixel := byte(0)
	// 	if p.bgShifterAttribHi&bitMux > 0 {
	// 		p1Pixel = 1
	// 	}

	// 	bgPixel = (p1Pixel << 1) | p0Pixel

	// 	bgPal0 := byte(0)
	// 	if p.bgShifterAttribLo&bitMux > 0 {
	// 		bgPal0 = 1
	// 	}
	// 	bgPal1 := byte(0)
	// 	if p.bgShifterAttribHi&bitMux > 0 {
	// 		bgPal1 = 1
	// 	}

	// 	bgPalette = (bgPal1 << 1) | bgPal0
	// }

	// Paint pixel
	// p.GetScreen().SetPixel(int(p.cycle)-1, int(p.scanLine), p.GetColorFromPaletteRAM(bgPalette, bgPixel))

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
