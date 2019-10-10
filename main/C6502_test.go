package main

import (
	"runtime/debug"
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
	memory := &Memory64k{}
	testBus := &Bus{testCPU, memory}
	cpu.ConnectBus(testBus)
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
	cpu.addressAbs = Word(0xf000)
	cpu.bus.Write(cpu.addressAbs, byte(0x3))

	ADC(cpu)

	assertFalse(t, cpu.StatusRegister(C))
	assertEqualsB(t, byte(0x5), cpu.a)
}

func TestOperationADCWithCarry(t *testing.T) {
	cpu := testCPU
	cpu.Reset()

	cpu.a = byte(0xFE)
	cpu.addressAbs = Word(0xf000)
	cpu.bus.Write(cpu.addressAbs, byte(0x3))

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
	cpu.addressAbs = Word(0xf000)
	cpu.bus.Write(cpu.addressAbs, byte(z))

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
	stackedStatus, e := cpu.Read(STACK + Word(stkp))
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
