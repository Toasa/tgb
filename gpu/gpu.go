package gpu

const w_resolution int = 160
const h_resolution int = 144

const (
	// (0xAARRGGBB: Alpha-Red-Green-Blue)
	BLACH      = 0xFF000000
	DARK_GLAY  = 0xFF555555
	LIGHT_GLAY = 0xFFAAAAAA
	WHITE      = 0xFFFFFFFF

	LCDC = 0xFF40
	STAT = 0xFF41
)

type GPU struct {
	Screen [w_resolution][h_resolution][4]int
	BackGround [256][256][4]int
}

func New() *GPU {
	return &GPU{}
}

func (gpu *GPU) LCDC() {
	// Bit 7 - LCD Control Operation
	// 	0: Stop completely (no picture on screen)
	// 	1: operation
	// Bit 6 - Window Tile Map Display Select
	// 	0: $9800-$9BFF
	// 	1: $9C00-$9FFF
	// Bit 5 - Window Display
	// 	0: off
	// 	1: on
	// Bit 4 - BG & Window Tile Data Select
	// 	0: $8800-$97FF
	// 	1: $8000-$8FFF <- Same area as OBJ
	// Bit 3 - BG Tile Map Display Select
	// 	0: $9800-$9BFF
	// 	1: $9C00-$9FFF
	// Bit 2 - OBJ (Sprite) Size
	// 	0: 8*8
	// 	1: 8*16 (width*height)
	// Bit 1 - OBJ (Sprite) Display
	// 	0: off
	// 	1: on
	// Bit 0 - BG & Window Display
	// 	0: off
	//  1: on
}

func (gpu *GPU) RenderScreen() {

}

func (gpu *GPU) UpdateGraphics(cycles int) {

}
