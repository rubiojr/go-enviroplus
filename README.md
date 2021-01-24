# Pimoroni Enviro+ Drivers

Go library to read data from [Pimoroni's Enviro+](https://learn.pimoroni.com/tutorial/sandyj/getting-started-with-enviro-plus) sensors.

## LTR559

Light/Proximity Sensor.

The driver is a port of [the Python driver](https://github.com/pimoroni/ltr559-python) from Pimoroni.

### Reading data from the sensor

```Go
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
```
