package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	file *os.File
)

func initFile() (*os.File, []byte) {
	file, _ = os.Open("resources/DRAKS0005.sl2")
	data = loadData()

	headerOk := isFileHeaderOk()
	if !headerOk {
		log.Fatalln("Header check not ok!")
	}

	return file, data
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

func loadData() (data []byte) {
	file.Seek(0, io.SeekStart)
	data, _ = ioutil.ReadAll(file)

	return data
}
