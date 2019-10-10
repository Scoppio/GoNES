package main

import (
	"reflect"
	"runtime"
)

// Flag : int that defines which register flag is active
type Flag int

// Word : two bytes
type Word uint16

// Reader : something that can read memory
type Reader interface {
	Read(address rune, readOnly bool) (data byte, err error)
}

// Writer : something that can write to memory
type Writer interface {
	Write(address rune, data byte) error
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

// CPU : Functions pertaining to a cpu
type CPU interface {
	Resetable
	Completable
	Clocker
	Reader
	Writer
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
