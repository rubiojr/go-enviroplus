// Read proximity and light from the LTR559 sensor
package main

import (
	"fmt"
	"time"

	"github.com/rubiojr/go-enviroplus/ltr559"
)

func main() {
	d := ltr559.New()

	fmt.Printf("Manufacturer ID:  0x%x\n", d.ManufacturerID())
	fmt.Printf("Part ID:          0x%x\n", d.PartID())

	for {
		fmt.Println("proximity: ", d.Proximity())
		fmt.Println("      lux: ", d.Lux())
		time.Sleep(1 * time.Second)
	}
}
