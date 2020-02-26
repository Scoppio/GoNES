package main

import (
	"bytes"
	"log"
)

const (
	// C : Carry flag
	C = 0
	// Z : Zero flag
	Z = 1
	// I : flag
	I = 2
	// D : Decimal flag
	D = 3
	// B : Break flag
	B = 4
	// U : unised flag
	U = 5
	// V : Overflow flag
	V = 6
	// N : Negative flag
	N = 7
	// Stack : Stack memory address
	Stack = Word(0x0100)
)

var (
	// OpCodesLookupTable : table with all the instructions of the 6502
	OpCodesLookupTable []Instruction
)

// CPU6502 : Struct that represents the 6502 chip
type CPU6502 struct {
	a, x, y, stkp, status, fetched, opcode, cycles byte
	pc, addressAbs, addressRel                     Word
	bus                                            *Bus
}

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

// CreateCPU : creates a new CPU
func CreateCPU() *CPU6502 {
	return &CPU6502{}
}

// StatusRegister : Checks the state of the given flag
func (c *CPU6502) StatusRegister(flag Flag) bool {
	val := c.status & (1 << uint(flag))
	return (val > 0)
}

// StatusRegisterAsWord : Checks the state of the given flag in a Word format
func (c *CPU6502) StatusRegisterAsWord(flag Flag) Word {
	return Word(c.status & (1 << uint(flag)))
}

// SetStatusRegisterFlag : Sets the state of the given flag
func (c *CPU6502) SetStatusRegisterFlag(flag Flag, val bool) {
	if val {
		c.status |= byte(1 << uint(flag))
	} else {
		c.status &= ^byte(1 << uint(flag))
	}
}

// SetFlagsZeroAndNegative : Sets flags Z and N for a given value
func (c *CPU6502) SetFlagsZeroAndNegative(val byte) {
	c.SetStatusRegisterFlag(Z, val == 0x00)
	c.SetStatusRegisterFlag(N, (val&0x80) != 0)
}

// ConnectBus : connects the CPU to the Bus
func (c *CPU6502) ConnectBus(bus *Bus) {
	c.bus = bus
}

// CPUReadStack : reads memory position at stack
func (c *CPU6502) CPUReadStack(address Word) (byte, error) {
	return c.bus.CPURead(Stack+address, false)
}

// CPURead : reads memory position from the CPU
func (c *CPU6502) CPURead(address Word) (byte, error) {
	return c.bus.CPURead(address, false)
}

// CPUWriteStack : writes memory position at stack
func (c *CPU6502) CPUWriteStack(address Word, data byte) error {
	return c.bus.CPUWrite(Stack+address, data)
}

// CPUWrite : writes memory position from the CPU
func (c *CPU6502) CPUWrite(address Word, data byte) error {
	return c.bus.CPUWrite(address, data)
}

// Clock : Does a single clock which will execute an instruction when reaches 0
func (c *CPU6502) Clock() {
	// execute
	if c.cycles == byte(0x00) {

		c.opcode, _ = c.CPURead(c.pc)
		c.SetStatusRegisterFlag(U, true)

		c.pc++

		c.cycles = OpCodesLookupTable[c.opcode].cycles
		additionalCycle1 := OpCodesLookupTable[c.opcode].addressmode(c)
		additionalCycle2 := OpCodesLookupTable[c.opcode].operate(c)
		c.cycles += (additionalCycle1 & additionalCycle2)
		// Always set the unused status flag bit to 1
		c.SetStatusRegisterFlag(U, true)
		OperationCount++
	}

	c.cycles--
}

// Complete : Checks if the cycle has reached 0
func (c *CPU6502) Complete() bool {
	return c.cycles == byte(0x00)
}

// Reset : Resets the CPU and puts it in a known state
func (c *CPU6502) Reset() {
	// Get start address
	c.addressAbs = 0xFFFC
	lo, _ := c.CPURead(c.addressAbs + 0)
	hi, _ := c.CPURead(c.addressAbs + 1)
	// Set it to PC
	c.pc = Word(hi)<<8 | Word(lo)

	c.a = 0
	c.x = 0
	c.y = 0
	c.stkp = 0xFD
	c.status = 0x00
	c.SetStatusRegisterFlag(U, true)

	// Clear internal helper variables
	c.addressRel = 0x0000
	c.addressAbs = 0x0000
	c.fetched = 0x00

	// Reset takes time
	c.cycles = 8
}

