// Driver for the LTR559 light and proximity sensor
package ltr559

import (
	"encoding/binary"
	"errors"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

var ch0_c = []int{17743, 42785, 5926, 0}
var ch1_c = []int{-11059, 19548, -1185, 0}

const integration_time = 50.0

const (
	ALS_CONTROL   = 0x80
	ALS_PS_STATUS = 0x8c
	ALS_DATA      = 0x88
	ALS_MEAS_RATE = 0x85
	ALS_THRESHOLD = 0x97

	MANUFACTURER_ID = 0x87
	PART_ID         = 0x86

	PS_OFFSET    = 0x94
	PS_DATA      = 0x8D
	PS_LED       = 0x82
	PS_N_PULSES  = 0x83
	PS_CONTROL   = 0x81
	PS_MEAS_RATE = 0x84
	PS_THRESHOLD = 0x90
)

// LTR559 is the sensor struct
type LTR559 struct {
	device     *i2c.Dev
	ps0        uint16
	als        []byte
	als0, als1 uint16
	lux        float64
	bus        i2c.BusCloser
}

// New create an LTR559 sensor ready for reading light/proximity values.
func New() (*LTR559, error) {
	d := &LTR559{}

	return d, d.init()
}

// Close closes the the i2c bus
func (s *LTR559) Close() {
	s.bus.Close()
}

func (s *LTR559) init() error {
	var err error
	if _, err = host.Init(); err != nil {
		return err
	}

	s.bus, err = i2creg.Open("")
	if err != nil {
		return err
	}

	// Dev is a valid conn.Conn.
	s.device = &i2c.Dev{Addr: 0x23, Bus: s.bus}

	// als_control sw_reset=1
	s.setRegister(ALS_CONTROL, bitsToBytes("00000010"))

	// ps_led current_ma: 50, duty_cycle: 1.0, pulse: 30khz
	s.setRegister(PS_LED, bitsToBytes("00011011"))

	// ps_n_pulses count 1
	s.setRegister(PS_N_PULSES, bitsToBytes("00001111"))

	// als_control mode=1 gain=4
	s.setRegister(ALS_CONTROL, bitsToBytes("00001001"))

	// ps_control active=1 saturation_indicator_enable=1
	s.setRegister(PS_CONTROL, bitsToBytes("00100011"))

	// ps_meas_rate rate_ms=100
	s.setRegister(PS_MEAS_RATE, bitsToBytes("00000010"))

	// als_meas_rate repeat_rate=50 integration_time_ms=50
	s.setRegister(ALS_MEAS_RATE, bitsToBytes("00001000"))

	// als_threshold lower=0x0000 upper=0xffff
	s.setRegister(ALS_THRESHOLD, []byte{0xFF, 0xFF, 0x00, 0x00})

	// ps_threshold lower=0x0000 upper=0xffff
	s.setRegister(PS_THRESHOLD, []byte{0xFF, 0xFF, 0x00, 0x00})

	// ps_offset offset=0
	s.setRegister(PS_OFFSET, []byte{0x00, 0x00})

	return nil
}

// Lux returns the ambient light value in lux.
func (s *LTR559) Lux() (float64, error) {
	err := s.updateSensor()
	if err != nil {
		return 0, err
	}

	return s.lux, nil
}

// Proximity returns the RAW proximity reading from the sensor.
func (s *LTR559) Proximity() (float64, error) {
	err := s.updateSensor()
	if err != nil {
		return 0, err
	}

	return float64(s.ps0), nil
}

func (s *LTR559) updateSensor() error {
	var err error
	b, err := s.getRegister(ALS_PS_STATUS, 1)
	if err != nil {
		return err
	}
	status := b[0]

	ps_int := (status&0x02 != 0x0) || (status&0x04 != 0)
	als_int := (status&0x08 != 0x0) || (status&0x04 != 0x0)

	if ps_int {
		res, err := s.getRegister(PS_DATA, 2)
		if err != nil {
			return err
		}
		s.ps0 = binary.LittleEndian.Uint16(res)
	}

	if als_int {
		s.als, err = s.getRegister(ALS_DATA, 4)
		if err != nil {
			return err
		}
		s.als0 = binary.LittleEndian.Uint16(s.als[0:2])
		s.als1 = binary.LittleEndian.Uint16(s.als[2:])

		var ratio uint16
		if s.als0+s.als1 > 0 {
			ratio = s.als1 * 100 / (s.als1 + s.als0)
		} else {
			ratio = 101
		}

		var ch_idx int
		if ratio < 45 {
			ch_idx = 0
		} else if ratio < 64 {
			ch_idx = 1
		} else if ratio < 85 {
			ch_idx = 2
		} else {
			ch_idx = 3
		}

		lux := float64((int(s.als0) * ch0_c[ch_idx]) - (int(s.als1) * ch1_c[ch_idx]))
		lux /= integration_time / 100.0
		lux /= 4
		lux /= 10000.0
		s.lux = lux
	}

	return nil
}

func (s *LTR559) ManufacturerID() ([]byte, error) {
	return s.getRegister(MANUFACTURER_ID, 1)
}

func (s *LTR559) PartID() (byte, error) {
	b, err := s.getRegister(PART_ID, 1)
	if err != nil {
		return 0, err
	}

	return b[0], nil
}

func (s *LTR559) sendData(d []byte) error {
	return s.device.Tx(d, nil)
}

func (s *LTR559) setRegister(addr byte, data []byte) {
	//fmt.Printf("0x%x: %b [0x%x]\n", addr, data, data)
	l := []byte{}
	l = append(l, addr)
	l = append(l, data...)
	s.sendData(l)
}

func (s *LTR559) getRegister(addr byte, count int) ([]byte, error) {
	read := make([]byte, count)
	err := s.device.Tx([]byte{addr}, read)

	return read, err
}

func bitsToBytes(s string) []byte {
	b := make([]byte, (len(s)+(8-1))/8)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '1' {
			panic(errors.New("value out of range"))
		}
		b[i>>3] |= (c - '0') << uint(7-i&7)
	}
	return b
}
