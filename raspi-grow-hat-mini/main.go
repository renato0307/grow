package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/grow/monitor-raspi4/pkg/grow"
	"github.com/grow/monitor-raspi4/pkg/options"
	"github.com/grow/monitor-raspi4/pkg/publish"
)

type MoistureReader interface {
	Close() error
	Name() string
	Read() float64
}

type Publisher func(name string, value float64) error

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))

	options, err := options.Get()
	if err != nil {
		slog.Error("invalid options")
		os.Exit(1)
	}
	slog.Info("plants configured", "plants", options.Plants)
	slog.Info("publishers configured", "publishers", options.Publishers)

	readers := setupReaders(options.Plants)
	defer func() {
		for _, r := range readers {
			r.Close()
		}
	}()
	publishers := setupPublishers(options)
	for {
		readAndPublish(readers, publishers, options.Frequency)
	}

}

func setupReaders(plants []options.Plant) []MoistureReader {
	readers := make([]MoistureReader, len(plants))
	for i := range plants {
		r, err := grow.NewGrowHatMoistureReader(plants[i].Name, plants[i].Connector)
		if err != nil {
			slog.Error("could not init reader", "plant", plants[i].Name, "error", err)
			os.Exit(1)
		}
		readers[i] = r
	}
	return readers
}

func setupPublishers(opt options.Options) []publish.Publisher {
	publishers := []publish.Publisher{}
	for _, pt := range opt.Publishers {
		switch pt {
		case options.Console:
			publishers = append(publishers, publish.NewConsolePublisher())
		case options.NATS:
			natsPub, err := publish.NewNATSPublisher(opt.NATS)
			if err != nil {
				slog.Error("could not init NATS publisher", "error", err)
				os.Exit(1)
			}
			publishers = append(publishers, natsPub)
		}
	}
	return publishers
}

func readAndPublish(readers []MoistureReader, publishers []publish.Publisher, frequency time.Duration) {
	for _, reader := range readers {
		reading := reader.Read()
		slog.Debug("reading", "name", reader.Name(), "value", reading)
		for _, publisher := range publishers {
			err := publisher(publish.Reading{
				Timestamp: time.Now(),
				Name:      reader.Name(),
				Value:     reading,
			})
			if err != nil {
				slog.Error("could not publish", "error", err)
			}
		}
	}
	time.Sleep(frequency)
}
