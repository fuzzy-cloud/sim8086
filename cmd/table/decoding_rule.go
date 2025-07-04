package main

import (
	cpu "cpu8086"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var mapStrToMnemonic = map[string]cpu.Mnemonic{
	"MOV": cpu.MOV,
}

func parseDecodingRule(raw string) (out cpu.DecodingRule, err error) {
	split := strings.Split(raw, " | ")
	if len(split) < 2 {
		err = fmt.Errorf("not enough bytes: %d", len(split))
		return
	}

	mnemonic, ok := mapStrToMnemonic[split[0]]
	if !ok {
		err = fmt.Errorf("invalid mnemonic: %s\n", split[0])
		return
	}
	out.Mnemonic = mnemonic

	split = split[1:]
	for byteIdx, rawByte := range split {
		rawParts := strings.Split(rawByte, " ")
		if len(rawParts) == 0 {
			continue
		}

		var (
			b       = &out.Bytes[byteIdx]
			shift   = 8
			partIdx = 0
		)

		for _, rawPart := range rawParts {
			p := &b.Parts[partIdx]
			p.Literal = -1

			switch rawPart {
			case "mod":
				shift -= 2
				p.Kind = cpu.PartMOD
				p.Mask = 0b11
				p.Shift = shift
			case "reg":
				shift -= 3
				p.Kind = cpu.PartREG
				p.Mask = 0b111
				p.Shift = shift
			case "rm":
				shift -= 3
				p.Kind = cpu.PartRM
				p.Mask = 0b111
				p.Shift = shift
			case "d":
				shift -= 1
				p.Kind = cpu.PartD
				p.Mask = 0b1
				p.Shift = shift
			case "w":
				shift -= 1
				p.Kind = cpu.PartW
				p.Mask = 0b1
				p.Shift = shift
			case "disp-lo":
				shift -= 8
				p.Kind = cpu.PartDISP_LO
				p.Mask = 0b11111111
				p.Shift = shift
			case "disp-hi":
				shift -= 8
				p.Kind = cpu.PartDISP_LO
				p.Mask = 0b11111111
				p.Shift = shift
			default:
				literal, parseErr := strconv.ParseInt(rawPart, 2, 16)
				if parseErr != nil {
					err = fmt.Errorf("failed to parse a literal: %s", rawPart)
					return
				}

				shift -= len(rawPart)
				p.Kind = cpu.PartLiteral
				p.Mask = 2<<len(rawPart) - 1
				p.Shift = shift
				p.Literal = literal
			}

			if shift < 0 {
				err = fmt.Errorf("shift is less than zero: %d", shift)
				return
			}

			b.NotEmpty = true
			p.NotEmpty = true

			partIdx++
		}
	}

	if c := countNotEmptyBytes(out.Bytes); c == 0 {
		err = errors.New("all bytes are empty")
	}

	return
}

func countNotEmptyBytes(bytes [6]cpu.ByteDecoding) (c int) {
	for i := range bytes {
		if bytes[i].NotEmpty {
			c++
		}
	}
	return
}
