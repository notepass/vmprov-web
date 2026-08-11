package main

import (
	"fmt"
	"log/slog"

	"github.com/notepass/vmprov-web/internal/config"
	"github.com/notepass/vmprov-web/internal/logger"
	"github.com/notepass/vmprov-web/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}

	log := logger.New(cfg.LogLevel)
	slog.SetDefault(log)

	log.Info("starting server", "port", cfg.ServerPort)

	srv := server.New(cfg, log)
	addr := fmt.Sprintf(":%d", cfg.ServerPort)

	if err := server.Start(srv, addr); err != nil {
		log.Error("server error", "error", err)
		return
	}
}
