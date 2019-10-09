package main

import (
	"bytes"
	"fmt"
	"log"
)

const (
	C     = 0
	Z     = 1
	I     = 2
	D     = 3
	B     = 4
	U     = 5
	V     = 6
	N     = 7
	STACK = Word(1 << 3)
)

var (
	clock_count = 0
)

type CPU6502 struct {
	a, x, y, stkp, status, fetched, opcode, cycles byte
	pc, address_abs, address_rel                   Word
	bus                                            *Bus
}

var OpCodesLookupTable []Instruction

func init() {
	OpCodesLookupTable = []Instruction{
		{"BRK", BRK, IMM, 7}, {"ORA", ORA, IZX, 6}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"???", NOP, IMP, 3}, {"ORA", ORA, ZP0, 3}, {"ASL", ASL, ZP0, 5}, {"???", XXX, IMP, 5}, {"PHP", PHP, IMP, 3}, {"ORA", ORA, IMM, 2}, {"ASL", ASL, IMP, 2}, {"???", XXX, IMP, 2}, {"???", NOP, IMP, 4}, {"ORA", ORA, ABS, 4}, {"ASL", ASL, ABS, 6}, {"???", XXX, IMP, 6},
		{"BPL", BPL, REL, 2}, {"ORA", ORA, IZY, 5}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"???", NOP, IMP, 4}, {"ORA", ORA, ZPX, 4}, {"ASL", ASL, ZPX, 6}, {"???", XXX, IMP, 6}, {"CLC", CLC, IMP, 2}, {"ORA", ORA, ABY, 4}, {"???", NOP, IMP, 2}, {"???", XXX, IMP, 7}, {"???", NOP, IMP, 4}, {"ORA", ORA, ABX, 4}, {"ASL", ASL, ABX, 7}, {"???", XXX, IMP, 7},
		{"JSR", JSR, ABS, 6}, {"AND", AND, IZX, 6}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"BIT", BIT, ZP0, 3}, {"AND", AND, ZP0, 3}, {"ROL", ROL, ZP0, 5}, {"???", XXX, IMP, 5}, {"PLP", PLP, IMP, 4}, {"AND", AND, IMM, 2}, {"ROL", ROL, IMP, 2}, {"???", XXX, IMP, 2}, {"BIT", BIT, ABS, 4}, {"AND", AND, ABS, 4}, {"ROL", ROL, ABS, 6}, {"???", XXX, IMP, 6},
		{"BMI", BMI, REL, 2}, {"AND", AND, IZY, 5}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"???", NOP, IMP, 4}, {"AND", AND, ZPX, 4}, {"ROL", ROL, ZPX, 6}, {"???", XXX, IMP, 6}, {"SEC", SEC, IMP, 2}, {"AND", AND, ABY, 4}, {"???", NOP, IMP, 2}, {"???", XXX, IMP, 7}, {"???", NOP, IMP, 4}, {"AND", AND, ABX, 4}, {"ROL", ROL, ABX, 7}, {"???", XXX, IMP, 7},
		{"RTI", RTI, IMP, 6}, {"EOR", EOR, IZX, 6}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"???", NOP, IMP, 3}, {"EOR", EOR, ZP0, 3}, {"LSR", LSR, ZP0, 5}, {"???", XXX, IMP, 5}, {"PHA", PHA, IMP, 3}, {"EOR", EOR, IMM, 2}, {"LSR", LSR, IMP, 2}, {"???", XXX, IMP, 2}, {"JMP", JMP, ABS, 3}, {"EOR", EOR, ABS, 4}, {"LSR", LSR, ABS, 6}, {"???", XXX, IMP, 6},
		{"BVC", BVC, REL, 2}, {"EOR", EOR, IZY, 5}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"???", NOP, IMP, 4}, {"EOR", EOR, ZPX, 4}, {"LSR", LSR, ZPX, 6}, {"???", XXX, IMP, 6}, {"CLI", CLI, IMP, 2}, {"EOR", EOR, ABY, 4}, {"???", NOP, IMP, 2}, {"???", XXX, IMP, 7}, {"???", NOP, IMP, 4}, {"EOR", EOR, ABX, 4}, {"LSR", LSR, ABX, 7}, {"???", XXX, IMP, 7},
		{"RTS", RTS, IMP, 6}, {"ADC", ADC, IZX, 6}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"???", NOP, IMP, 3}, {"ADC", ADC, ZP0, 3}, {"ROR", ROR, ZP0, 5}, {"???", XXX, IMP, 5}, {"PLA", PLA, IMP, 4}, {"ADC", ADC, IMM, 2}, {"ROR", ROR, IMP, 2}, {"???", XXX, IMP, 2}, {"JMP", JMP, IND, 5}, {"ADC", ADC, ABS, 4}, {"ROR", ROR, ABS, 6}, {"???", XXX, IMP, 6},
		{"BVS", BVS, REL, 2}, {"ADC", ADC, IZY, 5}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"???", NOP, IMP, 4}, {"ADC", ADC, ZPX, 4}, {"ROR", ROR, ZPX, 6}, {"???", XXX, IMP, 6}, {"SEI", SEI, IMP, 2}, {"ADC", ADC, ABY, 4}, {"???", NOP, IMP, 2}, {"???", XXX, IMP, 7}, {"???", NOP, IMP, 4}, {"ADC", ADC, ABX, 4}, {"ROR", ROR, ABX, 7}, {"???", XXX, IMP, 7},
		{"???", NOP, IMP, 2}, {"STA", STA, IZX, 6}, {"???", NOP, IMP, 2}, {"???", XXX, IMP, 6}, {"STY", STY, ZP0, 3}, {"STA", STA, ZP0, 3}, {"STX", STX, ZP0, 3}, {"???", XXX, IMP, 3}, {"DEY", DEY, IMP, 2}, {"???", NOP, IMP, 2}, {"TXA", TXA, IMP, 2}, {"???", XXX, IMP, 2}, {"STY", STY, ABS, 4}, {"STA", STA, ABS, 4}, {"STX", STX, ABS, 4}, {"???", XXX, IMP, 4},
		{"BCC", BCC, REL, 2}, {"STA", STA, IZY, 6}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 6}, {"STY", STY, ZPX, 4}, {"STA", STA, ZPX, 4}, {"STX", STX, ZPY, 4}, {"???", XXX, IMP, 4}, {"TYA", TYA, IMP, 2}, {"STA", STA, ABY, 5}, {"TXS", TXS, IMP, 2}, {"???", XXX, IMP, 5}, {"???", NOP, IMP, 5}, {"STA", STA, ABX, 5}, {"???", XXX, IMP, 5}, {"???", XXX, IMP, 5},
		{"LDY", LDY, IMM, 2}, {"LDA", LDA, IZX, 6}, {"LDX", LDX, IMM, 2}, {"???", XXX, IMP, 6}, {"LDY", LDY, ZP0, 3}, {"LDA", LDA, ZP0, 3}, {"LDX", LDX, ZP0, 3}, {"???", XXX, IMP, 3}, {"TAY", TAY, IMP, 2}, {"LDA", LDA, IMM, 2}, {"TAX", TAX, IMP, 2}, {"???", XXX, IMP, 2}, {"LDY", LDY, ABS, 4}, {"LDA", LDA, ABS, 4}, {"LDX", LDX, ABS, 4}, {"???", XXX, IMP, 4},
		{"BCS", BCS, REL, 2}, {"LDA", LDA, IZY, 5}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 5}, {"LDY", LDY, ZPX, 4}, {"LDA", LDA, ZPX, 4}, {"LDX", LDX, ZPY, 4}, {"???", XXX, IMP, 4}, {"CLV", CLV, IMP, 2}, {"LDA", LDA, ABY, 4}, {"TSX", TSX, IMP, 2}, {"???", XXX, IMP, 4}, {"LDY", LDY, ABX, 4}, {"LDA", LDA, ABX, 4}, {"LDX", LDX, ABY, 4}, {"???", XXX, IMP, 4},
		{"CPY", CPY, IMM, 2}, {"CMP", CMP, IZX, 6}, {"???", NOP, IMP, 2}, {"???", XXX, IMP, 8}, {"CPY", CPY, ZP0, 3}, {"CMP", CMP, ZP0, 3}, {"DEC", DEC, ZP0, 5}, {"???", XXX, IMP, 5}, {"INY", INY, IMP, 2}, {"CMP", CMP, IMM, 2}, {"DEX", DEX, IMP, 2}, {"???", XXX, IMP, 2}, {"CPY", CPY, ABS, 4}, {"CMP", CMP, ABS, 4}, {"DEC", DEC, ABS, 6}, {"???", XXX, IMP, 6},
		{"BNE", BNE, REL, 2}, {"CMP", CMP, IZY, 5}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"???", NOP, IMP, 4}, {"CMP", CMP, ZPX, 4}, {"DEC", DEC, ZPX, 6}, {"???", XXX, IMP, 6}, {"CLD", CLD, IMP, 2}, {"CMP", CMP, ABY, 4}, {"NOP", NOP, IMP, 2}, {"???", XXX, IMP, 7}, {"???", NOP, IMP, 4}, {"CMP", CMP, ABX, 4}, {"DEC", DEC, ABX, 7}, {"???", XXX, IMP, 7},
		{"CPX", CPX, IMM, 2}, {"SBC", SBC, IZX, 6}, {"???", NOP, IMP, 2}, {"???", XXX, IMP, 8}, {"CPX", CPX, ZP0, 3}, {"SBC", SBC, ZP0, 3}, {"INC", INC, ZP0, 5}, {"???", XXX, IMP, 5}, {"INX", INX, IMP, 2}, {"SBC", SBC, IMM, 2}, {"NOP", NOP, IMP, 2}, {"???", SBC, IMP, 2}, {"CPX", CPX, ABS, 4}, {"SBC", SBC, ABS, 4}, {"INC", INC, ABS, 6}, {"???", XXX, IMP, 6},
		{"BEQ", BEQ, REL, 2}, {"SBC", SBC, IZY, 5}, {"???", XXX, IMP, 2}, {"???", XXX, IMP, 8}, {"???", NOP, IMP, 4}, {"SBC", SBC, ZPX, 4}, {"INC", INC, ZPX, 6}, {"???", XXX, IMP, 6}, {"SED", SED, IMP, 2}, {"SBC", SBC, ABY, 4}, {"NOP", NOP, IMP, 2}, {"???", XXX, IMP, 7}, {"???", NOP, IMP, 4}, {"SBC", SBC, ABX, 4}, {"INC", INC, ABX, 7}, {"???", XXX, IMP, 7}}

}

