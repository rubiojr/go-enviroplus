// Based on code from Mark Hansen:
//   https://github.com/mhansen/breathe/blob/master/breathe.go
//
// Pimoroni's driver used as a reference also:
//  https://github.com/pimoroni/pms5003-python
//
// Binary breathe reads air quality data from a PMS5003 chip, exporting the data over prometheus HTTP.
//
// PMS5003 datasheet: http://www.aqmd.gov/docs/default-source/aq-spec/resources-page/plantower-pms5003-manual_v2-3.pdf
package pms5003

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/rs/zerolog"
	"golang.org/x/sys/unix"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

const (
	magic1 = 0x42 // :)
	magic2 = 0x4d
)

// PMS5003 wraps an air quality packet, as documented in https://cdn-shop.adafruit.com/product-files/3686/plantower-pms5003-manual_v2-3.pdf
type PMS5003 struct {
	Length         uint16
	Pm10Std        uint16
	Pm25Std        uint16
	Pm100Std       uint16
	Pm10Env        uint16
	Pm25Env        uint16
	Pm100Env       uint16
	Particles3um   uint16
	Particles5um   uint16
	Particles10um  uint16
	Particles25um  uint16
	Particles50um  uint16
	Particles100um uint16
	Unused         uint16
	Checksum       uint16
}

type Device struct {
	pinEnable, pinReset gpio.PinIO
	rw                  io.ReadWriteCloser
	serialPort          string
	log                 zerolog.Logger
	lastRead            *PMS5003
}

// New device with custom options.
//
// https://pinout.xyz/pinout/enviro_plus#
//
// resetPin: module reset signal pin
// enablePin: enable/disable the module
// serialPort: usually /dev/ttyAMA0 on a Raspberry PI
func NewWithOpts(resetPin, enablePin, serialPort string) (*Device, error) {
	dev := &Device{}

	if _, err := host.Init(); err != nil {
		return dev, err
	}

	dev.pinEnable = gpioreg.ByName(enablePin)
	if err := dev.pinEnable.Out(gpio.High); err != nil {
		return nil, err
	}

	dev.pinReset = gpioreg.ByName(resetPin)
	if err := dev.pinReset.Out(gpio.High); err != nil {
		return nil, err
	}

	dev.serialPort = serialPort

	dev.log = zerolog.New(os.Stderr).With().Timestamp().Logger()
	dev.log = dev.log.Level(zerolog.InfoLevel)

	dev.lastRead = &PMS5003{}
	return dev, nil
}

// Return the last value read from the sensor
//
// Call StartReading() first to start reading values from the sensor.
func (dev *Device) LastValue() *PMS5003 {
	return dev.lastRead
}

// New device with sane default values for Enviro+ with PMS5003
// from Plantower.
func New() (*Device, error) {
	return NewWithOpts("GPIO27", "GPIO22", "/dev/ttyAMA0")
}

func (dev *Device) EnableDebugging() {
	dev.log = dev.log.Level(zerolog.DebugLevel)
}

// Start reading values from the serial port.
//
// The call will loop indefinitely and store the values read.
// Use LastValue() to retrieve the last value read.
func (dev *Device) StartReading() error {
	options := serial.OpenOptions{
		PortName:              dev.serialPort,
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       1,
		InterCharacterTimeout: 4,
	}

	var err error
	dev.rw, err = serial.Open(options)
	if err != nil {
		return err
	}
	defer dev.rw.Close()

	dev.reset()

	for {
		dev.log.Print("Attempting to read.")
		err := dev.readPMS()
		if err != nil {
			dev.log.Printf("readPMS: %v\n", err)
			dev.reset()
			continue
		}
		dev.log.Printf("pms = %+v\n", dev.lastRead)
		if !dev.lastRead.valid() {
			dev.log.Print("pms is not valid. Ignoring...")
			continue
		}
	}
}

func (p *PMS5003) valid() bool {
	if p.Length != 28 {
		return false
	}
	return true
}

func (dev *Device) readPMS() error {
	if err := dev.awaitMagic(); err != nil {
		// Read errors are likely unrecoverable - just quit and restart.
		dev.log.Error().Err(err).Msgf("awaitMagic: %v", err)
		return err
	}
	buf := make([]byte, 30)
	n, err := io.ReadFull(dev.rw, buf)
	if err != nil {
		// Read errors are likely unrecoverable - just quit and restart.
		dev.log.Error().Err(err).Msgf("readfull: %v", err)
		return err
	}
	if n != 30 {
		return fmt.Errorf("too few bytes read: want %d got %d", 30, n)
	}

	var sum uint16 = uint16(magic1) + uint16(magic2)
	for i := 0; i < 28; i++ {
		sum += uint16(buf[i])
	}

	bufR := bytes.NewReader(buf)
	binary.Read(bufR, binary.BigEndian, dev.lastRead)

	if sum != dev.lastRead.Checksum {
		// This error is recoverable
		return fmt.Errorf("checksum: got %v want %v", sum, dev.lastRead)
	}

	return nil
}

func (dev *Device) awaitMagic() error {
	dev.log.Print("Awaiting magic... ")
	var b1 byte
	b2, err := dev.pop()
	if err != nil {
		return err
	}
	for {
		b1 = b2
		b2, err = dev.pop()
		if err != nil {
			return err
		}
		if b1 == magic1 && b2 == magic2 {
			// found magic
			return nil
		}
	}
}

func (dev *Device) pop() (byte, error) {
	b := make([]byte, 1)
	_, err := dev.rw.Read(b)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// Discards data written to the port but not transmitted,
// or data received but not read
// https://github.com/tarm/serial/blob/master/serial_linux.go
func (dev *Device) flushSerial() error {
	const TCFLSH = 0x540B
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr((dev.rw).(*os.File).Fd()),
		uintptr(TCFLSH),
		uintptr(unix.TCIOFLUSH),
	)

	if errno == 0 {
		return nil
	}
	return errno
}

// reset the PMS5003 module
func (dev *Device) reset() error {
	if err := dev.pinReset.Out(gpio.Low); err != nil {
		return err
	}

	dev.flushSerial()
	time.Sleep(100 * time.Millisecond)

	if err := dev.pinReset.Out(gpio.High); err != nil {
		return err
	}

	return nil
}
