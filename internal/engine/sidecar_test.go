package engine

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStartSidecarRequiresAbsolutePath(t *testing.T) {
	_, err := StartSidecar(context.Background(), "swarm-engine", "", nil)
	if !errorsIs(err, ErrProtocol) {
		t.Fatalf("expected ErrProtocol, got %v", err)
	}
}

func TestStartSidecarRequiresAbsoluteDownloadRoot(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	helper := filepath.Join(t.TempDir(), "engine-helper")
	if err := os.Symlink(exe, helper); err != nil {
		t.Skipf("cannot create helper symlink: %v", err)
	}
	_, err = StartSidecar(context.Background(), helper, "relative", nil)
	if !errorsIs(err, ErrProtocol) {
		t.Fatalf("expected ErrProtocol, got %v", err)
	}
}

func TestSidecarHandshakeAndCalls(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	helper := filepath.Join(t.TempDir(), "engine-helper")
	if err := os.Symlink(exe, helper); err != nil {
		t.Skipf("cannot create helper symlink: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sidecar, err := StartSidecar(ctx, helper, t.TempDir(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer sidecar.Close()

	caps := sidecar.Capabilities(ctx)
	if !caps.Ready || caps.Name != "test-engine" || caps.Version != "1.0" {
		t.Fatalf("unexpected capabilities: %+v", caps)
	}
	if err := sidecar.Add(ctx, AddRequest{ID: "abc", MagnetURI: "magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567"}); err != nil {
		t.Fatal(err)
	}
	if err := sidecar.AddTorrent(ctx, AddTorrentRequest{ID: "torrent", Metainfo: []byte("d4:infodee")}); err != nil {
		t.Fatal(err)
	}
	if err := sidecar.Pause(ctx, "abc"); err != nil {
		t.Fatal(err)
	}
	if err := sidecar.Resume(ctx, "abc"); err != nil {
		t.Fatal(err)
	}
	if err := sidecar.Remove(ctx, "abc", false); err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	if filepath.Base(os.Args[0]) == "engine-helper" {
		runEngineHelper()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func runEngineHelper() {
	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	for scanner.Scan() {
		var req sidecarRequest
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			return
		}
		resp := sidecarResponse{ID: req.ID, Result: json.RawMessage(`{}`)}
		switch req.Method {
		case "hello":
			result := helloResult{
				ProtocolVersion: SidecarProtocolVersion,
				Capabilities: Capabilities{
					Ready:     true,
					Name:      "test-engine",
					Version:   "1.0",
					Supported: []Capability{CapabilityMagnet, CapabilityTorrentV1, CapabilityTorrentV2, CapabilityHybrid},
				},
			}
			raw, _ := json.Marshal(result)
			resp.Result = raw
		case "add", "add_torrent", "pause", "resume", "remove":
			resp.Result = json.RawMessage(`{}`)
		case "shutdown":
			_ = encoder.Encode(resp)
			return
		default:
			resp.Error = &sidecarError{Code: "unsupported", Message: fmt.Sprintf("unknown method %s", req.Method)}
			resp.Result = nil
		}
		if err := encoder.Encode(resp); err != nil {
			return
		}
	}
}

func errorsIs(err, target error) bool {
	for err != nil {
		if err == target {
			return true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := err.(unwrapper)
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}
