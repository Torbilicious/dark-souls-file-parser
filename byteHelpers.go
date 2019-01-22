package main

import (
	"encoding/binary"
	"unicode/utf16"
	"unicode/utf8"
)

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
