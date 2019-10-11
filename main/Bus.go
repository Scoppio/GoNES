package main

import (
	"errors"
	"fmt"
	"runtime/debug"
)

var (
	// ClockCount : Used for debug, counts total clocks so far
	ClockCount = 0
)

// Bus Databus and things connected to it
type Bus struct {
	cpu  *CPU6502
	ppu  *PPU2C02
	cart *Cartridge
	ram  [2 * 1024]byte
}

// CreateBus : creates a new bus
func CreateBus(cpu *CPU6502, ppu *PPU2C02) *Bus {
	bus := &Bus{cpu, ppu, nil, [2 * 1024]byte{}}
	cpu.ConnectBus(bus)
	ppu.ConnectBus(bus)
	return bus
}

// PreLoadMemory : inserts data into memory using string format
// inserted data must be in hexadecimal writen as a string
// and they may have space after each 2 bytes
// example : "A9 0F 8D 15 40 60" or "A90F8D154060"
func (b *Bus) PreLoadMemory(offset rune, data string) {
	nOffset := offset
	for i := 0; i < len(data); i += 2 {
		b.CPUWrite(nOffset, ByteToHex(data[i])<<4|ByteToHex(data[i+1]))
		if i+2 < len(data) && data[i+2] == byte(' ') {
			i++
		}
		nOffset++
	}
}

// SetCodeEntry : Set the address that starts your program
func (b *Bus) SetCodeEntry(address rune) {
	b.CPUWrite(0xFFFC, byte(address))
	b.CPUWrite(0xFFFD, byte(address>>8))
}

func (b *Bus) CPURead(address rune, readOnly bool) (byte, error) {
	var d byte = 0x00
	var e error = nil
	if data, ok := b.cart.CPURead(address); ok {
		d = data
	} else if address >= 0x0000 && address <= 0x1FFF {
		d = b.ram[address&0x07FF]
	} else if address >= 0x2000 && address <= 0x3FFF {
		d, e = b.ppu.CPURead(address&0x0007, readOnly)
	} else if address == 0xFFFC || address == 0xFFFD {
		d = b.ram[address&0x07FF]
	} else {
		e = errors.New(fmt.Sprintln("tried to access index out of range", Hex(uint32(address), 4)))
		fmt.Print(string(debug.Stack()))
	}
	return d, e
}

func (b *Bus) CPUWrite(address rune, data byte) error {
	var e error = nil
	if ok := b.cart.CPUWrite(address, data); ok {
		//
	} else if address >= 0x0000 && address <= 0x1fff {
		b.ram[address&0x07FF] = data
	} else if address >= 0x2000 && address <= 0x3FFF {
		e = b.ppu.CPUWrite(address&0x0007, data)
	} else if address == 0xFFFC || address == 0xFFFD {
		b.ram[address&0x07FF] = data
	} else {
		e = errors.New(fmt.Sprintln("tried to access index out of range", Hex(uint32(address), 4)))
		fmt.Print(string(debug.Stack()))
	}
	return e
}

// Clock : Bus clock implementation pulses the clock to all things attached to it
func (b *Bus) Clock() {
	b.cpu.Clock()
}

func (b *Bus) ExecutOperation() {
	b.Clock()
	for !b.cpu.Complete() {
		b.Clock()
	}
	b.Clock()

	for b.cpu.Complete() {
		b.Clock()
	}
}

func (b *Bus) Reset() {
	b.cpu.Reset()
	ClockCount = 0
}

func (b *Bus) InsertCartridge(c *Cartridge) {
	b.cart = c
	b.ppu.InsertCartridge(c)
}
