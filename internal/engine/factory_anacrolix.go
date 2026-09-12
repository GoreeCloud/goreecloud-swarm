//go:build anacrolix_engine

package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	atorrent "github.com/anacrolix/torrent"
)

const anacrolixVersion = "v1.61.0"

type anacrolixEngine struct {
	client    *atorrent.Client
	mu        sync.RWMutex
	transfers map[string]*atorrent.Torrent
}

func newConfigured(cfg Config) (Engine, error) {
	switch cfg.Kind {
	case KindUnavailable:
		return Unavailable{}, nil
	case KindAnacrolix:
		return newAnacrolix(cfg)
	default:
		return nil, unsupportedEngine(cfg.Kind)
	}
}

func newAnacrolix(cfg Config) (Engine, error) {
	if cfg.DownloadDir == "" {
		return nil, ErrDownloadDirRequired
	}

	root, err := filepath.Abs(cfg.DownloadDir)
	if err != nil {
		return nil, fmt.Errorf("resolve download directory: %w", err)
	}
	if err := os.MkdirAll(root, 0o700); err != nil {
		return nil, fmt.Errorf("create download directory: %w", err)
	}

	clientConfig := atorrent.NewDefaultClientConfig()
	clientConfig.DataDir = root
	client, err := atorrent.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("initialize anacrolix torrent client: %w", err)
	}

	return &anacrolixEngine{
		client:    client,
		transfers: make(map[string]*atorrent.Torrent),
	}, nil
}

func (e *anacrolixEngine) Capabilities(context.Context) Capabilities {
	return Capabilities{
		Ready:   true,
		Name:    "anacrolix/torrent",
		Version: anacrolixVersion,
		Supported: []Capability{
			CapabilityMagnet,
			CapabilityTorrentV1,
			CapabilityTorrentV2,
			CapabilityHybrid,
			CapabilityDHT,
			CapabilityPEX,
			CapabilityEncryption,
		},
		Limitation: "initial adapter: file import, persistence, storage deletion, network lock, proxy policy, and runtime platform integrations are not yet implemented",
	}
}

func (e *anacrolixEngine) Add(ctx context.Context, req AddRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if req.ID == "" {
		return errors.New("engine transfer id is required")
	}
	if req.MagnetURI == "" {
		return errors.New("magnet URI is required")
	}

	e.mu.Lock()
	if _, exists := e.transfers[req.ID]; exists {
		e.mu.Unlock()
		return ErrDuplicateTransferID
	}
	e.transfers[req.ID] = nil
	e.mu.Unlock()

	t, err := e.client.AddMagnet(req.MagnetURI)
	if err != nil {
		e.mu.Lock()
		delete(e.transfers, req.ID)
		e.mu.Unlock()
		return fmt.Errorf("add magnet to anacrolix client: %w", err)
	}
	if err := ctx.Err(); err != nil {
		t.Drop()
		e.mu.Lock()
		delete(e.transfers, req.ID)
		e.mu.Unlock()
		return err
	}

	t.DownloadAll()
	e.mu.Lock()
	e.transfers[req.ID] = t
	e.mu.Unlock()
	return nil
}

func (e *anacrolixEngine) Pause(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	t, err := e.lookup(id)
	if err != nil {
		return err
	}
	t.DisallowDataDownload()
	t.DisallowDataUpload()
	return nil
}

func (e *anacrolixEngine) Resume(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	t, err := e.lookup(id)
	if err != nil {
		return err
	}
	t.AllowDataDownload()
	t.AllowDataUpload()
	t.DownloadAll()
	return nil
}

func (e *anacrolixEngine) Remove(ctx context.Context, id string, deleteData bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if deleteData {
		return ErrDeleteDataUnsupported
	}
	t, err := e.lookup(id)
	if err != nil {
		return err
	}
	t.Drop()
	e.mu.Lock()
	delete(e.transfers, id)
	e.mu.Unlock()
	return nil
}

func (e *anacrolixEngine) Close() error {
	errs := e.client.Close()
	e.mu.Lock()
	clear(e.transfers)
	e.mu.Unlock()
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func (e *anacrolixEngine) lookup(id string) (*atorrent.Torrent, error) {
	e.mu.RLock()
	t, ok := e.transfers[id]
	e.mu.RUnlock()
	if !ok || t == nil {
		return nil, ErrTransferNotFound
	}
	return t, nil
}
