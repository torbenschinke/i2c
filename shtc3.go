package i2c

import (
	"fmt"
	"github.com/torbenschinke/i2c/sensor/shtc3"
	"periph.io/x/conn/v3/i2c"
	"time"
)

type shtc3x struct {
	t        shtc3.Temperature
	h        shtc3.Humidity
	dev      dev
	id       shtc3.Ident
	lastPoll time.Time
}

func newShtc3x(bus i2c.BusCloser) (*shtc3x, error) {
	s := &shtc3x{
		dev: dev{
			addr: shtc3.I2CAddr,
			bus:  bus,
		},
	}

	id, err := shtc3.ReadID(s.dev)
	if err != nil {
		return nil, fmt.Errorf("cannot read shtc3 id: %w", err)
	}

	if !id.SHTC3() {
		return nil, fmt.Errorf("i2c responded but not shtc3 identifier: %v", id)
	}

	s.id = id

	return s, nil
}

func (s *shtc3x) ID() ID {
	return ID(s.id)
}

func (s *shtc3x) Poll(d *Dispatcher) {
	if s.lastPoll.IsZero() || s.lastPoll.After(d.Time().Add(24*time.Hour)) {
		if err := shtc3.SoftReset(s.dev); err != nil {
			d.OnError(s.ID(), fmt.Errorf("cannot soft reset sensor shtc3x: %w", err))
			return
		}
	}

	temp, err := shtc3.ReadTemperature(s.dev)
	if err != nil {
		d.OnError(s.ID(), fmt.Errorf("cannot read sensor shtc3x temp: %w", err))
		return
	}

	hum, err := shtc3.ReadHumidity(s.dev)
	if err != nil {
		d.OnError(s.ID(), fmt.Errorf("cannot read sensor shtc3x hum: %w", err))
		return
	}

	d.OnTemperature(s.ID(), T(temp.Celsius()*1000))
	d.OnHumidity(s.ID(), RH(hum.Relative()*1000))

	s.lastPoll = d.Time()
}
