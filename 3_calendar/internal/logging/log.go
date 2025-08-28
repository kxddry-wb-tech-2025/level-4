package logging

import (
	appLog "calendar/internal/models/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap.Logger.
type Logger struct {
	log *zap.Logger
}

// New creates a new Logger instance.
func New(env string) *Logger {
	var cfg zap.Config

	cfg = zap.NewProductionConfig()
	cfg.DisableStacktrace = true
	cfg.DisableCaller = true
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch env {
	case "dev", "local":
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	case "prod":
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		cfg.DisableStacktrace = true
		cfg.DisableCaller = true
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	zl := zap.Must(cfg.Build())
	return &Logger{log: zl}
}

// Listen consumes log entries from the channel and writes them using Zap.
func (l *Logger) Listen(entries <-chan appLog.Entry) {
	go func() {
		for e := range entries {
			l.writeEntry(e)
		}
	}()
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) writeEntry(e appLog.Entry) {
	fields := fieldsFromEntry(e)
	switch e.Level {
	case appLog.LevelDebug:
		l.log.Debug(e.Message, fields...)
	case appLog.LevelInfo:
		l.log.Info(e.Message, fields...)
	case appLog.LevelWarn:
		l.log.Warn(e.Message, fields...)
	case appLog.LevelError:
		l.log.Error(e.Message, fields...)
	default:
		l.log.Info(e.Message, fields...)
	}
}

func fieldsFromEntry(e appLog.Entry) []zap.Field {
	fields := make([]zap.Field, 0, len(e.Meta)+1)
	if e.Error != nil {
		fields = append(fields, zap.Error(e.Error))
	}
	for k, v := range e.Meta {
		switch val := v.(type) {
		case string:
			fields = append(fields, zap.String(k, val))
		case bool:
			fields = append(fields, zap.Bool(k, val))
		case int:
			fields = append(fields, zap.Int(k, val))
		case int32:
			fields = append(fields, zap.Int32(k, val))
		case int64:
			fields = append(fields, zap.Int64(k, val))
		case uint:
			fields = append(fields, zap.Uint(k, val))
		case uint32:
			fields = append(fields, zap.Uint32(k, val))
		case uint64:
			fields = append(fields, zap.Uint64(k, val))
		case float32:
			fields = append(fields, zap.Float32(k, val))
		case float64:
			fields = append(fields, zap.Float64(k, val))
		default:
			fields = append(fields, zap.Any(k, v))
		}
	}
	return fields
}
