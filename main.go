package main

import (
	"encoding/binary"
	"fmt"
	"gopkg.in/restruct.v1"
	"io/ioutil"
	"log"
	"os"
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
	file                *os.File
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

func main() {
	file, _ = os.Open("resources/DRAKS0005.sl2")
	defer file.Close()
	data, _ = ioutil.ReadAll(file)

	amount := getAmountOfSlots()
	fmt.Printf("Amount of slots: %d\n\n", amount)

	for slotIndex := 0; slotIndex < amount; slotIndex++ {
		offset := BlockIndex + BlockSize*slotIndex
		deaths := readInt(offset+DeathsOffset, 4)

		bytes := make([]byte, 24)
		realOffset := offset + NameOffset

		fmt.Printf("offset: %d\n", realOffset)

		restruct.Unpack(data[realOffset:], binary.LittleEndian, &bytes)

		bytes = sliceBytesToCorrectLength(bytes)

		name := UTF16BytesToString(bytes, binary.LittleEndian)

		fmt.Printf("name: %s\n", name)
		fmt.Printf("deaths: %d\n\n", deaths)
	}
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

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func getAmountOfSlots() int {
	amount := readInt(SlotsAmountOffset, 4)

	headers := make([]SlotHeader, amount)
	restruct.Unpack(data[SlotsMetadataOffset:], binary.LittleEndian, &headers)

	var counter = 0
	for _, header := range headers {

		file.Seek(int64(header.BlockStartOffset)+int64(BlockDataOffset), 0)
		inByte := readNextBytes(file, 1)

		if inByte[0] != 0 {
			counter++
		}
	}

	return counter
}
