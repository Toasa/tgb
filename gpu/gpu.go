package gpu

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth = 160
	winHeight = 144
	 
	PIXEL_SIZE = 5

	// (0xAARRGGBB: Alpha-Red-Green-Blue)
	BLACH      = 0xFF000000
	DARK_GLAY  = 0xFF555555
	LIGHT_GLAY = 0xFFAAAAAA
	WHITE      = 0xFFFFFFFF

	LCDC = 0xFF40
	STAT = 0xFF41
)

type GPU struct {
	Screen [winWidth][winHeight][4]int
	BackGround [256][256][4]int
	Title string

	Window *sdl.Window
	Surface *sdl.Surface
}

func New() *GPU {
	return &GPU{
		Title: "test",
	}
}

func (gpu *GPU) Init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		gpu.Title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth * PIXEL_SIZE,
		winHeight * PIXEL_SIZE,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}

	gpu.Window = window

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	gpu.Surface = surface
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
	for x := 0; x < winWidth; x++ {
		for y := 0; y < winHeight; y++ {
			rect := sdl.Rect{
				X: int32(x) * PIXEL_SIZE,
				Y: int32(y) * PIXEL_SIZE,
				W: PIXEL_SIZE,
				H: PIXEL_SIZE,
			}
			gpu.Surface.FillRect(&rect, toColor(gpu.Screen[x][y]))
		}
	}
	gpu.Window.UpdateSurface()
}

func (gpu *GPU) UpdateGraphics(cycles int) {
	gpu.RenderScreen()
}

func toColor(scrn [4]int) uint32 {
	return uint32(scrn[3]) << 24 | uint32(scrn[2]) << 16 | uint32(scrn[1]) << 8 | uint32(scrn[0])
}