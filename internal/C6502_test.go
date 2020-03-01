package main

import (
	"strconv"
	"testing"
)

var (
	testCPU *CPU6502
)

const (
	ZeroB byte = 0
	ZeroW Word = 0
)

func init() {
	testCPU = &CPU6502{}
	ppu := &PPU2C02{}

	testBus := CreateBus(testCPU, ppu)
	cpu.ConnectBus(testBus)
}

func TestReset(t *testing.T) {
	cpu := testCPU
	cpu.Reset()
	t.Logf("status = 0b%s", strconv.FormatInt(int64(cpu.status), 2))

	assertEqualsB(t, ZeroB, cpu.x)
	assertEqualsB(t, ZeroB, cpu.y)
	assertEqualsB(t, ZeroB, cpu.a)
	assertEqualsB(t, byte(0xFD), cpu.stkp)
	assertFalse(t, cpu.StatusRegister(C))
	assertFalse(t, cpu.StatusRegister(Z))
	assertFalse(t, cpu.StatusRegister(I))
	assertFalse(t, cpu.StatusRegister(D))
	assertFalse(t, cpu.StatusRegister(B))
	assertTrue(t, cpu.StatusRegister(U))
	assertFalse(t, cpu.StatusRegister(V))
	assertFalse(t, cpu.StatusRegister(N))
	assertEqualsW(t, ZeroW, cpu.addressAbs)
	assertEqualsW(t, ZeroW, cpu.addressRel)
	assertEqualsW(t, ZeroW, cpu.pc)
	assertEqualsB(t, ZeroB, cpu.fetched)
	assertEqualsB(t, byte(8), cpu.cycles)
}

func TestOperationADC(t *testing.T) {

	cpu := testCPU
	cpu.Reset()

	cpu.a = byte(0x2)
	cpu.addressAbs = Word(0x0100)
	cpu.bus.CPUWrite(cpu.addressAbs, byte(0x3))

	ADC(cpu)

	assertFalse(t, cpu.StatusRegister(C))
	assertEqualsB(t, byte(0x5), cpu.a)
}

func TestOperationADCWithCarry(t *testing.T) {
	cpu := testCPU
	cpu.Reset()

	cpu.a = byte(0xFE)
	cpu.addressAbs = Word(0x0100)
	cpu.bus.CPUWrite(cpu.addressAbs, byte(0x3))

	ADC(cpu)

	assertTrue(t, cpu.StatusRegister(C))
	assertFalse(t, cpu.StatusRegister(V))
	assertEqualsB(t, byte(0x1), cpu.a)
}

func TestOperationADCWithOverflow(t *testing.T) {
	cpu := testCPU
	cpu.Reset()
	v := -10
	z := -127
	cpu.a = byte(v)
	cpu.addressAbs = Word(0x0100)
	cpu.bus.CPUWrite(cpu.addressAbs, byte(z))

	ADC(cpu)

	assertTrue(t, cpu.StatusRegister(C))
	assertTrue(t, cpu.StatusRegister(V))
}

func TestOperationPHP(t *testing.T) {
	cpu := testCPU
	cpu.Reset()

	oldStatus := byte(0x30)

	stkp := cpu.stkp
	PHP(cpu)
	stackedStatus, e := cpu.CPURead(Stack + Word(stkp))
	assertNil(t, e)
	assertEqualsB(t, oldStatus, stackedStatus)
}

func TestOperationJSR(t *testing.T) {
	cpu := testCPU
	cpu.Reset()
	cpu.addressAbs = Word(0xABCD)
	cpu.pc = Word(2)

	JSR(cpu)

	assertEqualsW(t, Word(0xABCD), cpu.pc)
}
