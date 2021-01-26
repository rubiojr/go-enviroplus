package main

import (
	"fmt"
	"time"

	"github.com/rubiojr/go-enviroplus/mics6814"
)

func main() {
	dev, err := mics6814.New()
	if err != nil {
		panic(err)
	}
	defer dev.Halt()

	go func() {
		dev.StartReading()
	}()

	for {
		fmt.Printf("Oxidising: %.2f\n", dev.LastValue().Oxidising)
		fmt.Printf("Reducing:  %.2f\n", dev.LastValue().Reducing)
		fmt.Printf("NH3:       %.2f\n", dev.LastValue().NH3)
		time.Sleep(1 * time.Second)
	}
}
