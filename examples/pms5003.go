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
