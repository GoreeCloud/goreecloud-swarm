package storage

import (
	"errors"
	"path/filepath"
	"strings"
)

var ErrUnsafePath = errors.New("unsafe torrent path")

// SafeJoin resolves an untrusted torrent-relative path inside root.
// It rejects absolute paths, parent traversal, empty roots, and any resolved
// result that escapes the authorized storage boundary.
func SafeJoin(root, untrusted string) (string, error) {
	root = strings.TrimSpace(root)
	untrusted = strings.TrimSpace(untrusted)
	if root == "" || untrusted == "" {
		return "", ErrUnsafePath
	}
	if filepath.IsAbs(untrusted) || filepath.VolumeName(untrusted) != "" {
		return "", ErrUnsafePath
	}

	clean := filepath.Clean(untrusted)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", ErrUnsafePath
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", ErrUnsafePath
	}
	candidate := filepath.Join(absRoot, clean)
	rel, err := filepath.Rel(absRoot, candidate)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", ErrUnsafePath
	}
	return candidate, nil
}
