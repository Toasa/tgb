package cpu

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

var ramOrigin *[0x10000]byte

func New(org *[0x10000]byte) *CPU {
	ramOrigin = org
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
	decodedInst := cpu.decodeAndExecute(inst, operands)
	return opcodeCycles[inst]
}

func (cpu *CPU) fetch() (opcode, []uint8) {
	inst := opcode(cpu.read(cpu.pc))
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

	case 0x06: // LD B, n
		cpu.b = operands[0]

	case 0x0A: // LD A, [BC]
		cpu.a = cpu.read(cpu.bc())

	case 0x0E: // LD E, n
		cpu.c = operands[0]

	case 0x11: // LD DE, nn
		lsb := operands[0]
		msb := operands[1]
		cpu.set_de(u8tou16(lsb, msb))

	case 0x12: // LD [DE], A
		cpu.write(cpu.de(), cpu.a)

	case 0x16: // LD D, n
		cpu.d = operands[0]

	case 0x18: // JR r
		r := operands[0]
		cpu.pc += uint16(int16(r))

	case 0x1A: // LD A, [DE]
		cpu.a = cpu.read(cpu.de())

	case 0x1E: // LD E, n
		cpu.e = operands[0]

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

	case 0x26: // LD H, n
		cpu.h = operands[0]

	case 0x28: // JR Z, r
		r := operands[0]
		if cpu.isZeroFlag() {
			cpu.pc += uint16(int16(r))
		}

	case 0x2A: // LDI A, [HL+]
		cpu.a = cpu.read(cpu.hl())
		cpu.set_hl(cpu.hl() + 1)

	case 0x2E: // LD L, n
		cpu.l = operands[0]

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

	case 0x3A: // LDD A, [HL-]
		cpu.a = cpu.read(cpu.hl())
		cpu.set_hl(cpu.hl() - 1)



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
		cpu.modifyFlags(int(cpu.a) + int(cpu.b), "+")
		cpu.a += cpu.b

	case 0x81: // ADD A, C
		cpu.modifyFlags(int(cpu.a) + int(cpu.c), "+")
		cpu.a += cpu.c

	case 0x82: // ADD A, D
		cpu.modifyFlags(int(cpu.a) + int(cpu.d), "+")
		cpu.a += cpu.d

	case 0x83: // ADD A, E
		cpu.modifyFlags(int(cpu.a) + int(cpu.e), "+")
		cpu.a += cpu.e

	case 0x84: // ADD A, H
		cpu.modifyFlags(int(cpu.a) + int(cpu.h), "+")
		cpu.a += cpu.h

	case 0x85: // ADD A, L
		cpu.modifyFlags(int(cpu.a) + int(cpu.l), "+")
		cpu.a += cpu.l

	case 0x86: // ADD A, [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlags(int(cpu.a) + int(n), "+")
		cpu.a += n

	case 0x87: // ADD A, A
		cpu.modifyFlags(int(cpu.a) + int(cpu.a), "+")
		cpu.a += cpu.a

	case 0x88: // ADC A, B
		cpu.modifyFlags(int(cpu.a) + int(cpu.b) + int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.b + cpu.getCarryFlag()

	case 0x89: // ADC A, C
		cpu.modifyFlags(int(cpu.a) + int(cpu.c) + int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.c + cpu.getCarryFlag()

	case 0x8A: // ADC A, D
		cpu.modifyFlags(int(cpu.a) + int(cpu.d) + int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.d + cpu.getCarryFlag()

	case 0x8B: // ADC A, E
		cpu.modifyFlags(int(cpu.a) + int(cpu.e) + int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.e + cpu.getCarryFlag()

	case 0x8C: // ADC A, H
		cpu.modifyFlags(int(cpu.a) + int(cpu.h) + int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.h + cpu.getCarryFlag()

	case 0x8D: // ADC A, L
		cpu.modifyFlags(int(cpu.a) + int(cpu.l) + int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.l + cpu.getCarryFlag()

	case 0x8E: // ADC A, [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlags(int(cpu.a) + int(n) + int(cpu.getCarryFlag()), "+")
		cpu.a += n + cpu.getCarryFlag()

	case 0x8F: // ADC A, A
		cpu.modifyFlags(int(cpu.a) + int(cpu.a) + int(cpu.getCarryFlag()), "+")
		cpu.a += cpu.a + cpu.getCarryFlag()

	case 0x90: // SUB B
		cpu.modifyFlags(int(cpu.a) - int(cpu.b), "-")
		cpu.a -= cpu.b

	case 0x91: // SUB C
		cpu.modifyFlags(int(cpu.a) - int(cpu.c), "-")
		cpu.a -= cpu.c

	case 0x92: // SUB D
		cpu.modifyFlags(int(cpu.a) - int(cpu.d), "-")
		cpu.a -= cpu.d

	case 0x93: // SUB E
		cpu.modifyFlags(int(cpu.a) - int(cpu.e), "-")
		cpu.a -= cpu.e

	case 0x94: // SUB H
		cpu.modifyFlags(int(cpu.a) - int(cpu.h), "-")
		cpu.a -= cpu.h

	case 0x95: // SUB L
		cpu.modifyFlags(int(cpu.a) - int(cpu.l), "-")
		cpu.a -= cpu.l

	case 0x96: // SUB [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlags(int(cpu.a) - int(n), "-")
		cpu.a -= n

	case 0x97: // SUB A
		cpu.modifyFlags(int(cpu.a) - int(cpu.a), "-")
		cpu.a -= cpu.a

	case 0x98: // SBC A, B
		cpu.modifyFlags(int(cpu.a) - int(cpu.b) - int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.b + cpu.getCarryFlag()

	case 0x99: // SBC A, C
		cpu.modifyFlags(int(cpu.a) - int(cpu.c) - int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.c + cpu.getCarryFlag()

	case 0x9A: // SBC A, D
		cpu.modifyFlags(int(cpu.a) - int(cpu.d) - int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.d + cpu.getCarryFlag()

	case 0x9B: // SBC A, E
		cpu.modifyFlags(int(cpu.a) - int(cpu.e) - int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.e + cpu.getCarryFlag()

	case 0x9C: // SBC A, H
		cpu.modifyFlags(int(cpu.a) - int(cpu.h) - int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.h + cpu.getCarryFlag()

	case 0x9D: // SBC A, L
		cpu.modifyFlags(int(cpu.a) - int(cpu.l) - int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.l + cpu.getCarryFlag()

	case 0x9E: // SBC A, [HL]
		n := cpu.read(cpu.hl())
		cpu.modifyFlags(int(cpu.a) - int(n) - int(cpu.getCarryFlag()), "-")
		cpu.a -= n + cpu.getCarryFlag()

	case 0x9F: // SBC A, A
		cpu.modifyFlags(int(cpu.a) - int(cpu.a) - int(cpu.getCarryFlag()), "-")
		cpu.a -= cpu.a + cpu.getCarryFlag()

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
			cpu.pushCurrentPC()
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xC5: // PUSH BC
		cpu.write(cpu.sp-1, cpu.b)
		cpu.write(cpu.sp-2, cpu.c)
		cpu.sp -= 2

	case 0xC6: // ADD A, n
		n := operands[0]
		cpu.modifyFlags(int(cpu.a) + int(n), "+")
		cpu.a += n

	case 0xC7: // RST 0x00
		cpu.pushCurrentPC()
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

	case 0xCC: // CALL Z, nn
		lsb := operands[0]
		msb := operands[1]
		if cpu.isZeroFlag() {
			cpu.pushCurrentPC()
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xCD: // CALL nn
		l := operands[0]
		m := operands[1]
		nn := u8tou16(l, m)
		cpu.pushCurrentPC()
		cpu.pc = nn

	case 0xCE: // ADC A, n
		n := operands[0]
		cpu.modifyFlags(int(cpu.a) + int(n) + int(cpu.getCarryFlag()), "+")
		cpu.a += n + cpu.getCarryFlag()

	case 0xCF: // RST 0x08
		cpu.pushCurrentPC()
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

	case 0xD4: // CALL NC, nn
		lsb := operands[0]
		msb := operands[1]
		if !cpu.isCarryFlag() {
			cpu.pushCurrentPC()
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xD5: // PUSH DE
		cpu.write(cpu.sp-1, cpu.d)
		cpu.write(cpu.sp-2, cpu.e)
		cpu.sp -= 2

	case 0xD6: // SUB n
		n := operands[0]
		cpu.modifyFlags(int(cpu.a) - int(n), "-")
		cpu.a -= n

	case 0xD7: // RST 0x10
		cpu.pushCurrentPC()
		cpu.pc = 0x0010

	case 0xD8: // RET C
		if cpu.isCarryFlag() {
			cpu.popPreservedPC()
		}

	case 0xD9: // RETI
		cpu.popPreservedPC()
		cpu.setIMEFlag()

	case 0xDA: // JP C, nn
		lsb := operands[0]
		msb := operands[1]
		if cpu.isCarryFlag() {
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xDC: // CALL C, nn
		lsb := operands[0]
		msb := operands[1]
		if cpu.isCarryFlag() {
			cpu.pushCurrentPC()
			cpu.pc = u8tou16(lsb, msb)
		}

	case 0xDE: // SBC A, n
		n := operands[0]
		cpu.modifyFlags(int(cpu.a) - int(n) - int(cpu.getCarryFlag()), "-")
		cpu.a -= n + cpu.getCarryFlag()

	case 0xDF: // RST 0x18
		cpu.pushCurrentPC()
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

	case 0xE5: // PUSH HL
		cpu.write(cpu.sp-1, cpu.h)
		cpu.write(cpu.sp-2, cpu.l)
		cpu.sp -= 2

	case 0xE7: // RST 0x20
		cpu.pushCurrentPC()
		cpu.pc = 0x0020

	case 0xE9: // JP HL
		cpu.pc = cpu.hl()

	case 0xEA: // LD [nn], A
		lsb := operands[0]
		msb := operands[1]
		addr := u8tou16(lsb, msb)
		cpu.write(addr, cpu.a)

	case 0xEF: // RST 0x28
		cpu.pushCurrentPC()
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

	case 0xF5: // PUSH AF
		cpu.write(cpu.sp-1, cpu.a)
		cpu.write(cpu.sp-2, cpu.f)
		cpu.sp -= 2

	case 0xF7: // RST 0x30
		cpu.pushCurrentPC()
		cpu.pc = 0x0030

	case 0xF9: // LD SP, HL
		cpu.sp = cpu.hl()

	case 0xFA: // LD A, [nn]
		lsb := operands[0]
		msb := operands[1]
		cpu.a = cpu.read(u8tou16(lsb, msb))

	case 0xFF: // RST 0x38
		cpu.pushCurrentPC()
		cpu.pc = 0x0038





	// HALT
	// STOP
	
	case 0xF3: // DI(Disable interrupt)
		cpu.clearIMEFlag()

	
	case 0xFB: // EI
		IME_scheduled = true
	

	
	case 0x3F: // CCF
		cpu.clearSubFlag()
		cpu.clearHalfCarryFlag()
		if cpu.isCarryFlag() {
			cpu.clearCarryFlag()
		} else {
			cpu.setCarryFlag()
		}

	// DAA
	case 0x27:
		// Z, C = star?
		cpu.clearHalfCarryFlag()
	// CPL
	case 0x2F:
		cpu.a ^= 0xFF
		cpu.setSubFlag()
		cpu.setHalfCarryFlag()
}

func (cpu *CPU) read(addr uint16) uint8 {
	return (*ramOrigin)[addr]
}

func (cpu *CPU) write(addr uint16, val uint8) {
	*ramOrigin[addr] = val
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

func (cpu *CPU) pushCurrentPC() {
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



func (cpu *CPU) ld_rr_nn(inst opcode, nn uint16) {
	switch inst {
	case 0x01:
		cpu.set_bc(nn)
	case 0x11:
		cpu.set_de(nn)
	case 0x21:
		cpu.set_hl(nn)
	case 0x31:
		cpu.sp = nn
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
	} else if 15 < res && res <= 255 {
		cpu.clearZeroFlag()
		cpu.setHalfCarryFlag()
		cpu.clearCarryFlag()
	} else if 255 < res {
		cpu.clearZeroFlag()
		cpu.clearHalfCarryFlag()
		cpu.setCarryFlag()
	}
}


