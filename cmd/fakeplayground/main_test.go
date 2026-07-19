package main

import (
	"net"
	"testing"
)

func TestListenLoopbackAllowsEphemeralPort(t *testing.T) {
	listener, err := listenLoopback("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	if !addr.IP.IsLoopback() || addr.Port == 0 {
		t.Fatalf("listener address = %s, want loopback with assigned port", addr)
	}
}

func TestListenLoopbackRefusesWildcard(t *testing.T) {
	listener, err := listenLoopback("0.0.0.0:0")
	if err == nil {
		listener.Close()
		t.Fatal("listenLoopback succeeded for a non-loopback address")
	}
}
