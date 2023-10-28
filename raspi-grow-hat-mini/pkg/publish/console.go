package publish

import (
	"fmt"
	"log/slog"
)

func NewConsolePublisher() func(Reading) error {
	return func(r Reading) error {
		slog.Info("reading", "plant", r.Name, "value", fmt.Sprintf("%.15f", r.Value))
		return nil
	}
}
