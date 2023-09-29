package main

import (
	"flag"
	"log/slog"
	"os"

	_ "github.com/lib/pq"

	"github.com/isutare412/tasks/api/internal/config"
	"github.com/isutare412/tasks/api/internal/log"
	"github.com/isutare412/tasks/api/internal/wire"
)

var configPath = flag.String("config", "config.yaml", "YAML config file path")

func init() {
	flag.Parse()
}

func main() {
	cfg, err := config.LoadValidated(*configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	log.Init(cfg.ToLogConfig())

	slog.Debug("loaded config", "config", cfg)

	components, err := wire.NewComponents(cfg)
	if err != nil {
		slog.Error("failed to wire components", "error", err)
		os.Exit(1)
	}

	if err := components.Initialize(); err != nil {
		slog.Error("failed to initialize components", "error", err)
		os.Exit(1)
	}

	components.Run()
	components.Shutdown()
}
