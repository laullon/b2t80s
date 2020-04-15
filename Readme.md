
# b2t80s
Z80 Based Computers Emulator (back to the 80's)

## Download

<https://github.com/laullon/b2t80s/releases/latest>

## Usage

```
  -bp string
        Breakpoints [0xXXXX[,0xXXXX,...]]
  -cpuprofile file
        write cpu profile to file
  -debug
        shows debugger
  -dskA string
        disc file to load on drive A
  -memprofile file
        write memory profile to file
  -mode string
        Spectrum model to emulate [48k|128k|plus3|cpc464|cpc6128] (default "48k")
  -slow
        Real Spectrum loading process
  -tap string
        tap file to load
  -z80 string
        z80 file to load
```

## Build and Run 

### Dependencies
```
go get -u github.com/go-bindata/go-bindata/...
$HOME/go/bin/go-bindata -pkg data -o data/data.go data/...
```

### Run
```
go run main.go --mode 48k -tap "./games/ManicMiner.tap"
```

## links

### ZX

* <https://stackoverflow.com/questions/1215777/writing-a-graphical-z80-emulator-in-c-or-c>
* tests: <http://mdfs.net/Software/Z80/Exerciser/>
* The Complete Spectrum ROM Disassembly: <https://skoolkid.github.io/rom/maps/all.html#0038>
* SPECTRUM 128 ROM 0 DISASSEMBLY <http://www.matthew-wilson.net/spectrum/rom/128_ROM0.html?LMCL=aH_qpw&LMCL=L7lymk#L1F45>
* Roms: <http://www.shadowmagic.org.uk/spectrum/roms.html>
* Contention Test Success <http://www.zxdesign.info/testSuccess.shtml>
* Video Parameters <http://www.zxdesign.info/vidparam.shtml>

### CPC

* cpc6128 rom: <http://cpctech.cpc-live.com/docs/os.asm>
* <http://cpctech.cpc-live.com/docs/basic.asm>
* <http://cpctech.cpc-live.com/docs/amsdos.asm>

### Z80

* <https://www.chibiakumas.com/z80/>

### my

* int.asm <https://gist.github.com/laullon/9928e27738df3c5a194d92c7b2977710>

## ZexDoc

```
zmac --zmac zexdocsmall.asm
go test -v -timeout 999m github.com/laullon/b2t80s/emulator -run TestZEXDoc
```
