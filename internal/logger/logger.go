package logger

import (
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultFilename   = "logs/app.log"
	DefaultMaxSize    = 100 // megabytes
	DefaultMaxAge     = 30  // days
	DefaultMaxBackups = 10
	DefaultCompress   = true
)

type Level = zapcore.Level

const (
	DebugLevel Level = zapcore.DebugLevel
	InfoLevel  Level = zapcore.InfoLevel
	WarnLevel  Level = zapcore.WarnLevel
	ErrorLevel Level = zapcore.ErrorLevel
	FatalLevel Level = zapcore.FatalLevel
	PanicLevel Level = zapcore.PanicLevel
)

var (
	defaultLogger *zap.Logger
)

// InitLogger configures and returns a zap logger
func InitLogger(opts ...Option) {
	// Default config
	config := &Config{
		level:         zapcore.InfoLevel,
		consoleOutput: true,
		fileOutput:    true,
		filename:      "logs/app.log",
		maxSize:       100, // megabytes
		maxAge:        30,  // days
		maxBackups:    10,
		compress:      true,
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	var cores []zapcore.Core

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Configure console output
	if config.consoleOutput {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			config.level,
		)
		cores = append(cores, consoleCore)
	}

	// Configure file output with rotation
	if config.fileOutput {
		// Ensure logs directory exists
		if err := os.MkdirAll(filepath.Dir(config.filename), 0755); err != nil {
			panic(err)
		}

		fileWriter := &lumberjack.Logger{
			Filename:   config.filename,
			MaxSize:    config.maxSize,
			MaxAge:     config.maxAge,
			MaxBackups: config.maxBackups,
			Compress:   config.compress,
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(fileWriter),
			config.level,
		)
		cores = append(cores, fileCore)
	}

	// Create multi-core logger
	core := zapcore.NewTee(cores...)
	defaultLogger = zap.New(core, zap.AddCaller())
}

// Config holds logger configuration
type Config struct {
	level         zapcore.Level
	writer        io.Writer
	consoleOutput bool
	fileOutput    bool
	filename      string
	maxSize       int  // megabytes
	maxAge        int  // days
	maxBackups    int  // number of backups to keep
	compress      bool // compress rotated files
}

// Option is a function that configures the logger
type Option func(*Config)

// WithLevel sets the logging level
func WithLevel(level Level) Option {
	return func(c *Config) {
		c.level = level
	}
}

// WithOutput sets the output writer
func WithOutput(w io.Writer) Option {
	return func(c *Config) {
		c.writer = w
	}
}

// WithConsoleOutput enables/disables console output
func WithConsoleOutput(enabled bool) Option {
	return func(c *Config) {
		c.consoleOutput = enabled
	}
}

// WithFileOutput enables/disables file output
func WithFileOutput(enabled bool) Option {
	return func(c *Config) {
		c.fileOutput = enabled
	}
}

// WithFilename sets the log filename
func WithFilename(filename string) Option {
	return func(c *Config) {
		c.filename = filename
	}
}

// WithRotationConfig sets the log rotation configuration
func WithRotationConfig(maxSize, maxAge, maxBackups int, compress bool) Option {
	return func(c *Config) {
		c.maxSize = maxSize
		c.maxAge = maxAge
		c.maxBackups = maxBackups
		c.compress = compress
	}
}

// Logger returns the configured zap.Logger instance
func Logger() *zap.Logger {
	return defaultLogger
}
