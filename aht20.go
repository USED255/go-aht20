// https://github.com/tinygo-org/drivers/tree/release/aht20 的 https://github.com/d2r2/go-i2c 移植
package aht20

// 移植自 https://github.com/tinygo-org/drivers/tree/release/aht20
// 参考了 https://github.com/Chouffy/python_sensor_aht20
// 和 https://github.com/d2r2/go-bh1750

// 许可: https://github.com/tinygo-org/drivers/blob/81bc1bcad1862f719556d3c7ff411c3f56b143dd/LICENSE
/* Copyright (c) 2018-2022 The TinyGo Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */

// 相比原代码, 添加了包装, 做了必要的修改, 添加了一些 log.

import (
	"time"

	"github.com/d2r2/go-i2c"
)

// Device wraps an I2C connection to an AHT20 device.
type Device struct {
	bus *i2c.I2C
	// Address  uint16 // 注意, d2r2/go-i2c 需要在 i2c.NewI2C() 时传入地址.
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

// Configure the device
func (d *Device) Configure() {
	// Check initialization state
	status := d.Status()
	if status&STATUS_CALIBRATED == 1 {
		lg.Debug("AHT20 is initialized")
		// Device is initialized
		return
	}

	// Force initialization
	lg.Debug("Initialization AHT20")
	d.tx([]byte{CMD_INITIALIZE, 0x08, 0x00}, nil)
	time.Sleep(10 * time.Millisecond)
}

// Reset the device
func (d *Device) Reset() {
	lg.Debug("Reset sensor...")
	d.tx([]byte{CMD_SOFTRESET}, nil)
}

// Status of the device
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
		if data[0]&STATUS_CALIBRATED != 0 && data[0]&STATUS_BUSY == 0 {
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

func (d *Device) DeciRelHumidity() int32 {
	return (int32(d.humidity) * 1000) / 0x100000
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
