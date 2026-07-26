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
		{
			name:  "usual",
			input: "8080",
			want:  8080,
		},
		{
			name:  "trims space",
			input: " 443\n",
			want:  443,
		},
		{
			name:  "maximum",
			input: "65535",
			want:  65_535,
		},
		{
			name:    "zero",
			input:   "0",
			wantErr: true,
		},
		{
			name:    "above maximum",
			input:   "65536",
			wantErr: true,
		},
		{
			name:    "not a number",
			input:   "http",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parsePort(c.input)
			if (err != nil) != c.wantErr {
				t.Errorf(
					"error = %v; want error = %v",
					err,
					c.wantErr,
				)
			}
			if got != c.want {
				t.Errorf(
					"parsePort(%q) = %d; want %d",
					c.input,
					got,
					c.want,
				)
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
