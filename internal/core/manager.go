package core

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/GoreeCloud/goreecloud-swarm/internal/engine"
	"github.com/GoreeCloud/goreecloud-swarm/internal/protocol"
)

var ErrTransferNotFound = errors.New("transfer not found")

type Manager struct {
	engine    engine.Engine
	mu        sync.RWMutex
	transfers map[string]Transfer
}

func NewManager(e engine.Engine) *Manager {
	return &Manager{engine: e, transfers: make(map[string]Transfer)}
}

func (m *Manager) EngineCapabilities(ctx context.Context) engine.Capabilities {
	return m.engine.Capabilities(ctx)
}

func (m *Manager) AddMagnet(ctx context.Context, raw string) (Transfer, error) {
	caps := m.engine.Capabilities(ctx)
	if !caps.Ready {
		return Transfer{}, engine.ErrUnavailable
	}
	if !caps.Supports(engine.CapabilityMagnet) {
		return Transfer{}, engine.ErrUnsupported
	}

	magnet, err := protocol.ParseMagnet(raw)
	if err != nil {
		return Transfer{}, err
	}

	id, err := randomID()
	if err != nil {
		return Transfer{}, err
	}

	now := time.Now().UTC()
	transfer := Transfer{
		ID:        id,
		Name:      magnet.DisplayName,
		State:     StatePending,
		AddedAt:   now,
		UpdatedAt: now,
	}

	if err := m.engine.Add(ctx, engine.AddRequest{ID: id, MagnetURI: magnet.Raw}); err != nil {
		transfer.State = StateError
		transfer.LastError = err.Error()
		transfer.UpdatedAt = time.Now().UTC()
		return transfer, err
	}

	m.mu.Lock()
	m.transfers[id] = transfer
	m.mu.Unlock()
	return transfer, nil
}

func (m *Manager) Get(id string) (Transfer, error) {
	m.mu.RLock()
	transfer, ok := m.transfers[id]
	m.mu.RUnlock()
	if !ok {
		return Transfer{}, ErrTransferNotFound
	}
	return transfer, nil
}

func (m *Manager) Pause(ctx context.Context, id string) (Transfer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	transfer, ok := m.transfers[id]
	if !ok {
		return Transfer{}, ErrTransferNotFound
	}
	if err := m.engine.Pause(ctx, id); err != nil {
		return transfer, err
	}
	transfer.State = StatePaused
	transfer.LastError = ""
	transfer.UpdatedAt = time.Now().UTC()
	m.transfers[id] = transfer
	return transfer, nil
}

func (m *Manager) Resume(ctx context.Context, id string) (Transfer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	transfer, ok := m.transfers[id]
	if !ok {
		return Transfer{}, ErrTransferNotFound
	}
	if err := m.engine.Resume(ctx, id); err != nil {
		return transfer, err
	}
	transfer.State = StatePending
	transfer.LastError = ""
	transfer.UpdatedAt = time.Now().UTC()
	m.transfers[id] = transfer
	return transfer, nil
}

func (m *Manager) Remove(ctx context.Context, id string, deleteData bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.transfers[id]; !ok {
		return ErrTransferNotFound
	}
	if err := m.engine.Remove(ctx, id, deleteData); err != nil {
		return err
	}
	delete(m.transfers, id)
	return nil
}

func (m *Manager) List() []Transfer {
	m.mu.RLock()
	out := make([]Transfer, 0, len(m.transfers))
	for _, transfer := range m.transfers {
		out = append(out, transfer)
	}
	m.mu.RUnlock()

	sort.Slice(out, func(i, j int) bool { return out[i].AddedAt.Before(out[j].AddedAt) })
	return out
}

func randomID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
