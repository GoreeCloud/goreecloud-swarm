package api

import (
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
