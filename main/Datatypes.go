package main

import (
	"reflect"
	"runtime"
)

type Flag int
type Word uint16

type Reader interface {
	Read(address rune, readOnly bool) (data byte, err error)
}

type Writer interface {
	Write(address rune, data byte) error
}

type Completable interface {
	Complete() bool
}

type Clocker interface {
	Clock(n byte)
}

type Resetable interface {
	Reset()
}

type CPU interface {
	Resetable
	Completable
	Clocker
	Reader
	Writer
}

type Operate func(c *CPU6502) byte

type Addressing func(c *CPU6502) byte

func FnEquals(fn interface{}, target interface{}) bool {
	return fnName(fn) == fnName(target)
}

func fnName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return name
}

type Instruction struct {
	name        string
	operate     Operate
	addressmode Addressing
	cycles      byte
}
