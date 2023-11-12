package options

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/grow/monitor-ghm/pkg/grow"
	"github.com/spf13/pflag"
)

const (
	Console         = "console" // Console publisher
	NATS            = "nats"    // NATS publisher
	DefaultNATSURL  = "nats://192.168.1.2:4222"
	SensorSeparator = "|"
	MaxMoisture     = 6.5
	MinMoisture     = 25.5
)

var DefaultSensors = []string{
	fmt.Sprintf("%s%s%d", "espadas", SensorSeparator, grow.Moisture1),
	fmt.Sprintf("%s%s%d", "abacateiro", SensorSeparator, grow.Moisture2),
	fmt.Sprintf("%s%s%d", "pilea", SensorSeparator, grow.Moisture3),
}

type Sensors struct {
	Name        string
	Connector   int
	MaxMoisture float64
	MinMoisture float64
}

type NATSConfig struct {
	URL string

	StreamName    string
	StreamSubject string
}

type Options struct {
	Frequency  time.Duration
	NATS       NATSConfig
	Publishers []string
	Sensors    []Sensors
	LogLevel   *slog.LevelVar
}

func Get() (Options, error) {

	opt := Options{}
	var sensors []string
	var logLevelValue string

	pflag.DurationVar(&opt.Frequency, "readings-frequency", 5*time.Minute, "How frequently data is read from the sensors")
	pflag.StringArrayVar(&opt.Publishers, "publisher", []string{NATS}, "Which data publishers to use like console and nats")
	pflag.StringVar(&opt.NATS.URL, "nats-url", DefaultNATSURL, "NATS URL to publish the messages")
	pflag.StringVar(&opt.NATS.StreamName, "nats-stream", "PlantReadings", "NATS stream name to publish messages")
	pflag.StringVar(&opt.NATS.StreamSubject, "nats-stream-sub", "PlantReadings.home", "NATS stream subject name to publish messages")
	pflag.StringArrayVar(&sensors, "sensor", DefaultSensors, `List of sensors in the "<name>,<sensor-pin>" format`)
	pflag.StringVar(&logLevelValue, "log-level", "info", "Changes the log level like info, warn, error, and debug")

	pflag.Parse()

	levelVar := &slog.LevelVar{}
	err := levelVar.UnmarshalText([]byte(logLevelValue))
	if err != nil {
		return opt, fmt.Errorf("invalid log level value: %s", logLevelValue)
	}
	opt.LogLevel = levelVar

	for _, s := range sensors {
		sensorCfg := strings.Split(s, SensorSeparator)
		if len(sensorCfg) < 2 || len(sensorCfg) > 4 {
			return opt, fmt.Errorf("invalid sensor value: %s", s)
		}
		connector, err := strconv.Atoi(sensorCfg[1])
		if err != nil {
			return opt, fmt.Errorf("invalid connector value: %s", sensorCfg[1])
		}

		minMoisture := MinMoisture
		maxMoisture := MaxMoisture
		if len(sensorCfg) >= 3 {
			minMoisture, err = strconv.ParseFloat(sensorCfg[2], 64)
			if err != nil {
				return opt, fmt.Errorf("invalid mininum moisture value: %s", sensorCfg[2])
			}
		}
		if len(sensorCfg) == 4 {
			maxMoisture, err = strconv.ParseFloat(sensorCfg[3], 64)
			if err != nil {
				return opt, fmt.Errorf("invalid mininum moisture value: %s", sensorCfg[2])
			}
		}
		opt.Sensors = append(opt.Sensors, Sensors{
			Name:        sensorCfg[0],
			Connector:   connector,
			MaxMoisture: minMoisture,
			MinMoisture: maxMoisture,
		})
	}

	return opt, nil
}
