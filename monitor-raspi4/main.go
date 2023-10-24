package main

import (
	"fmt"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

var count int64
var timeLastReading = time.Now()
var reading float64

func handler(evt gpiod.LineEvent) {
	count += 1
	timeElapsed := float64(time.Now().UnixNano()-timeLastReading.UnixNano()) / 1000000000
	if timeElapsed >= 1.0 {
		reading = float64(count) / timeElapsed
		count = 0
		timeLastReading = time.Now()
	}
}

func main() {
	c, _ := gpiod.NewChip("gpiochip0")
	l, _ := c.RequestLine(rpi.J8p16,
		gpiod.WithRisingEdge,
		gpiod.WithEventHandler(handler))
	defer l.Close()

	for {
		fmt.Printf("%.15f\n", reading)
		time.Sleep(1 * time.Second)
	}
}
