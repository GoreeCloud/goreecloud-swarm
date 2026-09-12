package engine

import (
	"context"
	"errors"
)

var ErrUnavailable = errors.New("transfer engine unavailable")

type Capability string

const (
	CapabilityMagnet     Capability = "magnet"
	CapabilityTorrentV1  Capability = "torrent_v1"
	CapabilityTorrentV2  Capability = "torrent_v2"
	CapabilityHybrid     Capability = "torrent_hybrid"
	CapabilityDHT        Capability = "dht"
	CapabilityPEX        Capability = "pex"
	CapabilityLPD        Capability = "lpd"
	CapabilityWebSeeds   Capability = "web_seeds"
	CapabilityIPv6       Capability = "ipv6"
	CapabilityEncryption Capability = "peer_encryption"
)

type Capabilities struct {
	Ready      bool         `json:"ready"`
	Name       string       `json:"name"`
	Version    string       `json:"version,omitempty"`
	Supported  []Capability `json:"supported"`
	Limitation string       `json:"limitation,omitempty"`
}

type AddRequest struct {
	ID        string
	MagnetURI string
}

type Engine interface {
	Capabilities(ctx context.Context) Capabilities
	Add(ctx context.Context, req AddRequest) error
	Pause(ctx context.Context, id string) error
	Resume(ctx context.Context, id string) error
	Remove(ctx context.Context, id string, deleteData bool) error
	Close() error
}

type Unavailable struct{}

func (Unavailable) Capabilities(context.Context) Capabilities {
	return Capabilities{
		Ready:      false,
		Name:       "unavailable",
		Supported:  []Capability{},
		Limitation: "no BitTorrent engine adapter is connected",
	}
}

func (Unavailable) Add(context.Context, AddRequest) error      { return ErrUnavailable }
func (Unavailable) Pause(context.Context, string) error        { return ErrUnavailable }
func (Unavailable) Resume(context.Context, string) error       { return ErrUnavailable }
func (Unavailable) Remove(context.Context, string, bool) error { return ErrUnavailable }
func (Unavailable) Close() error                               { return nil }
