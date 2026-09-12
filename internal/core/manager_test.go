package core

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-swarm/internal/engine"
)

func TestUnavailableEngineDoesNotPersistTransfer(t *testing.T) {
	m := NewManager(engine.Unavailable{})
	_, err := m.AddMagnet(context.Background(), "magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567&dn=Example")
	if !errors.Is(err, engine.ErrUnavailable) {
		t.Fatalf("error = %v, want ErrUnavailable", err)
	}
	if got := len(m.List()); got != 0 {
		t.Fatalf("persisted %d transfers despite engine rejection", got)
	}
}
