package main

import (
	"flag"
	"fmt"
	"github.com/torbenschinke/i2c"
	"time"
)

func main() {
	iv := flag.Duration("interval", time.Second*10, "interval to poll sensors")
	flag.Parse()

	poller := i2c.NewPolling(*iv)
	poller.Dispatcher().Register(printObserver{})

	select {}
}

type printObserver struct {
}

func (p printObserver) OnTemperature(id i2c.ID, t i2c.T) {
	fmt.Printf("%#x: %.2f°\n", id, t.Celsius())
}

func (p printObserver) OnHumidity(id i2c.ID, h i2c.RH) {
	fmt.Printf("%#x: %.2f%%\n", id, h.Humidity())
}

func (p printObserver) OnError(id i2c.ID, err error) {
	fmt.Printf("%#x: %v\n", id, err)
}
