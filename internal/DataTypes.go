package main

import (
	"reflect"
	"runtime"
)

// Flag : int that defines which register flag is active
type Flag int

// Register : int that defines which register is accessed
type Register int

// Word : uint 16
type Word uint16

// CPUReader : something that can read memory
type CPUReader interface {
	CPURead(address Word, readOnly bool) (data byte, err error)
}

// CPUWriter : something that can write to memory
type CPUWriter interface {
	CPUWrite(address Word, data byte) error
}

// PPUReader : something that can read memory
type PPUReader interface {
	PPURead(address Word, readOnly bool) (data byte, err error)
}

// PPUWriter : something that can write to memory
type PPUWriter interface {
	PPUWrite(address Word, data byte) error
}

// Completable : if the total number of cycles completed
type Completable interface {
	Complete() bool
}

// Clocker : a clockable thingy
type Clocker interface {
	Clock(n byte)
}

// Resetable : a resetable thingy
type Resetable interface {
	Reset()
}

// Accessable : Accessable
type Accessable interface {
	CPUWriter
	CPUReader
	PPUWriter
	PPUReader
}

// Operate : opcode operation
type Operate func(c *CPU6502) byte

// Addressing : opcode addressing mode
type Addressing func(c *CPU6502) byte

// FnEquals : checks if both functions have the same name
func FnEquals(fn interface{}, target interface{}) bool {
	return fnName(fn) == fnName(target)
}

// fnName : returns the name of the function
func fnName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return name
}

// Instruction : definition of an opcode/instruction
type Instruction struct {
	name        string
	operate     Operate
	addressmode Addressing
	cycles      byte
}
