package main

import "testing"

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
