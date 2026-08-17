package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/notepass/vmprov-web/internal/config"
)

// RouteRegistrar registers API routes on the Echo instance.
type RouteRegistrar interface {
	RegisterRoutes(e *echo.Echo)
}

// New creates a new HTTP server using Echo.
// Handlers implementing RouteRegistrar are registered in the given order.
func New(cfg *config.Config, logger *slog.Logger, registrars ...RouteRegistrar) *http.Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	srv := e.Server
	srv.ReadTimeout = 15 * time.Second
	srv.WriteTimeout = 15 * time.Second

	for _, r := range registrars {
		r.RegisterRoutes(e)
	}

	return srv
}

// Start runs the server with graceful shutdown handling.
// onShutdown is called after the HTTP server has been shut down,
// allowing cleanup of resources like database connections.
func Start(srv *http.Server, addr string, onShutdown func()) error {
	srv.Addr = addr
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	if onShutdown != nil {
		onShutdown()
	}

	slog.Info("server stopped")
	return nil
}
