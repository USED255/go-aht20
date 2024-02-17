package aht20

import "errors"

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
