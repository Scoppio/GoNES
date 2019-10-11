package main

import "fmt"

type PPU2C02 struct {
	bus          *Bus
	cart         *Cartridge
	nameTable    [2][1024]byte
	paletteTable [32]byte
}

// ConnectBus : connects the CPU to the Bus
func (p *PPU2C02) ConnectBus(bus *Bus) {
	p.bus = bus
}

func CreatePPU() *PPU2C02 {
	return &PPU2C02{nil, nil, [2][1024]byte{}, [32]byte{}}
}

func (p *PPU2C02) PPURead(address rune, readOnly bool) (byte, error) {
	var data byte = 0
	address &= 0x3FFF

	if mappedAddress, ok := p.cart.CPURead(address); ok {
		fmt.Println(mappedAddress)
	}

	return data, nil
}

func (p *PPU2C02) PPUWrite(address rune, data byte) error {
	address &= 0x3FFF

	if p.cart.CPUWrite(address, data) {

	}

	return nil
}

func (p *PPU2C02) CPURead(address rune, readOnly bool) (byte, error) {
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

func (p *PPU2C02) CPUWrite(address rune, data byte) error {

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
func (p *PPU2C02) Clock(n byte) {}

func (p *PPU2C02) InsertCartridge(c *Cartridge) {
	p.cart = c
}
