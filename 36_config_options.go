//go:build ignore

package main

import "fmt"

// Rust reaches for a builder plus #[derive(Default)]. Go's version: a plain
// struct whose defaults live in New, plus functional options — closures
// that mutate a *Config, applied over those defaults.
type Config struct {
	Host    string
	Port    int
	Verbose bool
}

// Option mutates a Config in place, like one ConfigBuilder::with_x call.
type Option func(*Config)

func WithHost(host string) Option {
	return func(config *Config) { config.Host = host }
}

func WithPort(port int) Option {
	return func(config *Config) { config.Port = port }
}

func WithVerbose() Option {
	return func(config *Config) { config.Verbose = true }
}

// New sets explicit defaults, then applies each option in order. Verbose's
// zero value (false) is already the default, so nothing requests it: in Go
// the zero value often IS the default.
func New(opts ...Option) Config {
	config := Config{
		Host: "localhost",
		Port: 8080,
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}

func main() {
	defaults := New()
	custom := New(WithHost("example.com"), WithPort(443), WithVerbose())

	fmt.Printf("defaults: %+v\n", defaults)
	fmt.Printf("custom:   %+v\n", custom)
}
