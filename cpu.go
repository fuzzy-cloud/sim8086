package cpu

import (
	"fmt"
	"strings"
)

type opcode byte

const (
	opcodeInvalid opcode = iota
	MOV
)

type register byte

const (
	registerInvalid register = iota
	AL
	CL
	DL
	BL
	AH
	CH
	DH
	BH
	AX
	CX
	DX
	BX
	SP
	BP
	SI
	DI
)

var opcodeToString = map[opcode]string{
	opcodeInvalid: "INVALID",
	MOV:           "mov",
}

var registerToString = map[register]string{
	registerInvalid: "INVALID",
	AL:              "al",
	CL:              "cl",
	DL:              "dl",
	BL:              "bl",
	AH:              "ah",
	CH:              "ch",
	DH:              "dh",
	BH:              "bh",
	AX:              "ax",
	CX:              "cx",
	DX:              "dx",
	BX:              "bx",
	SP:              "sp",
	BP:              "bp",
	SI:              "si",
	DI:              "di",
}

var mapModeAndWidthToRegister = [...][2]register{
	0b000: {AL, AX},
	0b001: {CL, CX},
	0b010: {DL, DX},
	0b011: {BL, BX},
	0b100: {AH, SP},
	0b101: {CH, BP},
	0b110: {DH, SI},
	0b111: {BH, DI},
}

type decodeResult struct {
	op opcode
	// "Direction" bit. Equals to 1 when src is specified in REG field (and 0 for dst)
	d uint8
	// "Size" bit. 0 — we're working with bytes. 1 — with words.
	w uint8
	// "Memory mode" bits. Indicates whether one operand is in memory or both are registers
	mod uint8
	// "Register" bits. Used to identify which register to use (usually)
	reg uint8
	// "Register/memory" bits. Used to indentify which register to use or used in EAC
	rm uint8
}

func disassemble(stream []byte) (string, error) {
	if len(stream) < 2 {
		return "", fmt.Errorf("byte stream should at least have 2 bytes, got — %d", len(stream))
	}

	var out strings.Builder
	out.WriteString("bits 16\n")

	for ip := 0; ip < len(stream); {
		res := decode(stream[ip:])
		ip += 2

		switch res.mod {
		case 0b11:
			operand1 := mapModeAndWidthToRegister[res.reg][res.w]
			operand2 := mapModeAndWidthToRegister[res.rm][res.w]

			opStr := opcodeToString[res.op]
			dstStr := registerToString[operand2]
			srcStr := registerToString[operand1]

			fmt.Fprintf(&out, "\n%s %s, %s", opStr, dstStr, srcStr)
		default:
			return "", fmt.Errorf("unsupported MOD value: %b", res.mod)
		}
	}

	return out.String(), nil
}

func decode(stream []byte) (res decodeResult) {
	// NOTE: An instruction could from 1 to 6 byte in length
	b1 := stream[0]

	switch {
	case b1>>2 == 0b100010:
		b2 := stream[1]

		// b1
		res.op = MOV
		res.d = b1 >> 1 & 0x1
		res.w = b1 & 0x1

		// b2
		res.mod = b2 >> 6
		res.reg = b2 >> 3 & 0b111
		res.rm = b2 & 0b111
	default:
		panic("unsupported instruction")
	}

	return res
}
