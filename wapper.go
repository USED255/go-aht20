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

// 来自 https://github.com/d2r2/go-bh1750/blob/master/logger.go

// You can manage verbosity of log output
// in the package by changing last parameter value.
var lg = logger.NewPackageLogger("aht20",
	logger.DebugLevel,
	// logger.InfoLevel,
)

// New creates a new AHT20 connection. The I2C bus must already be configured.
func NewAHT20(bus *i2c.I2C) Device {
	aht20 := New(bus)
	//aht20.Reset()
	aht20.Configure()

	return aht20
}

func (d *Device) ReadRelHumidity() (float32, error) {
	err := d.Read()
	if err != nil {
		return 0, err
	}
	return d.RelHumidity(), nil
}

func (d *Device) ReadCelsius() (float32, error) {
	err := d.Read()
	if err != nil {
		return 0, err
	}
	return d.Celsius(), nil
}