// InterruptRequest : Sets the system in a state to execute code from an interruption
func (c *CPU6502) InterruptRequest() {
	if !c.StatusRegister(I) {
		c.CPUWriteStack(Word(c.stkp), byte((c.pc>>8)&0x00ff))
		c.stkp--
		c.CPUWriteStack(Word(c.stkp), byte(c.pc&0x00FF))
		c.stkp--

		c.SetStatusRegisterFlag(B, false)
		c.SetStatusRegisterFlag(U, true)
		c.SetStatusRegisterFlag(I, true)
		c.CPUWriteStack(Word(c.stkp), c.status)
		c.stkp--

		c.addressAbs = 0xFFFE
		lo, _ := c.CPURead(c.addressAbs + 0)
		hi, _ := c.CPURead(c.addressAbs + 1)
		c.pc = Word(hi)<<8 | Word(lo)

		c.cycles = 7
	}
}

// NonMaskableInterruptRequest : Sets the system in a state to execute code from an interruption
func (c *CPU6502) NonMaskableInterruptRequest() {
	c.CPUWriteStack(Word(c.stkp), byte((c.pc>>8)&0x00ff))
	c.stkp--
	c.CPUWriteStack(Word(c.stkp), byte(c.pc&0x00FF))
	c.stkp--

	c.SetStatusRegisterFlag(B, false)
	c.SetStatusRegisterFlag(U, true)
	c.SetStatusRegisterFlag(I, true)
	c.CPUWriteStack(Word(c.stkp), c.status)
	c.stkp--

	c.addressAbs = 0xFFFA
	lo, _ := c.CPURead(c.addressAbs + 0)
	hi, _ := c.CPURead(c.addressAbs + 1)
	c.pc = Word(hi)<<8 | Word(lo)

	c.cycles = 8
}

func (c *CPU6502) fetch() byte {
	if !FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		f, _ := c.CPURead(c.addressAbs)
		c.fetched = f
	}
	return c.fetched
}

// ADDRESSING MODES

// IMP : Implicit address
func IMP(c *CPU6502) byte {
	c.fetched = c.a
	return 0
}

// ZP0 : Zero Page Addressing
func ZP0(c *CPU6502) byte {
	add, _ := c.CPURead(c.pc)
	c.addressAbs = Word(add)
	c.pc++
	c.addressAbs &= 0x00FF
	return 0
}

// ZPY : Zero Page Adressing with Y
func ZPY(c *CPU6502) byte {
	add, _ := c.CPURead(c.pc)
	c.addressAbs = Word(add) + Word(c.y)
	c.pc++
	c.addressAbs &= 0x00FF
	return 0
}

// ABS : Absolute addressing
func ABS(c *CPU6502) byte {
	lo, _ := c.CPURead(c.pc)
	c.pc++
	hi, _ := c.CPURead(c.pc)
	c.pc++
	c.addressAbs = (Word(hi) << 8) | Word(lo)
	return 0
}

// ABY : Absolute addressing with Y offset
func ABY(c *CPU6502) byte {
	lo, _ := c.CPURead(c.pc)
	c.pc++
	hi, _ := c.CPURead(c.pc)
	c.pc++

	c.addressAbs = (Word(hi) << 8) | Word(lo)
	c.addressAbs += Word(c.y)
	// If Y added to address overflows, then page has
	// changed, therefore one extra cycle is needed
	if (c.addressAbs & 0xFF00) != (Word(hi) << 8) {
		return 1
	}

	return 0
}

// IZX : Indirect Zero Page with X offset
func IZX(c *CPU6502) byte {
	t, _ := c.CPURead(c.pc)
	c.pc++
	lo, _ := c.CPURead((Word(t) + Word(c.x)) & 0x00ff)
	hi, _ := c.CPURead((Word(t) + (Word(c.x) + 1)) & 0x00ff)

	c.addressAbs = (Word(hi) << 8) | Word(lo)
	return 0
}

// IMM : Immediate addressing
func IMM(c *CPU6502) byte {
	c.addressAbs = c.pc
	c.pc++
	return 0
}

