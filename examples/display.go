package main

import (
	"image/color"
	"time"

	"github.com/rubiojr/go-enviroplus/display"
)

func main() {
	d, err := display.Init()
	if err != nil {
		panic(err)
	}

	err = d.PowerOn()
	if err != nil {
		panic(err)
	}

	err = d.FillScreen(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	err = d.PowerOff()
	if err != nil {
		panic(err)
	}
}
