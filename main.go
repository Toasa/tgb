package main

import (
	"fmt"
	"tgb/gb"
)

func main() {
	filename := "roms/Tetris.gb"
	gb := gb.New(filename)
	gb.Update()
	
}