func (c *CPU6502) StatusRegister(flag Flag) bool {
	val := c.status & (1 << flag)
	return (val != 0)
}

func (c *CPU6502) StatusRegisterAsWord(flag Flag) Word {
	return Word(c.status & (1 << flag))
}

func (c *CPU6502) SetStatusRegisterFlag(flag Flag, val bool) {
	if val {
		c.status |= byte(1 << flag)
	} else {
		c.status &= ^byte(1 << flag)
	}
}

func (c *CPU6502) SetFlagsZeroAndNegative(val byte) {
	c.SetStatusRegisterFlag(Z, val == 0x00)
	c.SetStatusRegisterFlag(N, val&0x80 != 0)
}

func (c *CPU6502) ConnectBus(bus *Bus) {
	c.bus = bus
}

func (c *CPU6502) Read(address Word) (data byte, err error) {
	return c.bus.Read(address, false)
}

func (c *CPU6502) Write(address Word, data byte) error {
	return c.bus.Write(address, data)
}

func (c *CPU6502) Clock() {
	// execute
	if c.cycles == byte(0x00) {
		var e error
		c.opcode, e = c.Read(c.pc)
		if e != nil {
			log.Fatalf("Error when trying to access address %d, error %s", c.pc, e)
		}
		c.SetStatusRegisterFlag(U, true)
		c.pc = c.pc + 1
		c.cycles = OpCodesLookupTable[c.opcode].cycles
		additionalCycle := OpCodesLookupTable[c.opcode].addressmode(c)
		additionalCycle2 := OpCodesLookupTable[c.opcode].operate(c)
		c.cycles += (additionalCycle & additionalCycle2)
		// Always set the unused status flag bit to 1
		c.SetStatusRegisterFlag(U, true)
	}
	clock_count++
	c.cycles--
}

