package cpu

import (
	"fmt"
	"log"
	"tgb/interrupt"
	"tgb/memory"
	// "tgb/interrupt"
)

type CPU struct {
	a uint8 // general register
	f uint8 // flag register
	b uint8 // general register
	c uint8 // general register
	d uint8 // general register
	e uint8 // general register
	h uint8 // general register
	l uint8 // general register

	sp uint16 // stack pointer
	pc uint16 // program counter

	ime   byte   // Interrupt Master Enable Flag
	ie    uint16 // Interrupt Enable
	iFlag uint16 // Interrupt Flag

	cycle int
}

type opcode uint8

func NewCPUinBoot() *CPU {
	cpu := &CPU {
		pc: 0x0000,
		cycle: 0,
	}
	return cpu
}

func NewCPU() *CPU {
	// only normal GB
	cpu := &CPU{
		a:     0x01,
		f:     0xB0,
		b:     0x00,
		c:     0x13,
		d:     0x00,
		e:     0xD8,
		h:     0x01,
		l:     0x4D,
		sp:    0xFFFE,
		pc:    0x0100,
		cycle: 0,
	}

	return cpu
}

// (fetch - decode - execute) 1 cycle
func (cpu *CPU) Step() int {
	inst, operands := cpu.fetch()
	cpu.decodeAndExecute(inst, operands)
	return opcodeCycles[inst]
}

func (cpu *CPU) fetch() (opcode, []uint8) {
	inst := opcode(cpu.read(cpu.pc))
	fmt.Printf("%s  ", opcodeLabels[inst])

	cpu.pc++

	l := operandLength(inst)
	operands := cpu.fetchOperands(l)

	return inst, operands
}

func (cpu *CPU) fetchOperands(len int) []uint8 {
	operands := []uint8{}
	for i := 0; i < len; i++ {
		operand := cpu.read(cpu.pc)
		cpu.pc++
		operands = append(operands, operand)
	}
	return operands
}

