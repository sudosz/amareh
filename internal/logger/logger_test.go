// I write this file for debug purpose, because I want to compare the performance of different loggers
// I'm not sure if this is the best way to do it, but it works for now
package logger_test

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	testMessage = "test message"
	testString  = "four score and seven years"
	testInt     = 2023
	testFloat   = 3.14159
	testBool    = true
)

type concurrentBuffer struct {
	sync.Mutex
	bytes.Buffer
}

func (b *concurrentBuffer) Write(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	return b.Buffer.Write(p)
}

func newZerologLogger(w io.Writer) zerolog.Logger {
	return zerolog.New(w).
		With().
		Timestamp().
		Caller().
		Logger()
}

func newZerologConsoleLogger(w io.Writer) zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}).With().Timestamp().Caller().Logger()
}

func newZapLogger(w io.Writer) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.RFC3339TimeEncoder

	return zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(w),
		zap.InfoLevel,
	))
}

func newZapConsoleLogger(w io.Writer) *zap.Logger {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(w),
		zap.InfoLevel,
	))
}

func newSlogLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func newSlogConsoleLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func BenchmarkZerologLogger(b *testing.B) {
	logger := newZerologLogger(io.Discard)
	b.ReportAllocs()

	for b.Loop() {
		logger.Info().Msg(testMessage)
	}
}

func BenchmarkZapLogger(b *testing.B) {
	logger := newZapLogger(io.Discard)
	b.ReportAllocs()

	for b.Loop() {
		logger.Info(testMessage)
	}
}

func BenchmarkSlogLogger(b *testing.B) {
	logger := newSlogLogger(io.Discard)
	b.ReportAllocs()

	for b.Loop() {
		logger.Info(testMessage)
	}
}

func BenchmarkZerologLoggerWithConsoleOutput(b *testing.B) {
	logger := newZerologConsoleLogger(os.Stdout)
	b.ReportAllocs()

	for b.Loop() {
		logger.Info().Msg(testMessage)
	}
}

func BenchmarkZapLoggerWithConsoleOutput(b *testing.B) {
	logger := newZapConsoleLogger(os.Stdout)
	b.ReportAllocs()

	for b.Loop() {
		logger.Info(testMessage)
	}
}

func BenchmarkSlogLoggerWithConsoleOutput(b *testing.B) {
	logger := newSlogConsoleLogger(os.Stdout)
	b.ReportAllocs()

	for b.Loop() {
		logger.Info(testMessage)
	}
}

func BenchmarkZerologLoggerWithBuffer(b *testing.B) {
	buf := &concurrentBuffer{}
	logger := newZerologLogger(buf)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().
				Str("string", testString).
				Int("integer", testInt).
				Float64("float", testFloat).
				Bool("bool", testBool).
				Msg(testMessage)
		}
	})
}

func BenchmarkZapLoggerWithBuffer(b *testing.B) {
	buf := &concurrentBuffer{}
	logger := newZapLogger(buf)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(testMessage,
				zap.String("string", testString),
				zap.Int("integer", testInt),
				zap.Float64("float", testFloat),
				zap.Bool("bool", testBool),
			)
		}
	})
}

func BenchmarkSlogLoggerWithBuffer(b *testing.B) {
	buf := &concurrentBuffer{}
	logger := newSlogLogger(buf)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(testMessage,
				"string", testString,
				"integer", testInt,
				"float", testFloat,
				"bool", testBool,
			)
		}
	})
}
