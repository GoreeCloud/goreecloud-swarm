//go:build anacrolix_engine

package engine

import (
	"errors"
	"testing"
)

func TestAnacrolixRequiresDownloadDirectory(t *testing.T) {
	_, err := NewConfigured(Config{Kind: KindAnacrolix})
	if !errors.Is(err, ErrDownloadDirRequired) {
		t.Fatalf("error = %v, want ErrDownloadDirRequired", err)
	}
}