func (c *CPU6502) Complete() bool {
	return c.cycles == byte(0x00)
}

func (c *CPU6502) Reset() {
	c.a = 0
	c.x = 0
	c.y = 0
	c.stkp = 0xFD
	c.status = 0x00
	c.SetStatusRegisterFlag(U, true)
	c.address_abs = 0xFFFC
	lo, e := c.Read(c.address_abs + 0)
	if addressingError(e) {
		//
	}
	hi, e2 := c.Read(c.address_abs + 1)
	if addressingError(e2) {
		//
	}
	c.pc = Word(hi)<<8 | Word(lo)
	c.address_rel = 0x0000
	c.address_abs = 0x0000
	c.fetched = 0x00

	c.cycles = 8
}

func (c *CPU6502) InterruptRequest() {
	if !c.StatusRegister(I) {
		c.Write(STACK+Word(c.stkp), byte((c.pc>>8)&0x00ff))
		c.stkp--
		c.Write(STACK+Word(c.stkp), byte(c.pc&0x00FF))
		c.stkp--

		c.SetStatusRegisterFlag(B, false)
		c.SetStatusRegisterFlag(U, true)
		c.SetStatusRegisterFlag(I, true)
		c.Write(STACK+Word(c.stkp), c.status)
		c.stkp--

		c.address_abs = 0xFFFE
		lo, e := c.Read(c.address_abs + 0)
		if addressingError(e) {
			//
		}
		hi, e2 := c.Read(c.address_abs + 1)
		if addressingError(e2) {
			//
		}
		c.pc = Word(hi)<<8 | Word(lo)

		c.cycles = 7
	}
}

func (c *CPU6502) NonMaskableInterruptRequest() {
	c.Write(STACK+Word(c.stkp), byte((c.pc>>8)&0x00ff))
	c.stkp--
	c.Write(STACK+Word(c.stkp), byte(c.pc&0x00FF))
	c.stkp--

	c.SetStatusRegisterFlag(B, false)
	c.SetStatusRegisterFlag(U, true)
	c.SetStatusRegisterFlag(I, true)
	c.Write(STACK+Word(c.stkp), c.status)
	c.stkp--

	c.address_abs = 0xFFFA
	lo, e := c.Read(c.address_abs + 0)
	if addressingError(e) {
		//
	}
	hi, e2 := c.Read(c.address_abs + 1)
	if addressingError(e2) {
		//
	}
	c.pc = Word(hi)<<8 | Word(lo)

	c.cycles = 8
}

func (c *CPU6502) fetch() byte {
	if !FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		f, e := c.Read(c.address_abs)
		if !addressingError(e) {
			//
		}
		c.fetched = f
	}
	return c.fetched
}

// ADDRESSING MODES

// Implicit address
func IMP(c *CPU6502) byte {
	c.fetched = c.a
	return 0
}

// Zero Page Addressing
func ZP0(c *CPU6502) byte {
	add, e := c.Read(c.pc)
	if !addressingError(e) {
		add = 0
	}
	c.address_abs = Word(add)
	c.pc++
	c.address_abs &= 0x00FF
	return 0
}

// Zero Page Adressing with Y
func ZPY(c *CPU6502) byte {
	add, e := c.Read(c.pc)
	if !addressingError(e) {
		add = 0
	}
	add += c.y
	c.address_abs = Word(add)
	c.pc++
	c.address_abs &= 0x00FF
	return 0
}

