package plugin

import (
	"context"
	"github.com/lzap/cborpc/log"
	"github.com/rs/zerolog"
)

type appLoggerImpl struct {
	Logger *zerolog.Logger
}

func (a appLoggerImpl) Msgf(level log.Level, format string, values ...any) {
	switch level {
	case log.TRC:
		a.Logger.Trace().Msgf(format, values...)
	case log.DBG:
		a.Logger.Debug().Msgf(format, values...)
	case log.INF:
        a.Logger.Info().Msgf(format, values...)
	case log.WRN:
        a.Logger.Warn().Msgf(format, values...)
	case log.ERR:
        a.Logger.Error().Msgf(format, values...)
	}
}

func WithPluginLogger(ctx context.Context, zl *zerolog.Logger) context.Context {
	return log.ContextWithLogger(ctx, &appLoggerImpl{Logger: zl})
}
