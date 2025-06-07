package cpu

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type mnemonic byte

const (
	mnemonicInvalid mnemonic = iota
	MOV
	ADD
	SUB
	CMP
	JNZ
	JE
	JL
	JLE
	JB
	JBE
	JP
	JO
	JS
	JNE
	JNL
	JG
	JNB
	JA
	JNP
	JNO
	JNS
	LOOP
	LOOPZ
	LOOPNZ
	JCXZ
)

var mnemonicToString = [...]string{
	mnemonicInvalid: "INVALID",
	MOV:             "mov",
	ADD:             "add",
	SUB:             "sub",
	CMP:             "cmp",
	JNZ:             "jnz",
	JE:              "je",
	JL:              "jl",
	JLE:             "jle",
	JB:              "jb",
	JBE:             "jbe",
	JP:              "jp",
	JO:              "jo",
	JS:              "js",
	JNE:             "jne",
	JNL:             "jnl",
	JG:              "jg",
	JNB:             "jnb",
	JA:              "ja",
	JNP:             "jnp",
	JNO:             "jno",
	JNS:             "jns",
	LOOP:            "loop",
	LOOPZ:           "loopz",
	LOOPNZ:          "loopnz",
	JCXZ:            "jcxz",
}

func (o mnemonic) String() string {
	return mnemonicToString[o]
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

var REGTable = [...][2]register{
	0b000: {AL, AX},
	0b001: {CL, CX},
	0b010: {DL, DX},
	0b011: {BL, BX},
	0b100: {AH, SP},
	0b101: {CH, BP},
	0b110: {DH, SI},
	0b111: {BH, DI},
}

var EACTable = [...][2]register{
	0b000: {BX, SI},
	0b001: {BX, DI},
	0b010: {BP, SI},
	0b011: {BP, DI},
	0b100: {SI},
	0b101: {DI},
	0b110: {BP},
	0b111: {BX},
}

// 0b111 — reg1 + reg2 + disp
// 0b101 — reg1 + disp
// 0b100 — reg1
// 0b000 — DIRECT ACCESS
var EACFormTable = [...][3]uint8{
	0b000: {0b110, 0b111, 0b111},
	0b001: {0b110, 0b111, 0b111},
	0b010: {0b110, 0b111, 0b111},
	0b011: {0b110, 0b111, 0b111},
	0b100: {0b100, 0b101, 0b101},
	0b101: {0b100, 0b101, 0b101},
	0b110: {0b000, 0b101, 0b101},
	0b111: {0b100, 0b101, 0b101},
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
		switch inst.mnemonic {
		case
			JNZ,
			JE,
			JL,
			JLE,
			JB,
			JBE,
			JP,
			JO,
			JS,
			JNE,
			JNL,
			JG,
			JNB,
			JA,
			JNP,
			JNO,
			JNS,
			LOOP,
			LOOPZ,
			LOOPNZ,
			JCXZ:
			if inst.jump+2 > 0 {
				fmt.Fprintf(&out, "%s $+%d+0", inst.mnemonic, inst.jump+2)
			} else if inst.jump+2 == 0 {
				fmt.Fprintf(&out, "%s $+0", inst.mnemonic)
			} else {
				fmt.Fprintf(&out, "%s $%d+0", inst.mnemonic, inst.jump+2)
			}
		case MOV:
			if inst.dst.kind == opKindEAC && inst.src.kind == opKindImm {
				if inst.src.imm.word {
					fmt.Fprintf(&out, "%s %s, word %s", inst.mnemonic, inst.dst, inst.src)
					break
				}
				fmt.Fprintf(&out, "%s %s, byte %s", inst.mnemonic, inst.dst, inst.src)
				break
			}
			fmt.Fprintf(&out, "%s %s, %s", inst.mnemonic, inst.dst, inst.src)
		default:
			if inst.dst.kind == opKindEAC && inst.src.kind == opKindImm {
				if inst.src.imm.word {
					fmt.Fprintf(&out, "%s word %s, %s", inst.mnemonic, inst.dst, inst.src)
					break
				}
				fmt.Fprintf(&out, "%s byte %s, %s", inst.mnemonic, inst.dst, inst.src)
				break
			}
			fmt.Fprintf(&out, "%s %s, %s", inst.mnemonic, inst.dst, inst.src)
		}
	}

	return out.String(), outErr
}

