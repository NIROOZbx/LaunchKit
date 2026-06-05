package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Lauchkit/LaunchKit/shared/config"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

func NewLogger(cfg *config.LogConfig) zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}

	var writers []io.Writer
	writers = append(writers, consoleWriter)

	// Only adds file writer if a file path is configured
	if cfg.File != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.File), 0755); err != nil {
			panic("Failed to create log directory: " + err.Error())
		}

		fileWriter := &lumberjack.Logger{
			Filename:   cfg.File,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   true,
		}
		writers = append(writers, fileWriter)
	}

	level, _ := zerolog.ParseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	return zerolog.New(io.MultiWriter(writers...)).
		With().
		Timestamp().
		Caller().
		Logger()
}