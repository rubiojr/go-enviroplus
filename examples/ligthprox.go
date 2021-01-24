// Read proximity and light from the LTR559 sensor
package main

import (
	"fmt"
	"time"

	"github.com/rubiojr/go-enviroplus/ltr559"
)

func main() {
	d, _ := ltr559.New()

	mid, _ := d.ManufacturerID()
	pid, _ := d.PartID()

	fmt.Printf("Manufacturer ID:  0x%x\n", mid)
	fmt.Printf("Part ID:          0x%x\n", pid)

	for {
		prox, _ := d.Proximity()
		lux, _ := d.Lux()
		fmt.Println("proximity: ", prox)
		fmt.Println("      lux: ", lux)
		time.Sleep(1 * time.Second)
	}
}
