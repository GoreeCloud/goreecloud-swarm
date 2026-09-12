package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoreeCloud/goreecloud-swarm/internal/core"
	"github.com/GoreeCloud/goreecloud-swarm/internal/engine"
)

func testServer() *Server {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewServer(core.NewManager(engine.Unavailable{}), logger)
}

type acceptingAPIEngine struct{}

func (acceptingAPIEngine) Capabilities(context.Context) engine.Capabilities {
	return engine.Capabilities{Ready: true, Name: "test", Version: "1", Supported: []engine.Capability{engine.CapabilityMagnet, engine.CapabilityTorrentV1}}
}
func (acceptingAPIEngine) Add(context.Context, engine.AddRequest) error               { return nil }
func (acceptingAPIEngine) AddTorrent(context.Context, engine.AddTorrentRequest) error { return nil }
func (acceptingAPIEngine) Pause(context.Context, string) error                        { return nil }
func (acceptingAPIEngine) Resume(context.Context, string) error                       { return nil }
func (acceptingAPIEngine) Remove(context.Context, string, bool) error                 { return nil }
func (acceptingAPIEngine) Close() error                                               { return nil }

func acceptingTestServer() *Server {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewServer(core.NewManager(acceptingAPIEngine{}), logger)
}

func TestHealthIsTruthfullyDegradedWithoutEngine(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	testServer().Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "degraded" || body["ready"] != false {
		t.Fatalf("unexpected health response: %s", rec.Body.String())
	}
}

func TestAddTransferReportsEngineUnavailable(t *testing.T) {
	body := `{"magnet_uri":"magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transfers", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testServer().Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestSecurityHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	testServer().Handler().ServeHTTP(rec, req)
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q", got)
	}
}

func TestInspectMetainfoV1(t *testing.T) {
	info := "d6:lengthi4e4:name4:test12:piece lengthi16384e6:pieces20:aaaaaaaaaaaaaaaaaaaae"
	body := "d4:info" + info + "e"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/metainfo/inspect", strings.NewReader(body))
	rec := httptest.NewRecorder()
	testServer().Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var meta map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &meta); err != nil {
		t.Fatal(err)
	}
	if meta["version"] != "v1" || meta["name"] != "test" {
		t.Fatalf("unexpected metainfo: %s", rec.Body.String())
	}
}

func TestInspectMetainfoRejectsInvalidData(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/metainfo/inspect", strings.NewReader("not-bencode"))
	rec := httptest.NewRecorder()
	testServer().Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAddTorrentTransfer(t *testing.T) {
	srv := acceptingTestServer()
	info := "d6:lengthi4e4:name4:test12:piece lengthi16384e6:pieces20:aaaaaaaaaaaaaaaaaaaae"
	body := "d4:info" + info + "e"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transfers/torrent", strings.NewReader(body))
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
	var transfer core.Transfer
	if err := json.Unmarshal(rec.Body.Bytes(), &transfer); err != nil {
		t.Fatal(err)
	}
	if transfer.Name != "test" {
		t.Fatalf("name = %q", transfer.Name)
	}
}

func TestTransferLifecycleEndpoints(t *testing.T) {
	srv := acceptingTestServer()
	addBody := `{"magnet_uri":"magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567&dn=Example"}`
	addReq := httptest.NewRequest(http.MethodPost, "/api/v1/transfers", strings.NewReader(addBody))
	addReq.Header.Set("Content-Type", "application/json")
	addRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(addRec, addReq)
	if addRec.Code != http.StatusAccepted {
		t.Fatalf("add status = %d body=%s", addRec.Code, addRec.Body.String())
	}
	var transfer core.Transfer
	if err := json.Unmarshal(addRec.Body.Bytes(), &transfer); err != nil {
		t.Fatal(err)
	}

	pauseReq := httptest.NewRequest(http.MethodPost, "/api/v1/transfers/"+transfer.ID+"/pause", nil)
	pauseRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(pauseRec, pauseReq)
	if pauseRec.Code != http.StatusOK {
		t.Fatalf("pause status = %d body=%s", pauseRec.Code, pauseRec.Body.String())
	}
	var paused core.Transfer
	if err := json.Unmarshal(pauseRec.Body.Bytes(), &paused); err != nil {
		t.Fatal(err)
	}
	if paused.State != core.StatePaused {
		t.Fatalf("pause state = %q", paused.State)
	}

	resumeReq := httptest.NewRequest(http.MethodPost, "/api/v1/transfers/"+transfer.ID+"/resume", nil)
	resumeRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(resumeRec, resumeReq)
	if resumeRec.Code != http.StatusOK {
		t.Fatalf("resume status = %d body=%s", resumeRec.Code, resumeRec.Body.String())
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/transfers/"+transfer.ID+"?delete_data=false", nil)
	deleteRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(deleteRec, deleteReq)
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d body=%s", deleteRec.Code, deleteRec.Body.String())
	}

	missingReq := httptest.NewRequest(http.MethodPost, "/api/v1/transfers/"+transfer.ID+"/pause", nil)
	missingRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d body=%s", missingRec.Code, missingRec.Body.String())
	}
}

func TestRemoveRejectsInvalidDeleteData(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/transfers/example?delete_data=maybe", nil)
	rec := httptest.NewRecorder()
	acceptingTestServer().Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", rec.Code, rec.Body.String())
	}
}
