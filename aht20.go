// https://github.com/tinygo-org/drivers/tree/release/aht20 的 https://github.com/d2r2/go-i2c 移植
package aht20

// 移植自 https://github.com/tinygo-org/drivers/tree/release/aht20
// 参考了 https://github.com/Chouffy/python_sensor_aht20
// 和 https://github.com/d2r2/go-bh1750
// 许可: https://github.com/tinygo-org/drivers/blob/81bc1bcad1862f719556d3c7ff411c3f56b143dd/LICENSE

import (
	"time"

	"github.com/d2r2/go-i2c"
)

// Device wraps an I2C connection to an AHT20 device.
type Device struct {
	bus      *i2c.I2C
	humidity uint32
	temp     uint32
}

// New creates a new AHT20 connection. The I2C bus must already be
// configured.
//
// This function only creates the Device object, it does not touch the device.
func New(bus *i2c.I2C) Device {
	return Device{
		bus: bus,
	}
}

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
	d.tx([]byte{CMD_INITIALIZE, 0x08, 0x00}, nil)
	time.Sleep(10 * time.Millisecond)
}

// Reset the AHT20
func (d *Device) Reset() {
	lg.Debug("Reset sensor...")
	d.tx([]byte{CMD_SOFTRESET}, nil)
	time.Sleep(20 * time.Millisecond)
}

// Status of the AHT20
func (d *Device) Status() byte {
	data := []byte{0}

	d.tx([]byte{CMD_STATUS}, data)

	return data[0]
}

// Read the temperature and humidity
//
// The actual temperature and humidity are stored
// and can be accessed using `Temp` and `Humidity`.
func (d *Device) Read() error {
	d.tx([]byte{CMD_TRIGGER, 0x33, 0x00}, nil)

	data := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	lg.Debug("Reading...")
	for retry := 0; retry < 3; retry++ {
		time.Sleep(80 * time.Millisecond)
		err := d.tx(nil, data)
		if err != nil {
			return err
		}

		// If measurement complete, store values
		if data[0]&0x04 != 0 && data[0]&0x80 == 0 {
			d.humidity = uint32(data[1])<<12 | uint32(data[2])<<4 | uint32(data[3])>>4
			d.temp = (uint32(data[3])&0xF)<<16 | uint32(data[4])<<8 | uint32(data[5])
			return nil
		}
		lg.Debug("retry...")
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
