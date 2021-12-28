package i2c

import (
	"fmt"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"time"
)

type Polling struct {
	busses     []i2c.BusCloser
	sensors    []Sensor
	dispatcher *Dispatcher
	ticker     *time.Ticker
	done       chan bool
}

// NewPolling tries to identify the correct i2c bus and the connected devices and immediately starts polling.
// It does never fail, and instead just returns an instance which has no sensors.
func NewPolling(interval time.Duration) *Polling {
	polling := &Polling{
		dispatcher: &Dispatcher{
			listener: map[int]interface{}{},
		},
		ticker: time.NewTicker(interval),
		done:   make(chan bool),
	}

	if _, err := host.Init(); err != nil {
		fmt.Printf("cannot init periph.io host: %v\n", err)
	}

	for _, ref := range i2creg.All() {
		if ref.Number < 0 {
			continue
		}

		bus, err := ref.Open()
		if err != nil {
			fmt.Printf("cannot open available i2c bus: %v\n", err)
			continue
		}

		polling.busses = append(polling.busses, bus)
		if shtc3Sensor, err := newShtc3x(bus); err == nil {
			fmt.Printf("found shtc3 sensor on i2c %v.%v\n", ref.Number, shtc3Sensor.ID())
			polling.sensors = append(polling.sensors, shtc3Sensor)
		}
	}

	go polling.poll()

	return polling
}

func (p *Polling) poll() {
	for {
		select {
		case <-p.done:
			return
		case t := <-p.ticker.C:
			p.dispatcher.lastTickTime = t
			for _, sensor := range p.sensors {
				sensor.Poll(p.dispatcher)
			}
		}
	}
}

func (p *Polling) Close() error {
	p.ticker.Stop()
	p.done <- true
	return nil
}
