// Driver to read read the MICS6814 via an ads1015 ADC
//
// Reference driver: https://github.com/pimoroni/enviroplus-python/blob/master/library/enviroplus/gas.py
package mics6814

import (
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ads1x15"
	"periph.io/x/host/v3"
)

type Device struct {
	oxPin     ads1x15.PinADC
	redPin    ads1x15.PinADC
	nh3Pin    ads1x15.PinADC
	lastRead  Readings
	waitTime  time.Duration
	heaterPin gpio.PinIO
	log       zerolog.Logger
}

type Readings struct {
	Oxidising float64
	Reducing  float64
	NH3       float64
}

// Return the last value read from the sensor
//
// Call StartReading() first to start reading values from the sensor.
func (dev *Device) LastValue() Readings {
	return dev.lastRead
}

func (dev *Device) StartReading() {
	defer dev.Halt()

	for {
		ox, err := dev.readPin(dev.oxPin)
		if err != nil {
			dev.log.Error().Err(err)
		}

		red, err := dev.readPin(dev.redPin)
		if err != nil {
			dev.log.Error().Err(err)
		}

		nh3, _ := dev.readPin(dev.nh3Pin)
		if err != nil {
			dev.log.Error().Err(err)
		}

		dev.lastRead = Readings{
			Oxidising: ox,
			Reducing:  red,
			NH3:       nh3,
		}
		time.Sleep(dev.waitTime)
	}
}

type Opts struct {
	Wait       time.Duration // wait time between readings
	I2cAddress byte
}

var DefaultOpts = &Opts{
	Wait:       1 * time.Second,
	I2cAddress: 0x49,
}

func New() (*Device, error) {
	return NewWithOpts(*DefaultOpts)
}

func NewWithOpts(opts Opts) (*Device, error) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}

	dev := &Device{}
	dev.waitTime = opts.Wait
	dev.log = zerolog.New(os.Stderr).With().Timestamp().Logger()
	dev.log = dev.log.Level(zerolog.InfoLevel)

	// https://pinout.xyz/pinout/enviro_plus
	// Gas Heater En
	dev.heaterPin = gpioreg.ByName("GPIO24")
	if err := dev.heaterPin.Out(gpio.High); err != nil {
		return nil, err
	}

	// Create a new ADS1015 ADC.
	dopts := &ads1x15.DefaultOpts
	dopts.I2cAddress = 0x49
	adc, err := ads1x15.NewADS1015(bus, dopts)
	if err != nil {
		return nil, err
	}

	dev.oxPin, err = adc.PinForChannel(ads1x15.Channel0, 1*physic.Volt, 1*physic.Hertz, ads1x15.BestQuality)
	if err != nil {
		return dev, err
	}

	dev.redPin, err = adc.PinForChannel(ads1x15.Channel1, 1*physic.Volt, 1*physic.Hertz, ads1x15.BestQuality)
	if err != nil {
		return dev, err
	}

	dev.nh3Pin, err = adc.PinForChannel(ads1x15.Channel2, 1*physic.Volt, 1*physic.Hertz, ads1x15.BestQuality)
	if err != nil {
		return dev, err
	}

	return dev, nil
}

func (dev *Device) Halt() {
	dev.oxPin.Halt()
	dev.redPin.Halt()
	dev.nh3Pin.Halt()
	dev.heaterPin.Out(gpio.Low)
}

func (dev *Device) readPin(pin ads1x15.PinADC) (float64, error) {
	reading, err := pin.Read()
	if err != nil {
		return 0, err
	}
	v := float64(reading.V) / float64(physic.Volt)
	return float64(v*56000) / (3.3 - v), nil
}
