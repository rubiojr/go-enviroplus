package display

import (
	"image"
	"image/color"
	"io"
	"log"
	"sync"

	"github.com/asssaf/st7735-go/st7735"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

type Rotation uint8

const (
	NO_ROTATION  Rotation = 0
	ROTATION_90  Rotation = 1 // 90 degrees clock-wise rotation
	ROTATION_180 Rotation = 2
	ROTATION_270 Rotation = 3
	WIDTH        int      = 80
	HEIGHT       int      = 160
)

var once sync.Once
var display *Display

type Display struct {
	p   spi.PortCloser
	dev *st7735.Dev
}

func init() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
	}
}

func Init() (*Display, error) {
	var err error
	once.Do(func() {
		display = &Display{}
		display.p, err = spireg.Open("SPI0.1")
		if err != nil {
			return
		}
		display.dev, err = st7735.New(display.p.(spi.Port), gpioreg.ByName("GPIO9"), nil, nil, &st7735.DefaultOpts)
	})

	return display, err
}

func (d *Display) Close() {
}

func (d *Display) DrawImage(reader io.Reader) {
}

func (d *Display) DrawRAW(img image.Image) {
}

//func (d *Display) Rotate(rotation Rotation) {
//	d.dev.SetRotation(st7789.Rotation(rotation))
//}

func (d *Display) FillScreen(c color.RGBA) error {
	bounds := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{int(WIDTH), HEIGHT}}
	img := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			img.Set(x, y, c)
		}
	}
	return d.dev.DisplayImage(0, 0, img)
}

func (d *Display) SetPixel(x int16, y int16, c color.RGBA) {
}

// PowerOff the display
func (d *Display) PowerOff() error {
	dev, err := d.bldev()
	if err != nil {
		return err
	}
	dev.SetBacklight(false)

	return nil
}

// PowerOn the display
func (d *Display) PowerOn() error {
	dev, err := d.bldev()
	if err != nil {
		return err
	}
	dev.SetBacklight(true)

	return nil
}

func (d *Display) bldev() (*st7735.Dev, error) {
	backlightPin := gpioreg.ByName("GPIO12")
	p, err := spireg.Open("SPI0.0")
	if err != nil {
		return nil, err
	}
	defer p.Close()

	return st7735.New(p.(spi.Port), gpioreg.ByName("GPIO9"), nil, backlightPin, &st7735.DefaultOpts)
}
