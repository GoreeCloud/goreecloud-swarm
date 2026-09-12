package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoreeCloud/goreecloud-swarm/internal/api"
	"github.com/GoreeCloud/goreecloud-swarm/internal/core"
	"github.com/GoreeCloud/goreecloud-swarm/internal/engine"
)

const defaultListenAddr = "127.0.0.1:8080"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	addr := os.Getenv("SWARM_LISTEN_ADDR")
	if addr == "" {
		addr = defaultListenAddr
	}

	if !isLoopbackAddress(addr) {
		logger.Error("refusing non-loopback bind before remote authentication is implemented", "listen_addr", addr)
		os.Exit(2)
	}

	transferEngine, err := openTransferEngine(context.Background(), os.Getenv("SWARM_ENGINE_PATH"))
	if err != nil {
		logger.Error("Swarm engine startup failed", "error", err)
		os.Exit(2)
	}
	manager := core.NewManager(transferEngine)
	server := api.NewServer(manager, logger)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		caps := transferEngine.Capabilities(context.Background())
		logger.Info("Swarm service starting", "listen_addr", addr, "engine", caps.Name, "engine_version", caps.Version, "engine_ready", caps.Ready, "status", "development")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Swarm service failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Swarm service shutdown failed", "error", err)
	}
	if err := transferEngine.Close(); err != nil {
		logger.Error("Swarm engine shutdown failed", "error", err)
	}
}

func isLoopbackAddress(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func openTransferEngine(ctx context.Context, executable string) (engine.Engine, error) {
	if executable == "" {
		return engine.Unavailable{}, nil
	}
	return engine.StartSidecar(ctx, executable, os.Stderr)
}
