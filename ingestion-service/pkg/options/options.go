package options

import (
	"fmt"
	"log/slog"

	"github.com/spf13/pflag"
)

const (
	DefaultNATSURL       = "nats://192.168.1.129:4222,nats://192.168.1.130:4222,nats://192.168.1.131:4222"
	DefaultPrometheusURL = "http://localhost:9090/api/v1/write"
)

type NATSConfig struct {
	StreamName    string
	StreamSubject string
	URL           string
}

type PrometheusConfig struct {
	URL string
}

type Options struct {
	LogLevel   *slog.LevelVar
	NATS       NATSConfig
	Prometheus PrometheusConfig
	ProbesAddr string
}

func Get() (Options, error) {

	opt := Options{}
	var logLevelValue string

	pflag.StringVar(&opt.NATS.StreamName, "nats-stream", "PlantReadings", "NATS stream name to publish messages")
	pflag.StringVar(&opt.NATS.StreamSubject, "nats-stream-sub", "PlantReadings.home", "NATS stream subject name to publish messages")
	pflag.StringVar(&opt.NATS.URL, "nats-url", DefaultNATSURL, "NATS URL to publish the messages")
	pflag.StringVar(&opt.Prometheus.URL, "prom-url", DefaultPrometheusURL, "Prometheus URL to send metrics")
	pflag.StringVar(&opt.ProbesAddr, "probes-addr", ":8222", "The bind address for health and readiness probes")
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