type instruction struct {
	mnemonic mnemonic
	jump     int8
	dst      operand
	src      operand
}

type operand struct {
	kind operandKind
	reg  register
	imm  struct {
		val  int16
		word bool
	}
	eac struct {
		form     uint8
		reg1     register
		reg2     register
		dispOrDA int16
	}
}

type operandKind int

const (
	opKindReg operandKind = iota + 1
	opKindImm
	opKindEAC
)

func operandReg(reg register) (o operand) {
	o.kind = opKindReg
	o.reg = reg
	return
}

func operandImm(val int16, word bool) (o operand) {
	o.kind = opKindImm
	o.imm.val = val
	o.imm.word = word
	return
}

func operandEAC(form uint8, disp int16, regs ...register) (o operand) {
	o.kind = opKindEAC
	o.eac.form = form
	o.eac.dispOrDA = disp
	if len(regs) > 1 {
		o.eac.reg1 = regs[0]
		o.eac.reg2 = regs[1]
	}
	if len(regs) > 0 {
		o.eac.reg1 = regs[0]
	}
	return
}

func (o operand) String() string {
	switch o.kind {
	case opKindReg:
		return registerToString[o.reg]
	case opKindImm:
		return strconv.Itoa(int(o.imm.val))
	case opKindEAC:
		switch o.eac.form {
		case 0b000:
			return fmt.Sprintf("[%d]", o.eac.dispOrDA)
		case 0b100:
			return fmt.Sprintf("[%s]", o.eac.reg1)
		case 0b110:
			return fmt.Sprintf("[%s + %s]", o.eac.reg1, o.eac.reg2)
		case 0b101:
			if o.eac.dispOrDA < 0 {
				return fmt.Sprintf("[%s - %d]", o.eac.reg1, -o.eac.dispOrDA)
			} else if o.eac.dispOrDA > 0 {
				return fmt.Sprintf("[%s + %d]", o.eac.reg1, o.eac.dispOrDA)
			}
			return fmt.Sprintf("[%s]", o.eac.reg1)
		case 0b111:
			if o.eac.dispOrDA < 0 {
				return fmt.Sprintf("[%s + %s - %d]", o.eac.reg1, o.eac.reg2, -o.eac.dispOrDA)
			} else if o.eac.dispOrDA > 0 {
				return fmt.Sprintf("[%s + %s + %d]", o.eac.reg1, o.eac.reg2, o.eac.dispOrDA)
			}
			return fmt.Sprintf("[%s + %s]", o.eac.reg1, o.eac.reg2)
		default:
			panic(fmt.Sprintf("invalid form of EAC: %d", o.eac.form))
		}
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
		isEAC          bool
		isSrcImmediate bool
		isAddrToAcc    bool
		isAccToAddr    bool
		isImmToAcc     bool
		isJump         bool
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
		// TODO: sign bit
		s = -1
	)

	// NOTE: An instruction could from 1 to 6 byte in length
	b1 := stream[n]
	n++

	switch {
	// MOVs
	case b1>>2 == 0b100010:
		inst.mnemonic = MOV

		// b1
		d = int(b1 >> 1 & 0b1)
		w = int(b1 & 0b1)

		// b2
		b2 := stream[n]
		n++

		mod = int(b2 >> 6)
		reg = int(b2 >> 3 & 0b111)
		rm = int(b2 & 0b111)
	case b1>>4 == 0b1011:
		inst.mnemonic = MOV

		// b1
		w = int(b1 >> 3 & 0b1)
		reg = int(b1 & 0b111)

		// NOTE: knowledge encoded into this specific instruction: we should read data and put into src
		if w == 0 {
			readByteData = true
		} else {
			readWordData = true
		}
		isSrcImmediate = true
	case b1>>1 == 0b1100011:
		inst.mnemonic = MOV

		// b1
		w = int(b1 & 0b1)

		// b2
		b2 := stream[n]
		n++

		// NOTE: knowledge encoded into this specific instruction:
		// the dst should be register or memory and the src — immediate
		mod = int(b2 >> 6)
		rm = int(b2 & 0b111)

		// NOTE: knowledge encoded into this specific instruction: we should read data and put into src
		if w == 0 {
			readByteData = true
		} else {
			readWordData = true
		}
		isSrcImmediate = true
	case b1>>1 == 0b1010000:
		inst.mnemonic = MOV

		// b1
		w = int(b1 & 0b1)

		// NOTE: knowledge encoded into this specific instruction: we should read addr into src
		if w == 0 {
			readByteData = true
		} else {
			readWordData = true
		}
		isAddrToAcc = true
	case b1>>1 == 0b1010001:
		inst.mnemonic = MOV

		// b1
		w = int(b1 & 0b1)

		// NOTE: knowledge encoded into this specific instruction: we should read addr into dst
		if w == 0 {
			readByteData = true
		} else {
			readWordData = true
		}
		isAccToAddr = true

	// ADDs
	case b1>>2 == 0:
		inst.mnemonic = ADD

		// b1
		d = int(b1 >> 1 & 0b1)
		w = int(b1 & 0b1)

		// b2
		b2 := stream[n]
		n++

		mod = int(b2 >> 6)
		reg = int(b2 >> 3 & 0b111)
		rm = int(b2 & 0b111)
	case b1>>2 == 0b100000 && (stream[n]>>3&0b111) == 0b000:
		inst.mnemonic = ADD

		// b1
		s = int(b1 >> 1 & 0b1)
		w = int(b1 & 0b1)

		// b2
		b2 := stream[n]
		n++

		mod = int(b2 >> 6)
		rm = int(b2 & 0b111)

		// NOTE: knowledge encoded into this specific instruction: we should read data and put into src
		// FIXME
		if (w == 0 && s == 0) || (w == 1 && s == 1) {
			readByteData = true
		} else {
			readWordData = true
		}
		isSrcImmediate = true
	case b1>>1 == 0b10:
		inst.mnemonic = ADD

		// b1
		w = int(b1 & 0b1)

		// NOTE: knowledge encoded into this specific instruction: we should read addr into dst
		if w == 0 {
			readByteData = true
		} else {
			readWordData = true
		}
		isImmToAcc = true

	// SUBs
	case b1>>2 == 0b1010:
		inst.mnemonic = SUB

		// b1
		d = int(b1 >> 1 & 0b1)
		w = int(b1 & 0b1)

		// b2
		b2 := stream[n]
		n++

		mod = int(b2 >> 6)
		reg = int(b2 >> 3 & 0b111)
		rm = int(b2 & 0b111)
	case b1>>2 == 0b100000 && (stream[n]>>3&0b111) == 0b101:
		inst.mnemonic = SUB

		// b1
		s = int(b1 >> 1 & 0b1)
		w = int(b1 & 0b1)

		// b2
		b2 := stream[n]
		n++

		mod = int(b2 >> 6)
		rm = int(b2 & 0b111)

		// NOTE: knowledge encoded into this specific instruction: we should read data and put into src
		// FIXME
		if (w == 0 && s == 0) || (w == 1 && s == 1) {
			readByteData = true
		} else {
			readWordData = true
		}
		isSrcImmediate = true
	case b1>>1 == 0b10110:
		inst.mnemonic = SUB

		// b1
		w = int(b1 & 0b1)

		// NOTE: knowledge encoded into this specific instruction: we should read addr into dst
		if w == 0 {
			readByteData = true
		} else {
			readWordData = true
		}
		isImmToAcc = true

	// CMPs
	case b1>>2 == 0b1110:
		inst.mnemonic = CMP

		// b1
		d = int(b1 >> 1 & 0b1)
		w = int(b1 & 0b1)

		// b2
		b2 := stream[n]
		n++

		mod = int(b2 >> 6)
		reg = int(b2 >> 3 & 0b111)
		rm = int(b2 & 0b111)
	case b1>>2 == 0b100000 && (stream[n]>>3&0b111) == 0b111:
		inst.mnemonic = CMP

		// b1
		s = int(b1 >> 1 & 0b1)
		w = int(b1 & 0b1)

		// b2
		b2 := stream[n]
		n++

		mod = int(b2 >> 6)
		rm = int(b2 & 0b111)

		// NOTE: knowledge encoded into this specific instruction: we should read data and put into src
		// FIXME
		if (w == 0 && s == 0) || (w == 1 && s == 1) {
			readByteData = true
		} else {
			readWordData = true
		}
		isSrcImmediate = true
	case b1>>1 == 0b11110:
		inst.mnemonic = CMP

		// b1
		w = int(b1 & 0b1)

		// NOTE: knowledge encoded into this specific instruction: we should read addr into dst
		if w == 0 {
			readByteData = true
		} else {
			readWordData = true
		}
		isImmToAcc = true

	// JMPs
	default:
		jumps := map[byte]mnemonic{
			0b01110100: JE,
			0b01111100: JL,
			0b01111110: JLE,
			0b01110010: JB,
			0b01110110: JBE,
			0b01111010: JP,
			0b01110000: JO,
			0b01111000: JS,
			0b01110101: JNE,
			0b01111101: JNL,
			0b01111111: JG,
			0b01110011: JNB,
			0b01110111: JA,
			0b01111011: JNP,
			0b01110001: JNO,
			0b01111001: JNS,
			0b11100010: LOOP,
			0b11100001: LOOPZ,
			0b11100000: LOOPNZ,
			0b11100011: JCXZ,
		}

		if mnemonic, ok := jumps[b1]; ok {
			inst.mnemonic = mnemonic
			// NOTE: knowledge encoded into this specific instruction
			readByteData = true
			isJump = true
			break
		}

		err = fmt.Errorf("unsupported instruction: %b", b1)
		return
	}

	if s == 1 && w == 0 {
		err = errors.New("impossible combination s = 1 and w = 0")
		return
	}

	switch mod {
	case 0b00:
		isEAC = true
		if rm == 0b110 {
			readWordDisp = true
		}
	case 0b01:
		isEAC = true
		readByteDisp = true
	case 0b10:
		isEAC = true
		readWordDisp = true
	case 0b11:
		// NOTE: MOD can be use to identify a register of dst when src is immediate
		// HACK: I'm not sure it will work for all cases
		if !isSrcImmediate {
			isRegToReg = true
		}
	}

	var disp int16
	if readByteDisp {
		disp = int16(int8(stream[n]))
		n++
	}
	if readWordDisp {
		disp = int16(binary.LittleEndian.Uint16(stream[n:]))
		n += 2
	}

	var data int16
	if readByteData {
		data = int16(int8(stream[n]))
		n++
	}
	if readWordData {
		data = int16(binary.LittleEndian.Uint16(stream[n:]))
		n += 2
	}

	var eacForm uint8
	if mod != 0b11 && rm != -1 {
		eacForm = EACFormTable[rm][mod]
	}

	switch {
	// regToReg
	case isRegToReg:
		operand1 := operandReg(REGTable[rm][w])
		operand2 := operandReg(REGTable[reg][w])

		if d == 0 {
			inst.dst = operand1
			inst.src = operand2
		} else {
			inst.dst = operand2
			inst.src = operand1
		}
	// immToMemEAC
	case isEAC && isSrcImmediate:
		regs := EACTable[rm]

		inst.dst = operandEAC(eacForm, disp, regs[0], regs[1])
		inst.src = operandImm(data, w == 1)
	// memToReg or regToMem
	case isEAC:
		regs := EACTable[rm]

		operand1 := operandEAC(eacForm, disp, regs[0], regs[1])
		operand2 := operandReg(REGTable[reg][w])

		if d == 0 {
			inst.dst = operand1
			inst.src = operand2
		} else {
			inst.dst = operand2
			inst.src = operand1
		}
	// immToReg
	case isSrcImmediate:
		// HACK: it is sooo gross
		if reg == -1 && mod == 0b11 {
			inst.dst = operandReg(REGTable[rm][w])
		} else {
			inst.dst = operandReg(REGTable[reg][w])
		}
		inst.src = operandImm(data, w == 1)
	// accToAddr
	case isAccToAddr:
		inst.dst = operandEAC(eacForm, data)
		inst.src = operandReg(AX)
	// addrToAcc
	case isAddrToAcc:
		inst.dst = operandReg(AX)
		inst.src = operandEAC(eacForm, data)
	case isImmToAcc:
		if w == 1 {
			inst.dst = operandReg(AX)
			inst.src = operandImm(data, true)
		} else {
			inst.dst = operandReg(AL)
			inst.src = operandImm(data, false)
		}
	case isJump:
		inst.jump = int8(data)
	}

	return
}
