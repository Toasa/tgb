package timer

import (
	"tgb/memory"
	"tgb/interrupt"
)

type Timer struct {
	Cycle int

	// 16bit counter
	// Upper 8bit of this counter is exactly DIV timer.
	InternalCounter int
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

// TAC = 00のとき、タイマー割り込みは1秒間に4096回起こる。
// つまりTIMAは一秒間に4096 * 256 = 1048576回インクリメントが起こった

func New() *Timer {
	return &Timer{
		Cycle: 0,
		InternalCounter: 0,
	}
}

func (t *Timer) UpdateTimers(cycles int) {

	t.Cycle += cycles
	if t.Cycle >= 4 {
		t.Cycle -= 4

		t.InternalCounter++
	}

	if t.isInternalCounterOverflow() {
		t.InternalCounter -= 65536
	}

	// internalCounterを更新したときに、DIVも更新
	t.reloadDIV()

	// TIMA の更新
	write(TIMA, read(TIMA)+1)

	if t.IsTIMAOverflowed() {

		t.loadTMA()
		t.timerInterruptRequest()
	}
	
}

func read(addr uint16) uint8 {
	return memory.Read(addr)
}

func write(addr uint16, val uint8) {
	memory.Write(addr, val)
}

func DividerRegister(cycles int) {

}

func (t *Timer) IsTIMAOverflowed() bool {
	return false
}

func (t *Timer) timerInterruptRequest() {
	interrupt.SetIF_TimerFlag()
}

func (t *Timer) isInternalCounterOverflow() bool {
	return t.InternalCounter >= 65536
}

func (t *Timer) reloadDIV() {
	currentDIVVal := uint8(t.InternalCounter >> 8)
	write(DIV, currentDIVVal)
}

func (t *Timer) loadTMA() {
	write(TIMA, read(TMA))
}

