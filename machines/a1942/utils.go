package a1942

import (
	"archive/zip"

	"github.com/laullon/b2t80s/utils"
)

func loadRom(fileName string) []byte {
	zipFile := "/Users/glaullon/go/src/github.com/laullon/b2t80s/games/1942.zip"
	var mem []byte
	zf, err := zip.OpenReader(zipFile)
	if err != nil {
		panic(err)
	}

	for _, file := range zf.File {
		if file.Name == fileName {
			mem = utils.ReadZipFile(file)
		}
	}

	err = zf.Close()
	if err != nil {
		panic(err)
	}

	println("loaded rom:", fileName, "size:", len(mem))
	if len(mem) == 0 {
		panic("rom not found")
	}

	return mem
}
