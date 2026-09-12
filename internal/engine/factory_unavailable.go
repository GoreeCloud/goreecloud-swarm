//go:build !anacrolix_engine

package engine

import "fmt"

func newConfigured(cfg Config) (Engine, error) {
	switch cfg.Kind {
	case KindUnavailable:
		return Unavailable{}, nil
	case KindAnacrolix:
		return nil, fmt.Errorf("%w: rebuild with -tags anacrolix_engine", ErrEngineNotBuilt)
	default:
		return nil, unsupportedEngine(cfg.Kind)
	}
}
