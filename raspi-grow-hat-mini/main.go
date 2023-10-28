package main

import (
	"log"
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

func main() {
	options, err := options.Get()
	if err != nil {
		log.Fatal("invalid options")
	}
	log.Println("(i) running for", options.Plants)
	log.Println("(i) publishing to", options.Publishers)

	readers := setupReaders(options.Plants)
	defer func() {
		for _, r := range readers {
			r.Close()
		}
	}()
	publishers := setupPublishers(options.Publishers)
	for {
		readAndPublish(readers, publishers, options.Frequency)
	}

}

func setupReaders(plants []options.Plant) []MoistureReader {
	readers := make([]MoistureReader, len(plants))
	for i := range plants {
		r, err := grow.NewGrowHatMoistureReader(plants[i].Name, plants[i].Connector)
		if err != nil {
			log.Fatalf("could not init reader for %s: %w", plants[i].Name, err)
		}
		readers[i] = r
	}
	return readers
}

func setupPublishers(publisherTypes []string) []publish.Publisher {
	publishers := []publish.Publisher{}
	for _, pt := range publisherTypes {
		switch pt {
		case options.Console:
			publishers = append(publishers, publish.NewConsolePublisher())
		case options.NATS:
			publishers = append(publishers, publish.NewNATSPublisher())
		}
	}
	return publishers
}

func readAndPublish(readers []MoistureReader, publishers []publish.Publisher, frequency time.Duration) {
	for _, reader := range readers {
		reading := reader.Read()
		log.Println("read", reading)
		for _, publish := range publishers {
			publish(reader.Name(), reading)
		}
	}
	time.Sleep(frequency)
}
