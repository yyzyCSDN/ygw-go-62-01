package main

import (
	"flag"
	"time"
)

// Config carries the runtime options of the column store server.
type Config struct {
	Listen    string
	DataDir   string
	Retention time.Duration
	Chunk     int
}

// DefaultConfig returns the standard development configuration.
func DefaultConfig() Config {
	return Config{
		Listen:    "127.0.0.1:18080",
		DataDir:   ".columnstore-data",
		Retention: 24 * time.Hour,
		Chunk:     4,
	}
}

// ParseConfig reads flags and applies overrides on top of the defaults.
func ParseConfig(args []string) (Config, error) {
	cfg := DefaultConfig()
	flags := flag.NewFlagSet("columnstore", flag.ContinueOnError)
	flags.StringVar(&cfg.Listen, "listen", cfg.Listen, "http listen address")
	flags.StringVar(&cfg.DataDir, "data", cfg.DataDir, "segment data directory")
	flags.DurationVar(&cfg.Retention, "retention", cfg.Retention, "block retention window")
	flags.IntVar(&cfg.Chunk, "chunk", cfg.Chunk, "vectorized scan chunk size")
	err := flags.Parse(args)
	return cfg, err
}
