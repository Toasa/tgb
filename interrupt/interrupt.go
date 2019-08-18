package interrupt

import (
	"tgb/memory"
)

// Interrupt Master Enable Flag
// false - Disable all Interrupts
// true - Enable all Interrupts that are enabled in IE Register (FFFF)
var IME bool

const (
	// FFFF - IE - Interrupt Enable (R/W)
	//  Bit 0: V-Blank  Interrupt Enable  (INT 40h)  (1=Enable)
	//  Bit 1: LCD STAT Interrupt Enable  (INT 48h)  (1=Enable)
	//  Bit 2: Timer    Interrupt Enable  (INT 50h)  (1=Enable)
	//  Bit 3: Serial   Interrupt Enable  (INT 58h)  (1=Enable)
	//  Bit 4: Joypad   Interrupt Enable  (INT 60h)  (1=Enable)
	IE = 0xFFFF

	// FF0F - IF - Interrupt Flag (R/W)
	//  Bit 0: V-Blank  Interrupt Request (INT 40h)  (1=Request)
	//  Bit 1: LCD STAT Interrupt Request (INT 48h)  (1=Request)
	//  Bit 2: Timer    Interrupt Request (INT 50h)  (1=Request)
	//  Bit 3: Serial   Interrupt Request (INT 58h)  (1=Request)
	//  Bit 4: Joypad   Interrupt Request (INT 60h)  (1=Request)
	IF = 0xFF0F
)

func DoInterrupt() {

}

func SetIMEFlag() {
	IME = true
}

func ClearIMEFlag() {
	IME = false
}

func CheckInterruptVector() {

}

func CheckInterrupts() bool {
	if isAllSet(read(IE)) && isAllSet(read(IF)) && IME == true {
		return true
	}
	return false
}

func isAllSet(flag uint8) bool {
	return 0xFF == flag
}

func DoInterrupts() {

}

func write(addr uint16, val uint8) {
	memory.Write(addr, val)
}

func read(addr uint16) uint8 {
	return memory.Read(addr)
}

func SetIF_VBlankFlag() {
	write(IF, read(IF) | 0x01)
}

func SetIF_LCDFlag() {
	write(IF, read(IF) | 0x02)
}

func SetIF_TimerFlag() {
	write(IF, read(IF) | 0x04)
}

func SetIF_SerialFlag() {
	write(IF, read(IF) | 0x08)
}

func SetIF_JoypadFlag() {
	write(IF, read(IF) | 0x10)
}

func ClearIF_VBlankFlag() {
	write(IF, read(IF) & 0xFE)
}

func ClearIF_LCDFlag() {
	write(IF, read(IF) & 0xFD)
}

func ClearIF_TimerFlag() {
	write(IF, read(IF) & 0xFB)
}

func ClearIF_SerialFlag() {
	write(IF, read(IF) & 0xF7)
}

func ClearIF_JoypadFlag() {
	write(IF, read(IF) & 0xEF)
}