package main

import (
	"bytes"
	"fmt"
)

func ByteToHex(val byte) byte {
	//For uppercase A-F letters:
	if val < 58 {
		return val - 48
	} else {
		return val - 55
	}
}

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
