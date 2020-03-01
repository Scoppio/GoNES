package main

import (
	"fmt"
	"time"
)

var (
	nes          *Bus
	cpu          *CPU6502
	mapAsm       map[Word]string
	frames       = 0
	second       = time.Tick(time.Second)
	emulationRun = false
	residualTime = 0.0
	elapsedTime  = 0.0
	lastUpdate   = time.Now()
)

// SetRom : Put a ROM on the memory of the Nes Emulator
func SetRom(rom string) {
	nes.InsertCartridge(LoadCartridge(rom))
	nes.Reset()
	mapAsm = cpu.Disassemble(0x0000, 0xFFFF)
	filename := time.Now().Format("2006-01-02_15:04:05")
	WriteDisassemble(mapAsm, "../output/disasemble_"+filename+".txt")
}

func init() {
	nes = CreateBus(CreateCPU(), CreatePPU())
	cpu = nes.cpu
	nes.Reset()
}

func tick() {
	nes.ExecuteOperation()
}

func reset() {
	nes.Reset()
	for !nes.cpu.Complete() {
		nes.Clock()
	}
}

func testCode() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Assert failed:", r)
		}
	}()
	// run Program until it sends the success msg
	for {
		tick()
		fmt.Println(cpu.pc, cpu.a, cpu.x, cpu.y, cpu.opcode, cpu.status, cpu.stkp)
	}
}
