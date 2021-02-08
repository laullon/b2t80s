package mappers

import (
	"fmt"

	"github.com/laullon/b2t80s/utils"
)

type nesFile struct {
	header  header
	trainer []byte
	prg     []byte
	chr     []byte
}

type header struct {
	mapper  byte
	prgSize byte
	chrSize byte
	flags   []byte
}

func loadFile(fileName string) *nesFile {
	data := utils.ReadFile(fileName)
	if string(data[:3]) != "NES" {
		panic(-1)
	}
	file := &nesFile{}
	file.header.prgSize = data[4]
	file.header.chrSize = data[5]
	file.header.flags = data[6:11]
	file.header.mapper = (file.header.flags[1] & 0xf0) | ((file.header.flags[0] & 0xf0) >> 4)

	fmt.Printf("file: %v \n", file)

	data = data[16:]
	if file.header.flags[0]&0b00000100 != 0 {
		file.trainer = data[:512]
		data = data[512:]
	}

	file.prg = data[:0x4000*uint32(file.header.prgSize)]
	data = data[0x4000*uint32(file.header.prgSize):]

	file.chr = data[:0x2000*uint32(file.header.chrSize)]
	data = data[0x2000*uint32(file.header.chrSize):]

	if len(data) != 0 {
		panic(-1)
	}

	return file
}
