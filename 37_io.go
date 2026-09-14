//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// upperReader wraps an io.Reader, upper-casing bytes.
type upperReader struct {
	inner io.Reader
}

func (reader upperReader) Read(buffer []byte) (int, error) {
	n, err := reader.inner.Read(buffer)
	// ASCII-safe; UTF-8 needs rune-aware decoding.
	for i := 0; i < n; i++ {
		buffer[i] = byte(unicode.ToUpper(rune(buffer[i])))
	}
	return n, err
}

func main() {
	// strings.NewReader reads; bytes.Buffer writes.
	source := upperReader{
		inner: strings.NewReader("hello, io"),
	}
	var sink bytes.Buffer

	// io.Copy streams reader to writer until EOF.
	n, err := io.Copy(&sink, source)
	if err != nil {
		panic(err)
	}

	fmt.Printf("copied %d bytes: %s\n", n, sink.String())
}
