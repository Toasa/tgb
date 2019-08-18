package memory

var Data [0x10000]uint8

func Write(addr uint16, val uint8) {
	// Unused memory area in GB
	if 0xFEA0 <= addr && addr <= 0xFEFF {
		return
	}

	Data[addr] = val
}

func Read(addr uint16) uint8 {
	// Unused memory area in GB
	if 0xFEA0 <= addr && addr <= 0xFEFF {
		return 0x00
	}

	return Data[addr]
}
