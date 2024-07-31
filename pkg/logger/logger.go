package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel string `koanf:"APP_LOGGER_LEVEL"`
}

type Logger interface {
	Info(msg string)
	Track(resource, ip string, status int)
	Debug(err error)
	Error(err error)
}

type Log struct {
	zerolog.Logger
}

func New(lvl string) *Log {
	setLoggerLevel(lvl)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	return &Log{
		Logger: zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

func setLoggerLevel(lvl string) {
	var l zerolog.Level

	switch lvl {
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)
}

func (l *Log) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *Log) Track(resource, ip string, status int) {
	l.Logger.Info().
		Str("resource", resource).
		Str("ip", ip).
		Int("status", status).
		Send()
}

func (l *Log) Debug(err error) {
	l.Logger.Debug().Err(err).Send()
}

func (l *Log) Error(err error) {
	l.Logger.Error().Err(err).Send()
}
