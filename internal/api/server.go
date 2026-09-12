package api

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/GoreeCloud/goreecloud-swarm/internal/core"
	"github.com/GoreeCloud/goreecloud-swarm/internal/engine"
	"github.com/GoreeCloud/goreecloud-swarm/internal/protocol"
)

type Server struct {
	manager *core.Manager
	logger  *slog.Logger
	mux     *http.ServeMux
}

type healthResponse struct {
	Status    string              `json:"status"`
	Ready     bool                `json:"ready"`
	Service   componentHealth     `json:"service"`
	Engine    engine.Capabilities `json:"engine"`
	Timestamp time.Time           `json:"timestamp"`
}

type componentHealth struct {
	Ready bool   `json:"ready"`
	Name  string `json:"name"`
}

type addTransferRequest struct {
	MagnetURI string `json:"magnet_uri"`
}

func NewServer(manager *core.Manager, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Server{manager: manager, logger: logger, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return withSecurityHeaders(s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/v1/health", s.health)
	s.mux.HandleFunc("GET /api/v1/capabilities", s.capabilities)
	s.mux.HandleFunc("GET /api/v1/transfers", s.listTransfers)
	s.mux.HandleFunc("POST /api/v1/metainfo/inspect", s.inspectMetainfo)
	s.mux.HandleFunc("POST /api/v1/transfers", s.addTransfer)
	s.mux.HandleFunc("POST /api/v1/transfers/{id}/pause", s.pauseTransfer)
	s.mux.HandleFunc("POST /api/v1/transfers/{id}/resume", s.resumeTransfer)
	s.mux.HandleFunc("DELETE /api/v1/transfers/{id}", s.removeTransfer)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	caps := s.manager.EngineCapabilities(r.Context())
	status := "healthy"
	if !caps.Ready {
		status = "degraded"
	}
	writeJSON(w, http.StatusOK, healthResponse{
		Status:    status,
		Ready:     caps.Ready,
		Service:   componentHealth{Ready: true, Name: "swarmd"},
		Engine:    caps,
		Timestamp: time.Now().UTC(),
	})
}

func (s *Server) capabilities(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.manager.EngineCapabilities(r.Context()))
}

func (s *Server) listTransfers(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"transfers": s.manager.List()})
}

func (s *Server) inspectMetainfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const maxMetainfoBytes = 16 << 20
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxMetainfoBytes))
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "metainfo_too_large", "torrent metainfo exceeds the development inspection limit")
			return
		}
		writeError(w, http.StatusBadRequest, "invalid_metainfo", "torrent metainfo could not be read")
		return
	}
	meta, err := protocol.ParseTorrent(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_metainfo", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, meta)
}

func (s *Server) addTransfer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	var req addTransferRequest
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body is invalid")
		return
	}
	if strings.TrimSpace(req.MagnetURI) == "" {
		writeError(w, http.StatusBadRequest, "invalid_request", "magnet_uri is required")
		return
	}

	transfer, err := s.manager.AddMagnet(r.Context(), req.MagnetURI)
	if err != nil {
		if errors.Is(err, engine.ErrUnavailable) {
			writeError(w, http.StatusServiceUnavailable, "engine_unavailable", "transfer engine is not available")
			return
		}
		if errors.Is(err, engine.ErrUnsupported) {
			writeError(w, http.StatusUnprocessableEntity, "capability_unsupported", "transfer engine does not support magnet ingestion")
			return
		}
		writeError(w, http.StatusBadRequest, "invalid_transfer", err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, transfer)
}

func (s *Server) pauseTransfer(w http.ResponseWriter, r *http.Request) {
	transfer, err := s.manager.Pause(r.Context(), r.PathValue("id"))
	if err != nil {
		s.writeLifecycleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, transfer)
}

func (s *Server) resumeTransfer(w http.ResponseWriter, r *http.Request) {
	transfer, err := s.manager.Resume(r.Context(), r.PathValue("id"))
	if err != nil {
		s.writeLifecycleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, transfer)
}

func (s *Server) removeTransfer(w http.ResponseWriter, r *http.Request) {
	deleteData := false
	if raw := r.URL.Query().Get("delete_data"); raw != "" {
		var err error
		deleteData, err = strconv.ParseBool(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "delete_data must be true or false")
			return
		}
	}
	if err := s.manager.Remove(r.Context(), r.PathValue("id"), deleteData); err != nil {
		s.writeLifecycleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) writeLifecycleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrTransferNotFound):
		writeError(w, http.StatusNotFound, "transfer_not_found", "transfer was not found")
	case errors.Is(err, engine.ErrUnavailable):
		writeError(w, http.StatusServiceUnavailable, "engine_unavailable", "transfer engine is not available")
	default:
		writeError(w, http.StatusConflict, "operation_failed", err.Error())
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSONStatus(w, status, map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	writeJSONStatus(w, status, v)
}

func writeJSONStatus(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
