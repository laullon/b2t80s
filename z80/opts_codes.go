package z80

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/laullon/b2t80s/data"
	"github.com/laullon/b2t80s/emulator"
	"github.com/laullon/b2t80s/utils"
)

type opcodeData struct {
	ids     string
	str     string
	tStates string
	length  string
}

var _Instructions = make([]emulator.Instruction, 0xff+1)
var _CBInstructions = make([]emulator.Instruction, 0xff+1)
var _EDInstructions = make([]emulator.Instruction, 0xff+1)
var _FDInstructions = make([]emulator.Instruction, 0xff+1)
var _FDCBInstructions = make([]emulator.Instruction, 0xff+1)
var _DDInstructions = make([]emulator.Instruction, 0xff+1)
var _DDCBInstructions = make([]emulator.Instruction, 0xff+1)

func GetOpCode(mem []byte) (emulator.Instruction, error) {
	opCode := mem[0]
	var ins emulator.Instruction
	switch opCode {
	case 0xcb:
		ins = _CBInstructions[mem[1]]
	case 0xed:
		ins = _EDInstructions[mem[1]]
	case 0xdd:
		byte1 := mem[1]
		if byte1 == 0xcb {
			ins = _DDCBInstructions[mem[3]]
		} else {
			ins = _DDInstructions[mem[1]]
		}
		if !ins.Valid { // DD with out IX, use next byte - zexdoc
			ins = _Instructions[byte1]
		}
	case 0xfd:
		byte1 := mem[1]
		if byte1 == 0xcb {
			ins = _FDCBInstructions[mem[3]]
		} else {
			ins = _FDInstructions[mem[1]]
		}
		if !ins.Valid { // DD with out IY, use next byte - zexdoc
			ins = _Instructions[byte1]
		}
	default:
		ins = _Instructions[opCode]
	}

	if !ins.Valid {
		return ins, NotSupported(fmt.Sprintf("%02X %02X %02X %02X", mem[0], mem[1], mem[2], mem[3]))
	}

	ins.Mem = mem[:ins.Length]
	return ins, nil
}

func LoadOPCodess() {
	z80ops, err := data.Asset("data/z80ops.csv")
	if err != nil {
		panic("data/z80ops.csv not found")
	}

	var allOpcodeData []*opcodeData
	r := csv.NewReader(bytes.NewReader(z80ops))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		allOpcodeData = append(allOpcodeData, newOpcodeData(record))
	}

	// "BIT b,r","8","2","CB 40+8*b+rb"
	rb := []string{"B", "C", "D", "E", "H", "L", "(HL)", "A"}
	for _, opc := range allOpcodeData {
		if strings.Index(opc.ids, "rb") != -1 {
			for i, b := range rb {
				newOpc := opc.copy()
				newOpc.str = strings.ReplaceAll(newOpc.str, "r", b)
				newOpc.ids = strings.ReplaceAll(newOpc.ids, "rb", strconv.Itoa(i))
				allOpcodeData = append(allOpcodeData, newOpc)
			}
		}
	}

	// str, tStates, length, ids
	// "BIT b,(IX+N)","20","4","DD CB XX 46+8*b"
	for _, opc := range allOpcodeData {
		if strings.Index(opc.ids, "b") != -1 {
			for b := 0; b < 8; b++ {
				newOpc := opc.copy()
				newOpc.str = strings.ReplaceAll(newOpc.str, "b", strconv.Itoa(b))
				newOpc.ids = strings.ReplaceAll(newOpc.ids, "b", strconv.Itoa(b))
				allOpcodeData = append(allOpcodeData, newOpc)
			}
		}
	}

	for _, opc := range allOpcodeData {
		opc.ids = strings.Replace(opc.ids, " XX", "", -1)
		if strings.IndexAny(opc.ids, "rb") == -1 {
			ids := strings.Split(opc.ids, " ")
			switch len(ids) {
			case 1:
				id0 := utils.Eval(ids[0])
				if (_Instructions[id0].Valid) && (opc.str != _Instructions[id0].Opcode) {
					// log.Printf("'%v' ignored (good:%v)", opc.str, _Instructions[id0].Opcode)
				} else {
					_Instructions[id0] = newInstruction(int32(id0), opc)
				}

			case 2:
				id0 := utils.Eval(ids[0])
				id1 := utils.Eval(ids[1])
				id := int32(id0<<8 | id1)
				if id0 == 0xcb {
					_CBInstructions[id1] = newInstruction(id, opc)
				} else if id0 == 0xfd {
					_FDInstructions[id1] = newInstruction(id, opc)
				} else if id0 == 0xdd {
					_DDInstructions[id1] = newInstruction(id, opc)
				} else if id0 == 0xed {
					_EDInstructions[id1] = newInstruction(id, opc)
				} else {
					panic(fmt.Sprintf("%02X %02X (%s %s) %s", id0, id1, ids[0], ids[1], opc.str))
				}
			case 3:
				id0 := utils.Eval(ids[0])
				id1 := utils.Eval(ids[1])
				id2 := utils.Eval(ids[2])
				id := int32(id0<<16 | id1<<8 | id2)
				if id0 == 0xdd && id1 == 0xcb {
					_DDCBInstructions[id2] = newInstruction(id, opc)
				} else if id0 == 0xfd && id1 == 0xcb {
					_FDCBInstructions[id2] = newInstruction(id, opc)
				} else {
					panic(fmt.Sprintf("%02X %02X %02X (%s %s %s) %s", id0, id1, id2, ids[0], ids[1], ids[2], opc.str))
				}
			}
		}
	}
}

func newOpcodeData(record []string) *opcodeData {
	return &opcodeData{
		str:     record[0],
		tStates: record[1],
		length:  record[2],
		ids:     record[3],
	}
}

func (opc *opcodeData) copy() *opcodeData {
	return &opcodeData{
		str:     opc.str,
		tStates: opc.tStates,
		length:  opc.length,
		ids:     opc.ids,
	}
}

func newInstruction(id int32, opc *opcodeData) emulator.Instruction {
	l, err := strconv.ParseUint(opc.length, 10, 16)
	if err != nil {
		panic(err)
	}

	ts, err := strconv.ParseInt(strings.Split(opc.tStates, "/")[0], 10, 16)
	if err != nil {
		panic(err)
	}

	return emulator.Instruction{
		Instruction: id,
		Opcode:      opc.str,
		// TODO Tstates and altTstates
		Tstates: uint(ts),
		Length:  uint16(l),
		Valid:   true,
	}
}

type NotSupportedError struct {
	op string
}

func (e *NotSupportedError) Error() string { return fmt.Sprintf("opt code '%s' not supported\n", e.op) }

func NotSupported(op string) *NotSupportedError {
	return &NotSupportedError{op: op}
}

type NotImplementedError struct {
	op string
}

func (e *NotImplementedError) Error() string {
	return fmt.Sprintf("opt '%s' not supported\n", e.op)
}

func NotImplemented(op string) *NotImplementedError {
	return &NotImplementedError{op: op}
}
