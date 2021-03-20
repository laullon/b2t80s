package mappers

import (
	"fmt"

	"github.com/laullon/b2t80s/utils"
)

type gbFile struct {
	header header
	data   []byte
}

type header struct {
	mapper  byte
	name    string
	romSize byte
	ramSize byte
}

func loadFile(fileName string) *gbFile {
	data := utils.ReadFile(fileName)
	file := &gbFile{}
	file.data = data

	file.header.mapper = data[0x0147]
	file.header.romSize = data[0x0148]
	file.header.ramSize = data[0x0149]
	file.header.name = string(data[0x0134:0x0143])

	fmt.Printf("file: %v \n", file)
	return file
}

func (f gbFile) String() string {
	return fmt.Sprintf("name:%s mapper:%d romSize:%d ramSize:%d", f.header.name, f.header.mapper, f.header.romSize, f.header.ramSize)
}
