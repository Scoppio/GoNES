package main

// Mapper : interface to implement mappers
type Mapper interface {
	CPUMapRead(address Word) (uint32, bool)
	CPUMapWrite(address Word) (uint32, bool)
	PPUMapRead(address Word) (uint32, bool)
	PPUMapWrite(address Word) (uint32, bool)
}

// Mapper000 : default mapper
type Mapper000 struct {
	PGRBanks byte
	CHABanks byte
}

// CPUMapRead : Reads data from CPU
func (m *Mapper000) CPUMapRead(address Word) (uint32, bool) {
	if address >= 0x8000 && address <= 0xFFFF {
		var mappedAddress uint32
		if m.PGRBanks > 1 {
			mappedAddress = 0x00007FFF
		} else {
			mappedAddress = 0x00003FFF
		}
		mappedAddress &= uint32(address)
		return mappedAddress, true
	}
	return 0, false
}

// CPUMapWrite : map the write process to the correct address
func (m *Mapper000) CPUMapWrite(address Word) (uint32, bool) {
	if address >= 0x8000 && address <= 0xFFFF {
		var mappedAddress uint32
		if m.PGRBanks > 1 {
			mappedAddress = 0x00007FFF
		} else {
			mappedAddress = 0x00003FFF
		}
		mappedAddress &= uint32(address)
		return mappedAddress, true
	}
	return 0, false
}

// PPUMapRead : map the read process to the correct address
func (m *Mapper000) PPUMapRead(address Word) (uint32, bool) {
	if address >= 0x0000 && address <= 0x1FFF {

		return uint32(address), true
	}
	return 0, false
}

// PPUMapWrite : map the write process to the correct address
func (m *Mapper000) PPUMapWrite(address Word) (uint32, bool) {
	return 0, false
}

// Reset : resets mapper - can implement things like change starting position to run other games
func (m *Mapper000) Reset() {
	// do nothing
}
