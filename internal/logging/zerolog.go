package logging

import (
	"fmt"
	"os"
	"time"
	"unicode/utf8"

	"github.com/hardcaporg/hardcap/internal/config"
	"github.com/hardcaporg/hardcap/internal/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func truncateText(str string, length int) string {
	if length <= 0 {
		return str
	}

	if utf8.RuneCountInString(str) <= length {
		return str
	}

	trimmed := []rune(str)[:length]

	if trimmed[0] == '"' {
		return string(trimmed) + "...\""
	} else {
		return string(trimmed) + "..."
	}
}

func decorate(l zerolog.Logger) zerolog.Logger {
	logger := l.With().Timestamp().
		Str("hostname", version.Hostname)

	if version.BuildCommit != "" {
		logger = logger.Str("version", version.BuildCommit)
	}

	return logger.Logger()
}

func Initialize() {
	level, err := zerolog.ParseLevel(config.Logging.Level)
	if err != nil {
		panic(fmt.Errorf("cannot parse log level '%s': %w", config.Logging.Level, err))
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = decorate(log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Kitchen,
		FormatFieldValue: func(i interface{}) string {
			return truncateText(fmt.Sprintf("%s", i), config.Logging.MaxField)
		},
	}))
}
