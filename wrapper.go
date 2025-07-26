package aht20

import (
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

func (d *Device) controllerTransmit(w []byte) error {
	_, err := d.bus.WriteBytes(w)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) controllerReceive(r []byte) error {
	_, err := d.bus.ReadBytes(r)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) tx(w, r []byte) error {
	if len(w) > 0 {
		if err := d.controllerTransmit(w); nil != err {
			return err
		}
	}

	if len(r) > 0 {
		if err := d.controllerReceive(r); nil != err {
			return err
		}
	}

	return nil
}

// NewAHT20 return new sensor instance.
func NewAHT20(bus *i2c.I2C) Device {
	aht20 := New(bus)
	//aht20.Reset()
	aht20.Configure()

	return aht20
}

func (d *Device) ReadRelativeHumidity() (float32, error) {
	err := d.Read()
	if err != nil {
		return 0, err
	}
	return d.RelHumidity(), nil
}

// ReadTemperatureC reads and calculates temrature in C (celsius).
func (d *Device) ReadTemperatureC() (float32, error) {
	err := d.Read()
	if err != nil {
		return 0, err
	}
	return d.Celsius(), nil
}

// 来自  https://github.com/d2r2/go-bh1750/blob/master/logger.go
// 许可: https://github.com/d2r2/go-bh1750/blob/master/LICENSE
/* MIT License

Copyright (c) 2018 Denis Dyakov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE. */

// You can manage verbosity of log output
// in the package by changing last parameter value.
var lg = logger.NewPackageLogger("aht20",
	logger.DebugLevel,
	// logger.InfoLevel,
)
