package files

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type DSK interface {
	SeekSector(chrn []byte) *Sector
	SeekTrack(side, track int)
	ActualSector() *Sector
	NextSector() *Sector
	GetID(side int) []byte
}

type dsk struct {
	sides      byte
	trackSizes []uint16
	tracks     [][]*track
	flipped    bool
	head       *head
}

type head struct {
	side, track, sector int
}

type track struct {
	side       int
	sectorSize int
	sectors    []*Sector
}

type Sector struct {
	CHRN     []byte
	ST1, ST2 byte
	Data     []byte
}

func (s *Sector) String() string {
	return fmt.Sprintf("chrn:%v ST1:%d ST2:%d data:%v", s.CHRN, s.ST1, s.ST2, len(s.Data))
}

func LoadDsk(fileName string) DSK {
	fi, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	file := make([]byte, fi.Size()+1)
	l, err := f.Read(file)
	if err != nil {
		panic(err)
	}
	file = file[:l]

	// fmt.Printf("Loading disc '%s'\n", fileName)
	t := string(file[0:21])
	// fmt.Printf("- Disc type '%s'\n", t)

	tracks := file[0x30]
	dsk := &dsk{
		sides:  file[0x31],
		tracks: make([][]*track, file[0x31]),
		head:   &head{0, 0, 0},
	}

	switch {
	case strings.HasPrefix(t, "MV"):
		tracks := int(file[0x30]) * int(file[0x31])
		trackSize := int(file[0x32]) | int(file[0x33])<<8
		file = file[0x100:]
		for idx := 0; idx < tracks; idx++ {
			track := readTrack(file[:trackSize])
			// fmt.Printf("(%d)%+v\n", trackSize, track)
			dsk.tracks[track.side] = append(dsk.tracks[track.side], track)
			file = file[trackSize:]
		}

	case strings.HasPrefix(t, "EXTENDED CPC DSK"):
		tracks := int(tracks) * int(dsk.sides)
		for idx := 0; idx < tracks; idx++ {
			ts := uint16(file[0x34+idx]) << 8
			dsk.trackSizes = append(dsk.trackSizes, ts)
		}
		file = file[0x100:]
		for idx := 0; idx < tracks; idx++ {
			track := readTrack(file[:dsk.trackSizes[idx]])
			dsk.tracks[track.side] = append(dsk.tracks[track.side], track)
			file = file[dsk.trackSizes[idx]:]
		}
	default:
		panic(t)
	}
	// fmt.Printf("%v\n", dsk)
	return dsk
}

func (dsk *dsk) ActualSector() *Sector {
	return dsk.tracks[dsk.head.side][dsk.head.track].sectors[dsk.head.sector]
}

func (dsk *dsk) GetID(side int) []byte {
	dsk.head.side = side
	return dsk.tracks[dsk.head.side][dsk.head.track].sectors[dsk.head.sector].CHRN
}

func (dsk *dsk) SeekTrack(side, track int) {
	dsk.head.side = side
	dsk.head.track = track
	dsk.head.sector = 0
}

func (dsk *dsk) SeekSector(chrn []byte) *Sector {
	dsk.head.side = int(chrn[1])
	dsk.head.track = int(chrn[0])
	for idx, sector := range dsk.tracks[dsk.head.side][dsk.head.track].sectors {
		if reflect.DeepEqual(sector.CHRN, chrn) {
			dsk.head.sector = idx
			return sector
		}
	}
	// println(dsk.String())
	panic(fmt.Sprintf("Sector '%v' not found", chrn))
	// fmt.Printf("Sector '%v' not found", chrn)
	// return nil
}

func (dsk *dsk) NextSector() *Sector {
	dsk.head.sector++
	if dsk.head.sector > len(dsk.tracks[dsk.head.side][dsk.head.track].sectors) {
		dsk.head.sector = 0
		dsk.head.track++
		if dsk.head.track > len(dsk.tracks[dsk.head.side]) {
			dsk.head.track = 0
			dsk.head.side = 1 - dsk.head.side
		}
	}
	return dsk.ActualSector()
}

func readTrack(file []byte) *track {
	if len(file) == 0 {
		return &track{
			sectorSize: 0,
		}
	}
	header := string(file[0:0x0b])
	if !strings.HasPrefix(header, "Track-Info") {
		panic("header")
	}

	t := &track{
		sectorSize: 128 << int(file[0x14]),
		side:       int(file[0x11]),
	}
	sn := int(file[0x15])
	for i := 0; i < sn; i++ {
		sector := &Sector{
			CHRN: file[0x18+8*i : 0x18+8*i+4],
			ST1:  file[0x18+8*i+4],
			ST2:  file[0x18+8*i+5],
			Data: file[0x100+t.sectorSize*i : 0x100+t.sectorSize*i+t.sectorSize],
		}
		t.sectors = append(t.sectors, sector)
	}
	return t
}

func (t *track) String() string {
	str := fmt.Sprintf("[Track] side:%d sectorSize:%d sectors:%d", t.side, t.sectorSize, len(t.sectors))
	return str
}

func (d *dsk) String() string {
	str := fmt.Sprintf("[dsk] sides:%d tracks:%v", d.sides, len(d.tracks))
	for s, side := range d.tracks {
		for t, track := range side {
			str = fmt.Sprintf("%s\n\t%v", str, track)
			for c, sector := range track.sectors {
				str = fmt.Sprintf("%s\n\t\t %v (%d.%d.%d)", str, sector, s, t, c)
			}
		}
	}
	return str
}
