package main

import (
	"bytes"
	"fmt"
	"os"
)

// ByteToHex : Converts letters to its equivalent hexadecimal code
// example: "A" becomes byte(0xA)
func ByteToHex(val byte) byte {
	//For uppercase A-F letters:
	if val < 58 {
		return val - 48
	}
	return val - 55
}

// Hex : converts a uint32 to its hexadecimal representation with leading zeroes
func Hex(n uint32, d int) string {
	s := fmt.Sprintf("%X", n)
	var b bytes.Buffer
	// b.WriteString("0x")
	if len(s) < d {
		for i := d - len(s); i > 0; i-- {
			b.WriteByte('0')
		}
	}
	b.WriteString(s)
	return b.String()
}

// WriteDisassemble : WriteDisassemble
func WriteDisassemble(disassemble map[Word]string, filename string) {
	f, _ := os.Create(filename)
	defer f.Close()
	var i Word
	for i < 0xFFFF {
		if v, ok := disassemble[i]; ok {
			f.WriteString(v + "\n")
		}
		i++
	}
}
