package i2c

import (
	"fmt"
	"github.com/torbenschinke/i2c/sensor/sht3x"
	"github.com/torbenschinke/i2c/sensor/shtc3"
	"log"
	"math/rand"
	"periph.io/x/conn/v3/i2c"
	"time"
)

type sht3 struct {
	t        shtc3.Temperature
	h        shtc3.Humidity
	dev      dev
	id       int32
	lastPoll time.Time
}

func findSht3(bus i2c.BusCloser) ([]*sht3, error) {
	var sensors []*sht3
	for _, addr := range sht3x.I2CAddrs {
		s, err := newSht3x(bus, addr)
		if err != nil {
			log.Printf("probing sht3x on %#x failed (%v)\n", addr, err)
		} else {
			sensors = append(sensors, s)
		}
	}

	return sensors, nil
}

func newSht3x(bus i2c.BusCloser, addr uint16) (*sht3, error) {
	s := &sht3{
		dev: dev{
			addr: addr,
			bus:  bus,
		},
		id: rand.Int31(),
	}

	err := sht3x.SoftReset(s.dev)
	if err != nil {
		return nil, fmt.Errorf("cannot reset sht3x: %w", err)
	}

	if err := sht3x.ClearStatus(s.dev); err != nil {
		return nil, fmt.Errorf("cannot clear status sht3x: %w", err)
	}

	status, err := sht3x.ReadStatus(s.dev) // seems to read as 0x8010 without clearing
	if err != nil {
		return nil, fmt.Errorf("cannot read status sht3x: %w", err)
	}

	if status != 0x0 {
		return nil, fmt.Errorf("unexpected status: %#x", status)
	}

	return s, nil
}

func (s *sht3) Poll(d *Dispatcher) {
	if s.lastPoll.IsZero() || s.lastPoll.After(d.Time().Add(24*time.Hour)) {
		if err := sht3x.SoftReset(s.dev); err != nil {
			d.OnError(ID(s.id), fmt.Errorf("cannot soft reset sensor sht3x: %w", err))
			return
		}
	}

	temp, hum, err := sht3x.ReadTempHum(s.dev)
	if err != nil {
		d.OnError(ID(s.id), fmt.Errorf("cannot read sensor shtc3x hum: %w", err))
		return
	}

	d.OnTemperature(ID(s.id), T(temp.Celsius()*1000))
	d.OnHumidity(ID(s.id), RH(hum.Relative()*1000))
}
