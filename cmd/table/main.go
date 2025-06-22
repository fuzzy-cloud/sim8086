package main

import (
	cpu "cpu8086"
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed table.sim8086
var rawTable string

var mapStrToMnemonic = map[string]cpu.Mnemonic{
	"MOV": cpu.MOV,
}

type Key uint8

const (
	KeyMnemonic Key = iota + 1
	KeyMOD
	KeyRM
	KeyREG
	KeyD
	KeyW
)

type Rule struct {
	Mnemoic cpu.Mnemonic
	Bytes   [6]Byte
}

type Byte struct {
	Valid bool
	Parts [3]Part
}

type Part struct {
	Valid bool
	Key   Key
	Mask  int
	Shift int
	Value int
}

func main() {
	for line := range strings.Lines(rawTable) {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, ";") {
			continue
		}

		var rule Rule
		var ok bool

		split := strings.Split(line, " | ")
		switch len(split) {
		case 0:
			panic(fmt.Sprintf("no mnemonic: %s\n", line))
		case 1:
			panic(fmt.Sprintf("no first byte: %s\n", line))
		}

		// mnemonic
		{
			rule.Mnemoic, ok = mapStrToMnemonic[split[0]]
			if !ok {
				panic(fmt.Sprintf("invalid mnemonic: %s\n", split[0]))
			}
		}

		// 1st byte
		var b1 Byte
		{
			split = strings.Split(split[1], " ")
			if len(split) == 0 {
				panic("empty 1st byte")
			}

			value, err := strconv.ParseInt(split[0], 2, 16)
			if err != nil {
				panic(err)

			}

			b1.Valid = true
			b1.Parts[0] = Part{
				Valid: true,
				Key:   KeyMnemonic,
				Mask:  1<<(8-len(split[0])) - 1,
				Shift: 8 - len(split[0]),
				Value: int(value),
			}
		}

		fmt.Printf("%+v\n", rule)
	}
}
