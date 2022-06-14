package aht20

// 移植自 https://github.com/tinygo-org/drivers/tree/release/aht20

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

// 用于 github.com/d2r2/go-i2c 的包装
type AHT20 struct {
	bus      *i2c.I2C
	humidity uint32
	temp     uint32
}

func (i *AHT20) controllerTransmit(w []byte) error {
	_, err := i.bus.WriteBytes(w)
	if err != nil {
		return err
	}
	return nil
}
func (i *AHT20) controllerReceive(r []byte) error {
	_, err := i.bus.ReadBytes(r)
	if err != nil {
		return err
	}
	return nil
}
func (i *AHT20) Tx(w, r []byte) error {
	if len(w) > 0 {
		if err := i.controllerTransmit(w); nil != err {
			return err
		}
	}

	if len(r) > 0 {
		if err := i.controllerReceive(r); nil != err {
			return err
		}
	}

	return nil
}

// New creates a new AHT20 connection. The I2C bus must already be
// configured.
//
// This function only creates the AHT20 object, it does not touch the AHT20.
func NewAHT20(bus *i2c.I2C) AHT20 {
	aht20 := AHT20{
		bus: bus,
	}
	//aht20.Reset()
	aht20.Configure()
	return aht20
}

// Configure the AHT20
func (d *AHT20) Configure() {
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
func (d *AHT20) Reset() {
	lg.Debug("Reset sensor...")
	d.Tx([]byte{CMD_SOFTRESET}, nil)
	time.Sleep(20 * time.Millisecond)
}

// Status of the AHT20
func (d *AHT20) Status() byte {
	data := []byte{0}

	d.Tx([]byte{CMD_STATUS}, data)

	return data[0]
}

// Read the temperature and humidity
//
// The actual temperature and humidity are stored
// and can be accessed using `Temp` and `Humidity`.
func (d *AHT20) Read() error {
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

func (d *AHT20) ReadWithRetry(i int) error {
	err := d.Read()
	for {
		if err == nil {
			return nil
		}
		lg.Debug("Retry...")
		if i == 0 {
			return err
		}
		err = d.Read()
		i = i - 1
	}
}

func (d *AHT20) RawHumidity() uint32 {
	return d.humidity
}

func (d *AHT20) RawTemp() uint32 {
	return d.temp
}

func (d *AHT20) RelHumidity() float32 {
	return (float32(d.humidity) * 100) / 0x100000
}

// Temperature in degrees celsius
func (d *AHT20) Celsius() float32 {
	return (float32(d.temp*200.0) / 0x100000) - 50
}

// https://github.com/d2r2/go-bh1750/blob/master/logger.go 抄来的
// You can manage verbosity of log output
// in the package by changing last parameter value.
var lg = logger.NewPackageLogger("aht20",
	logger.DebugLevel,
	// logger.InfoLevel,
)
