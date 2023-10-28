package publish

import "log"

func NewConsolePublisher() Publisher {
	return func(name string, value float64) error {
		log.Printf("%s reads %.15f\n", name, value)
		return nil
	}
}
