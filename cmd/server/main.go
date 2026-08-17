package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/notepass/vmprov-web/internal/api"
	"github.com/notepass/vmprov-web/internal/config"
	"github.com/notepass/vmprov-web/internal/db"
	"github.com/notepass/vmprov-web/internal/libvirt"
	"github.com/notepass/vmprov-web/internal/logger"
	"github.com/notepass/vmprov-web/internal/repository"
	"github.com/notepass/vmprov-web/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel)
	slog.SetDefault(log)

	// Initialize database
	adaptor, err := db.NewAdaptor(db.Config{
		ConnString:      cfg.DBConnString,
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime,
	})
	if err != nil {
		log.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	// Run migrations
	if err := db.RunMigrations(adaptor); err != nil {
		log.Error("failed to run migrations", "error", err)
		adaptor.Close()
		os.Exit(1)
	}

	// Health check
	if err := adaptor.HealthCheck(); err != nil {
		log.Error("database health check failed", "error", err)
		adaptor.Close()
		os.Exit(1)
	}

	log.Info("database connected")

	// Create repositories
	dialect := repository.DialectMySQL
	if _, ok := adaptor.(*db.PostgresAdaptor); ok {
		dialect = repository.DialectPostgres
	}

	userRepo := repository.NewUserRepo(adaptor.DB(), dialect, log)
	templateRepo := repository.NewTemplateRepo(adaptor.DB(), dialect, log)
	auditLogRepo := repository.NewAuditLogRepo(adaptor.DB(), dialect, log)
	connRepo := repository.NewLibvirtConnectionRepo(adaptor.DB(), dialect, log)

	// Ensure the libvirt known_hosts file exists.
	knownHostsFile := expandHome(cfg.LibvirtKnownHostsFile)
	if err := ensureKnownHostsFile(knownHostsFile); err != nil {
		log.Error("failed to ensure known_hosts file", "path", knownHostsFile, "error", err)
		adaptor.Close()
		os.Exit(1)
	}

	connHandler := api.NewLibvirtConnectionHandler(
		connRepo,
		libvirt.New(),
		time.Duration(cfg.LibvirtConnectTimeout)*time.Second,
		knownHostsFile,
		log,
	)

	log.Info("starting server", "port", cfg.ServerPort)

	srv := server.New(cfg, log, connHandler)
	addr := fmt.Sprintf(":%d", cfg.ServerPort)

	if err := server.Start(srv, addr, func() {
		log.Info("closing database connection...")
		if err := adaptor.Close(); err != nil {
			log.Error("failed to close database", "error", err)
			return
		}
		log.Info("database connection closed")
	}); err != nil {
		log.Error("server error", "error", err)
		os.Exit(1)
	}

	_ = userRepo
	_ = templateRepo
	_ = auditLogRepo
}

// expandHome replaces a leading ~/ with the user's home directory.
func expandHome(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil || home == "" {
			return path
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~"))
	}
	return path
}

// ensureKnownHostsFile creates the known_hosts file and its parent
// directory if they do not exist.
func ensureKnownHostsFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err == nil {
		f.Close()
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(path, nil, 0o600)
}
