package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/grow/service-ingestion/pkg/options"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Reading struct {
	Timestamp time.Time
	Name      string
	Value     string
}

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

		reading := Reading{}
		err := json.Unmarshal(msg.Data(), &reading)
		if err != nil {
			slog.Warn("message with invalid format - ignoring it", "error", err)
			msg.Ack()
			return
		}
		value, err := strconv.ParseFloat(reading.Value, 64)
		if err != nil {
			slog.Warn("message with invalid format - ignoring it", "error", err)
			msg.Ack()
			return
		}

		slog.Debug("writing to influxdb", "name", reading.Name, "value", value, "ts", reading.Timestamp)
		err = writeToInfluxDB(reading.Name, reading.Timestamp, value, options.InfluxDB)
		if err != nil {
			slog.Warn("could not write reading", "error", err)
			return
		}
		slog.Debug("writing done", "name", reading.Name, "value", value, "ts", reading.Timestamp)
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

func writeToInfluxDB(name string, timestamp time.Time, value float64, influxDB options.InfluxDBConfig) error {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(influxDB.URL, influxDB.Token)

	//Wwrite point asynchronously
	writeAPI := client.WriteAPIBlocking(influxDB.Organization, influxDB.Bucket)

	// Create point
	p := influxdb2.NewPointWithMeasurement("soil_moisture").
		AddTag("name", name).
		AddField("value", value).
		SetTime(timestamp)

	return writeAPI.WritePoint(context.Background(), p)
}
