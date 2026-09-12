//go:build !anacrolix_engine

package engine

import (
	"errors"
	"testing"
)

func TestNewConfiguredDefaultsToUnavailable(t *testing.T) {
	got, err := NewConfigured(Config{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got.(Unavailable); !ok {
		t.Fatalf("engine type = %T, want Unavailable", got)
	}
}

func TestNewConfiguredRejectsEngineMissingFromBuild(t *testing.T) {
	_, err := NewConfigured(Config{Kind: KindAnacrolix, DownloadDir: "/tmp/swarm"})
	if !errors.Is(err, ErrEngineNotBuilt) {
		t.Fatalf("error = %v, want ErrEngineNotBuilt", err)
	}
}

func TestNewConfiguredRejectsUnknownEngine(t *testing.T) {
	_, err := NewConfigured(Config{Kind: "unknown"})
	if !errors.Is(err, ErrUnsupportedEngine) {
		t.Fatalf("error = %v, want ErrUnsupportedEngine", err)
	}
}