// Absolute addressing
func ABS(c *CPU6502) byte {
	lo, e := c.Read(c.pc)
	if !addressingError(e) {
		lo = 0
	}
	c.pc++
	hi, e2 := c.Read(c.pc)
	if !addressingError(e2) {
		hi = 0
	}
	c.pc++
	c.address_abs = (Word(hi) << 8) | Word(lo)
	return 0
}

// Absolute addressing with Y offset
func ABY(c *CPU6502) byte {
	lo, e := c.Read(c.pc)
	if !addressingError(e) {
		lo = 0
	}
	c.pc++
	hi, e2 := c.Read(c.pc)
	if !addressingError(e2) {
		hi = 0
	}
	c.pc++
	c.address_abs = (Word(hi) << 8) | Word(lo)
	c.address_abs += Word(c.y)
	// If Y added to address overflows, then page has
	// changed, therefore one extra cycle is needed
	if (c.address_abs & 0xFF00) != (Word(hi) << 8) {
		return 1
	}

	return 0
}

// Indirect Zero Page with X offset
func IZX(c *CPU6502) byte {
	t, e := c.Read(c.pc)
	if !addressingError(e) {
		t = 0
	}
	c.pc++

	lo, e2 := c.Read((Word(t) + Word(c.x)) & 0x00ff)
	if !addressingError(e2) {
		lo = 0
	}

	hi, e2 := c.Read((Word(t) + (Word(c.x) + 1)) & 0x00ff)
	if !addressingError(e2) {
		hi = 0
	}

	c.address_abs = (Word(hi) << 8) | Word(lo)
	return 0
}

// Immediate addressing
func IMM(c *CPU6502) byte {
	c.address_abs = c.pc
	c.pc++
	return 0
}

// Zero page addressing with X offset
func ZPX(c *CPU6502) byte {
	add, e := c.Read(c.pc)
	if !addressingError(e) {
		add = 0
	}
	add += c.x
	c.address_abs = Word(add)
	c.pc++
	c.address_abs &= 0x00FF
	return 0
}

// Relative addressing
func REL(c *CPU6502) byte {
	add, e := c.Read(c.pc)
	if !addressingError(e) {
		add = 0
	}
	c.pc++
	c.address_rel = Word(add)
	if c.address_rel&0x80 != 0 {
		c.address_rel += 0xFF00
	}

	return 0
}

// Absolute addressing with X
func ABX(c *CPU6502) byte {
	lo, e := c.Read(c.pc)
	if !addressingError(e) {
		lo = 0
	}
	c.pc++
	hi, e2 := c.Read(c.pc)
	if !addressingError(e2) {
		hi = 0
	}
	c.pc++
	c.address_abs = (Word(hi) << 8) | Word(lo)
	c.address_abs += Word(c.x)
	// If X added to address overflows, then page has
	// changed, therefore one extra cycle is needed
	if (c.address_abs & 0xFF00) != (Word(hi) << 8) {
		return 1
	}

	return 0
}

// Indirect Addressing
func IND(c *CPU6502) byte {
	ptr_lo, e := c.Read(c.pc)
	if !addressingError(e) {
		ptr_lo = 0
	}
	c.pc++
	ptr_hi, e2 := c.Read(c.pc)
	if !addressingError(e2) {
		ptr_hi = 0
	}
	c.pc++
	ptr := (Word(ptr_hi) << 8) | Word(ptr_lo)
	lo, e3 := c.Read(ptr + 0)
	if !addressingError(e3) {
		lo = 0
	}
	// Page Boundary Bug
	readAddress := ptr
	if ptr_lo == 0x00ff {
		readAddress &= 0x00FF
	} else {
		readAddress++
	}

	hi, e4 := c.Read(readAddress)
	if !addressingError(e4) {
		hi = 0
	}
	c.address_abs = (Word(hi) << 8) | Word(lo)
	return 0
}

// Indirect Zero Page with Y
func IZY(c *CPU6502) byte {
	t, e := c.Read(c.pc)
	if !addressingError(e) {
		t = 0
	}
	c.pc++

	lo, e2 := c.Read(Word(t) & 0x00ff)
	if !addressingError(e2) {
		lo = 0
	}

	hi, e2 := c.Read((Word(t) + 1) & 0x00ff)
	if !addressingError(e2) {
		hi = 0
	}

	c.address_abs = (Word(hi) << 8) | Word(lo)
	c.address_abs += Word(c.y)
	// If Y added to address overflows, then page has
	// changed, therefore one extra cycle is needed
	if (c.address_abs & 0xFF00) != (Word(hi) << 8) {
		return 1
	}
	return 0
}

func addressingError(e error) bool {
	if e != nil {
		log.Fatalf("%s", e)
		return true
	}
	return false
}

func logError(e error) {
	if e != nil {
		log.Panicf("%s", e)
	}
	e = nil
}

// Begin OPCODES //

// Invalid OpCode
func XXX(c *CPU6502) byte {
	//
	return 0
}

