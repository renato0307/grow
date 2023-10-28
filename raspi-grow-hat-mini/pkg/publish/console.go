package publish

import (
	"fmt"
	"log/slog"
)

func NewConsolePublisher() Publisher {
	return func(name string, value float64) error {
		slog.Info("reading", "plant", name, "value", fmt.Sprintf("%.15f", value))
		return nil
	}
}
