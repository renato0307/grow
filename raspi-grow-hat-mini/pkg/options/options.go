package options

import (
	"time"

	"github.com/grow/monitor-raspi4/pkg/grow"
)

const (
	Console = "console" // Console publisher
	NATS    = "nats"    // NATS publisher
)

type Plant struct {
	Name      string
	Connector int
}

type NATSConfig struct {
	URL string

	StreamName    string
	StreamSubject string
}

type Options struct {
	Frequency  time.Duration
	NATS       NATSConfig
	Plants     []Plant
	Publishers []string
}

func Get() (Options, error) {
	return Options{
		Frequency:  5 * time.Second,
		Publishers: []string{Console, NATS},
		NATS: NATSConfig{
			URL:           "nats://192.168.1.129:4222,nats://192.168.1.130:4222,nats://192.168.1.131:4222",
			StreamName:    "PlantReadings",
			StreamSubject: "PlantReadings.new",
		},
		Plants: []Plant{
			{
				Name:      "espadas",
				Connector: grow.Moisture1,
			},
			{
				Name:      "abacateiro",
				Connector: grow.Moisture2,
			},
			{
				Name:      "pilea peperomioides",
				Connector: grow.Moisture3,
			},
		},
	}, nil
}
