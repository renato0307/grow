package publish

import "time"

type Reading struct {
	Timestamp time.Time
	Name      string
	Value     float64
}

type Publisher func(Reading) error