// ZPX : Zero page addressing with X offset
func ZPX(c *CPU6502) byte {
	add, _ := c.CPURead(c.pc)
	c.addressAbs = Word(add) + Word(c.x)
	c.pc++
	c.addressAbs &= 0x00FF
	return 0
}

// REL : Relative addressing
func REL(c *CPU6502) byte {
	add, _ := c.CPURead(c.pc)
	c.addressRel = Word(add)
	c.pc++
	if c.addressRel&0x80 != 0 {
		c.addressRel |= 0xFF00
	}

	return 0
}

// ABX : Absolute addressing with X
func ABX(c *CPU6502) byte {
	lo, _ := c.CPURead(c.pc)
	c.pc++
	hi, _ := c.CPURead(c.pc)

	c.pc++
	c.addressAbs = (Word(hi) << 8) | Word(lo)
	c.addressAbs += Word(c.x)
	// If X added to address overflows, then page has
	// changed, therefore one extra cycle is needed
	if (c.addressAbs & 0xFF00) != (Word(hi) << 8) {
		return 1
	}

	return 0
}

// IND : Indirect Addressing
func IND(c *CPU6502) byte {
	pointerLow, _ := c.CPURead(c.pc)
	c.pc++
	pointerHigh, _ := c.CPURead(c.pc)
	c.pc++
	ptr := (Word(pointerHigh) << 8) | Word(pointerLow)
	lo, _ := c.CPURead(ptr + 0)

	// Page Boundary Bug
	readAddress := ptr

	if pointerLow == 0x00ff {
		readAddress &= 0x00FF
	} else {
		readAddress++
	}

	hi, _ := c.CPURead(readAddress)
	c.addressAbs = (Word(hi) << 8) | Word(lo)
	return 0
}

// IZY : Indirect Zero Page with Y
func IZY(c *CPU6502) byte {
	t, _ := c.CPURead(c.pc)
	c.pc++

	lo, _ := c.CPURead(Word(t) & 0x00ff)
	hi, _ := c.CPURead((Word(t) + 1) & 0x00ff)

	c.addressAbs = (Word(hi) << 8) | Word(lo)
	c.addressAbs += Word(c.y)
	// If Y added to address overflows, then page has
	// changed, therefore one extra cycle is needed
	if (c.addressAbs & 0xFF00) != (Word(hi) << 8) {
		return 1
	}

	return 0
}

// Begin OPCODES //

// XXX : Invalid OpCode
func XXX(c *CPU6502) byte {
	//
	return 0
}

// ADC : Add memory to accumulator with carry
// Function -> A = A + M
// Flags -> C, Z, N, V
func ADC(c *CPU6502) byte {
	f := Word(c.fetch())
	flagVal := c.StatusRegisterAsWord(C)

	temp := Word(c.a) + f + flagVal
	overflows := (^(Word(c.a) ^ Word(c.fetched)) & (Word(c.a) ^ Word(temp))) & 0x0080
	c.SetStatusRegisterFlag(C, temp > 255)
	c.SetStatusRegisterFlag(Z, (temp&0xFF00) == 0)
	c.SetStatusRegisterFlag(N, (temp&0x0080) != 0)
	c.SetStatusRegisterFlag(V, overflows != 0)
	c.a = byte(temp & 0x00FF)
	return 1
}

// BCS : Branch if Carry
func BCS(c *CPU6502) byte {
	if c.StatusRegister(C) {
		jump(c)
	}
	return 0
}

// BNE : Branch if result not zero
func BNE(c *CPU6502) byte {
	if !c.StatusRegister(Z) {
		jump(c)
	}
	return 0
}

func jump(c *CPU6502) {
	c.cycles++
	c.addressAbs = c.pc + c.addressRel
	if (c.addressAbs & 0xFF00) != (c.pc & 0xFF00) {
		c.cycles++
	}

	c.pc = c.addressAbs
}

// BVS : Branch on overflow set
func BVS(c *CPU6502) byte {
	if c.StatusRegister(V) {
		jump(c)
	}
	return 0
}

// CLV : Clear Overflow Flag
func CLV(c *CPU6502) byte {
	c.SetStatusRegisterFlag(V, false)
	return 0
}

