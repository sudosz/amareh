package main

import (
	stdlog "log"

	"github.com/sudosz/amareh/i18n"
	"github.com/sudosz/amareh/internal/config"
	"github.com/sudosz/amareh/internal/logger"
)

func main() {
	cfg, err := config.LoadConfig("configs")
	if err != nil {
		stdlog.Fatalf("failed to load config: %v", err)
	}

	logger.InitLogger(
		logger.WithLevel(logger.DebugLevel),
		logger.WithConsoleOutput(true),
		logger.WithFileOutput(true),
		logger.WithFilename(cfg.LogDirectory+"/bot.log"),
	)

	log := logger.Logger()
	log.Info(i18n.T("hello.from", map[string]any{"name": "bot"}))
}
