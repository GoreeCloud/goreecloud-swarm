package main

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-swarm/internal/engine"
)

func TestIsLoopbackAddress(t *testing.T) {
	for _, addr := range []string{"127.0.0.1:8080", "[::1]:8080", "localhost:8080"} {
		if !isLoopbackAddress(addr) {
			t.Fatalf("expected %q to be loopback", addr)
		}
	}
	for _, addr := range []string{"0.0.0.0:8080", "192.0.2.10:8080", "bad-address"} {
		if isLoopbackAddress(addr) {
			t.Fatalf("expected %q to be rejected", addr)
		}
	}
}

func TestOpenTransferEngineDefaultsUnavailable(t *testing.T) {
	e, err := openTransferEngine(context.Background(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	if caps := e.Capabilities(context.Background()); caps.Ready || caps.Name != "unavailable" {
		t.Fatalf("unexpected capabilities: %+v", caps)
	}
}

func TestOpenTransferEngineRejectsRelativeExecutable(t *testing.T) {
	_, err := openTransferEngine(context.Background(), "swarm-engine", "")
	if !errors.Is(err, engine.ErrProtocol) {
		t.Fatalf("expected ErrProtocol, got %v", err)
	}
}