// DEC : Decrement value at memory location
func DEC(c *CPU6502) byte {
	c.fetch()
	temp := c.fetched - 1
	c.CPUWrite(c.addressAbs, temp)
	c.SetFlagsZeroAndNegative(temp)
	return 0
}

// INC : Increase value at memory location
func INC(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.fetched) + 1
	c.CPUWrite(c.addressAbs, byte(temp&0x00FF))
	c.SetFlagsZeroAndNegative(byte(temp & 0x00FF))
	return 0
}

// JSR : Jump to Sub-Routine
// Function -> Push PC to Stack, then pc = address
func JSR(c *CPU6502) byte {
	c.pc--
	c.CPUWriteStack(Word(c.stkp), byte(c.pc>>8))
	c.stkp--
	c.CPUWriteStack(Word(c.stkp), byte(c.pc))
	c.stkp--
	c.pc = c.addressAbs
	return 0
}

// LSR : Shift one bit right - memory or accumulator
func LSR(c *CPU6502) byte {
	c.fetch()
	c.SetStatusRegisterFlag(C, (c.fetched&0x01) == 1)
	temp := c.fetched >> 1
	c.SetFlagsZeroAndNegative(temp)
	if FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		c.a = temp
	} else {
		c.CPUWrite(c.addressAbs, temp)
	}
	return 0
}

// PHP : Instruction: Push Status Register to Stack
// Function:    status -> Stack
// Note:        Break flag is set to 1 before push
func PHP(c *CPU6502) byte {
	c.SetStatusRegisterFlag(B, true)
	c.SetStatusRegisterFlag(U, true)
	c.CPUWriteStack(Word(c.stkp), c.status)
	c.SetStatusRegisterFlag(B, false)
	c.SetStatusRegisterFlag(U, false)
	c.stkp--
	return 0
}

// ROR : Rotate one bit Right
func ROR(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.fetched)>>1 | c.StatusRegisterAsWord(C)<<7
	c.SetStatusRegisterFlag(C, (temp&0x0001) != 0)
	c.SetFlagsZeroAndNegative(byte(temp & 0x00FF))

	if FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		c.a = byte(temp & 0x00FF)
	} else {
		c.CPUWrite(c.addressAbs, byte(temp&0x00FF))
	}
	return 0
}

// SEC : Instruction: Set Carry Flag
// Function:    C = 1
func SEC(c *CPU6502) byte {
	c.SetStatusRegisterFlag(C, true)
	return 0
}

// STX : Instruction: Store X Register at Address
// Function:    M = X
func STX(c *CPU6502) byte {
	c.CPUWrite(c.addressAbs, c.x)
	return 0
}

// TSX : Instruction: Transfer Stack Pointer to X Register
// Function:    X = Stack pointer
// Flags Out:   N, Z
func TSX(c *CPU6502) byte {
	c.x = c.stkp
	c.SetFlagsZeroAndNegative(c.x)
	return 0
}

// AND : AND operation
// Function:    A = A & M
// Flags Out:   N, Z
func AND(c *CPU6502) byte {
	c.a = c.a & c.fetch()
	c.SetFlagsZeroAndNegative(c.a)
	return 1
}

// BEQ : Branch if Equal
func BEQ(c *CPU6502) byte {
	if c.StatusRegister(Z) {
		jump(c)
	}
	return 0
}

// BPL : Branch if Positive
func BPL(c *CPU6502) byte {
	if !c.StatusRegister(N) {
		jump(c)
	}
	return 0
}

// CLC : Clear Carry Bit
func CLC(c *CPU6502) byte {
	c.SetStatusRegisterFlag(C, false)
	return 0
}

// CMP : Compare Accumulator
func CMP(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.a) - Word(c.fetched)
	c.SetStatusRegisterFlag(C, c.a >= c.fetched)
	c.SetFlagsZeroAndNegative(byte(temp))
	return 1
}

// DEX : Decrement value at X Register
func DEX(c *CPU6502) byte {
	c.x--
	c.SetFlagsZeroAndNegative(c.x)
	return 0
}

// INX : Increment value at X Register
func INX(c *CPU6502) byte {
	c.x++
	c.SetFlagsZeroAndNegative(c.x)
	return 0
}

