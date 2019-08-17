package gb

import (
	"fmt"
	"time"
	"errors"
	"io/ioutil"
	"tgb/cpu"
	"tgb/gpu"
	"tgb/interrupt"
	"tgb/memory"
	"tgb/timer"
)

type GB struct {
	CPU     *cpu.CPU
	GPU     *gpu.GPU
	ROM     []byte
	RomInfo Rom_info
	Memory *[0x10000]uint8

	// 4.194304MHz / 256 = 16.384KHz
	timer *timer.Timer

	current_cycle int
}

type Rom_info struct {
	entryPoint           []byte
	nintendoLogo         []byte
	Title                string
	CGBFlag              bool
	newLicenseeCode      []byte
	SGBFlag              bool
	cartridgeType        string
	romSize              int
	ramSize              int
	destinationCode      byte
	oldLicenseeCode      byte
	maskROMVersionNumber byte
}

var cartridgeTypeMap map[byte]string = map[byte]string{
	0x00: "ROM ONLY",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BATTERY",
	0x05: "MBC2",
	0x06: "MBC2+BATTERY",
	0x08: "ROM+RAM",
	0x09: "ROM+RAM+BATTERY",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BATTERY",
	0x0F: "MBC3+TIMER+BATTERY",
	0x10: "MBC3+TIMER+RAM+BATTERY",
	0x11: "MBC3",
	0x12: "MBC3+RAM",
	0x13: "MBC3+RAM+BATTERY",
	0x15: "MBC4",
	0x16: "MBC4+RAM",
	0x17: "MBC4+RAM+BATTERY",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BATTERY",
	0x1C: "MBC5+RUMBLE",
	0x1D: "MBC5+RUMBLE+RAM",
	0x1E: "MBC5+RUMBLE+RAM+BATTERY",
	0xFC: "POCKET CAMERA",
	0xFD: "BANDAI TAMA5",
	0xFE: "HuC3",
	0xFF: "HuC1+RAM+BATTERY",
}

var romSizeMap map[byte]int = map[byte]int{
	0x00: 32768,
	0x01: 65536,
	0x02: 131072,
	0x03: 262144,
	0x04: 524288,
	0x05: 1048576,
	0x06: 2097152,
	0x07: 4194304,
	0x52: 1153434,
	0x53: 1258292,
	0x54: 1572864,
}

var ramSizeMap map[byte]int = map[byte]int{
	0x00: 0,
	0x01: 2048,
	0x02: 8192,
	0x03: 32768,
}

func New(filename string) *GB {
	var gb *GB

	rom, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	ri := getRomInfo(rom)

	gb = &GB{
		CPU:     cpu.NewCPUinBoot(),
		ROM:     rom,
		RomInfo: ri,
		Memory: &memory.Data,
	}

	switch ri.cartridgeType {
	case "ROM ONLY":
		gb.loadROMtoRAM(rom, 0x8000)
	default:
	}
	return gb
}

func (gb *GB) Boot() error {
	// 0x0104 - 0x0133
	gb.displayLogo()
	result := gb.checkLogoArea()
	if !result {
		return errors.New("Compared the ROM and internal memory in 0x0104-0x0133, but didn't get exactly coincide.")
	}

	// Bootときに設定した各レジスタを初期値に上書きする
	gb.CPU = cpu.NewCPU()
	gb.setMemoryValueInBoot()
	return nil
}

func (gb *GB) write(addr uint16, val uint8) {
	memory.Write(addr, val)
}

func (gb *GB) read(addr uint16) uint8 {
	return memory.Read(addr)
}

func (gb *GB) loadROMtoRAM(rom []byte, size int) {
	// rom sizeが0x8000B限定の場合
	for i := 0; i < size; i++ {
		gb.write(uint16(i), rom[i])
	}
}

