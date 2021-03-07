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

	vMirror   bool
	fourPages bool

	flags []byte
}

func loadFile(fileName string) *nesFile {
	data := utils.ReadFile(fileName)
	if string(data[:3]) != "NES" {
		panic(-1)
	}
	file := &nesFile{}
	file.header.prgSize = data[4]
	file.header.chrSize = data[5]

	file.header.vMirror = data[6]&0x01 == 1
	file.header.fourPages = data[6]&0x08 != 0

	file.header.mapper = (data[7] & 0xf0) | ((data[6] & 0xf0) >> 4)

	file.header.flags = data[6:11]

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

func (f nesFile) String() string {
	return fmt.Sprintf("type:%d prg:%d chr:%d vMirror:%v", f.header.mapper, f.header.prgSize, f.header.chrSize, f.header.vMirror)
}