// LDA : Instruction: Load The Accumulator
// Function:    A = M
// Flags Out:   N, Z
func LDA(c *CPU6502) byte {
	c.a = c.fetch()
	c.SetFlagsZeroAndNegative(c.a)
	return 1
}

// NOP : No Operation
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

// PLA : Pop from Stack
func PLA(c *CPU6502) byte {
	c.stkp++

	c.a, _ = c.CPUReadStack(Word(c.stkp))
	c.SetFlagsZeroAndNegative(c.a)
	return 0
}

// RTI : Return from Interrupt
func RTI(c *CPU6502) byte {
	c.stkp++
	c.status, _ = c.CPUReadStack(Word(c.stkp))

	c.status &= ^(byte(1) << B)
	c.status &= ^(byte(1) << U)

	c.stkp++
	lo, _ := c.CPUReadStack(Word(c.stkp))
	c.stkp++
	hi, _ := c.CPUReadStack(Word(c.stkp))
	c.pc = (Word(hi) << 8) | Word(lo)
	return 0
}

// SED : Instruction: Set Decimal Flag
// Function:    D = 1
func SED(c *CPU6502) byte {
	c.SetStatusRegisterFlag(D, true)
	return 0
}

// STY : Instruction: Store Y Register at Address
// Function:    M = Y
func STY(c *CPU6502) byte {
	c.CPUWrite(c.addressAbs, c.y)
	return 0
}

// TXA : Instruction: Transfer X Register to Accumulator
// Instruction: Transfer X Register to Accumulator
// Function:    A = X
// Flags Out:   N, Z
func TXA(c *CPU6502) byte {
	c.a = c.x
	c.SetFlagsZeroAndNegative(c.a)
	return 0
}

// ASL : Instruction: Arithmetic Shift Left
// Function:    A = C <- (A << 1) <- 0
// Flags Out:   N, Z, C
func ASL(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.fetched) << 1
	c.SetStatusRegisterFlag(C, (temp&0xFF00) != 0)
	c.SetFlagsZeroAndNegative(byte(temp))
	if FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		c.a = byte(temp & 0x00FF)
	} else {
		c.CPUWrite(c.addressAbs, byte(temp&0x00FF))
	}

	return 0
}

// BIT : Test Bits in memory with accumulator
func BIT(c *CPU6502) byte {
	c.fetch()
	temp := c.a & c.fetched
	c.SetStatusRegisterFlag(Z, (temp&0x00FF) == 0x00)
	c.SetStatusRegisterFlag(N, c.fetched&(1<<7) != 0)
	c.SetStatusRegisterFlag(V, c.fetched&(1<<6) != 0)
	return 0
}

// BRK : Instruction: Break
// Function: Program Sourced Interrupt
func BRK(c *CPU6502) byte {
	c.pc++

	c.SetStatusRegisterFlag(I, true)
	c.CPUWrite(Word(0x0100)+Word(c.stkp), byte((c.pc>>8)&0x00FF))
	c.stkp--
	c.CPUWrite(Word(0x0100)+Word(c.stkp), byte(c.pc&0x00FF))
	c.stkp--

	c.SetStatusRegisterFlag(B, true)
	c.CPUWrite(Word(0x0100)+Word(c.stkp), c.status)
	c.stkp--
	c.SetStatusRegisterFlag(B, false)
	lo, _ := c.CPURead(0xFFFE)
	hi, _ := c.CPURead(0xFFFF)

	c.pc = Word(hi)<<8 | Word(lo)
	return 0
}

// CLD : Clear Decimal Register
func CLD(c *CPU6502) byte {
	c.SetStatusRegisterFlag(D, false)
	return 0
}

// CPX : Compare X register
func CPX(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.x) - Word(c.fetched)
	c.SetStatusRegisterFlag(C, c.x >= c.fetched)
	c.SetFlagsZeroAndNegative(byte(temp))
	return 0
}

// DEY : Decrement value from Y register
func DEY(c *CPU6502) byte {
	c.y--
	c.SetFlagsZeroAndNegative(c.y)
	return 0
}

// INY : Increment value at Y register
func INY(c *CPU6502) byte {
	c.y++
	c.SetFlagsZeroAndNegative(c.y)
	return 0
}

