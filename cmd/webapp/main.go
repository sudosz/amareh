package main

import (
	"github.com/sudosz/amareh/i18n"
	"github.com/sudosz/amareh/internal/logger"
)

func init() {
	logger.InitLogger(
		logger.WithLevel(logger.DebugLevel),
		logger.WithConsoleOutput(true),
		logger.WithFileOutput(true),
		logger.WithFilename("logs/webapp.log"),
	)
}

func main() {
	log := logger.Logger()
	log.Info(i18n.T("hello.from", map[string]any{"name": "webapp"}))
}
