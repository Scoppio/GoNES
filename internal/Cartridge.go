package main

import (
	"log"
	"os"
)

const (
	// Horizontal : Horizontal
	Horizontal = 0
	// Vertical : Vertical
	Vertical = 1
	// OnescreenLo : OnescreenLo
	OnescreenLo = 2
	// OnescreenHi : OnescreenHi
	OnescreenHi = 3
)

type header struct {
	name         [4]byte
	PGRRomBlocks byte
	CHARomBlocks byte
	mapper1      byte
	mapper2      byte
	PRGRamSize   byte
	TVSystem1    byte
	TVSystem2    byte
	unused       [5]byte
}

// Cartridge : struct that defines the Cart object
type Cartridge struct {
	bus       *Bus
	header    *header
	mapperID  byte
	mapper    *Mapper000
	PRGMemory []byte
	CHAMemory []byte
	PRGBanks  byte
	CHABanks  byte
	Mirror    int
}

// TestCartridge : handmade cart for testing
func TestCartridge(rom string, offset Word) *Cartridge {

	cartHeader := &header{
		[4]byte{'t', 'e', 's', 't'},
		1,
		1,

		0,
		0,
		0,

		0,
		0,
		[5]byte{}}

	mapperID := ((cartHeader.mapper2 >> 4) << 4) | (cartHeader.mapper1 >> 4)

	// Discover what kind of iNes file, hardcoded 1 for now
	var PRGMemory, CHAMemory []byte

	buf := make([]byte, int(cartHeader.PGRRomBlocks)*16384)
	PRGMemory = buf

	nOffset := 0
	for i := 0; i < len(rom); i += 2 {
		PRGMemory[nOffset] = ByteToHex(rom[i])<<4 | ByteToHex(rom[i+1])
		if i+2 < len(rom) && rom[i+2] == byte(' ') {
			i++
		}
		nOffset++
	}

	PRGMemory[0xFFFC&0x3FFF] = byte(offset)
	PRGMemory[0xFFFD&0x3FFF] = byte(offset >> 8)

	buf = make([]byte, int(cartHeader.CHARomBlocks)*8192)
	CHAMemory = buf

	cart := &Cartridge{nil, cartHeader, mapperID,
		&Mapper000{cartHeader.PGRRomBlocks, cartHeader.CHARomBlocks}, PRGMemory, CHAMemory, cartHeader.PGRRomBlocks, cartHeader.CHARomBlocks, Horizontal}

	return cart
}

// LoadCartridge : loads the cart after giving a filepath
func LoadCartridge(filepath string) *Cartridge {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var cartHeader *header

	bh := make([]byte, 16) // define buffer header

	head, err := file.Read(bh)
	if head > 0 {
		cartHeader = &header{
			[4]byte{bh[0], bh[1], bh[2], bh[3]},
			bh[4],
			bh[5],

			bh[6],
			bh[7],
			bh[8],

			bh[9],
			bh[10],
			[5]byte{}}
	}

	if cartHeader.mapper1&0x04 != 0 {
		file.Seek(512, 1)
	}

	mapperID := ((cartHeader.mapper2 >> 4) << 4) | (cartHeader.mapper1 >> 4)
	mirror := Horizontal
	if cartHeader.mapper1&0x01 > 0 {
		mirror = Vertical
	}

	// Discover what kind of iNes file, hardcoded 1 for now
	fileType := 1
	var PRGMemory, CHAMemory []byte
	if fileType == 0 {
		// Not implemented yet
	} else if fileType == 1 {

		buf := make([]byte, int(cartHeader.PGRRomBlocks)*16384) // define your buffer size here.
		_, e := file.Read(buf)
		logError(e)
		PRGMemory = buf
		if cartHeader.CHARomBlocks == 0 {
			buf = make([]byte, 8192) // define your buffer size here.
			_, e = file.Read(buf)
			logError(e)
			CHAMemory = buf
		} else {
			buf = make([]byte, int(cartHeader.CHARomBlocks)*8192) // define your buffer size here.
			_, e = file.Read(buf)
			logError(e)
			CHAMemory = buf
		}

	} else if fileType == 2 {
		// Not implemented yet
	}

	cart := &Cartridge{nil, cartHeader, mapperID, &Mapper000{cartHeader.PGRRomBlocks, cartHeader.CHARomBlocks}, PRGMemory, CHAMemory, cartHeader.PGRRomBlocks, cartHeader.CHARomBlocks, mirror}

	return cart
}

// CPURead : allows the reading of data by the CPU
func (c *Cartridge) CPURead(address Word) (byte, bool) {
	if mappedAddress, ok := c.mapper.CPUMapRead(address); ok {
		return c.PRGMemory[mappedAddress], true
	}
	return 0, false
}

// CPUWrite : allows the CPU to write data
func (c *Cartridge) CPUWrite(address Word, data byte) bool {
	if mappedAddress, ok := c.mapper.CPUMapWrite(address); ok {
		c.PRGMemory[mappedAddress] = data
		return true
	}
	return false
}

// PPURead : reads PPU data
func (c *Cartridge) PPURead(address Word) (byte, bool) {
	if mappedAddress, ok := c.mapper.PPUMapRead(address); ok {
		return c.CHAMemory[mappedAddress], true
	}
	return 0, false
}

// PPUWrite : writes data to PPU
func (c *Cartridge) PPUWrite(address Word, data byte) bool {
	if mappedAddress, ok := c.mapper.PPUMapWrite(address); ok {
		c.CHAMemory[mappedAddress] = data
		return true
	}
	return false
}

// Reset : reset process
func (c *Cartridge) Reset() {
	if c != nil && c.mapper != nil {
		c.mapper.Reset()
	}
}
