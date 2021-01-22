package atetris

import (
	"archive/zip"

	"github.com/laullon/b2t80s/utils"
)

func loadRom(fileName string) []byte {
	zipFile := "/Users/glaullon/go/src/github.com/laullon/b2t80s/games/atetris.zip"
	var mem []byte
	zf, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}

	for _, file := range zf.File {
		println(file.Name)
		if file.Name == fileName {
			mem = utils.ReadZipFile(file)
		}
	}

	err = zf.Close()
	if err != nil {
		panic(err)
	}
	return mem
}