// Add memory to accumulator with carry
// Function -> A = A + M
// Flags -> C, Z, N, V
func ADC(c *CPU6502) byte {
	c.fetch()
	flagVal := c.StatusRegisterAsWord(C)

	temp := Word(c.a) + Word(c.fetched) + flagVal
	overflows := (^(Word(c.a) ^ Word(c.fetched)) & (Word(c.a) ^ Word(temp))) & 0x0080
	c.SetStatusRegisterFlag(C, temp > 255)
	c.SetStatusRegisterFlag(Z, (temp&0xFF00) == 0)
	c.SetStatusRegisterFlag(N, (temp&0x0080) != 0)
	c.SetStatusRegisterFlag(V, overflows != 0)
	c.a = byte(temp & 0x00FF)
	return 1
}

// Branch if Carry
func BCS(c *CPU6502) byte {
	if c.StatusRegister(C) {
		c.cycles++
		c.address_abs = c.pc + c.address_rel
		if (c.address_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.address_abs
	}
	return 0
}

// Branch if result not zero
func BNE(c *CPU6502) byte {
	if !c.StatusRegister(Z) {
		c.cycles++
		c.address_abs = c.pc + c.address_rel
		if (c.address_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.address_abs
	}
	return 0
}

// Branch on overflow set
func BVS(c *CPU6502) byte {
	if c.StatusRegister(V) {
		c.cycles++
		c.address_abs = c.pc + c.address_rel
		if (c.address_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.address_abs
	}
	return 0
}

// Clear Overflow Flag
func CLV(c *CPU6502) byte {
	c.SetStatusRegisterFlag(V, false)
	return 0
}

// Decrement value at memory location
func DEC(c *CPU6502) byte {
	c.fetch()
	temp := c.fetched - 1
	c.Write(c.address_abs, temp)
	c.SetFlagsZeroAndNegative(temp)
	return 0
}

// Increase value at memory location
func INC(c *CPU6502) byte {
	c.fetch()
	temp := c.fetched + 1
	c.Write(c.address_abs, temp)
	c.SetFlagsZeroAndNegative(temp)
	return 0
}

// Jump to Sub-Routine
// Function -> Push PC to stack, then pc = address
func JSR(c *CPU6502) byte {
	c.pc--
	c.Write(STACK+Word(c.stkp), byte(c.pc>>8))
	c.stkp--
	c.Write(STACK+Word(c.stkp), byte(c.pc&0x00ff))

	c.pc = c.address_abs
	return 0
}

// Shift one bit right - memory or accumulator
func LSR(c *CPU6502) byte {
	c.fetch()
	c.SetStatusRegisterFlag(C, (c.fetched&0x01) == 1)
	temp := c.fetched >> 1
	c.SetFlagsZeroAndNegative(temp)
	if FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		c.a = temp
	} else {
		c.Write(c.address_abs, temp)
	}
	return 0
}

// Instruction: Push Status Register to Stack
// Function:    status -> stack
// Note:        Break flag is set to 1 before push
func PHP(c *CPU6502) byte {
	c.SetStatusRegisterFlag(B, true)
	c.SetStatusRegisterFlag(U, true)
	c.Write(STACK+Word(c.stkp), c.status)
	c.SetStatusRegisterFlag(B, false)
	c.SetStatusRegisterFlag(U, false)
	c.stkp--
	return 0
}

// Rotate one bit Right
func ROR(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.fetched)>>1 | c.StatusRegisterAsWord(C)<<7
	c.SetFlagsZeroAndNegative(byte(temp & 0x00FF))
	c.SetStatusRegisterFlag(C, (temp&0xFF00) != 0)

	if FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		c.a = byte(temp & 0x00FF)
	} else {
		c.Write(c.address_abs, byte(temp&0x00FF))
	}
	return 0
}

// Instruction: Set Carry Flag
// Function:    C = 1
func SEC(c *CPU6502) byte {
	c.SetStatusRegisterFlag(C, true)
	return 0
}

// Instruction: Store X Register at Address
// Function:    M = X
func STX(c *CPU6502) byte {
	c.Write(c.address_abs, c.x)
	return 0
}

// Instruction: Transfer Stack Pointer to X Register
// Function:    X = stack pointer
// Flags Out:   N, Z
func TSX(c *CPU6502) byte {
	c.x = c.stkp
	c.SetFlagsZeroAndNegative(c.x)
	return 0
}

// AND operation
func AND(c *CPU6502) byte {
	c.fetch()
	c.a = c.a & c.fetched
	c.SetStatusRegisterFlag(Z, c.a == 0x00)
	c.SetStatusRegisterFlag(N, c.a&0x80 != 0)
	return 1
}

