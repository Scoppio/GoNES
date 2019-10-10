package main

// Memory64k : Memory with 64k bytes
type Memory64k struct {
	mem [64 * 1024]byte
}

// Reset : clears memory
func (m *Memory64k) Reset() {
	m.mem = [64 * 1024]byte{}
}

// PreLoadMemory : inserts data into memory using string format
// inserted data must be in hexadecimal writen as a string
// and they may have space after each 2 bytes
// example : "A9 0F 8D 15 40 60" or "A90F8D154060"
func (m *Memory64k) PreLoadMemory(offset Word, data string) {
	nOffset := offset
	for i := 0; i < len(data); i += 2 {
		m.mem[nOffset] = ByteToHex(data[i])<<4 | ByteToHex(data[i+1])
		if i+2 < len(data) && data[i+2] == byte(' ') {
			i++
		}
		nOffset++
	}
}

// SetCodeEntry : Set the address that starts your program
func (m *Memory64k) SetCodeEntry(address Word) {
	m.mem[0xFFFC] = byte(address)
	m.mem[0xFFFD] = byte(address >> 8)
}
