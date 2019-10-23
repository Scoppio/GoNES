package main

import (
	"testing"
)

func TestFnEquals(t *testing.T) {
	fn := IMM

	if got := FnEquals(fn, IMM); got != true {
		t.Errorf("Expected: %t, got: %t", true, got)
	}
}

func TestLoopyRegister(t *testing.T) {

	l := CreateLoopyRegister()
	want := Word(0x000)
	if l.getAddress() != want {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}
	// func (l *loopyRegister) getAddress() Word {
	// 	return Word(l.coarseX) | Word(l.coarseY)<<5 | Word(l.nametableX)<<10 | Word(l.nametableY)<<11 | Word(l.fineY)<<12
	// }
	l = CreateLoopyRegister()
	l.coarseX = byte(0x02)
	want = 0x0002
	if l.getAddress() != Word(0x0002) {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

	l = CreateLoopyRegister()
	l.coarseY = byte(0x01)
	want = 0x0020
	if l.getAddress() != Word(0x0020) {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

	l = CreateLoopyRegister()
	l.nametableX = byte(0x01)
	want = 0x0400
	if l.getAddress() != want {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

	l = CreateLoopyRegister()
	l.nametableY = byte(0x01)
	want = 0x0800
	if l.getAddress() != want {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

	l = CreateLoopyRegister()
	l.fineY = byte(0x01)
	want = 0x1000
	if l.getAddress() != want {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}
}

func TestLoopyRegisterAddition(t *testing.T) {

	l := CreateLoopyRegister()
	l.increment()
	want := Word(0x001)
	if l.getAddress() != want && l.coarseX != byte(want) {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

	l = CreateLoopyRegister()
	l.add(0x0020)
	want = Word(0x0020)
	if l.getAddress() != want && l.coarseY != byte(0x01) {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

	l = CreateLoopyRegister()
	l.add(0x0400)
	want = 0x0400
	if l.getAddress() != want && l.nametableX != byte(0x01) {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

	l = CreateLoopyRegister()
	l.add(0x0800)
	want = 0x0800
	if l.getAddress() != want && l.nametableY != byte(0x01) {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

	l = CreateLoopyRegister()
	l.set(0x1000)
	want = 0x1000
	if l.getAddress() != want && l.fineY != byte(0x01) && l.coarseX != 0 && l.coarseY != 0 && l.nametableX != 0 && l.nametableY != 0 {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

	l = CreateLoopyRegister()
	l.coarseY = 0x01
	l.increment()
	want = 0x0021
	if l.getAddress() != want && l.coarseX != byte(0x01) {
		t.Errorf("Expected: %x, got: %x", want, l.getAddress())
	}

}
