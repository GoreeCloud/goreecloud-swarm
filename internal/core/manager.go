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