// LDX : Instruction: Load The X Register
// Function:    X = M
// Flags Out:   N, Z
func LDX(c *CPU6502) byte {
	c.x = c.fetch()
	c.SetFlagsZeroAndNegative(c.x)
	return 1
}

// ORA : Instruction: Bitwise Logic OR
// Function:    A = A | M
// Flags Out:   N, Z
func ORA(c *CPU6502) byte {
	c.fetch()
	c.a |= c.fetched
	c.SetFlagsZeroAndNegative(c.a)
	return 1
}

// PLP : Instruction: Pop Status Register off Stack
// Function:    Status <- Stack
func PLP(c *CPU6502) byte {
	c.stkp++
	c.status, _ = c.CPUReadStack(Word(c.stkp))
	c.SetStatusRegisterFlag(U, true)
	return 0
}

// RTS : Return from sub routine
func RTS(c *CPU6502) byte {
	c.stkp++
	lo, _ := c.CPUReadStack(Word(c.stkp))
	c.stkp++
	hi, _ := c.CPUReadStack(Word(c.stkp))
	c.pc = Word(hi)<<8 | Word(lo)
	c.pc++
	return 0
}

// SEI : Instruction: Set Interrupt Flag / Enable Interrupts
// Function:    I = 1
func SEI(c *CPU6502) byte {
	c.SetStatusRegisterFlag(I, true)
	return 0
}

// TAX : Instruction: Transfer Accumulator to Y Register
// Function:    X = A
// Flags Out:   N, Z
func TAX(c *CPU6502) byte {
	c.x = c.a
	c.SetFlagsZeroAndNegative(c.x)
	return 0
}

// TXS : Instruction: Transfer Stack Pointer to X Register
// Function:    Stack pointer = X
func TXS(c *CPU6502) byte {
	c.stkp = c.x
	return 0
}

// BCC : Branch if Carry Clear
func BCC(c *CPU6502) byte {
	if !c.StatusRegister(C) {
		jump(c)
	}
	return 0
}

// BMI : Branch if Negative
func BMI(c *CPU6502) byte {
	if c.StatusRegister(N) {
		jump(c)
	}
	return 0
}

// BVC : Branch if Overflows
func BVC(c *CPU6502) byte {
	if !c.StatusRegister(V) {
		jump(c)
	}
	return 0
}

// CLI : Clear Interrup Flag
func CLI(c *CPU6502) byte {
	c.SetStatusRegisterFlag(I, false)
	return 0
}

// CPY : Compare Y Register
func CPY(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.y) - Word(c.fetched)
	c.SetStatusRegisterFlag(C, c.y >= c.fetched)
	c.SetFlagsZeroAndNegative(byte(temp))
	return 0
}

// EOR : Bitwise XOR
func EOR(c *CPU6502) byte {
	c.fetch()
	c.a = c.a ^ c.fetched
	c.SetFlagsZeroAndNegative(c.a)
	return 1
}

// JMP : Jump to location
// Function -> pc = address
func JMP(c *CPU6502) byte {
	c.pc = c.addressAbs
	return 0
}

// LDY : Instruction: Load The Y Register
// Function:    Y = M
// Flags Out:   N, Z
func LDY(c *CPU6502) byte {
	c.y = c.fetch()
	c.SetFlagsZeroAndNegative(c.y)
	return 1
}

// PHA : Push accumulator to Stack
func PHA(c *CPU6502) byte {
	c.CPUWriteStack(Word(c.stkp), c.a)
	c.stkp--
	return 0
}

// ROL : Rotate One Bit Left (memory or accumulator)
func ROL(c *CPU6502) byte {
	c.fetch()
	temp := Word(c.fetched)<<1 | c.StatusRegisterAsWord(C)

	c.SetStatusRegisterFlag(C, (temp&0xFF00) != 0)
	c.SetFlagsZeroAndNegative(byte(temp & 0x00FF))
	if FnEquals(OpCodesLookupTable[c.opcode].addressmode, IMP) {
		c.a = byte(temp & 0x00FF)
	} else {
		c.CPUWrite(c.addressAbs, byte(temp&0x00FF))
	}
	return 0
}

// SBC : Subtract Operation
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

// STA : Store Accumulator at Address
func STA(c *CPU6502) byte {
	c.CPUWrite(c.addressAbs, c.a)
	return 0
}

