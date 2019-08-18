# わかっていないこと
* mac専用の画面を出力し、その上にピクセルを描画する
* 何から始めればよいのか

**********************************************************************************

msb: Most Significant Bit
lsb: Least Significant Bit
ex.
149 = 0b10010101
        ^      ^
       msb    lsb

msb: Most Significant Byteの場合もある

**********************************************************************************

* ROM内の256ByteのプログラムでGBは立ち上がる。
- はじめに、$104 - $133 = 48Byte、NintendoのLogo
- 音楽流れる
- 再度$104 - $133のメモリを読み、最初のbitと比較する。違っていたらhalt. 合ってたら続ける
- $134 - $14d = 26Byte

**********************************************************************************

clock数云々を気にする必要があるのは、スクリーン描画のため？

**********************************************************************************

MBC(Memory Bank Contoroller)
bank switchingにより、利用可能なアドレス空間を拡張する
MBC chipはカートリッジ内に存在し、GB本体の中にはない

以下の7種類（0x0147の場所にMBCの種類が書かれている）
- MBC なし (32 KB ROM のみ)
    メモリ0000-7FFFにROMがコピーされる
    オプションで A000-BFFF に、8KB までの容量の RAM を接続することができます。
- MBC1 (最大 2 MB ROM と、オプションで 32 KB RAM)
- MBC2 (最大 256 KB ROM と 512 x 4 ビットの RAM)
- MBC3 (最大 2 MB ROM と、オプションで 32 KB RAM と、タイマー)
- MBC5
- HuC1 (赤外線コントローラ付き MBC)
- Rumble


**********************************************************************************

    0xFFFF +--------------+ 64K
           +==============+
           |   working    | 56K
           |     RAM      | 
    0xC000 +--------------+ 48K
           | External RAM |
    0xA000 +--------------+ 40K
           |   GPU VRAM   |
    0x8000 +--------------+ 32K
           |   ROM BANK 1 |
           |              |
           |              |
    0x4000 +--------------+ 16K
           |   ROM BANK 0 |
           |              |
           +==============+
    0x0000 +--------------+ 0K

    Addresses       Name    Description        
    FFFFh           IE      Register Interrupt enable flags.
    FF80h – FFFEh   HRAM    Internal CPU RAM
    FF00h – FF7Fh   I/O     Registers I/O registers are mapped here.
    FEA0h – FEFFh   UNUSED  Description of the behaviour below.
    FE00h – FE9Fh   OAM     (Object Attribute Table) Sprite information table.
    E000h – FDFFh   ECHO    Description of the behaviour below.
    D000h – DFFFh   WRAMX   Work RAM, switchable (1-7) in GBC mode
    C000h – CFFFh   WRAM0   Work RAM.
    A000h – BFFFh   SRAM    External RAM in cartridge, often battery buffered.
    8000h – 9FFFh   VRAM    Video RAM, switchable (0-1) in GBC mode.
                            (Tile Data Tableは「8000h-8FFFh」or「8800h-97FFh」に格納される)
                            (Tile Data Tableの種類はLCDCにより選ばれる)
    4000h – 7FFFh   ROMX    Switchable ROM bank.
    0000h – 3FFFh   ROM0    Non-switchable ROM Bank.

**********************************************************************************

解読、ムスカの気持ち

**********************************************************************************
GBの種類
GB, CGB, SGB
**********************************************************************************

目標: 『Tetris』の起動
- SGB Flag: No
- Cartridge type: Rom Only
- Rom Size: 32KB (no Rom Banking)
- External RAM Size: None
- Header Checksum: 0x0B

**********************************************************************************

MBC(Memory Bank Controller)
- 利用可能なアドレスを拡張するために用いられる
- MBCはカートリッジの内部に存在し、GB本体にはない

Bank switching
- アドレス空間を1次元配列から多次元配列にする

- bank switching用のレジスタが必要
- GBのカートリッジにもMBCというチップが入っており、
- ROMバンク切り替え
- SRAMバンク切り替え
- 赤外線リンク
を提供していた

**********************************************************************************

Goライブラリ「go-sdl2」をテストする

**********************************************************************************

・どのようにCPUが動くか、opcode cycleとはなにか
・どのようにGPUはスクリーン上にpixelを描画するか、HBlankとVBlankとはなにか
・どのようにサウンドは生成されるか
・どのように割り込み処理が起こり、なぜそれが便利なのか

**********************************************************************************

Why did I spent ...の実装順
・すべてのCPUopcodeの実装
・メモリ[0x10000]byteの実装
・CPUタイミング

**********************************************************************************

var opcodes map[uint16]Opcode

type Opcode struct {
       label string

}

setA()
setB()
setSP()
getA()
getB()
getSP()の実装

**********************************************************************************

タイマーの実装

CPUからVRAMに直接アクセスできる

**********************************************************************************

MBC(Memory Bank Controller)
- 利用可能なアドレスを拡張するために用いられる
- MBCはカートリッジの内部に存在し、GB本体にはない

Bank switching
- アドレス空間を1次元配列から多次元配列にする

- bank switching用のレジスタが必要
- GBのカートリッジにもMBCというチップが入っており、
	- ROMバンク切り替え
	- SRAMバンク切り替え
	- 赤外線リンク
を提供していた

**********************************************************************************

- Clock: Oscillator clock frequency is 4194304 Hz (8388608 Hz in double speed mode).
- Cycle: CPU cycle frequency is 1048576 Hz (2097152 Hz in double speed mode).

Oscillator
    発振器
**********************************************************************************