package fdc

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type cmd func(*fdc765)

type fdc765cmd struct {
	length    int
	resLength int
	direction int
	handler   cmd
	args      []byte
	result    []byte
	data      []byte
}

func newFDC765cmd(length int, resLength int, direction int, handler cmd) *fdc765cmd {
	return &fdc765cmd{
		length:    length,
		resLength: resLength,
		direction: direction,
		handler:   handler,
	}
}

func (cmd *fdc765cmd) init() {
	cmd.args = nil
	cmd.data = nil
	cmd.result = make([]byte, cmd.resLength)
}

func (cmd *fdc765cmd) String() string {
	name := runtime.FuncForPC(reflect.ValueOf(cmd.handler).Pointer()).Name()
	name = name[strings.LastIndex(name, ".")+1:]
	return fmt.Sprintf("CMD: '%s' args:%v result:%v data:%d", name, cmd.args, cmd.result, len(cmd.data))
}
