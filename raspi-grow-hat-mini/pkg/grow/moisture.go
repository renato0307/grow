package grow

import (
	"fmt"
	"log"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

const (
	Moisture1 = rpi.J8p16
	Moisture2 = rpi.J8p24
	Moisture3 = rpi.J8p22
)

type GrowHatMoistureReader struct {
	count           int64
	name            string
	offset          int
	reading         float64
	timeLastReading time.Time

	chip *gpiod.Chip
	line *gpiod.Line
}

func NewGrowHatMoistureReader(name string, offset int) (*GrowHatMoistureReader, error) {
	log.Println("(i) initializing", name, offset)
	r := &GrowHatMoistureReader{
		name:            name,
		offset:          offset,
		timeLastReading: time.Now(),
	}

	chipName := "gpiochip0"
	c, err := gpiod.NewChip(chipName)
	if err != nil {
		return nil, fmt.Errorf("could not initialize chip %s: %w", chipName, err)
	}

	l, err := c.RequestLine(r.offset,
		gpiod.WithRisingEdge,
		gpiod.WithEventHandler(r.handler))
	if err != nil {
		return nil, fmt.Errorf("could not request line to %d: %w", r.offset, err)
	}

	r.chip = c
	r.line = l

	return r, nil
}

func (r *GrowHatMoistureReader) Read() float64 {
	return r.reading
}

func (r *GrowHatMoistureReader) Close() error {
	return r.line.Close()
}

func (r *GrowHatMoistureReader) Name() string {
	return r.name
}

func (r *GrowHatMoistureReader) handler(evt gpiod.LineEvent) {
	log.Println("handling", evt.Offset)
	r.count += 1
	timeElapsed := float64(time.Now().UnixNano()-r.timeLastReading.UnixNano()) / 1000000000
	if timeElapsed >= 1.0 {
		r.reading = float64(r.count) / timeElapsed
		r.count = 0
		r.timeLastReading = time.Now()
	}
}