// Branch if Equal
func BEQ(c *CPU6502) byte {
	if c.StatusRegister(Z) {
		c.cycles++
		c.address_abs = c.pc + c.address_rel
		if (c.address_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.address_abs
	}
	return 0
}

// Branch if Positive
func BPL(c *CPU6502) byte {
	if !c.StatusRegister(N) {
		c.cycles++
		c.address_abs = c.pc + c.address_rel
		if (c.address_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.address_abs
	}
	return 0
}

// Clear Carry Bit
func CLC(c *CPU6502) byte {
	c.SetStatusRegisterFlag(C, false)
	return 0
}

// Compare Accumulator
func CMP(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.a) - Word(c.fetched)
	c.SetStatusRegisterFlag(C, c.a >= c.fetched)
	c.SetStatusRegisterFlag(Z, (temp&0x00FF) == 0x0000)
	c.SetStatusRegisterFlag(N, temp&0x0080 != 0)
	return 1
}

// Decrement value at X Register
func DEX(c *CPU6502) byte {
	c.x--
	c.SetStatusRegisterFlag(Z, c.x == 0x00)
	c.SetStatusRegisterFlag(N, c.x&0x80 != 0)
	return 0
}

// Increment value at X Register
func INX(c *CPU6502) byte {
	c.x++
	c.SetStatusRegisterFlag(Z, c.x == 0x00)
	c.SetStatusRegisterFlag(N, c.x&0x80 != 0)
	return 0
}

// Instruction: Load The Accumulator
// Function:    A = M
// Flags Out:   N, Z
func LDA(c *CPU6502) byte {
	c.fetch()
	c.a = c.fetched
	c.SetStatusRegisterFlag(Z, c.a == 0x00)
	c.SetStatusRegisterFlag(N, c.a&0x80 != 0)
	return 1
}

// No Operation
// There are multiple NOP, therefore this code is not as
// simple as it could be
func NOP(c *CPU6502) byte {
	// Sadly not all NOPs are equal, Ive added a few here
	// based on https://wiki.nesdev.com/w/index.php/CPU_unofficial_opcodes
	// and will add more based on game compatibility, and ultimately
	// I'd like to cover all illegal opcodes too
	switch c.opcode {
	case 0x1C:
	case 0x3C:
	case 0x5C:
	case 0x7C:
	case 0xDC:
	case 0xFC:
		return 1
	}
	return 0
}

// Pop from stack
func PLA(c *CPU6502) byte {
	c.stkp++
	var e error
	c.a, e = c.Read(Word(0x1<<3) + Word(c.stkp))
	if addressingError(e) {
		log.Fatalf("Error %q", e)
	}
	c.SetStatusRegisterFlag(Z, c.a == 0x00)
	c.SetStatusRegisterFlag(N, c.a&0x80 != 0)
	return 0
}

// Return from Interrupt
func RTI(c *CPU6502) byte {
	c.stkp++
	var e error
	c.status, e = c.Read(STACK + Word(c.stkp))
	if addressingError(e) {
		//
	}
	c.status &= ^byte(B)
	c.status &= ^byte(C)

	c.stkp++
	lo, e2 := c.Read(STACK + Word(c.stkp))
	if addressingError(e2) {
		//
	}
	c.stkp++
	hi, e3 := c.Read(STACK + Word(c.stkp))
	if addressingError(e3) {
		//
	}
	c.pc = (Word(hi) << 8) | Word(lo)
	return 0
}

// Instruction: Set Decimal Flag
// Function:    D = 1
func SED(c *CPU6502) byte {
	c.SetStatusRegisterFlag(D, true)
	return 0
}

// Instruction: Store Y Register at Address
// Function:    M = Y
func STY(c *CPU6502) byte {
	c.Write(c.address_abs, c.y)
	return 0
}

// Instruction: Transfer X Register to Accumulator
// Function:    A = X
// Flags Out:   N, Z
func TXA(c *CPU6502) byte {
	c.a = c.x
	c.SetFlagsZeroAndNegative(c.a)
	return 0
}

// Instruction: Arithmetic Shift Left
// Function:    A = C <- (A << 1) <- 0
// Flags Out:   N, Z, C
func ASL(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.fetched) << 1
	c.SetStatusRegisterFlag(C, (temp&0xFF00) != 0)
	c.SetStatusRegisterFlag(Z, (temp&0x00FF) == 0x00)
	c.SetStatusRegisterFlag(N, temp&0x80 != 0)
	if FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		c.a = byte(temp & 0x00FF)
	} else {
		c.Write(c.address_abs, byte(temp&0x00FF))
	}

	return 0
}

// Test Bits in memory with accumulator
func BIT(c *CPU6502) byte {
	c.fetch()
	temp := c.a & c.fetched
	c.SetStatusRegisterFlag(Z, (temp&0x00FF) == 0x00)
	c.SetStatusRegisterFlag(N, c.fetched&(1<<7) != 0)
	c.SetStatusRegisterFlag(V, c.fetched&(1<<6) != 0)
	return 0
}

// Instruction: Break
// Function:    Program Sourced Interrupt
func BRK(c *CPU6502) byte {
	c.pc++

	c.SetStatusRegisterFlag(I, true)
	c.Write(Word(1<<3+c.stkp), byte((c.pc>>8)&0x00FF))
	c.stkp--
	c.Write(Word(1<<3+c.stkp), byte(c.pc&0x00FF))
	c.stkp--

	c.SetStatusRegisterFlag(B, true)
	c.Write(Word(1<<3+c.stkp), c.status)
	c.stkp--
	c.SetStatusRegisterFlag(B, false)
	lo, e := c.Read(0xFFFE)
	if !addressingError(e) {

	}
	hi, e2 := c.Read(0xFFFF)
	if !addressingError(e2) {

	}
	c.pc = Word(hi)<<8 | Word(lo)
	return 0
}

// Clear Decimal Register
func CLD(c *CPU6502) byte {
	c.SetStatusRegisterFlag(D, false)
	return 0
}

// Compare X register
func CPX(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.x) - Word(c.fetched)
	c.SetStatusRegisterFlag(C, c.x >= c.fetched)
	c.SetStatusRegisterFlag(Z, (temp&0x00FF) == 0x0000)
	c.SetStatusRegisterFlag(N, temp&0x0080 != 0)
	return 1
}

