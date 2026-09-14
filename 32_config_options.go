//go:build ignore

package main

import "fmt"

// Functional options: closures over New's defaults.
type Config struct {
	Host    string
	Port    int
	Verbose bool
}

// Option mutates a Config in place.
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

// New seeds defaults, then applies each option in order.
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
	custom := New(
		WithHost("example.com"),
		WithPort(443),
		WithVerbose())

	fmt.Printf("defaults: %+v\n", defaults)
	fmt.Printf("custom:   %+v\n", custom)
}
