package aht20

// 移植自 https://github.com/tinygo-org/drivers/tree/release/aht20
// 参考了 https://github.com/Chouffy/python_sensor_aht20
// 和 https://github.com/d2r2/go-bh1750

import (
	"errors"
	"time"

	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

const (
	CMD_INITIALIZE = 0xBE
	CMD_STATUS     = 0x71
	CMD_TRIGGER    = 0xAC
	CMD_SOFTRESET  = 0xBA

	STATUS_BUSY       = 0x80
	STATUS_CALIBRATED = 0x08
)

var (
	ErrBusy    = errors.New("AHT20 busy")
	ErrTimeout = errors.New("timeout")
)

// 来自 https://github.com/d2r2/go-bh1750/blob/master/logger.go

// You can manage verbosity of log output
// in the package by changing last parameter value.
var lg = logger.NewPackageLogger("aht20",
	logger.DebugLevel,
	// logger.InfoLevel,
)

type Device struct {
	bus      *i2c.I2C
	humidity uint32
	temp     uint32
}

// 用于 https://github.com/d2r2/go-i2c 的包装
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

func (d *Device) Tx(w, r []byte) error {
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

func NewAHT20(bus *i2c.I2C) Device {
	aht20 := Device{
		bus: bus,
	}
	//aht20.Reset()
	aht20.Configure()

	return aht20
}

// 以下来自 https://github.com/tinygo-org/drivers/blob/release/aht20/aht20.go
// 许可: https://github.com/tinygo-org/drivers/blob/81bc1bcad1862f719556d3c7ff411c3f56b143dd/LICENSE

// Configure the AHT20
func (d *Device) Configure() {
	// Check initialization state
	status := d.Status()
	if status&0x08 == 1 {
		lg.Debug("AHT20 is initialized")
		// AHT20 is initialized
		return
	}

	// Force initialization
	lg.Debug("Initialization AHT20")
	d.Tx([]byte{CMD_INITIALIZE, 0x08, 0x00}, nil)
	time.Sleep(10 * time.Millisecond)
}

// Reset the AHT20
func (d *Device) Reset() {
	lg.Debug("Reset sensor...")
	d.Tx([]byte{CMD_SOFTRESET}, nil)
	time.Sleep(20 * time.Millisecond)
}

// Status of the AHT20
func (d *Device) Status() byte {
	data := []byte{0}

	d.Tx([]byte{CMD_STATUS}, data)

	return data[0]
}

// Read the temperature and humidity
//
// The actual temperature and humidity are stored
// and can be accessed using `Temp` and `Humidity`.
func (d *Device) Read() error {
	d.Tx([]byte{CMD_TRIGGER, 0x33, 0x00}, nil)

	data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for retry := 0; retry < 3; retry++ {
		time.Sleep(80 * time.Millisecond)
		err := d.Tx(nil, data)
		if err != nil {
			return err
		}

		// If measurement complete, store values
		if data[0]&0x04 != 0 && data[0]&0x80 == 0 {
			d.humidity = uint32(data[1])<<12 | uint32(data[2])<<4 | uint32(data[3])>>4
			d.temp = (uint32(data[3])&0xF)<<16 | uint32(data[4])<<8 | uint32(data[5])
			return nil
		}
	}

	return ErrTimeout
}

func (d *Device) RawHumidity() uint32 {
	return d.humidity
}

func (d *Device) RawTemp() uint32 {
	return d.temp
}

func (d *Device) RelHumidity() float32 {
	return (float32(d.humidity) * 100) / 0x100000
}

// Temperature in degrees celsius
func (d *Device) Celsius() float32 {
	return (float32(d.temp*200.0) / 0x100000) - 50
}

// Temperature in mutiples of one tenth of a degree celsius
//
// Using this method avoids floating point calculations.
func (d *Device) DeciCelsius() int32 {
	return ((int32(d.temp) * 2000) / 0x100000) - 500
}
