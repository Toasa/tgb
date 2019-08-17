package memory

import (
	"tgb/interrupt"
)

var Data [0x10000]uint8

const (
	DIV  = 0xFF04
	TIMA = 0xFF05
	TMA  = 0xFF06
	TAC  = 0xFF07
)

func Write(addr uint16, val uint8) {
	// Unused memory area in GB
	if 0xFEA0 <= addr && addr <= 0xFEFF {
		return
	}

	if addr == 0xFF0F {
		interrupt.IF = val
		return
	} else if addr == 0xFFFF {
		interrupt.IE = val
		return
	}

	Data[addr] = val
}

func Read(addr uint16) uint8 {
	// Unused memory area in GB
	if 0xFEA0 <= addr && addr <= 0xFEFF {
		return 0x00
	}

	if addr == 0xFF0F {
		return interrupt.IF
	} else if addr == 0xFFFF {
		return interrupt.IE
	}

	return Data[addr]
}
