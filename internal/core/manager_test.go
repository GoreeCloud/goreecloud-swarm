package core

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-swarm/internal/engine"
)

type acceptingEngine struct {
	pauseErr  error
	resumeErr error
	removeErr error
}

func (acceptingEngine) Capabilities(context.Context) engine.Capabilities {
	return engine.Capabilities{Ready: true, Name: "test", Version: "1", Supported: []engine.Capability{engine.CapabilityMagnet}}
}
func (acceptingEngine) Add(context.Context, engine.AddRequest) error               { return nil }
func (acceptingEngine) AddTorrent(context.Context, engine.AddTorrentRequest) error { return nil }
func (e acceptingEngine) Pause(context.Context, string) error                      { return e.pauseErr }
func (e acceptingEngine) Resume(context.Context, string) error                     { return e.resumeErr }
func (e acceptingEngine) Remove(context.Context, string, bool) error {
	return e.removeErr
}
func (acceptingEngine) Close() error { return nil }

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

func TestPauseResumeRemoveTransitionsAfterEngineAcceptance(t *testing.T) {
	m := NewManager(acceptingEngine{})
	transfer, err := m.AddMagnet(context.Background(), "magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567&dn=Example")
	if err != nil {
		t.Fatal(err)
	}

	paused, err := m.Pause(context.Background(), transfer.ID)
	if err != nil {
		t.Fatal(err)
	}
	if paused.State != StatePaused {
		t.Fatalf("pause state = %q", paused.State)
	}

	resumed, err := m.Resume(context.Background(), transfer.ID)
	if err != nil {
		t.Fatal(err)
	}
	if resumed.State != StatePending {
		t.Fatalf("resume state = %q", resumed.State)
	}

	if err := m.Remove(context.Background(), transfer.ID, false); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Get(transfer.ID); !errors.Is(err, ErrTransferNotFound) {
		t.Fatalf("Get after remove error = %v", err)
	}
}

func TestRejectedPauseDoesNotChangeAuthoritativeState(t *testing.T) {
	pauseErr := errors.New("pause rejected")
	m := NewManager(acceptingEngine{pauseErr: pauseErr})
	transfer, err := m.AddMagnet(context.Background(), "magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567&dn=Example")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.Pause(context.Background(), transfer.ID); !errors.Is(err, pauseErr) {
		t.Fatalf("pause error = %v", err)
	}
	current, err := m.Get(transfer.ID)
	if err != nil {
		t.Fatal(err)
	}
	if current.State != StatePending {
		t.Fatalf("state changed after rejected pause: %q", current.State)
	}
}

func TestAddTorrentUsesParsedVersionCapability(t *testing.T) {
	m := NewManager(acceptingTorrentEngine{})
	info := "d6:lengthi4e4:name4:test12:piece lengthi16384e6:pieces20:aaaaaaaaaaaaaaaaaaaae"
	data := []byte("d4:info" + info + "e")
	transfer, err := m.AddTorrent(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if transfer.Name != "test" {
		t.Fatalf("name = %q", transfer.Name)
	}
}

type acceptingTorrentEngine struct{ acceptingEngine }

func (acceptingTorrentEngine) Capabilities(context.Context) engine.Capabilities {
	return engine.Capabilities{Ready: true, Name: "test", Version: "1", Supported: []engine.Capability{engine.CapabilityTorrentV1}}
}

func TestLifecycleMissingTransfer(t *testing.T) {
	m := NewManager(acceptingEngine{})
	if _, err := m.Pause(context.Background(), "missing"); !errors.Is(err, ErrTransferNotFound) {
		t.Fatalf("pause missing error = %v", err)
	}
	if _, err := m.Resume(context.Background(), "missing"); !errors.Is(err, ErrTransferNotFound) {
		t.Fatalf("resume missing error = %v", err)
	}
	if err := m.Remove(context.Background(), "missing", false); !errors.Is(err, ErrTransferNotFound) {
		t.Fatalf("remove missing error = %v", err)
	}
}
