package i2c

import "periph.io/x/conn/v3/i2c"

type Sensor interface {
	Poll(*Dispatcher)
}

type dev struct {
	addr uint16
	bus  i2c.BusCloser
}

func (s dev) Read(p []byte) (n int, err error) {
	return len(p), s.bus.Tx(s.addr, nil, p)
}

func (s dev) Write(p []byte) (n int, err error) {
	return len(p), s.bus.Tx(s.addr, p, nil)
}
