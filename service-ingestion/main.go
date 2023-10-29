package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grow/service-ingestion/pkg/options"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	options, err := options.Get()
	if err != nil {
		fmt.Println("invalid options:", err)
		os.Exit(1)
	}
	slogOptions := &slog.HandlerOptions{
		Level: options.LogLevel,
	}
	handler := slog.NewTextHandler(os.Stdout, slogOptions)
	slog.SetDefault(slog.New(handler))

	// connect to nats server
	nc, _ := nats.Connect(options.NATS.URL)

	// create jetstream context from nats connection
	js, _ := jetstream.New(nc)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create a consumer (this is an idempotent operation)
	cons, err := js.CreateConsumer(ctx, options.NATS.StreamName, jetstream.ConsumerConfig{
		Durable:   "PlantReadingsIngestion",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		slog.Error("could not create consumer", "error", err)
		os.Exit(1)
	}

	// consume messages from the consumer in callback
	cc, err := cons.Consume(func(msg jetstream.Msg) {
		slog.Debug("received jetstream message", "msg", string(msg.Data()))
		msg.Ack()
	})
	if err != nil {
		slog.Error("error consuming messages", "error", err)
		os.Exit(1)
	}
	defer cc.Stop()

	// waits for termination
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
