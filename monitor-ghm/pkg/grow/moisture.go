package grow

import (
	"fmt"
	"log/slog"
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

	minMoisture float64
	maxMoisture float64

	chip *gpiod.Chip
	line *gpiod.Line
}

func NewGrowHatMoistureReader(name string, offset int, minMoisture, maxMoisture float64) (*GrowHatMoistureReader, error) {
	slog.Debug("initializing reader", "name", name, "offset", offset)
	r := &GrowHatMoistureReader{
		name:            name,
		offset:          offset,
		timeLastReading: time.Now(),
		minMoisture:     minMoisture,
		maxMoisture:     maxMoisture,
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
	return (r.maxMoisture - r.reading) * 100 / (r.maxMoisture - r.minMoisture)
}

func (r *GrowHatMoistureReader) Close() error {
	return r.line.Close()
}

func (r *GrowHatMoistureReader) Name() string {
	return r.name
}

func (r *GrowHatMoistureReader) handler(evt gpiod.LineEvent) {
	slog.Debug("handling value", "offset", evt.Offset)
	r.count += 1
	timeElapsed := float64(time.Now().UnixNano()-r.timeLastReading.UnixNano()) / 1000000000
	if timeElapsed >= 1.0 {
		r.reading = float64(r.count) / timeElapsed
		r.count = 0
		r.timeLastReading = time.Now()
	}
}
