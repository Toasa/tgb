package main

import (
	"log"
	"tgb/gb"
)

func main() {
	filename := "roms/Tetris.gb"
	gb := gb.New(filename)
	err := gb.Boot()

	if err != nil {
		log.Println(err)
		return
	}

	gb.Run()
}
