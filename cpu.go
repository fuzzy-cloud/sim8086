package cpu

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type opcode byte

const (
	opcodeInvalid opcode = iota
	MOV
)

var opcodeToString = [...]string{
	opcodeInvalid: "INVALID",
	MOV:           "mov",
}

func (o opcode) String() string {
	return opcodeToString[o]
}

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

var registerToString = [...]string{
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

func (r register) String() string {
	return registerToString[r]
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

var mapRMToRegister = [...][2]register{
	0b000: {BX, SI},
	0b001: {BX, DI},
	0b010: {BP, SI},
	0b011: {BP, DI},
	0b100: {SI},
	0b101: {DI},
	0b110: {BP},
	0b111: {BX},
}

func disassemble(stream []byte) (string, error) {
	var (
		out    strings.Builder
		outErr error
	)

	out.WriteString("bits 16\n")

	for ip := 0; ip < len(stream); {
		inst, n, err := decode(stream[ip:])
		if err != nil {
			outErr = err
			break
		}
		ip += n

		fmt.Fprintf(&out, "\n")
		fmt.Fprintf(&out, "%s %s, %s", inst.opcode, inst.dst, inst.src)
	}

	return out.String(), outErr
}

type instruction struct {
	opcode opcode
	dst    operand
	src    operand
}

type operand struct {
	kind operandKind
	reg  register
	imm  int16
	da   uint16
	eac  struct {
		reg1 register
		reg2 register
		disp uint16
	}
}

type operandKind int

const (
	opKindReg operandKind = iota + 1
	opKindImm
	opKindEAC
	opKindDa
)

func operandReg(reg register) (o operand) {
	o.kind = opKindReg
	o.reg = reg
	return
}

func operandImm(imm int16) (o operand) {
	o.kind = opKindImm
	o.imm = imm
	return
}

func operandEAC(disp uint16, reg1, reg2 register) (o operand) {
	o.kind = opKindEAC
	o.eac.disp = disp
	o.eac.reg1 = reg1
	o.eac.reg2 = reg2
	return
}

func (o operand) String() string {
	switch o.kind {
	case opKindReg:
		return registerToString[o.reg]
	case opKindImm:
		return strconv.Itoa(int(o.imm))
	case opKindEAC:
		var out strings.Builder

		fmt.Fprintf(&out, "[%s", o.eac.reg1)
		if o.eac.reg2 != registerInvalid {
			fmt.Fprintf(&out, " + %s", o.eac.reg2)
		}
		if o.eac.disp != 0 {
			fmt.Fprintf(&out, " + %d", o.eac.disp)
		}
		fmt.Fprintf(&out, "]")

		return out.String()
	default:
		panic(fmt.Sprintf("unsupported operand kind: %d", o.kind))
	}
}

func decode(stream []byte) (inst instruction, n int, err error) {
	var (
		readByteDisp bool
		readWordDisp bool
		readByteData bool
		readWordData bool
	)

	var (
		isRegToReg     bool
		isImmToReg     bool
		isDirectAccess bool
		isEAC          bool
	)

	var (
		// "Direction" bit. Equals to 0 when src is specified in REG field (and 1 for dst)
		d = -1
		// "Size" bit. 0 — we're working with bytes. 1 — with words.
		w = -1
		// "Memory mode" bits. Indicates whether one operand is in memory or both are registers
		mod = -1
		// "Register" bits. Used to identify which register to use (usually)
		reg = -1
		// "Register/memory" bits. Used to indentify which register to use or used in EAC
		rm = -1
	)

	// NOTE: An instruction could from 1 to 6 byte in length
	b1 := stream[n]
	n++

	switch {
	case b1>>2 == 0b100010:
		inst.opcode = MOV

		// b1
		d = int(b1 >> 1 & 0x1)
		w = int(b1 & 0x1)

		// b2
		b2 := stream[n]
		n++

		mod = int(b2 >> 6)
		reg = int(b2 >> 3 & 0b111)
		rm = int(b2 & 0b111)
	case b1>>4 == 0b1011:
		inst.opcode = MOV

		// b1
		w = int(b1 >> 3 & 0b1)
		reg = int(b1 & 0b111)

		// NOTE: knowledge encoded into this specific instruction: we should read data
		if w == 0 {
			readByteData = true
		} else {
			readWordData = true
		}

		// NOTE: knowledge encoded into this specific instruction: the data is immediate value
		isImmToReg = true
	default:
		err = fmt.Errorf("unsupported instruction: %b", b1)
		return
	}

	switch mod {
	case 0b00:
		if rm == 0b110 {
			isDirectAccess = true
		} else {
			isEAC = true
		}
	case 0b01:
		isEAC = true
		readByteDisp = true
	case 0b10:
		isEAC = true
		readWordDisp = true
	case 0b11:
		isRegToReg = true
	}

	var disp uint16
	if readByteDisp {
		disp = uint16(stream[n])
		n++
	}
	if readWordDisp {
		disp = binary.NativeEndian.Uint16(stream[n:])
		n += 2
	}

	var data int16
	if readByteData {
		data = int16(int8(stream[n]))
		n++
	}
	if readWordData {
		data = int16(binary.NativeEndian.Uint16(stream[n:]))
		n += 2
	}

	var (
		isDirectAccessToReg = isDirectAccess && d == 1
		isRegToDirectAccess = isDirectAccess && d == 0
	)

	switch {
	case isRegToReg:
		operand1 := operandReg(mapModeAndWidthToRegister[rm][w])
		operand2 := operandReg(mapModeAndWidthToRegister[reg][w])

		if d == 0 {
			inst.dst = operand1
			inst.src = operand2
		} else {
			inst.dst = operand2
			inst.src = operand1
		}
	case isImmToReg:
		inst.dst = operandReg(mapModeAndWidthToRegister[reg][w])
		inst.src = operandImm(data)
	case isEAC:
		regs := mapRMToRegister[rm]

		operand1 := operandEAC(disp, regs[0], regs[1])
		operand2 := operandReg(mapModeAndWidthToRegister[reg][w])

		if d == 0 {
			inst.dst = operand1
			inst.src = operand2
		} else {
			inst.dst = operand2
			inst.src = operand1
		}
	case isDirectAccessToReg:
		panic("todo: DA to reg")
	case isRegToDirectAccess:
		panic("todo: reg to DA")
	}

	return
}
