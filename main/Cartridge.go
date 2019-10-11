package main

import (
	"log"
	"os"
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

type Cartridge struct {
	bus       *Bus
	header    *header
	mapperID  byte
	mapper    *Mapper000
	PRGMemory []byte
	CHAMemory []byte
	PRGBanks  byte
	CHABanks  byte
}

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
		buf = make([]byte, int(cartHeader.CHARomBlocks)*8192) // define your buffer size here.
		_, e = file.Read(buf)
		logError(e)
		CHAMemory = buf

	} else if fileType == 2 {
		// Not implemented yet
	}

	cart := &Cartridge{nil, cartHeader, mapperID, &Mapper000{cartHeader.PGRRomBlocks, cartHeader.CHARomBlocks}, PRGMemory, CHAMemory, cartHeader.PGRRomBlocks, cartHeader.CHARomBlocks}

	return cart
}

func (c *Cartridge) CPURead(address rune) (byte, bool) {
	if mappedAddress, ok := c.mapper.CPUMapRead(address); ok {
		return c.PRGMemory[mappedAddress], true
	}
	return 0, false
}

func (c *Cartridge) CPUWrite(address rune, data byte) bool {
	if mappedAddress, ok := c.mapper.CPUMapWrite(address); ok {
		c.PRGMemory[mappedAddress] = data
		return true
	}
	return false
}

func (c *Cartridge) PPURead(address rune) (byte, bool) {
	if mappedAddress, ok := c.mapper.PPUMapRead(address); ok {
		return c.CHAMemory[mappedAddress], true
	}
	return 0, false
}

func (c *Cartridge) PPUWrite(address rune, data byte) bool {
	if mappedAddress, ok := c.mapper.PPUMapWrite(address); ok {
		c.CHAMemory[mappedAddress] = data
		return true
	}
	return false
}
