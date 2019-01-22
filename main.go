package main

import (
	"encoding/binary"
	"gopkg.in/restruct.v1"
	"log"
	"time"
)

var (
	BlockSize           = 0x60190
	BlockIndex          = 0x2c0
	BlockDataOffset     = 0x14
	SlotsAmountOffset   = 0xC
	SlotsMetadataOffset = 0x40
	NameOffset          = 0x100
	DeathsOffset        = 0x1f128
	data                []byte
)

type SlotHeader struct {
	BlockMetadataHigh uint32
	BlockMetadataLow  uint32
	BlockSize         uint64
	BlockStartOffset  uint32
	BlockUnknownData1 uint32
	BlockSkipBytes    uint32
	EndOfBlock        uint32
}

type Player struct {
	name   string
	deaths int
}

func main() {
	file, _data := initFile()
	defer file.Close()

	data = _data

	start := time.Now()

	players := getPlayers()
	for _, player := range players {
		log.Printf("name:   %s", player.name)
		log.Printf("deaths: %d", player.deaths)
	}

	log.Printf("took %v\n", time.Since(start))
}

func getPlayers() []Player {
	amount := getAmountOfSlots()

	players := make([]Player, 0)

	for slotIndex := 0; slotIndex < amount; slotIndex++ {
		offset := BlockIndex + BlockSize*slotIndex
		deaths := readInt(offset+DeathsOffset, 4)

		bytes := make([]byte, 24)
		realOffset := offset + NameOffset

		restruct.Unpack(data[realOffset:], binary.LittleEndian, &bytes)

		bytes = sliceBytesToCorrectLength(bytes)

		name := UTF16BytesToString(bytes, binary.LittleEndian)

		players = append(players, Player{deaths: deaths, name: name})
	}

	return players
}

func getAmountOfSlots() int {
	amount := readInt(SlotsAmountOffset, 4)

	headers := make([]SlotHeader, amount)
	restruct.Unpack(data[SlotsMetadataOffset:], binary.LittleEndian, &headers)

	var counter = 0
	for _, header := range headers {

		offset := int(header.BlockStartOffset) + BlockDataOffset
		inByte := readNextBytes(offset, 1)

		if inByte[0] == 0 {
			return counter
		} else {
			counter++
		}
	}

	return counter
}
