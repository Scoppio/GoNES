package main

import (
	"errors"
)

type Bus struct {
	cpu *CPU6502
	ram *Memory64k
}

func (b *Bus) Read(address Word, readOnly bool) (data byte, err error) {
	var d byte = 0x00
	var e error = nil
	if address >= 0x0000 && address <= 0xffff {
		d = b.ram.mem[address]
	} else {
		e = errors.New("Tried to access index out of range!")
	}
	return d, e
}

func (b *Bus) Write(address Word, data byte) error {
	var e error = nil
	if address >= 0x0000 && address <= 0xffff {
		b.ram.mem[address] = data
	} else {
		e = errors.New("Tried to access index out of range!")
	}
	return e
}

func (b *Bus) Clock(n byte) {}

func (b *Bus) String() string {
	return "Bus with CPU and MemoryBank"
}
