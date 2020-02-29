package main

import (
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
	romname      = "color_test"
	lastUpdate   = time.Now()
)

func init() {
	nes = CreateBus(CreateCPU(), CreatePPU())
	nes.InsertCartridge(LoadCartridge("../test/roms/" + romname + ".nes"))
	cpu = nes.cpu
	nes.Reset()
	// mapAsm = cpu.Disassemble(0x0000, 0xFFFF)
	WriteDisassemble(mapAsm, "../output/"+romname+".txt")
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
