# Pimoroni Enviro+ Drivers

Go library to read data from [Pimoroni's Enviro+](https://learn.pimoroni.com/tutorial/sandyj/getting-started-with-enviro-plus) sensors.

⚠️ Experimental, API subject to change ⚠

## BME250

Package to read pressure, relative humidity and temperature sensors.

```Go
package main

import (
	"fmt"
	"log"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
	"periph.io/x/periph/host"
)

func main() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer b.Close()

	d, err := bmxx80.NewI2C(b, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}
	e := physic.Env{}
	if err := d.Sense(&e); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%8s %10s %9s\n", e.Temperature, e.Pressure, e.Humidity)
}
```

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
