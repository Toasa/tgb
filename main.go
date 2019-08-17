package main

import (
	"fmt"
	"tgb/gb"
)

func main() {
	filename := "roms/Tetris.gb"
	gb := gb.New(filename)
	err := gb.Boot()

	if err != nil {
		fmt.Println(err)
		return
	}

	gb.Update()
}
