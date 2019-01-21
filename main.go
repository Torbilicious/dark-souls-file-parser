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

type Player struct {
	name   string
	deaths int
}

func main() {
	initFile()
	defer file.Close()

	running := true
	counter := 0
	maxTurns := 1
	for running {
		start := time.Now()

		players := getPlayers()
		for _, player := range players {
			log.Printf("name:   %s", player.name)
			log.Printf("deaths: %d", player.deaths)
		}

		log.Printf("took %v\n", time.Since(start))

		reloadData()
		counter++
		running = counter < maxTurns
	}
}

func initFile() {
	//file, _ = os.Open("C:\\Users\\Torbe\\Documents\\nbgi\\darksouls\\torbilicious\\draks0005.sl2")
	//file, _ = os.Open("C:\\Users\\Torbe\\Documents\\nbgi\\darksouls\\gestirnmoewe302\\draks0005.sl2")
	file, _ = os.Open("resources/DRAKS0005-aes-encrypted.sl2")
	reloadData()

	//headerOk := isFileHeaderOk()
	//if !headerOk {
	//	log.Fatalln("Header check not ok!")
	//}
}

func isFileHeaderOk() bool {
	file.Seek(0, io.SeekStart)
	bytes := readNextBytes(4)
	fileType := string(bytes)

	file.Seek(0x18, io.SeekStart)
	bytes = readNextBytes(8)
	version := string(bytes)

	log.Printf("File header check:")
	log.Printf("fileType: %s", fileType)
	log.Printf("version:  %s", version)

	return fileType == "BND4" && version == "00000001"
}

func reloadData() {
	file.Seek(0, io.SeekStart)
	data, _ = ioutil.ReadAll(file)
}

func getPlayers() []Player {
	//amount := getAmountOfSlots()
	amount := 11

	players := make([]Player, 0)

	for slotIndex := 0; slotIndex < amount; slotIndex++ {
		//https://github.com/pawREP/Dark-Souls-Remastered-SL2-Unpacker

		begin := 704 + slotIndex*0x060030
		end := begin + 0x060020
		slotBytes := data[begin:end]

		var aesKey = []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}

		result := decrypt(aesKey, slotBytes)

		before := data[:begin]
		endOfData := data[end:]

		data = before
		data = append(data, result...)
		data = append(data, endOfData...)

		offset := BlockIndex + BlockSize*slotIndex
		deaths := readInt(offset+DeathsOffset, 4)

		bytes := make([]byte, 24)
		realOffset := offset + NameOffset

		restruct.Unpack(data[realOffset:], binary.LittleEndian, &bytes)

		bytes = sliceBytesToCorrectLength(bytes)

		name := UTF16BytesToString(bytes, binary.LittleEndian)

		player := Player{deaths: deaths, name: name}
		players = append(players, player)

		log.Printf("name:   %s", player.name)
		log.Printf("deaths: %d", player.deaths)
	}

	return players
}

func AppendByte(slice []byte, data ...byte) []byte {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) { // if necessary, reallocate
		// allocate double what's needed, for future growth.
		newSlice := make([]byte, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
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

func readNextBytes(number int) []byte {
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

		file.Seek(int64(header.BlockStartOffset)+int64(BlockDataOffset), io.SeekStart)
		inByte := readNextBytes(1)

		if inByte[0] != 0 {
			counter++
		}
	}

	return counter
}
