package main

import (
	"fmt"

	"github.com/rubiojr/go-enviroplus/bme280"
	"periph.io/x/conn/v3/physic"
)

func main() {
	dev, _ := bme280.New()
	r, _ := dev.Read()

	fmt.Println("Temperature: ", r.Temperature)
	fmt.Printf("Humidity:     %.0f\n", float64(r.Humidity)/float64(physic.PercentRH))
	fmt.Println("Pressure:    ", float64(r.Pressure)/float64(physic.Pascal))
}
