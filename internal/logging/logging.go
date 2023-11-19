package logging

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	logger zerolog.Logger
	Level  = zerolog.InfoLevel
)

func init() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger = zerolog.New(consoleWriter).With().Timestamp().Logger().Level(Level)
}

func GetLogger() zerolog.Logger {
	return logger
}

func WithVerbose() zerolog.Logger {
	return logger.With().Logger().Level(zerolog.DebugLevel)
}

func EnableDebugLogging() {
	logger = logger.Level(zerolog.DebugLevel)
}
