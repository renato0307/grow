package publish

import "log"

func NewNATSPublisher() Publisher {
	return func(name string, value float64) error {
		log.Printf("publishing %s:%.15f to NATS\n", name, value)
		return nil
	}
}
