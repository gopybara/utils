package logger

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg Config) (*zap.Logger, error) {
	hooks := make([]zap.Option, 0)
	if cfg.GetSentryDsn() != nil {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              *cfg.GetSentryDsn(),
			Environment:      cfg.GetEnvironment(),
			TracesSampleRate: 0.2,
			Debug:            false,
			Release:          cfg.GetCommitTag(),
		})
		if err != nil {
			return nil, fmt.Errorf("cannot init sentry client: %w", err)
		}

		hooks = append(hooks, zap.Hooks(func(entry zapcore.Entry) error {
			if cfg.GetSentryDsn() != nil && entry.Level != zapcore.PanicLevel && entry.Level != zapcore.FatalLevel && entry.Level != zapcore.ErrorLevel {
				return nil
			}

			event := sentry.Event{
				Message: entry.Message,
			}
			defer sentry.CaptureEvent(&event)

			switch entry.Level {
			case zapcore.PanicLevel, zapcore.FatalLevel:
				event.Level = sentry.LevelFatal
			case zapcore.ErrorLevel:
				event.Level = sentry.LevelError
			}

			return nil
		}))
	}

	return zap.NewProduction(hooks...)
}
