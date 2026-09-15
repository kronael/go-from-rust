//go:build ignore

package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var errPortRange = errors.New("port outside 1..65535")

func parsePort(input string) (int, error) {
	input = strings.TrimSpace(input)
	port, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("parse port: %w", err)
	}
	if port < 1 || port > 65_535 {
		return 0, errPortRange
	}
	return port, nil
}

func TestParsePort(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "usual", input: "8080", want: 8080},
		{name: "trims space", input: " 443\n", want: 443},
		{name: "maximum", input: "65535", want: 65_535},
		{name: "zero", input: "0", wantErr: true},
		{name: "too big", input: "65536", wantErr: true},
		{name: "not numeric", input: "http", wantErr: true},
	}

	// Rust's assert_eq! panics, so one mismatch ends the
	// test. Go splits that into two calls: Fatalf stops
	// this subtest, Errorf records and runs on. Choosing
	// between them is how a Go test reports every real
	// mismatch in one run instead of only the first.
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parsePort(c.input)

			// Wrong error state: comparing the value now
			// would only add noise, so stop here.
			if (err != nil) != c.wantErr {
				t.Fatalf("error = %v; want error = %v",
					err, c.wantErr)
			}
			if c.wantErr {
				return
			}

			// An independent check, so record and continue.
			if got != c.want {
				t.Errorf("parsePort(%q) = %d; want %d",
					c.input, got, c.want)
			}
		})
	}
}

func main() {
	// Tour adapter; real projects use *_test.go.
	testing.Main(
		regexp.MatchString,
		[]testing.InternalTest{
			{
				Name: "TestParsePort",
				F:    TestParsePort,
			},
		},
		nil,
		nil,
	)
}