// Decrement value from Y register
func DEY(c *CPU6502) byte {
	c.y--
	c.SetStatusRegisterFlag(Z, c.y == 0x00)
	c.SetStatusRegisterFlag(N, c.y&0x80 != 0)
	return 0
}

// Increment value at Y register
func INY(c *CPU6502) byte {
	c.y++
	c.SetStatusRegisterFlag(Z, c.y == 0x00)
	c.SetStatusRegisterFlag(N, c.y&0x80 != 0)
	return 0
}

// Instruction: Load The X Register
// Function:    X = M
// Flags Out:   N, Z
func LDX(c *CPU6502) byte {
	c.fetch()
	c.x = c.fetched
	c.SetFlagsZeroAndNegative(c.x)
	return 1
}

// Instruction: Bitwise Logic OR
// Function:    A = A | M
// Flags Out:   N, Z
func ORA(c *CPU6502) byte {
	c.fetch()
	c.a |= c.fetched
	c.SetFlagsZeroAndNegative(c.a)
	return 1
}

// Instruction: Pop Status Register off Stack
// Function:    Status <- stack
func PLP(c *CPU6502) byte {
	c.stkp++
	var e error
	c.status, e = c.Read(STACK + Word(c.stkp))
	if !addressingError(e) {
		//
	}
	c.SetStatusRegisterFlag(U, true)
	return 0
}

// Return from sub routine
func RTS(c *CPU6502) byte {
	c.stkp++

	lo, e := c.Read(STACK + Word(c.stkp))
	if !addressingError(e) {
		//
	}
	c.stkp++
	hi, e2 := c.Read(STACK + Word(c.stkp))
	if !addressingError(e2) {
		//
	}
	c.pc = (Word(hi) << 8) | Word(lo)
	c.pc++
	return 0
}

// Instruction: Set Interrupt Flag / Enable Interrupts
// Function:    I = 1
func SEI(c *CPU6502) byte {
	c.SetStatusRegisterFlag(I, true)
	return 0
}

// Instruction: Transfer Accumulator to Y Register
// Function:    Y = A
// Flags Out:   N, Z
func TAX(c *CPU6502) byte {
	c.y = c.a
	c.SetFlagsZeroAndNegative(c.y)
	return 0
}
func TXS(c *CPU6502) byte { return 0 }

// Branch if Carry Clear
func BCC(c *CPU6502) byte {
	if !c.StatusRegister(C) {
		c.cycles++
		c.address_abs = c.pc + c.address_rel
		if (c.address_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.address_abs
	}
	return 0
}

// Branch if Negative
func BMI(c *CPU6502) byte {
	if c.StatusRegister(C) {
		c.cycles++
		c.address_abs = c.pc + c.address_rel
		if (c.address_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.address_abs
	}
	return 0
}

// Branch if Overflows
func BVC(c *CPU6502) byte {
	if !c.StatusRegister(V) {
		c.cycles++
		c.address_abs = c.pc + c.address_rel
		if (c.address_abs & 0xFF00) != (c.pc & 0xFF00) {
			c.cycles++
		}

		c.pc = c.address_abs
	}
	return 0
}

// Clear Interrup Flag
func CLI(c *CPU6502) byte {
	c.SetStatusRegisterFlag(I, false)
	return 0
}

// Compare Y Register
func CPY(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.y) - Word(c.fetched)
	c.SetStatusRegisterFlag(C, c.y >= c.fetched)
	c.SetStatusRegisterFlag(Z, (temp&0x00FF) == 0x0000)
	c.SetStatusRegisterFlag(N, temp&0x0080 != 0)
	return 1
}

// Bitwise XOR
func EOR(c *CPU6502) byte {
	c.fetch()
	c.a = c.a ^ c.fetched
	c.SetStatusRegisterFlag(Z, c.a == 0x00)
	c.SetStatusRegisterFlag(N, c.a&0x80 != 0)
	return 1
}

// Jump to location
// Function -> pc = address
func JMP(c *CPU6502) byte {
	c.pc = c.address_abs
	return 0
}

// Instruction: Load The Y Register
// Function:    Y = M
// Flags Out:   N, Z
func LDY(c *CPU6502) byte {
	c.fetch()
	c.y = c.fetched

	c.SetStatusRegisterFlag(Z, c.y == 0x00)
	c.SetStatusRegisterFlag(N, c.y&0x80 != 0)
	return 1
}

// Push accumulator to Stack
func PHA(c *CPU6502) byte {
	c.Write(STACK+Word(c.stkp), c.a)
	c.stkp--
	return 0
}

// Rotate One Bit Left (memory or accumulator)
func ROL(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.fetched)<<1 | c.StatusRegisterAsWord(C)

	c.SetStatusRegisterFlag(C, (temp&0xFF00) != 0)
	c.SetFlagsZeroAndNegative(byte(temp & 0x00FF))
	if FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		c.a = byte(temp & 0x00FF)
	} else {
		c.Write(c.address_abs, byte(temp&0x00FF))
	}
	return 0
}

