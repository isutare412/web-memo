package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/isutare412/web-memo/api/internal/config"
	"github.com/isutare412/web-memo/api/internal/log"
	"github.com/isutare412/web-memo/api/internal/metric"
	"github.com/isutare412/web-memo/api/internal/wire"
)

var configPath = flag.String("configs", ".", "path to config directory")

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
	if cfg.Metrics.Enabled {
		metric.Init()
	}

	slog.Debug("loaded config", "config", cfg)

	app, err := wire.NewApp(cfg)
	if err != nil {
		slog.Error("failed to wire app", "error", err)
		os.Exit(1)
	}

	if err := app.Initialize(); err != nil {
		slog.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}

	app.Run()
	app.Shutdown()
}