func (cpu *CPU) decodeAndExecute(inst opcode, operands []uint8) {
	switch inst {
	case 0x00: // NOP
	case 0x01: // LD BC, nn
		lsb := operands[0]
		msb := operands[1]
		cpu.set_bc(u8tou16(lsb, msb))

	case 0x02: // LD [BC], A
		cpu.write(cpu.bc(), cpu.a)

	case 0x03: // INC BC
		n := cpu.read(cpu.bc())
		cpu.write(cpu.bc(), n+1)

	case 0x04: // INC B
		cpu.modifyFlagsInIncOP(cpu.b+1, "INC")
		cpu.b += 1

	case 0x05: // DEC B
		cpu.modifyFlagsInIncOP(cpu.b-1, "DEC")
		cpu.b -= 1

	case 0x06: // LD B, n
		cpu.b = operands[0]

	case 0x07: // RLCA
		cpu.rlca()

	case 0x08: // LD [nn], SP
		lsb := operands[0]
		msb := operands[1]
		cpu.write(u8tou16(lsb, msb), cpu.read(cpu.sp))

	case 0x09: // ADD HL, BC
		cpu.modifyFlagsAddHL(int(cpu.hl()) + int(cpu.bc()))
		// cpu.write(cpu.hl(), cpu.hl()+cpu.bc())

	case 0x0A: // LD A, [BC]
		cpu.a = cpu.read(cpu.bc())

	case 0x0B: // DEC BC
		n := cpu.read(cpu.bc())
		cpu.write(cpu.bc(), n-1)

	case 0x0C: // INC C
		cpu.modifyFlagsInIncOP(cpu.c+1, "INC")
		cpu.c += 1

	case 0x0D: // DEC C
		cpu.modifyFlagsInIncOP(cpu.c-1, "DEC")
		cpu.c -= 1

	case 0x0E: // LD E, n
		cpu.c = operands[0]

	case 0x0F: // RRCA     4 000c rotate akku right
		cpu.rrca()

	case 0x11: // LD DE, nn
		lsb := operands[0]
		msb := operands[1]
		cpu.set_de(u8tou16(lsb, msb))

	case 0x12: // LD [DE], A
		cpu.write(cpu.de(), cpu.a)

	case 0x13: // INC DE
		n := cpu.read(cpu.de())
		cpu.write(cpu.de(), n+1)

	case 0x14: // INC D
		cpu.modifyFlagsInIncOP(cpu.d+1, "INC")
		cpu.d += 1

	case 0x15: // DEC D
		cpu.modifyFlagsInIncOP(cpu.d-1, "DEC")
		cpu.d -= 1

	case 0x16: // LD D, n
		cpu.d = operands[0]

	case 0x17: // RLA
		cpu.rla()

	case 0x18: // JR r
		r := operands[0]
		cpu.pc += uint16(int16(r))

	case 0x19: // ADD HL, DE
		cpu.modifyFlagsAddHL(int(cpu.hl()) + int(cpu.de()))
		// cpu.write(cpu.hl(), cpu.hl()+cpu.de())

	case 0x1A: // LD A, [DE]
		cpu.a = cpu.read(cpu.de())

	case 0x1B: // DEC DE
		n := cpu.read(cpu.de())
		cpu.write(cpu.de(), n-1)

	case 0x1C: // INC E
		cpu.modifyFlagsInIncOP(cpu.e+1, "INC")
		cpu.e += 1

	case 0x1D: // DEC E
		cpu.modifyFlagsInIncOP(cpu.e-1, "DEC")
		cpu.e -= 1

	case 0x1E: // LD E, n
		cpu.e = operands[0]

	case 0x1F: // RRA     4 000c rotate akku right through carry
		cpu.rra()

	case 0x20: // JR NZ, r
		r := operands[0]
		if !cpu.isZeroFlag() {
			cpu.pc += uint16(int16(r))
		}

	case 0x21: // LD HL, nn
		lsb := operands[0]
		msb := operands[1]
		cpu.set_hl(u8tou16(lsb, msb))

	case 0x22: // LDI [HL+], A
		cpu.write(cpu.hl(), cpu.a)
		cpu.set_hl(cpu.hl() + 1)

	case 0x23: // INC HL
		n := cpu.read(cpu.hl())
		cpu.write(cpu.hl(), n+1)

	case 0x24: // INC H
		cpu.modifyFlagsInIncOP(cpu.h+1, "INC")
		cpu.h += 1

	case 0x25: // DEC H
		cpu.modifyFlagsInIncOP(cpu.h-1, "DEC")
		cpu.h -= 1

	case 0x26: // LD H, n
		cpu.h = operands[0]

	case 0x28: // JR Z, r
		r := operands[0]
		if cpu.isZeroFlag() {
			cpu.pc += uint16(int16(r))
		}

	case 0x29: // ADD HL, HL
		cpu.modifyFlagsAddHL(int(cpu.hl()) + int(cpu.hl()))
		//cpu.write(cpu.hl(), cpu.hl()+cpu.hl())

	case 0x2A: // LDI A, [HL+]
		cpu.a = cpu.read(cpu.hl())
		cpu.set_hl(cpu.hl() + 1)

	case 0x2B: // DEC HL
		n := cpu.read(cpu.hl())
		cpu.write(cpu.hl(), n-1)

	case 0x2C: // INC L
		cpu.modifyFlagsInIncOP(cpu.l+1, "INC")
		cpu.l += 1

	case 0x2D: // DEC L
		cpu.modifyFlagsInIncOP(cpu.l-1, "DEC")
		cpu.l -= 1

	case 0x2E: // LD L, n
		cpu.l = operands[0]

	case 0x2F: // CPL
		cpu.a ^= 0xFF
		cpu.setSubFlag()
		cpu.setHalfCarryFlag()

	case 0x30: // JR NC, r
		r := operands[0]
		if !cpu.isCarryFlag() {
			cpu.pc += uint16(int16(r))
		}

	case 0x31: // LD SP, nn
		lsb := operands[0]
		msb := operands[1]
		cpu.set_sp(u8tou16(lsb, msb))

	case 0x32: // LDD [HL-], A
		cpu.write(cpu.hl(), cpu.a)
		cpu.set_hl(cpu.hl() - 1)

	case 0x33: // INC SP
		n := cpu.read(cpu.sp)
		cpu.write(cpu.sp, n+1)

	case 0x34: // INC [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlagsInIncOP(n+1, "INC")
		cpu.write(cpu.hl(), n+1)

	case 0x35: // DEC [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlagsInIncOP(n-1, "DEC")
		cpu.write(cpu.hl(), n+1)

	case 0x36: // LD [HL], n
		cpu.write(cpu.hl(), operands[0])

	case 0x37: // SCF
		cpu.clearSubFlag()
		cpu.clearHalfCarryFlag()
		cpu.setCarryFlag()

	case 0x38: // JR C, r
		r := operands[0]
		if cpu.isCarryFlag() {
			cpu.pc += uint16(int16(r))
		}

	case 0x39: // ADD HL, SP
		cpu.modifyFlagsAddHL(int(cpu.hl()) + int(cpu.sp))
		// cpu.write(cpu.hl(), cpu.hl()+cpu.sp)

	case 0x3A: // LDD A, [HL-]
		cpu.a = cpu.read(cpu.hl())
		cpu.set_hl(cpu.hl() - 1)

	case 0x3B: // DEC SP
		n := cpu.read(cpu.sp)
		cpu.write(cpu.sp, n-1)

	case 0x3C: // INC A
		cpu.modifyFlagsInIncOP(cpu.a+1, "INC")
		cpu.a += 1

	case 0x3D: // DEC A
		cpu.modifyFlagsInIncOP(cpu.a-1, "DEC")
		cpu.a -= 1

	case 0x3E: // LD A, n
		cpu.a = operands[0]

	case 0x3F: // CCF
		cpu.clearSubFlag()
		cpu.clearHalfCarryFlag()
		if cpu.isCarryFlag() {
			cpu.clearCarryFlag()
		} else {
			cpu.setCarryFlag()
		}

	case 0x40: // LD B, B
		cpu.b = cpu.b

	case 0x41: // LD B, C
		cpu.b = cpu.c

	case 0x42: // LD B, D
		cpu.b = cpu.d

	case 0x43: // LD B, E
		cpu.b = cpu.e

	case 0x44: // LD B, H
		cpu.b = cpu.h

	case 0x45: // LD B, L
		cpu.b = cpu.l

	case 0x46: // LD r, [HL]
		val := cpu.read(cpu.hl())
		cpu.b = val

	case 0x47: // LD B, A
		cpu.b = cpu.a

	case 0x48: // LD C, B
		cpu.c = cpu.b

	case 0x49: // LD C, C
		cpu.c = cpu.c

	case 0x4A: // LD C, D
		cpu.c = cpu.d

	case 0x4B: // LD C, E
		cpu.c = cpu.e

	case 0x4C: // LD C, H
		cpu.c = cpu.h

	case 0x4D: // LD C, L
		cpu.c = cpu.l

	case 0x4E: // LD r, [HL]
		val := cpu.read(cpu.hl())
		cpu.c = val

	case 0x4F: // LD C, A
		cpu.c = cpu.a

	case 0x50: // LD D, B
		cpu.d = cpu.b

	case 0x51: // LD D, C
		cpu.d = cpu.c

	case 0x52: // LD D, D
		cpu.d = cpu.d

	case 0x53: // LD D, E
		cpu.d = cpu.e

	case 0x54: // LD D, H
		cpu.d = cpu.h

	case 0x55: // LD D, L
		cpu.d = cpu.l

	case 0x56: // LD r, [HL]
		val := cpu.read(cpu.hl())
		cpu.d = val

	case 0x57: // LD D, A
		cpu.d = cpu.a

	case 0x58: // LD E, B
		cpu.e = cpu.b

	case 0x59: // LD E, C
		cpu.e = cpu.c

	case 0x5A: // LD E, D
		cpu.e = cpu.d

	case 0x5B: // LD E, E
		cpu.e = cpu.e

	case 0x5C: // LD E, H
		cpu.e = cpu.h

	case 0x5D: // LD E, L
		cpu.e = cpu.l

	case 0x5E: // LD r, [HL]
		val := cpu.read(cpu.hl())
		cpu.e = val

	case 0x5F: // LD E, A
		cpu.e = cpu.a

	case 0x60: // LD H, B
		cpu.h = cpu.b

	case 0x61: // LD H, C
		cpu.h = cpu.c

	case 0x62: // LD H, D
		cpu.h = cpu.d

	case 0x63: // LD H, E
		cpu.h = cpu.e

	case 0x64: // LD H, H
		cpu.h = cpu.h

	case 0x65: // LD H, L
		cpu.h = cpu.l

	case 0x66: // LD r, [HL]
		val := cpu.read(cpu.hl())
		cpu.h = val

	case 0x67: // LD H, A
		cpu.h = cpu.a

	case 0x68: // LD L, B
		cpu.l = cpu.b

	case 0x69: // LD L, C
		cpu.l = cpu.c

	case 0x6A: // LD L, D
		cpu.l = cpu.d

	case 0x6B: // LD L, E
		cpu.l = cpu.e

	case 0x6C: // LD L, H
		cpu.l = cpu.h

	case 0x6D: // LD L, L
		cpu.l = cpu.l

	case 0x6E: // LD r, [HL]
		val := cpu.read(cpu.hl())
		cpu.l = val

	case 0x6F: // LD L, A
		cpu.l = cpu.a

	case 0x70: // LD [HL], B
		cpu.write(cpu.hl(), cpu.b)

	case 0x71: // LD [HL], C
		cpu.write(cpu.hl(), cpu.c)

	case 0x72: // LD [HL], D
		cpu.write(cpu.hl(), cpu.d)

	case 0x73: // LD [HL], E
		cpu.write(cpu.hl(), cpu.e)

	case 0x74: // LD [HL], H
		cpu.write(cpu.hl(), cpu.h)

	case 0x75: // LD [HL], L
		cpu.write(cpu.hl(), cpu.l)

	case 0x77: // LD [HL], A
		cpu.write(cpu.hl(), cpu.a)

	case 0x78: // LD A, B
		cpu.a = cpu.b

	case 0x79: // LD A, C
		cpu.a = cpu.c

	case 0x7A: // LD A, D
		cpu.a = cpu.d

	case 0x7B: // LD A, E
		cpu.a = cpu.e

	case 0x7C: // LD A, H
		cpu.a = cpu.h

	case 0x7D: // LD A, L
		cpu.a = cpu.l

	case 0x7E: // LD r, [HL]
		val := cpu.read(cpu.hl())
		cpu.a = val

	case 0x7F: // LD A, A
		cpu.a = cpu.a

	case 0x80: // ADD A, B
		cpu.modifyFlags(int(cpu.a)+int(cpu.b), "+")
		cpu.a += cpu.b

	case 0x81: // ADD A, C
		cpu.modifyFlags(int(cpu.a)+int(cpu.c), "+")
		cpu.a += cpu.c

	case 0x82: // ADD A, D
		cpu.modifyFlags(int(cpu.a)+int(cpu.d), "+")
		cpu.a += cpu.d

	case 0x83: // ADD A, E
		cpu.modifyFlags(int(cpu.a)+int(cpu.e), "+")
		cpu.a += cpu.e

	case 0x84: // ADD A, H
		cpu.modifyFlags(int(cpu.a)+int(cpu.h), "+")
		cpu.a += cpu.h

	case 0x85: // ADD A, L
		cpu.modifyFlags(int(cpu.a)+int(cpu.l), "+")
		cpu.a += cpu.l

	case 0x86: // ADD A, [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlags(int(cpu.a)+int(n), "+")
		cpu.a += n

	case 0x87: // ADD A, A
		cpu.modifyFlags(int(cpu.a)+int(cpu.a), "+")
		cpu.a += cpu.a

	case 0x88: // ADC A, B
		cpu.modifyFlags(int(cpu.a)+int(cpu.b)+int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.b + cpu.getCarryFlag()

	case 0x89: // ADC A, C
		cpu.modifyFlags(int(cpu.a)+int(cpu.c)+int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.c + cpu.getCarryFlag()

	case 0x8A: // ADC A, D
		cpu.modifyFlags(int(cpu.a)+int(cpu.d)+int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.d + cpu.getCarryFlag()

	case 0x8B: // ADC A, E
		cpu.modifyFlags(int(cpu.a)+int(cpu.e)+int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.e + cpu.getCarryFlag()

	case 0x8C: // ADC A, H
		cpu.modifyFlags(int(cpu.a)+int(cpu.h)+int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.h + cpu.getCarryFlag()

	case 0x8D: // ADC A, L
		cpu.modifyFlags(int(cpu.a)+int(cpu.l)+int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.l + cpu.getCarryFlag()

	case 0x8E: // ADC A, [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlags(int(cpu.a)+int(n)+int(cpu.getCarryFlag()), "+")
		cpu.a += n + cpu.getCarryFlag()

	case 0x8F: // ADC A, A
		cpu.modifyFlags(int(cpu.a)+int(cpu.a)+int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.a + cpu.getCarryFlag()

	case 0x90: // SUB B
		cpu.modifyFlags(int(cpu.a)-int(cpu.b), "-")
		cpu.a -= cpu.b

	case 0x91: // SUB C
		cpu.modifyFlags(int(cpu.a)-int(cpu.c), "-")
		cpu.a -= cpu.c

	case 0x92: // SUB D
		cpu.modifyFlags(int(cpu.a)-int(cpu.d), "-")
		cpu.a -= cpu.d

	case 0x93: // SUB E
		cpu.modifyFlags(int(cpu.a)-int(cpu.e), "-")
		cpu.a -= cpu.e

	case 0x94: // SUB H
		cpu.modifyFlags(int(cpu.a)-int(cpu.h), "-")
		cpu.a -= cpu.h

	case 0x95: // SUB L
		cpu.modifyFlags(int(cpu.a)-int(cpu.l), "-")
		cpu.a -= cpu.l

	case 0x96: // SUB [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlags(int(cpu.a)-int(n), "-")
		cpu.a -= n

	case 0x97: // SUB A
		cpu.modifyFlags(int(cpu.a)-int(cpu.a), "-")
		cpu.a -= cpu.a

	case 0x98: // SBC A, B
		cpu.modifyFlags(int(cpu.a)-int(cpu.b)-int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.b + cpu.getCarryFlag()

	case 0x99: // SBC A, C
		cpu.modifyFlags(int(cpu.a)-int(cpu.c)-int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.c + cpu.getCarryFlag()

	case 0x9A: // SBC A, D
		cpu.modifyFlags(int(cpu.a)-int(cpu.d)-int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.d + cpu.getCarryFlag()

	case 0x9B: // SBC A, E
		cpu.modifyFlags(int(cpu.a)-int(cpu.e)-int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.e + cpu.getCarryFlag()

	case 0x9C: // SBC A, H
		cpu.modifyFlags(int(cpu.a)-int(cpu.h)-int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.h + cpu.getCarryFlag()

	case 0x9D: // SBC A, L
		cpu.modifyFlags(int(cpu.a)-int(cpu.l)-int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.l + cpu.getCarryFlag()

	case 0x9E: // SBC A, [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlags(int(cpu.a)-int(n)-int(cpu.getCarryFlag()), "-")
		cpu.a -= n + cpu.getCarryFlag()

	case 0x9F: // SBC A, A
		cpu.modifyFlags(int(cpu.a)-int(cpu.a)-int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.a + cpu.getCarryFlag()

	case 0xA0: // AND B
		cpu.modifyFlagsInAndOP(cpu.a & cpu.b)
		cpu.a &= cpu.b

	case 0xA1: // AND C
		cpu.modifyFlagsInAndOP(cpu.a & cpu.c)
		cpu.a &= cpu.c

	case 0xA2: // AND D
		cpu.modifyFlagsInAndOP(cpu.a & cpu.d)
		cpu.a &= cpu.d

	case 0xA3: // AND E
		cpu.modifyFlagsInAndOP(cpu.a & cpu.e)
		cpu.a &= cpu.e

	case 0xA4: // AND H
		cpu.modifyFlagsInAndOP(cpu.a & cpu.h)
		cpu.a &= cpu.h

	case 0xA5: // AND L
		cpu.modifyFlagsInAndOP(cpu.a & cpu.l)
		cpu.a &= cpu.l

	case 0xA6: // AND [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlagsInAndOP(cpu.a & n)
		cpu.a &= n

	case 0xA7: // AND A
		cpu.modifyFlagsInAndOP(cpu.a & cpu.a)
		cpu.a &= cpu.a

	case 0xA8: // XOR B
		cpu.modifyFlagsInOrOP(cpu.a ^ cpu.b)
		cpu.a ^= cpu.b

	case 0xA9: // XOR C
		cpu.modifyFlagsInOrOP(cpu.a ^ cpu.c)
		cpu.a ^= cpu.c

	case 0xAA: // XOR D
		cpu.modifyFlagsInOrOP(cpu.a ^ cpu.d)
		cpu.a ^= cpu.d

	case 0xAB: // XOR E
		cpu.modifyFlagsInOrOP(cpu.a ^ cpu.e)
		cpu.a ^= cpu.e

	case 0xAC: // XOR H
		cpu.modifyFlagsInOrOP(cpu.a ^ cpu.h)
		cpu.a ^= cpu.h

	case 0xAD: // XOR L
		cpu.modifyFlagsInOrOP(cpu.a ^ cpu.l)
		cpu.a ^= cpu.l

	case 0xAE: // XOR [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlagsInOrOP(cpu.a ^ n)
		cpu.a ^= n

	case 0xAF: // XOR A
		cpu.modifyFlagsInOrOP(cpu.a ^ cpu.a)
		cpu.a ^= cpu.a

	case 0xB0: // OR B
		cpu.modifyFlagsInOrOP(cpu.a | cpu.b)
		cpu.a |= cpu.b

	case 0xB1: // OR C
		cpu.modifyFlagsInOrOP(cpu.a | cpu.c)
		cpu.a |= cpu.c

	case 0xB2: // OR D
		cpu.modifyFlagsInOrOP(cpu.a | cpu.d)
		cpu.a |= cpu.d

	case 0xB3: // OR E
		cpu.modifyFlagsInOrOP(cpu.a | cpu.e)
		cpu.a |= cpu.e

	case 0xB4: // OR H
		cpu.modifyFlagsInOrOP(cpu.a | cpu.h)
		cpu.a |= cpu.h

	case 0xB5: // OR L
		cpu.modifyFlagsInOrOP(cpu.a | cpu.l)
		cpu.a |= cpu.l

	case 0xB6: // OR [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlagsInOrOP(cpu.a | n)
		cpu.a |= n

	case 0xB7: // OR A
		cpu.modifyFlagsInOrOP(cpu.a | cpu.a)
		cpu.a |= cpu.a

	case 0xB8: // CP B
		cpu.modifyFlagsInCP(cpu.b)

	case 0xB9: // CP C
		cpu.modifyFlagsInCP(cpu.c)

	case 0xBA: // CP D
		cpu.modifyFlagsInCP(cpu.d)

	case 0xBB: // CP E
		cpu.modifyFlagsInCP(cpu.e)

	case 0xBC: // CP H
		cpu.modifyFlagsInCP(cpu.h)

	case 0xBD: // CP L
		cpu.modifyFlagsInCP(cpu.l)

	case 0xBE: // CP [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlagsInCP(n)

	case 0xBF: // CP A
		cpu.modifyFlagsInCP(cpu.a)

	case 0xC0: // RET NZ
		if !cpu.isZeroFlag() {
			cpu.popPreservedPC()
		}

	case 0xC1: // POP BC
		l := cpu.read(cpu.sp)
		h := cpu.read(cpu.sp - 1)
		cpu.set_bc(u8tou16(l, h))
		cpu.sp += 2

	case 0xC2: // JP NZ, nn
		lsb := operands[0]
		msb := operands[1]
		if !cpu.isZeroFlag() {
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xC3: // JP nn
		lsb := operands[0]
		msb := operands[1]
		cpu.pc = u8tou16(lsb, msb)

	case 0xC4: // CALL NZ, nn
		lsb := operands[0]
		msb := operands[1]
		if !cpu.isZeroFlag() {
			cpu.PushCurrentPC()
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xC5: // PUSH BC
		cpu.write(cpu.sp-1, cpu.b)
		cpu.write(cpu.sp-2, cpu.c)
		cpu.sp -= 2

	case 0xC6: // ADD A, n
		n := operands[0]
		cpu.modifyFlags(int(cpu.a)+int(n), "+")
		cpu.a += n

	case 0xC7: // RST 0x00
		cpu.PushCurrentPC()
		cpu.pc = 0x0000

	case 0xC8: // RET Z
		if cpu.isZeroFlag() {
			cpu.popPreservedPC()
		}

	case 0xC9: // RET
		cpu.popPreservedPC()

	case 0xCA: // JP Z, nn
		lsb := operands[0]
		msb := operands[1]
		if cpu.isZeroFlag() {
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xCB: // PREFIX CB
		cpu.executeCBInst()

	case 0xCC: // CALL Z, nn
		lsb := operands[0]
		msb := operands[1]
		if cpu.isZeroFlag() {
			cpu.PushCurrentPC()
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xCD: // CALL nn
		l := operands[0]
		m := operands[1]
		nn := u8tou16(l, m)
		cpu.PushCurrentPC()
		cpu.pc = nn

	case 0xCE: // ADC A, n
		n := operands[0]
		cpu.modifyFlags(int(cpu.a)+int(n)+int(cpu.getCarryFlag()), "+")
		cpu.a += n + cpu.getCarryFlag()

	case 0xCF: // RST 0x08
		cpu.PushCurrentPC()
		cpu.pc = 0x0008

	case 0xD0: // RET NC
		if !cpu.isCarryFlag() {
			cpu.popPreservedPC()
		}

	case 0xD1: // POP DE
		l := cpu.read(cpu.sp)
		h := cpu.read(cpu.sp - 1)
		cpu.set_de(u8tou16(l, h))
		cpu.sp += 2

	case 0xD2: // JP NC, nn
		lsb := operands[0]
		msb := operands[1]
		if !cpu.isCarryFlag() {
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xD3: // EMPTY
		invalidInst()

	case 0xD4: // CALL NC, nn
		lsb := operands[0]
		msb := operands[1]
		if !cpu.isCarryFlag() {
			cpu.PushCurrentPC()
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xD5: // PUSH DE
		cpu.write(cpu.sp-1, cpu.d)
		cpu.write(cpu.sp-2, cpu.e)
		cpu.sp -= 2

	case 0xD6: // SUB n
		n := operands[0]
		cpu.modifyFlags(int(cpu.a)-int(n), "-")
		cpu.a -= n

	case 0xD7: // RST 0x10
		cpu.PushCurrentPC()
		cpu.pc = 0x0010

	case 0xD8: // RET C
		if cpu.isCarryFlag() {
			cpu.popPreservedPC()
		}

	case 0xD9: // RETI
		cpu.popPreservedPC()
		interrupt.SetIMEFlag()

	case 0xDA: // JP C, nn
		lsb := operands[0]
		msb := operands[1]
		if cpu.isCarryFlag() {
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xDB: // EMPTY
		invalidInst()

	case 0xDC: // CALL C, nn
		lsb := operands[0]
		msb := operands[1]
		if cpu.isCarryFlag() {
			cpu.PushCurrentPC()
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xDD: // EMPTY
		invalidInst()

	case 0xDE: // SBC A, n
		n := operands[0]
		cpu.modifyFlags(int(cpu.a)-int(n)-int(cpu.getCarryFlag()), "-")
		cpu.a -= n + cpu.getCarryFlag()

	case 0xDF: // RST 0x18
		cpu.PushCurrentPC()
		cpu.pc = 0x0018

	case 0xE0: // LD [FF00+n], A
		addr := u8tou16(operands[0], 0xFF)
		cpu.write(addr, cpu.a)

	case 0xE1: // POP HL
		l := cpu.read(cpu.sp)
		h := cpu.read(cpu.sp - 1)
		cpu.set_hl(u8tou16(l, h))
		cpu.sp += 2

	case 0xE2: // LD [FF00+C], A
		addr := u8tou16(cpu.c, 0xFF)
		cpu.write(addr, cpu.a)

	case 0xE3: // EMPTY
		invalidInst()

	case 0xE4: // EMPTY
		invalidInst()

	case 0xE5: // PUSH HL
		cpu.write(cpu.sp-1, cpu.h)
		cpu.write(cpu.sp-2, cpu.l)
		cpu.sp -= 2

	case 0xE6: // AND n
		n := operands[0]
		cpu.modifyFlagsInAndOP(cpu.a & n)
		cpu.a &= n

	case 0xE7: // RST 0x20
		cpu.PushCurrentPC()
		cpu.pc = 0x0020

	case 0xE9: // JP HL
		cpu.pc = cpu.hl()

	case 0xEA: // LD [nn], A
		lsb := operands[0]
		msb := operands[1]
		addr := u8tou16(lsb, msb)
		cpu.write(addr, cpu.a)

	case 0xEB: // EMPTY
		invalidInst()

	case 0xEC: // EMPTY
		invalidInst()

	case 0xED: // EMPTY
		invalidInst()

	case 0xEE: // XOR n
		n := operands[0]
		cpu.modifyFlagsInOrOP(cpu.a ^ n)
		cpu.a ^= n

	case 0xEF: // RST 0x28
		cpu.PushCurrentPC()
		cpu.pc = 0x0028

	case 0xF0: // LD A, [FF00+n]
		addr := u8tou16(operands[0], 0xFF)
		cpu.a = cpu.read(addr)

	case 0xF1: // POP AF
		l := cpu.read(cpu.sp)
		h := cpu.read(cpu.sp - 1)
		cpu.set_af(u8tou16(l, h))
		cpu.sp += 2

	case 0xF2: // LD A, [FF00+C]
		addr := u8tou16(cpu.c, 0xFF)
		cpu.a = cpu.read(addr)

	case 0xF3: // DI
		interrupt.ClearIMEFlag()

	case 0xF4: // EMPTY
		invalidInst()

	case 0xF5: // PUSH AF
		cpu.write(cpu.sp-1, cpu.a)
		cpu.write(cpu.sp-2, cpu.f)
		cpu.sp -= 2

	case 0xF6: // OR n
		n := operands[0]
		cpu.modifyFlagsInOrOP(cpu.a | n)
		cpu.a |= n

	case 0xF7: // RST 0x30
		cpu.PushCurrentPC()
		cpu.pc = 0x0030

	case 0xF9: // LD SP, HL
		cpu.sp = cpu.hl()

	case 0xFA: // LD A, [nn]
		lsb := operands[0]
		msb := operands[1]
		cpu.a = cpu.read(u8tou16(lsb, msb))

	case 0xFB: // EI
		interrupt.SetIMEFlag()

	case 0xFC: // EMPTY
		invalidInst()

	case 0xFD: // EMPTY
		invalidInst()

	case 0xFE: // CP n
		n := operands[0]
		cpu.modifyFlagsInCP(n)

	case 0xFF: // RST 0x38
		cpu.PushCurrentPC()
		cpu.pc = 0x0038

	case 0x10: // STOP

	case 0x27: // DAA
		// Decimal adjust register A.
		// This instruction adjusts register A so that the
		// correct representation of Binary Coded Decimal (BCD)
		// is obtained.
		// Z, C = star?
		if cpu.a == 0 {
			cpu.setZeroFlag()
		} else {
			cpu.clearZeroFlag()
		}
		cpu.clearHalfCarryFlag()

	case 0x76: // HALT
		if interrupt.IME {
			// (IE & IF & 1F) != 0 となるまで、CPUは停止される
		} else {

		}

	case 0xE8: // ADD SP, r - PENDING
		r := operands[0]
		cpu.clearZeroFlag()
		cpu.clearSubFlag()

		cpu.sp += uint16(int16(r))

	case 0xF8: // LD HL, SP+r8 - PENDING
	}
}

func (cpu *CPU) read(addr uint16) uint8 {
	return memory.Read(addr)
}

func (cpu *CPU) write(addr uint16, val uint8) {
	memory.Write(addr, val)
}

func (cpu *CPU) af() uint16 {
	return u8tou16(cpu.a, cpu.f)
}

func (cpu *CPU) bc() uint16 {
	return u8tou16(cpu.b, cpu.c)
}

func (cpu *CPU) de() uint16 {
	return u8tou16(cpu.d, cpu.e)
}

func (cpu *CPU) hl() uint16 {
	return u8tou16(cpu.h, cpu.l)
}

func (cpu *CPU) setA(val uint8) {
	cpu.a = val
}

func (cpu *CPU) setB(val uint8) {
	cpu.b = val
}

func (cpu *CPU) setC(val uint8) {
	cpu.c = val
}

func (cpu *CPU) setD(val uint8) {
	cpu.d = val
}

func (cpu *CPU) setH(val uint8) {
	cpu.h = val
}

func (cpu *CPU) setL(val uint8) {
	cpu.l = val
}

func (cpu *CPU) set_af(val uint16) {
	cpu.a = msb(val)
	cpu.f = lsb(val)
}

func (cpu *CPU) set_bc(val uint16) {
	cpu.b = msb(val)
	cpu.c = lsb(val)
}

func (cpu *CPU) set_de(val uint16) {
	cpu.d = msb(val)
	cpu.e = lsb(val)
}

func (cpu *CPU) set_hl(val uint16) {
	cpu.h = msb(val)
	cpu.l = lsb(val)
}

func (cpu *CPU) set_sp(val uint16) {
	cpu.sp = val
}

// Most Significant Byte
func msb(bytes uint16) uint8 {
	return uint8(bytes >> 8)
}

// Least Significant Byte
func lsb(bytes uint16) uint8 {
	return uint8(bytes & 0x0F)
}

// little endian
func u8tou16(u8, v8 uint8) uint16 {
	return (uint16(v8) << 8) | uint16(u8)
}

func (cpu *CPU) isZeroFlag() bool {
	return cpu.f&0x80 == 1
}

func (cpu *CPU) isSubFlag() bool {
	return cpu.f&0x40 == 1
}

func (cpu *CPU) isHalfCarryFlag() bool {
	return cpu.f&0x20 == 1
}

func (cpu *CPU) isCarryFlag() bool {
	return cpu.f&0x10 == 1
}

func (cpu *CPU) setZeroFlag() {
	cpu.f |= 0x80
}

func (cpu *CPU) setSubFlag() {
	cpu.f |= 0x40
}

func (cpu *CPU) setHalfCarryFlag() {
	cpu.f |= 0x20
}

func (cpu *CPU) setCarryFlag() {
	cpu.f |= 0x10
}

// Carry Flagが立ってたら 1, そうでなければ 0 を返す
func (cpu *CPU) getCarryFlag() uint8 {
	return uint8((cpu.f >> 4) & 0x01)
}

func (cpu *CPU) clearZeroFlag() {
	cpu.f &= 0x70
}

func (cpu *CPU) clearSubFlag() {
	cpu.f &= 0xB0
}

func (cpu *CPU) clearHalfCarryFlag() {
	cpu.f &= 0xD0
}

func (cpu *CPU) clearCarryFlag() {
	cpu.f &= 0xE0
}

func (cpu *CPU) PushCurrentPC() {
	cpu.sp--
	cpu.write(cpu.sp, msb(cpu.pc))
	cpu.sp--
	cpu.write(cpu.sp, lsb(cpu.pc))
}

func (cpu *CPU) popPreservedPC() {
	lsb := cpu.read(cpu.sp)
	cpu.sp++
	msb := cpu.read(cpu.sp)
	cpu.sp++
	cpu.pc = u8tou16(lsb, msb)
}

// Z N H C
// Z 0 1 0
func (cpu *CPU) modifyFlagsInAndOP(res uint8) {
	if res == 0 {
		cpu.setZeroFlag()
	} else {
		cpu.clearZeroFlag()
	}
	cpu.clearSubFlag()
	cpu.setHalfCarryFlag()
	cpu.clearCarryFlag()
}

func (cpu *CPU) modifyFlagsInOrOP(res uint8) {
	if res == 0 {
		cpu.setZeroFlag()
	} else {
		cpu.clearZeroFlag()
	}
	cpu.clearSubFlag()
	cpu.clearHalfCarryFlag()
	cpu.clearCarryFlag()
}

func (cpu *CPU) modifyFlagsInCP(val uint8) {
	cpu.setSubFlag()

	if cpu.a == val {
		cpu.setZeroFlag()
	} else {
		cpu.clearZeroFlag()
	}

	if cpu.a < val {
		cpu.setCarryFlag()
	} else {
		cpu.clearCarryFlag()
	}

	if (cpu.a & 0x0F) < (val & 0x0F) {
		cpu.setHalfCarryFlag()
	} else {
		cpu.clearHalfCarryFlag()
	}
}

// Z 0 H -
func (cpu *CPU) modifyFlagsInIncOP(res uint8, op string) {
	if op == "INC" {
		cpu.clearSubFlag()
	} else {
		cpu.setSubFlag()
	}

	if res == 0 {
		cpu.setZeroFlag()
	} else {
		cpu.clearZeroFlag()
	}

	if res > 0x0F {
		cpu.setHalfCarryFlag()
	} else {
		cpu.clearHalfCarryFlag()
	}
}

func (cpu *CPU) modifyFlagsAddHL(val int) {
	cpu.clearSubFlag()

	if val <= 0xFFF {
		cpu.clearHalfCarryFlag()
		cpu.clearCarryFlag()
	} else if 0xFFF < val && val <= 0xFFFF {
		cpu.setHalfCarryFlag()
		cpu.clearCarryFlag()
	} else {
		cpu.setHalfCarryFlag()
		cpu.setCarryFlag()
	}
}

func (cpu *CPU) modifyFlags(res int, op string) {
	switch op {
	case "+":
		cpu.clearSubFlag()
	case "-":
		cpu.setSubFlag()
	}

	if res == 0 {
		cpu.setZeroFlag()
		cpu.clearHalfCarryFlag()
		cpu.clearCarryFlag()
	} else if 0x0F < res && res <= 0xFF {
		cpu.clearZeroFlag()
		cpu.setHalfCarryFlag()
		cpu.clearCarryFlag()
	} else if 0xFF < res {
		cpu.clearZeroFlag()
		cpu.clearHalfCarryFlag()
		cpu.setCarryFlag()
	}
}

func (cpu *CPU) rlca() {
	lShifted := cpu.a << 1

	// Aの最上位ビットが1の場合、桁あふれした1を最下位ビットにつける
	if cpu.a&0x80 == 0x80 {
		cpu.setCarryFlag()
		lShifted ^= 0x01

		// Aの最上位ビットが0の場合
	} else {
		cpu.clearCarryFlag()
	}

	cpu.clearZeroFlag()
	cpu.clearSubFlag()
	cpu.clearHalfCarryFlag()
	cpu.a = lShifted
}

// 最下位ビットが立っていたら、桁あふれした1を最下位ビットにつける
func (cpu *CPU) rla() {
	lShifted := cpu.a << 1

	// Aの最上位ビットが1の場合
	if cpu.a&0x80 == 0x80 {
		if cpu.isCarryFlag() {
			lShifted ^= 0x01
		}
		cpu.setCarryFlag()

		// Aの最上位ビットが0の場合
	} else {
		cpu.clearCarryFlag()
	}

	cpu.clearZeroFlag()
	cpu.clearSubFlag()
	cpu.clearHalfCarryFlag()
	cpu.a = lShifted
}

func (cpu *CPU) rrca() {
	rShifted := cpu.a >> 1

	// Aの最下位ビットが1の場合、桁あふれした1を最上位ビットにつける
	if cpu.a&0x01 == 0x01 {
		cpu.setCarryFlag()
		rShifted ^= 0x80

		// Aの最下位ビットが0の場合
	} else {
		cpu.clearCarryFlag()
	}

	cpu.clearZeroFlag()
	cpu.clearSubFlag()
	cpu.clearHalfCarryFlag()
	cpu.a = rShifted
}

func (cpu *CPU) rra() {
	rShifted := cpu.a >> 1

	// Aの最下位ビットが1の場合
	if cpu.a&0x01 == 0x01 {
		if cpu.isCarryFlag() {
			rShifted ^= 0x80
		}
		cpu.setCarryFlag()

		// Aの最下位ビットが0の場合
	} else {
		cpu.clearCarryFlag()
	}

	cpu.clearZeroFlag()
	cpu.clearSubFlag()
	cpu.clearHalfCarryFlag()
	cpu.a = rShifted
}

func (cpu *CPU) executeCBInst() {

}

func invalidInst() {
	log.Println("invalid instruction")
}
