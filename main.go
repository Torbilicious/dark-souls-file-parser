package main

import (
	"encoding/binary"
	"gopkg.in/restruct.v1"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
	"unicode/utf16"
	"unicode/utf8"
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
	file := initFile()
	defer file.Close()

	start := time.Now()

	players := getPlayers()
	for _, player := range players {
		log.Printf("name:   %s", player.name)
		log.Printf("deaths: %d", player.deaths)
	}

	log.Printf("took %v\n", time.Since(start))
}

func initFile() (file *os.File) {
	file, _ = os.Open("resources/DRAKS0005.sl2")
	reloadData(file)

	headerOk := isFileHeaderOk()
	if !headerOk {
		log.Fatalln("Header check not ok!")
	}

	return file
}

func isFileHeaderOk() bool {
	bytes := readNextBytes(0x0, 4)
	fileType := string(bytes)

	bytes = readNextBytes(0x18, 8)
	version := string(bytes)

	log.Printf("File header check:")
	log.Printf("fileType: %s", fileType)
	log.Printf("version:  %s", version)

	return fileType == "BND4" && version == "00000001"
}

func reloadData(file *os.File) {
	file.Seek(0, io.SeekStart)
	data, _ = ioutil.ReadAll(file)
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

func sliceBytesToCorrectLength(bytes []byte) []byte {
	out := make([]byte, 0)

	for i := 0; i < len(bytes); i += 2 {
		if bytes[i] == 0 && bytes[i+1] == 0 {
			return out
		}

		out = append(out, bytes[i])
		out = append(out, bytes[i+1])
	}

	return out
}

func UTF16BytesToString(b []byte, o binary.ByteOrder) string {
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}
	return string(utf16.Decode(utf))
}

func readInt(offset int, length int) int {
	return int(binary.LittleEndian.Uint32(data[offset : offset+length]))
}

func readNextBytes(offset int, number int) []byte {
	return data[offset : offset+number]
}

func getAmountOfSlots() int {
	amount := readInt(SlotsAmountOffset, 4)

	headers := make([]SlotHeader, amount)
	restruct.Unpack(data[SlotsMetadataOffset:], binary.LittleEndian, &headers)

	var counter = 0
	for _, header := range headers {

		offset := int(header.BlockStartOffset) + BlockDataOffset
		inByte := readNextBytes(offset, 1)

		if inByte[0] != 0 {
			counter++
		}
	}

	return counter
}
