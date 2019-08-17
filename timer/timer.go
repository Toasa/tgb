package timer

import (
	"tgb/memory"
)

type Timer struct {
	InternalTimer int

	DIV  int
	TIMA int
	TMA  int
}

const (
	CLOCK_CYCLE = 4194304
	// MACHINE_CYCLE = 1048576
	FRAMES_SECOND = 60

	// CYCLES_FRAME is cpu clock num for 1frame.
	CYCLES_FRAME = CLOCK_CYCLE / FRAMES_SECOND

	// Divider Register (R/W)
	// This register is incremented at rate of 16384Hz (~16779Hz on SGB).
	//     CPUclock speedは4MiHzなので、
	// Writing any value to this register resets it to 00h.
	// Note: The divider is affected by CGB double speed mode,
	// and will increment at 32768Hz in double speed.
	DIV = 0xFF04

	// Timer counter (R/W)
	// This timer is incremented by a clock frequency specified by the TAC register ($FF07).
	// When the value overflows (gets bigger than FFh) then it will be reset to the value specified in TMA (FF06),
	// and an interrupt will be requested (0x0050).
	//     割り込みがrequestされるのは、IEフラグが立っており、IMEとなっているとき
	//     その割り込みのアドレスは0x0050
	TIMA = 0xFF05

	// Timer Modulo (R/W)
	// When the TIMA overflows, this data will be loaded.
	TMA = 0xFF06

	// Timer Control (R/W)
	// Bit  2   - Timer Enable / Disable
	//         0: Stop Timer
	//         1: Start Timer
	// Bits 1-0 - Input Clock Select
	// 		   00: CPU Clock / 1024 (DMG, CGB:   4096 Hz, SGB:   ~4194 Hz)
	// 		   01: CPU Clock / 16   (DMG, CGB: 262144 Hz, SGB: ~268400 Hz)
	// 		   10: CPU Clock / 64   (DMG, CGB:  65536 Hz, SGB:  ~67110 Hz)
	// 		   11: CPU Clock / 256  (DMG, CGB:  16384 Hz, SGB:  ~16780 Hz)
	TAC = 0xFF07
)

func newTimer() *Timer {
	return &Timer{
		InternalTimer: 0,
		DIV:           0,
		TIMA:          0,
		TMA:           0,
	}
}

func (t *Timer) UpdateTimers(cycles int) {
	// t.dividerRegister(cycles)

	// t.timer = gb.getInputClock()

	// for cycles > 0 {
	// 	//gb.timer += 4

	// 	// TIMA の更新
	// 	t.write(TIMA, t.read(TIMA)+1)

	// 	if t.isTIMAOverflowed() {
	// 		t.write(TIMA, t.read(TMA))

	// 		t.timerInterruptRequest()
	// 	}
	// }
}

func (t *Timer) read(addr uint16) uint8 {
	return memory.Read(addr)
}

func (t *Timer) write(addr uint16, val uint8) {
	memory.Write(addr, val)
}

func DividerRegister(cycles int) {

}

func IsTIMAOverflowed() bool {
	return false
}

func TimerInterruptRequest() {

}
