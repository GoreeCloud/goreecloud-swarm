package engine

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const (
	SidecarProtocolVersion = 1
	maxSidecarMessageBytes = 24 << 20
)

var (
	ErrProtocol = errors.New("transfer engine protocol error")
)

type sidecarRequest struct {
	ID     uint64 `json:"id"`
	Method string `json:"method"`
	Params any    `json:"params,omitempty"`
}

type sidecarResponse struct {
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *sidecarError   `json:"error,omitempty"`
}

type sidecarError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type helloParams struct {
	ProtocolVersion int    `json:"protocol_version"`
	Client          string `json:"client"`
}

type helloResult struct {
	ProtocolVersion int          `json:"protocol_version"`
	Capabilities    Capabilities `json:"capabilities"`
}

type Sidecar struct {
	mu           sync.Mutex
	cmd          *exec.Cmd
	stdin        io.WriteCloser
	scanner      *bufio.Scanner
	nextID       uint64
	capabilities Capabilities
	closed       bool
}

func StartSidecar(ctx context.Context, executable string, stderr io.Writer) (*Sidecar, error) {
	if executable == "" || !filepath.IsAbs(executable) {
		return nil, fmt.Errorf("%w: engine executable must be an absolute path", ErrProtocol)
	}

	cmd := exec.Command(executable)
	cmd.Env = []string{fmt.Sprintf("SWARM_ENGINE_PROTOCOL_VERSION=%d", SidecarProtocolVersion)}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("create engine stdout pipe: %w", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("create engine stdin pipe: %w", err)
	}
	if stderr != nil {
		cmd.Stderr = stderr
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start engine sidecar: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 4096), maxSidecarMessageBytes)
	s := &Sidecar{cmd: cmd, stdin: stdin, scanner: scanner}

	handshakeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var hello helloResult
	if err := s.call(handshakeCtx, "hello", helloParams{ProtocolVersion: SidecarProtocolVersion, Client: "swarmd"}, &hello); err != nil {
		_ = s.killAndWait()
		return nil, fmt.Errorf("engine handshake: %w", err)
	}
	if hello.ProtocolVersion != SidecarProtocolVersion {
		_ = s.killAndWait()
		return nil, fmt.Errorf("%w: unsupported engine protocol version %d", ErrProtocol, hello.ProtocolVersion)
	}
	if !hello.Capabilities.Ready {
		_ = s.killAndWait()
		return nil, fmt.Errorf("%w: engine reported not ready", ErrUnavailable)
	}
	s.capabilities = hello.Capabilities
	return s, nil
}

func (s *Sidecar) Capabilities(context.Context) Capabilities {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.capabilities
}

func (s *Sidecar) Add(ctx context.Context, req AddRequest) error {
	return s.call(ctx, "add", req, nil)
}

func (s *Sidecar) AddTorrent(ctx context.Context, req AddTorrentRequest) error {
	return s.call(ctx, "add_torrent", req, nil)
}

func (s *Sidecar) Pause(ctx context.Context, id string) error {
	return s.call(ctx, "pause", map[string]string{"id": id}, nil)
}

func (s *Sidecar) Resume(ctx context.Context, id string) error {
	return s.call(ctx, "resume", map[string]string{"id": id}, nil)
}

func (s *Sidecar) Remove(ctx context.Context, id string, deleteData bool) error {
	return s.call(ctx, "remove", map[string]any{"id": id, "delete_data": deleteData}, nil)
}

func (s *Sidecar) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	s.mu.Unlock()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = s.call(shutdownCtx, "shutdown", nil, nil)
	_ = s.stdin.Close()

	done := make(chan error, 1)
	go func() { done <- s.cmd.Wait() }()
	select {
	case err := <-done:
		return err
	case <-shutdownCtx.Done():
		_ = s.cmd.Process.Kill()
		return <-done
	}
}

func (s *Sidecar) call(ctx context.Context, method string, params any, out any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed && method != "shutdown" {
		return ErrUnavailable
	}
	s.nextID++
	id := s.nextID
	payload, err := json.Marshal(sidecarRequest{ID: id, Method: method, Params: params})
	if err != nil {
		return fmt.Errorf("encode engine request: %w", err)
	}
	if len(payload) > maxSidecarMessageBytes {
		return fmt.Errorf("%w: request too large", ErrProtocol)
	}

	if _, err := s.stdin.Write(append(payload, '\n')); err != nil {
		return fmt.Errorf("write engine request: %w", err)
	}

	scanDone := make(chan bool, 1)
	go func() { scanDone <- s.scanner.Scan() }()
	select {
	case ok := <-scanDone:
		if !ok {
			if err := s.scanner.Err(); err != nil {
				return fmt.Errorf("read engine response: %w", err)
			}
			return fmt.Errorf("%w: engine closed response stream", ErrProtocol)
		}
	case <-ctx.Done():
		_ = s.cmd.Process.Kill()
		return ctx.Err()
	}

	var resp sidecarResponse
	if err := json.Unmarshal(s.scanner.Bytes(), &resp); err != nil {
		return fmt.Errorf("%w: invalid JSON response: %v", ErrProtocol, err)
	}
	if resp.ID != id {
		return fmt.Errorf("%w: response id %d does not match request id %d", ErrProtocol, resp.ID, id)
	}
	if resp.Error != nil {
		if resp.Error.Code == "unavailable" {
			return fmt.Errorf("%w: %s", ErrUnavailable, resp.Error.Message)
		}
		if resp.Error.Code == "unsupported" {
			return fmt.Errorf("%w: %s", ErrUnsupported, resp.Error.Message)
		}
		return fmt.Errorf("engine %s: %s", resp.Error.Code, resp.Error.Message)
	}
	if out != nil {
		if len(resp.Result) == 0 {
			return fmt.Errorf("%w: missing result for %s", ErrProtocol, method)
		}
		if err := json.Unmarshal(resp.Result, out); err != nil {
			return fmt.Errorf("%w: decode %s result: %v", ErrProtocol, method, err)
		}
	}
	return nil
}

func (s *Sidecar) killAndWait() error {
	if s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	return s.cmd.Wait()
}