// TAY : Instruction: Transfer Accumulator to Y Register
// Function:    Y = A
// Flags Out:   N, Z
func TAY(c *CPU6502) byte {
	c.y = c.a
	c.SetFlagsZeroAndNegative(c.y)
	return 0
}

// TYA : Instruction: Transfer Y Register to Accumulator
// Function:    A = Y
// Flags Out:   N, Z
func TYA(c *CPU6502) byte {
	c.a = c.y
	c.SetFlagsZeroAndNegative(c.a)
	return 0
}

func logError(e error) {
	if e != nil {
		log.Panicf("%s", e)
	}
	e = nil
}

// Disassemble : This is the disassembly function. Its workings are not required for emulation.
// It is merely a convenience function to turn the binary instruction code into
// human readable form. Its included as part of the emulator because it can take
// advantage of many of the CPUs internal operations to do this.
func (c *CPU6502) Disassemble(start, stop Word) map[Word]string {
	var lo byte = 0
	var hi byte = 0
	var lineAddr Word = 0
	var addr uint32 = uint32(start)
	bus := c.bus
	mapLines := make(map[Word]string)

	// Starting at the specified address we read an instruction
	// byte, which in turn yields information from the lookup table
	// as to how many additional bytes we need to read and what the
	// addressing mode is. I need this info to assemble human readable
	// syntax, which is different depending upon the addressing mode

	// As the instruction is decoded, a std::string is assembled
	// with the readable output
	for uint32(addr) <= uint32(stop) {
		lineAddr = Word(addr)
		var sInst bytes.Buffer
		var e error
		var opcode byte
		// Prefix line with instruction address
		sInst.WriteString("$")
		sInst.WriteString(Hex(uint32(addr), 4))
		sInst.WriteString(": ")

		// Read instruction, and get its readable name
		opcode, e = c.bus.CPURead(Word(addr), true)
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
			value, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			sInst.WriteString("#$")
			sInst.WriteString(Hex(uint32(value), 2))
			sInst.WriteString(" {IMM}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ZP0) {
			lo, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(Hex(uint32(lo), 2))
			sInst.WriteString(" {ZP0}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ZPX) {
			lo, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(Hex(uint32(lo), 2))
			sInst.WriteString(", X {ZPX}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ZPY) {
			lo, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(Hex(uint32(lo), 2))
			sInst.WriteString(", Y {ZPY}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, IZX) {
			lo, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(Hex(uint32(lo), 2))
			sInst.WriteString(", X {IZX}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, IZY) {
			lo, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			hi = 0x00

			sInst.WriteString("$")
			sInst.WriteString(Hex(uint32(lo), 2))
			sInst.WriteString(", Y {IZY}")

		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ABS) {
			lo, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			hi, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("$")
			sInst.WriteString(Hex(uint32(Word(hi)<<8|Word(lo)), 4))
			sInst.WriteString(" {ABS}")

		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ABX) {
			lo, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			hi, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("$")
			sInst.WriteString(Hex(uint32(Word(hi)<<8|Word(lo)), 4))
			sInst.WriteString(", X {ABX}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, ABY) {
			lo, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			hi, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("$")
			sInst.WriteString(Hex(uint32(Word(hi)<<8|Word(lo)), 4))
			sInst.WriteString(", Y {ABY}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, IND) {
			lo, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++
			hi, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("($")
			sInst.WriteString(Hex(uint32(Word(hi)<<8|Word(lo)), 4))
			sInst.WriteString(") {IND}")
		} else if FnEquals(OpCodesLookupTable[opcode].addressmode, REL) {
			var val byte
			val, e = bus.CPURead(Word(addr), true)
			logError(e)
			addr++

			sInst.WriteString("$")
			sInst.WriteString(Hex(uint32(val), 4))
			sInst.WriteString(" [$")
			sInst.WriteString(Hex(uint32(int32(addr)+int32(int8(val))), 4))
			sInst.WriteString("] {REL}")
		}

		// Add the formed string to a std::map, using the instruction's
		// address as the key. This makes it convenient to look for later
		// as the instructions are variable in length, so a straight up
		// incremental index is not sufficient.
		mapLines[lineAddr] = sInst.String()
	}

	return mapLines
}
