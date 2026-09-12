package storage

import (
	"path/filepath"
	"testing"
)

func TestSafeJoin(t *testing.T) {
	got, err := SafeJoin("/srv/swarm/downloads", "linux/image.iso")
	if err != nil {
		t.Fatalf("SafeJoin() error = %v", err)
	}
	want := filepath.Clean("/srv/swarm/downloads/linux/image.iso")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestSafeJoinRejectsTraversal(t *testing.T) {
	for _, value := range []string{"../secret", "../../etc/passwd", "/etc/passwd", ".."} {
		if _, err := SafeJoin("/srv/swarm/downloads", value); err == nil {
			t.Fatalf("expected %q to be rejected", value)
		}
	}
}
