package zx

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"testing"

	"fyne.io/fyne/canvas"
	"github.com/laullon/b2t80s/machines"
	"github.com/stretchr/testify/assert"
)

var borderTrix = new(string)

func TestBorderTrix(t *testing.T) {
	machines.TapFile = new(string)
	*machines.TapFile = "test/BorderTrix.tap"

	machines.LoadSlow = new(bool)
	machines.Debug = new(bool)

	zx := NewZX48K().(*zx)

	zx.ula.monitor = &dummyMonitor{}

	zx.Clock().RunFor(15)

	f, err := os.Open("test/BorderTrix_result.png")
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	result, diff, err := ImgCompare(img, zx.ula.display)
	assert.Equal(t, int64(0), result, err)

	if result != 0 {
		f, err = os.Create("test/BorderTrix_error.png")
		if err != nil {
			panic(err)
		}
		png.Encode(f, zx.ula.display)

		f, err = os.Create("test/BorderTrix_diff.png")
		if err != nil {
			panic(err)
		}
		png.Encode(f, diff)
	}
}

type dummyMonitor struct{}

func (m *dummyMonitor) Canvas() *canvas.Image { return nil }
func (m *dummyMonitor) FrameDone()            {}
func (m *dummyMonitor) FPS() float64          { return 0 }

func ImgCompare(img1, img2 image.Image) (int64, image.Image, error) {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	if bounds1 != bounds2 {
		return math.MaxInt64, nil, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)
	resultImg := image.NewRGBA(image.Rect(
		bounds1.Min.X,
		bounds1.Min.Y,
		bounds1.Max.X,
		bounds1.Max.Y,
	))
	draw.Draw(resultImg, resultImg.Bounds(), img1, image.Point{0, 0}, draw.Src)

	for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
		for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()

			diff := int64(sqDiffUInt32(r1, r2))
			diff += int64(sqDiffUInt32(g1, g2))
			diff += int64(sqDiffUInt32(b1, b2))
			diff += int64(sqDiffUInt32(a1, a2))

			if diff > 0 {
				accumError += diff
				resultImg.Set(
					bounds1.Min.X+x,
					bounds1.Min.Y+y,
					color.RGBA{R: 255, A: 255})
			}
		}
	}

	return int64(math.Sqrt(float64(accumError))), resultImg, nil
}

func sqDiffUInt32(x, y uint32) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}
