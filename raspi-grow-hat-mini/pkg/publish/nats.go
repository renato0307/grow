package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/grow/monitor-raspi4/pkg/options"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSPublisher struct {
	nc            *nats.Conn
	js            jetstream.JetStream
	streamSubject string
}

func NewNATSPublisher(config options.NATSConfig) (Publisher, error) {
	nc, err := nats.Connect(config.URL)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to nats %s: %w", config.URL, err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to jetstream: %w", err)
	}

	// this is an idempotent operation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     config.StreamName,
		Subjects: []string{config.StreamSubject},
		Replicas: 3,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create stream %s: %w", config.StreamName, err)
	}

	np := &NATSPublisher{
		nc:            nc,
		js:            js,
		streamSubject: config.StreamSubject,
	}

	return np.Publish, nil
}

func (np *NATSPublisher) Publish(r Reading) error {
	formattedValue := fmt.Sprintf("%.15f", r.Value)
	slog.Info("publishing to NATS", "name", r.Name, "value", formattedValue)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	data := map[string]string{
		"name":      r.Name,
		"value":     formattedValue,
		"timestamp": r.Timestamp.UTC().Format(time.RFC3339),
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could marshall data to send to jetstreams: %w", err)
	}

	ack, err := np.js.PublishMsg(ctx, &nats.Msg{
		Data:    rawData,
		Subject: np.streamSubject,
	})
	if err != nil {
		return fmt.Errorf("could not send message to %s: %w", np.streamSubject, err)
	}

	slog.Debug("published msg to jetstream", "sequence", ack.Sequence, "stream", ack.Stream)
	return nil
}
