package main

import (
	"log"
	"time"

	"github.com/grow/monitor-raspi4/pkg/grow"
)

type MoistureReader interface {
	Close() error
	Name() string
	Read() float64
}

func main() {
	plants := []struct {
		name      string
		connector int
	}{
		{
			name:      "espadas",
			connector: grow.Moisture1,
		},
		{
			name:      "abacateiro",
			connector: grow.Moisture2,
		},
		{
			name:      "pilea peperomioides",
			connector: grow.Moisture3,
		},
	}

	readers := make([]MoistureReader, len(plants))
	for i := range plants {
		r, err := grow.NewGrowHatMoistureReader(plants[i].name, plants[i].connector)
		if err != nil {
			log.Fatalf("could not init reader for %s: %w", plants[i].name, err)
		}
		defer r.Close()
		readers[i] = r
	}

	for {
		for _, r := range readers {
			log.Printf("%s reads %.15f\n", r.Name(), r.Read())
		}
		time.Sleep(1 * time.Second)
	}
}
