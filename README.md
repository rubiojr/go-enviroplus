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

### PMS5003

Particle concentration sensor.

```Go
package main

import (
	"fmt"

	"github.com/rubiojr/go-enviroplus/pms5003"
)

func main() {
	dev, err := pms5003.New()
	if err != nil {
		panic(err)
	}
	dev.StartReading(func(r *pms5003.PMS5003) {
		fmt.Println("-------")
		fmt.Println("PM1.0 ug/m3 (ultrafine):                        ", r.Pm10Std)
		fmt.Println("PM2.5 ug/m3 (combustion, organic comp, metals): ", r.Pm25Std)
		fmt.Println("PM10 ug/m3 (dust, pollen, mould spores):        ", r.Pm100Std)
		fmt.Println("PM1.0 ug/m3 (atmos env):                        ", r.Pm10Env)
		fmt.Println("PM2.5 ug/m3 (atmos env):                        ", r.Pm25Env)
		fmt.Println("PM10 ug/m3 (atmos env):                         ", r.Pm100Env)
		fmt.Println("0.3um 1 0.1L air:                               ", r.Particles3um)
		fmt.Println("0.5um 1 0.1L air:                               ", r.Particles5um)
		fmt.Println("1.0um 1 0.1L air:                               ", r.Particles10um)
		fmt.Println("2.5um 1 0.1L air:                               ", r.Particles25um)
		fmt.Println("5um 1 0.1L air:                                 ", r.Particles50um)
		fmt.Println("10um 1 0.1L air:                                ", r.Particles100um)
	})
}
```
