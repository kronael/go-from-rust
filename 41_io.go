//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// upperReader wraps any io.Reader and upper-cases bytes as they pass by.
// Implementing one method, Read, is enough to be an io.Reader, the same
// shape as Rust's Read trait; you compose by wrapping, like Read adapters.
type upperReader struct {
	inner io.Reader
}

func (reader upperReader) Read(buffer []byte) (int, error) {
	n, err := reader.inner.Read(buffer)
	// Input here is ASCII, so upper-casing each byte is safe. Real UTF-8
	// needs rune-aware decoding — per-byte upcasing mangles multi-byte runes.
	for i := 0; i < n; i++ {
		buffer[i] = byte(unicode.ToUpper(rune(buffer[i])))
	}
	return n, err
}

func main() {
	// strings.NewReader is an in-memory reader; bytes.Buffer is a writer.
	source := upperReader{inner: strings.NewReader("hello, io")}
	var sink bytes.Buffer

	// io.Copy pumps reader to writer until EOF, like std::io::copy. It streams
	// in chunks; the wrapper transforms each before it lands in the buffer.
	n, err := io.Copy(&sink, source)
	if err != nil {
		panic(err)
	}

	fmt.Printf("copied %d bytes: %s\n", n, sink.String())
}
