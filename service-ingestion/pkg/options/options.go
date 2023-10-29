package options

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

const (
	DefaultNATSURL     = "nats://192.168.1.129:4222,nats://192.168.1.130:4222,nats://192.168.1.131:4222"
	DefaultInfluxDBURL = "http://localhost:8086"
)

type NATSConfig struct {
	StreamName    string
	StreamSubject string
	URL           string
}

type InfluxDBConfig struct {
	Bucket       string
	Organization string
	Token        string
	URL          string
}

type Options struct {
	LogLevel *slog.LevelVar
	NATS     NATSConfig
	InfluxDB InfluxDBConfig
}

func Get() (Options, error) {

	opt := Options{}
	var logLevelValue string

	pflag.StringVar(&opt.NATS.StreamName, "nats-stream", "PlantReadings", "NATS stream name to publish messages")
	pflag.StringVar(&opt.NATS.StreamSubject, "nats-stream-sub", "PlantReadings.home", "NATS stream subject name to publish messages")
	pflag.StringVar(&opt.NATS.URL, "nats-url", DefaultNATSURL, "NATS URL to publish the messages")
	pflag.StringVar(&opt.InfluxDB.Bucket, "influxdb-bucket", "grow", "InfluxDB bucket to store metrics")
	pflag.StringVar(&opt.InfluxDB.URL, "influxdb-url", DefaultInfluxDBURL, "InfluxDB URL to send metrics")
	pflag.StringVar(&opt.InfluxDB.Token, "influxdb-token", os.Getenv("INFLUXDB_TOKEN"), "InfluxDB token to send metrics")
	pflag.StringVar(&opt.InfluxDB.Organization, "influxdb-org", "influxdata", "InfluxDB organization to send metrics")
	pflag.StringVar(&logLevelValue, "log-level", "info", "Changes the log level like info, warn, error, and debug")

	pflag.Parse()

	levelVar := &slog.LevelVar{}
	err := levelVar.UnmarshalText([]byte(logLevelValue))
	if err != nil {
		return opt, fmt.Errorf("invalid log level value: %s", logLevelValue)
	}
	opt.LogLevel = levelVar

	return opt, nil
}
