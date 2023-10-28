package publish

import (
	"fmt"
	"log/slog"
)

func NewNATSPublisher() Publisher {
	return func(name string, value float64) error {
		slog.Info("publising to NATS", "plant", name, "value", fmt.Sprintf("%.15f", value))
		return nil
	}
}