func getRomInfo(rom_data []byte) Rom_info {
	return Rom_info{
		entryPoint:   rom_data[0x0100:0x103],
		nintendoLogo: rom_data[0x0104:0x0133],
		Title:        fetchTitle(rom_data[0x0134:0x0143]),
		//CGBFlag: ,
		newLicenseeCode:      rom_data[0x0144:0x145],
		SGBFlag:              checkSGBFlag(rom_data[0x146]),
		cartridgeType:        cartridgeTypeMap[rom_data[0x147]],
		romSize:              romSizeMap[rom_data[0x148]],
		ramSize:              ramSizeMap[rom_data[0x149]],
		destinationCode:      rom_data[0x014A],
		oldLicenseeCode:      rom_data[0x014B],
		maskROMVersionNumber: rom_data[0x014C],
	}
}

func fetchTitle(bytes []byte) string {
	var i int = len(bytes) - 1
	for ; 0 < i; i-- {
		if bytes[i] != 0x00 {
			break
		}
	}
	title := string(bytes[:i+1])
	return title
}

func checkSGBFlag(b byte) bool {
	if b == 0x03 {
		return true
	}
	return false
}

// Should be called 60 times/second
func (gb *GB) Update() {
	for {
		if interrupt.CheckInterrupts() {
			// the CPU will push the current PC into the stack, will jump
			// to the corresponding interrupt vector and set IME to '0'.
			// If IME is '0', this won't happen.
			gb.CPU.PushCurrentPC()
			interrupt.CheckInterruptVector()
		}
		cycles := gb.CPU.Step()
		fmt.Println(cycles)
		time.Sleep(time.Millisecond * 100)

		gb.timer.UpdateTimers(cycles)
		gb.GPU.UpdateGraphics(cycles)
		interrupt.DoInterrupts()

		gb.current_cycle += cycles
		if gb.current_cycle >= timer.CYCLES_FRAME {
			gb.current_cycle -= timer.CYCLES_FRAME
			gb.GPU.RenderScreen()
		}
	}
}

func (gb *GB) getInputClock() int {
	t := gb.read(timer.TAC) & 0x11
	switch t {
	case 0x00:
		return 1024
	case 0x01:
		return 16
	case 0x10:
		return 64
	case 0x11:
		return 256
	}

	return -1
}

func (gb *GB) checkLogoArea() bool {
	// FIX IT!
	return true
}

func (gb *GB) displayLogo() {
	
}

func (gb *GB) setMemoryValueInBoot() {
	gb.write(0xFF05, 0x00) // TIMA
	gb.write(0xFF06, 0x00) // TMA
	gb.write(0xFF07, 0x00) // TAC
	gb.write(0xFF10, 0x80) // NR10
	gb.write(0xFF11, 0xBF) // NR11
	gb.write(0xFF12, 0xF3) // NR12
	gb.write(0xFF14, 0xBF) // NR14
	gb.write(0xFF16, 0x3F) // NR21
	gb.write(0xFF17, 0x00) // NR22
	gb.write(0xFF19, 0xBF) // NR24
	gb.write(0xFF1A, 0x7F) // NR30
	gb.write(0xFF1B, 0xFF) // NR31
	gb.write(0xFF1C, 0x9F) // NR32
	gb.write(0xFF1E, 0xBF) // NR33
	gb.write(0xFF20, 0xFF) // NR41
	gb.write(0xFF21, 0x00) // NR42
	gb.write(0xFF22, 0x00) // NR43
	gb.write(0xFF23, 0xBF) // NR30
	gb.write(0xFF24, 0x77) // NR50
	gb.write(0xFF25, 0xF3) // NR51
	gb.write(0xFF26, 0xF1) // NR52
	gb.write(0xFF40, 0x91) // LCDC
	gb.write(0xFF42, 0x00) // SCY
	gb.write(0xFF43, 0x00) // SCX
	gb.write(0xFF45, 0x00) // LYC
	gb.write(0xFF47, 0xFC) // BGP
	gb.write(0xFF48, 0xFF) // OBP0
	gb.write(0xFF49, 0xFF) // OBP1
	gb.write(0xFF4A, 0x00) // WY
	gb.write(0xFF4B, 0x00) // WX
	gb.write(0xFFFF, 0x00) // IE
}