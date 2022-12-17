package telemetry

import (
	"fmt"
	"os"
	"time"

	"github.com/UndeadDemidov/ya-pr-diplomb/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
)

type AppLogger struct {
	zerolog.Logger
	cfg *config.Logger
}

// NewAppLogger creates logger.
func NewAppLogger(cfg *config.Logger) *AppLogger {
	return &AppLogger{cfg: cfg}
}

// NewTestAppLogger creates logger for tests.
func NewTestAppLogger() AppLogger {
	return AppLogger{cfg: &config.Logger{Development: true, Level: "debug"}}
}

// For mapping config logger to app logger levels.
var loggerLevelMap = map[string]zerolog.Level{ //nolint:gochecknoglobals
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"panic": zerolog.PanicLevel,
	"fatal": zerolog.FatalLevel,
}

func (l *AppLogger) getLevel(cfg *config.Logger) zerolog.Level {
	level, exist := loggerLevelMap[cfg.Level]
	if !exist {
		return zerolog.DebugLevel
	}

	return level
}

// InitLogger initialization logger.
// ToDo заменить на endpoint, который меняет уровень логирования.
func (l *AppLogger) InitLogger() {
	zerolog.SetGlobalLevel(l.getLevel(l.cfg))
	if l.cfg.Development {
		l.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().Timestamp().Caller().Stack().Logger()
	} else {
		wr := diode.NewWriter(os.Stdout, 1000, 10*time.Millisecond, func(missed int) { //nolint:gomnd
			fmt.Printf("Logger Dropped %d messages", missed) //nolint:forbidigo
		})
		l.Logger = log.Output(wr).With().Timestamp().Logger()
	}
}

// Logger methods

func (l *AppLogger) Debug(args ...interface{}) {
	l.Logger.Debug().Msg(fmt.Sprint(args...))
}

func (l *AppLogger) Debugf(template string, args ...interface{}) {
	l.Logger.Debug().Msg(fmt.Sprintf(template, args...))
}

func (l *AppLogger) Info(args ...interface{}) {
	l.Logger.Info().Msg(fmt.Sprint(args...))
}

func (l *AppLogger) Infof(template string, args ...interface{}) {
	l.Logger.Info().Msg(fmt.Sprintf(template, args...))
}

func (l *AppLogger) Warn(args ...interface{}) {
	l.Logger.Warn().Msg(fmt.Sprint(args...))
}

func (l *AppLogger) Warnf(template string, args ...interface{}) {
	l.Logger.Warn().Msg(fmt.Sprintf(template, args...))
}

func (l *AppLogger) Error(args ...interface{}) {
	l.Logger.Error().Msg(fmt.Sprint(args...))
}

func (l *AppLogger) Errorf(template string, args ...interface{}) {
	l.Logger.Error().Msg(fmt.Sprintf(template, args...))
}

func (l *AppLogger) Panic(args ...interface{}) {
	l.Logger.Panic().Msg(fmt.Sprint(args...))
}

func (l *AppLogger) Panicf(template string, args ...interface{}) {
	l.Logger.Panic().Msg(fmt.Sprintf(template, args...))
}

func (l *AppLogger) Fatal(args ...interface{}) {
	l.Logger.Fatal().Msg(fmt.Sprint(args...))
}

func (l *AppLogger) Fatalf(template string, args ...interface{}) {
	l.Logger.Fatal().Msg(fmt.Sprintf(template, args...))
}

func (l *AppLogger) Print(args ...interface{}) {
	l.Logger.Info().Msg(fmt.Sprint(args...))
}

func (l *AppLogger) Printf(format string, args ...interface{}) {
	l.Logger.Info().Msg(fmt.Sprintf(format, args...))
}

func (l *AppLogger) Println(args ...interface{}) {
	l.Print(args...)
}
