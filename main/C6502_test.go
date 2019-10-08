package main

import (
	"runtime/debug"
	"strconv"
	"testing"
)

var t_bus *Bus
var cpu *CPU6502

const (
	ZERO_B byte = 0
	ZERO_W Word = 0
)

func init() {
	cpu = &CPU6502{}
	memory := &Memory64k{}
	t_bus = &Bus{cpu, memory}
	cpu.ConnectBus(t_bus)
}

func assertNil(t *testing.T, got interface{}) {
	if got != nil {
		t.Errorf("Expected variable to be: nil")
		t.Log(string(debug.Stack()))
	}
}

func assertTrue(t *testing.T, got bool) {
	if !got {
		t.Errorf("Expected variable to be: %t, got: %t", true, got)
		t.Log(string(debug.Stack()))
	}
}

func assertFalse(t *testing.T, got bool) {
	if got {
		t.Errorf("Expected variable to be: %t, got: %t", false, got)
		t.Log(string(debug.Stack()))
	}
}

func assertEqualsB(t *testing.T, expect byte, got byte) {
	if expect != got {
		t.Errorf("Expected variable to be: %x, got: %x", expect, got)
		t.Log(string(debug.Stack()))
	}
}

func assertEqualsW(t *testing.T, expect Word, got Word) {
	if expect != got {
		t.Errorf("Expected variable to be: %x, got: %x", expect, got)
		t.Log(string(debug.Stack()))
	}
}

func TestReset(t *testing.T) {
	cpu.Reset()

	t.Logf("status = 0b%s", strconv.FormatInt(int64(cpu.status), 2))

	assertEqualsB(t, ZERO_B, cpu.x)
	assertEqualsB(t, ZERO_B, cpu.y)
	assertEqualsB(t, ZERO_B, cpu.a)
	assertEqualsB(t, byte(0xFD), cpu.stkp)
	assertFalse(t, cpu.StatusRegister(C))
	assertFalse(t, cpu.StatusRegister(Z))
	assertFalse(t, cpu.StatusRegister(I))
	assertFalse(t, cpu.StatusRegister(D))
	assertFalse(t, cpu.StatusRegister(B))
	assertTrue(t, cpu.StatusRegister(U))
	assertFalse(t, cpu.StatusRegister(V))
	assertFalse(t, cpu.StatusRegister(N))
	assertEqualsW(t, ZERO_W, cpu.address_abs)
	assertEqualsW(t, ZERO_W, cpu.address_rel)
	assertEqualsW(t, ZERO_W, cpu.pc)
	assertEqualsB(t, ZERO_B, cpu.fetched)
	assertEqualsB(t, byte(8), cpu.cycles)
}

func TestOperationADC(t *testing.T) {
	cpu.Reset()

	cpu.a = byte(0x2)
	cpu.address_abs = Word(0xf000)
	t_bus.Write(cpu.address_abs, byte(0x3))

	ADC(cpu)

	assertFalse(t, cpu.StatusRegister(C))
	assertEqualsB(t, byte(0x5), cpu.a)
}

func TestOperationADCWithCarry(t *testing.T) {
	cpu.Reset()

	cpu.a = byte(0xFE)
	cpu.address_abs = Word(0xf000)
	t_bus.Write(cpu.address_abs, byte(0x3))

	ADC(cpu)

	assertTrue(t, cpu.StatusRegister(C))
	assertFalse(t, cpu.StatusRegister(V))
	assertEqualsB(t, byte(0x1), cpu.a)
}

func TestOperationADCWithOverflow(t *testing.T) {
	cpu.Reset()
	v := -10
	z := -127
	cpu.a = byte(v)
	cpu.address_abs = Word(0xf000)
	t_bus.Write(cpu.address_abs, byte(z))

	ADC(cpu)

	assertTrue(t, cpu.StatusRegister(C))
	assertTrue(t, cpu.StatusRegister(V))
}

func TestOperationPHP(t *testing.T) {
	cpu.Reset()

	old_status := byte(0x30)

	stkp := cpu.stkp
	PHP(cpu)
	stacked_status, e := cpu.Read(STACK + Word(stkp))
	assertNil(t, e)
	assertEqualsB(t, old_status, stacked_status)
}

func TestOperationJSR(t *testing.T) {
	cpu.Reset()
	cpu.address_abs = Word(0xABCD)
	cpu.pc = Word(2)

	JSR(cpu)

	assertEqualsW(t, Word(0xABCD), cpu.pc)
}
