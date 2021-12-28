package i2c

import (
	"fmt"
	"sync"
	"time"
)

// T is the temperature in milli-Celsius (1°C = 1000T)
type T int32

// Celsius returns just the degree value.
func (t T) Celsius() float64 {
	return float64(t) / 1000
}

func (t T) String() string {
	return fmt.Sprintf("%3.2f°C", t.Celsius())
}

// RH is the relative humidity in milli-RH (1% = 1000RH)
type RH int32

// Humidity returns the value in relative percent.
func (v RH) Humidity() float64 {
	return float64(v) / 1000
}

func (v RH) String() string {
	return fmt.Sprintf("%3.2f%%", v.Humidity())
}

// ID is just some sensor identifier.
type ID int32

type TemperatureObserver interface {
	OnTemperature(id ID, t T)
}

type HumidityObserver interface {
	OnHumidity(id ID, h RH)
}

type ErrorObserver interface {
	OnError(id ID, err error)
}

type Dispatcher struct {
	mutex        sync.RWMutex
	listener     map[int]interface{}
	lastHnd      int
	lastTickTime time.Time
}

// OnTemperature dispatches a new temperature event.
func (d *Dispatcher) OnTemperature(id ID, t T) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, l := range d.listener {
		if cb, ok := l.(TemperatureObserver); ok {
			cb.OnTemperature(id, t)
		}
	}
}

// OnHumidity dispatches a new humidity event.
func (d *Dispatcher) OnHumidity(id ID, h RH) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, l := range d.listener {
		if cb, ok := l.(HumidityObserver); ok {
			cb.OnHumidity(id, h)
		}
	}
}

// OnError dispatches a new error event.
func (d *Dispatcher) OnError(id ID, err error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, l := range d.listener {
		if cb, ok := l.(ErrorObserver); ok {
			cb.OnError(id, err)
		}
	}
}

func (d *Dispatcher) Register(observer interface{}) (hnd int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.lastHnd++
	d.listener[d.lastHnd] = observer
	return d.lastHnd
}

func (d *Dispatcher) Unregister(hnd int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.listener, hnd)
}

func (d *Dispatcher) Time() time.Time {
	return d.lastTickTime
}

type PrintObserver struct {
}

func (p PrintObserver) OnTemperature(id ID, t T) {
	fmt.Printf("%v -> %v\n", id, t.String())
}

func (p PrintObserver) OnError(id ID, err error) {
	fmt.Printf("%v -> %v\n", id, err)
}

func (p PrintObserver) OnHumidity(id ID, h RH) {
	fmt.Printf("%v -> %v\n", id, h.String())
}
