package main

type Mapper interface {
	CPUMapRead(address rune) (uint32, bool)
	CPUMapWrite(address rune) (uint32, bool)
	PPUMapRead(address rune) (uint32, bool)
	PPUMapWrite(address rune) (uint32, bool)
}

type Mapper000 struct {
	PGRBanks byte
	CHABanks byte
}

func (m *Mapper000) CPUMapRead(address rune) (uint32, bool) {
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

func (m *Mapper000) CPUMapWrite(address rune) (uint32, bool) {
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

func (m *Mapper000) PPUMapRead(address rune) (uint32, bool) {
	if address >= 0x0000 && address <= 0x1FFF {

		return uint32(address), true
	}
	return 0, false
}

func (m *Mapper000) PPUMapWrite(address rune) (uint32, bool) {
	return 0, false
}
