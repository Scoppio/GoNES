package main

import (
	"fmt"
	"image/color"
	"math/rand"
)

type PPU2C02 struct {
	bus                *Bus
	cart               *Cartridge
	nameTable          [2][1024]byte
	paletteTable       [32]byte
	paletteScreen      [64]*color.RGBA
	spriteScreen       *Sprite    // Sprite(256, 240)
	spriteNameTable    [2]*Sprite // Sprite (128, 128), Sprite(128, 128)
	spritePatternTable [2]*Sprite // Sprite(256, 240), Sprite(256, 240)
	frameComplete      bool
	scanLine           int16
	cycle              int16
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

func (p *PPU2C02) GetPatternTable(i byte) *Sprite {
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
		[0x40]*color.RGBA{
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
		},
		CreateSprite(256, 240),
		[2]*Sprite{
			CreateSprite(256, 240),
			CreateSprite(256, 240)},
		[2]*Sprite{
			CreateSprite(256, 240),
			CreateSprite(256, 240)},
		false,
		0,
		0}
}

func (p *PPU2C02) PPURead(address Word, readOnly bool) (byte, error) {
	var data byte = 0
	address &= 0x3FFF

	if mappedAddress, ok := p.cart.CPURead(address); ok {
		fmt.Println(mappedAddress)
	}

	return data, nil
}

func (p *PPU2C02) PPUWrite(address Word, data byte) error {
	address &= 0x3FFF

	if p.cart.CPUWrite(address, data) {

	}

	return nil
}

func (p *PPU2C02) CPURead(address Word, readOnly bool) (byte, error) {
	var data byte = 0

	switch address {
	case 0x0000:
		break
	case 0x0001:
		break
	case 0x0002:
		break
	case 0x0003:
		break
	case 0x0004:
		break
	case 0x0005:
		break
	case 0x0006:
		break
	case 0x0007:
		break
	}
	return data, nil
}

func (p *PPU2C02) CPUWrite(address Word, data byte) error {

	switch address {
	case 0x0000:
		break
	case 0x0001:
		break
	case 0x0002:
		break
	case 0x0003:
		break
	case 0x0004:
		break
	case 0x0005:
		break
	case 0x0006:
		break
	case 0x0007:
		break
	}

	return nil
}

// Clock : Bus clock implementation pulses the clock to all things attached to it
func (p *PPU2C02) Clock() {

	c := byte(0x30)
	if rand.Int()%2 == 0 {
		c = byte(0x3F)
	}

	p.GetScreen().SetPixel(int(p.cycle)-1, int(p.scanLine), p.paletteScreen[c])

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

func (p *PPU2C02) InsertCartridge(c *Cartridge) {
	p.cart = c
}
