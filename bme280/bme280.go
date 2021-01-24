// Thin wrapper around perip.io bme280
package bme280

import (
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

type BME280 struct {
	bus    i2c.BusCloser
	device *bmxx80.Dev
}

func New() (*BME280, error) {
	var err error
	dev := &BME280{}
	if _, err = host.Init(); err != nil {
		return dev, err
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	dev.bus, err = i2creg.Open("")
	if err != nil {
		return dev, err
	}

	dev.device, err = bmxx80.NewI2C(dev.bus, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		return dev, err
	}

	return dev, nil
}

type Readings struct {
	Temperature physic.Temperature
	Pressure    physic.Pressure
	Humidity    physic.RelativeHumidity
}

func (dev *BME280) Read() (*Readings, error) {
	e := physic.Env{}
	if err := dev.device.Sense(&e); err != nil {
		return nil, err
	}

	return &Readings{
		Temperature: e.Temperature,
		Pressure:    e.Pressure,
		Humidity:    e.Humidity,
	}, nil
}