// Subtract Operation
func SBC(c *CPU6502) byte {
	c.fetch()
	flagVal := c.StatusRegisterAsWord(C)
	value := Word(c.fetched) ^ 0x00FF
	temp := Word(c.a) + value + flagVal
	c.SetStatusRegisterFlag(C, (temp&0xFF00) != 0)
	c.SetStatusRegisterFlag(V, ((temp^Word(c.a))&(temp^value)&0x0080) != 0)
	c.SetStatusRegisterFlag(Z, (temp&0xFF00) == 0)
	c.SetStatusRegisterFlag(N, (temp&0x0080) != 0)
	c.a = byte(temp & 0x00FF)
	return 1
}

// Store Accumulator at Address
func STA(c *CPU6502) byte {
	c.Write(c.address_abs, c.a)
	return 0
}

// Instruction: Transfer Accumulator to Y Register
// Function:    Y = A
// Flags Out:   N, Z
func TAY(c *CPU6502) byte {
	c.y = c.a
	c.SetFlagsZeroAndNegative(c.y)
	return 0
}

// Instruction: Transfer Y Register to Accumulator
// Function:    A = Y
// Flags Out:   N, Z
func TYA(c *CPU6502) byte {
	c.a = c.y
	c.SetFlagsZeroAndNegative(c.a)
	return 0
}

// This is the disassembly function. Its workings are not required for emulation.
// It is merely a convenience function to turn the binary instruction code into
// human readable form. Its included as part of the emulator because it can take
// advantage of many of the CPUs internal operations to do this.
func (c *CPU6502) Disassemble(start, stop Word) map[Word]string {
	var value byte = 0
	var lo byte = 0
	var hi byte = 0
	var line_addr Word = 0
	var addr uint32 = uint32(start)
	bus := c.bus
	mapLines := make(map[Word]string)

	// A convenient utility to convert variables into
	// hex strings because "modern C++"'s method with
	// streams is atrocious
	hex := func(n uint32, d int) string {
		s := fmt.Sprintf("%X", n)
		var b bytes.Buffer
		if len(s) < d {
			for i := d - len(s); i > 0; i-- {
				b.WriteByte('0')
			}
		}
		b.WriteString(s)
		return b.String()
	}

	// Starting at the specified address we read an instruction
	// byte, which in turn yields information from the lookup table
	// as to how many additional bytes we need to read and what the
	// addressing mode is. I need this info to assemble human readable
	// syntax, which is different depending upon the addressing mode

	// As the instruction is decoded, a std::string is assembled
	// with the readable output
	for addr <= uint32(stop) {
		line_addr = Word(addr)
		var sInst bytes.Buffer
		var e error
		var opcode byte
		// Prefix line with instruction address
		sInst.WriteString("$")
		sInst.WriteString(hex(addr, 4))
		sInst.WriteString(": ")

		// Read instruction, and get its readable name
		opcode, e = c.bus.Read(Word(addr), true)
		logError(e)
		addr++

		sInst.WriteString(OpCodesLookupTable[opcode].name)
		sInst.WriteByte(' ')

		// Get oprands from desired locations, and form the
		// instruction based upon its addressing mode. These
		// routines mimmick the actual fetch routine of the
		// 6502 in order to get accurate data as part of the
		// instruction
		if FnEquals(OpCodesLookupTable[opcode].addressmode, IMP) {
			sInst.WriteString("{IMP}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, IMM) {
			var value byte
			value, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			sInst.WriteString("#$")
			sInst.WriteString(hex(uint32(value), 2))
			sInst.WriteString(" {IMM}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ZP0) {
			lo, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(hex(uint32(lo), 2))
			sInst.WriteString(" {ZP0}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ZPX) {
			lo, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(hex(uint32(lo), 2))
			sInst.WriteString(", X {ZPX}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ZPY) {
			lo, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(hex(uint32(lo), 2))
			sInst.WriteString(", Y {ZPY}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, IZX) {
			lo, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(hex(uint32(lo), 2))
			sInst.WriteString(", X {IZX}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, IZY) {
			lo, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(hex(uint32(lo), 2))
			sInst.WriteString(", Y {IZY}")

		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ABS) {
			lo, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			hi, e = bus.Read(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("$")
			sInst.WriteString(hex(uint32(Word(hi)<<8|Word(lo)), 4))
			sInst.WriteString(" {ABS}")

		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ABX) {
			lo, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			hi, e = bus.Read(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("$")
			sInst.WriteString(hex(uint32(Word(hi)<<8|Word(lo)), 4))
			sInst.WriteString(", X {ABX}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ABY) {
			lo, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			hi, e = bus.Read(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("$")
			sInst.WriteString(hex(uint32(Word(hi)<<8|Word(lo)), 4))
			sInst.WriteString(", Y {ABY}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, IND) {
			lo, e = bus.Read(Word(addr), true)
			logError(e)
			addr++
			hi, e = bus.Read(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("($")
			sInst.WriteString(hex(uint32(Word(hi)<<8|Word(lo)), 4))
			sInst.WriteString(") {IND}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, REL) {
			var val byte
			val, e = bus.Read(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("$")
			sInst.WriteString(hex(uint32(val), 4))
			sInst.WriteString(" [$")
			sInst.WriteString(hex(addr+uint32(value), 4))
			sInst.WriteString("] {REL}")
		}

		// Add the formed string to a std::map, using the instruction's
		// address as the key. This makes it convenient to look for later
		// as the instructions are variable in length, so a straight up
		// incremental index is not sufficient.
		mapLines[line_addr] = sInst.String()
	}

	return mapLines
}
