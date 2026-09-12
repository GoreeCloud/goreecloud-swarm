package core

import "time"

type TransferState string

const (
	StatePending        TransferState = "pending"
	StateDownloading    TransferState = "downloading"
	StateSeeding        TransferState = "seeding"
	StatePaused         TransferState = "paused"
	StateChecking       TransferState = "checking"
	StateQueued         TransferState = "queued"
	StateStalled        TransferState = "stalled"
	StateStorageMissing TransferState = "storage_missing"
	StateNetworkBlocked TransferState = "network_blocked"
	StateError          TransferState = "error"
	StateCompleted      TransferState = "completed"
)

type Transfer struct {
	ID        string        `json:"id"`
	Name      string        `json:"name,omitempty"`
	State     TransferState `json:"state"`
	AddedAt   time.Time     `json:"added_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	LastError string        `json:"last_error,omitempty"`
}
