package engine

import (
	"errors"
	"fmt"
	"strings"
)

const (
	KindUnavailable = "unavailable"
	KindAnacrolix   = "anacrolix"
)

var (
	ErrEngineNotBuilt        = errors.New("requested transfer engine is not included in this build")
	ErrUnsupportedEngine     = errors.New("unsupported transfer engine")
	ErrDownloadDirRequired   = errors.New("download directory is required for the selected transfer engine")
	ErrTransferNotFound      = errors.New("engine transfer not found")
	ErrDuplicateTransferID   = errors.New("engine transfer id already exists")
	ErrDeleteDataUnsupported = errors.New("engine data deletion is not implemented")
)

type Config struct {
	Kind        string
	DownloadDir string
}

func NewConfigured(cfg Config) (Engine, error) {
	cfg.Kind = strings.ToLower(strings.TrimSpace(cfg.Kind))
	cfg.DownloadDir = strings.TrimSpace(cfg.DownloadDir)
	if cfg.Kind == "" {
		cfg.Kind = KindUnavailable
	}
	return newConfigured(cfg)
}

func unsupportedEngine(kind string) error {
	return fmt.Errorf("%w: %s", ErrUnsupportedEngine, kind)
}
