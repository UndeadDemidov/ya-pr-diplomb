package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/UndeadDemidov/ya-pr-diplomb/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

// Logger methods interface.
type Logger interface { //nolint:interfacebloat
	InitLogger()
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type apiLogger struct {
	cfg    *config.App
	logger *zerolog.Logger
}

// NewAPILogger creates logger.
func NewAPILogger(cfg *config.App) *apiLogger { //nolint:revive
	return &apiLogger{cfg: cfg}
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

func (l *apiLogger) getLevel(cfg *config.App) zerolog.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zerolog.DebugLevel
	}

	return level
}

// InitLogger initialization logger.
func (l *apiLogger) InitLogger() {
	zerolog.SetGlobalLevel(l.getLevel(l.cfg))

	var logger zerolog.Logger

	if l.cfg.Logger.Development {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().Timestamp().Caller().Stack().Logger()
	} else {
		wr := diode.NewWriter(os.Stdout, 1000, 10*time.Millisecond, func(missed int) { //nolint:gomnd
			fmt.Printf("Logger Dropped %d messages", missed) //nolint:forbidigo
		})
		logger = zerolog.New(wr).With().Timestamp().Logger()
	}

	l.logger = &logger
}

// Logger methods

func (l *apiLogger) Debug(args ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprint(args...))
}

func (l *apiLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprintf(template, args...))
}

func (l *apiLogger) Info(args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprint(args...))
}

func (l *apiLogger) Infof(template string, args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprintf(template, args...))
}

func (l *apiLogger) Warn(args ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprint(args...))
}

func (l *apiLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprintf(template, args...))
}

func (l *apiLogger) Error(args ...interface{}) {
	l.logger.Error().Msg(fmt.Sprint(args...))
}

func (l *apiLogger) Errorf(template string, args ...interface{}) {
	l.logger.Error().Msg(fmt.Sprintf(template, args...))
}

func (l *apiLogger) Panic(args ...interface{}) {
	l.logger.Panic().Msg(fmt.Sprint(args...))
}

func (l *apiLogger) Panicf(template string, args ...interface{}) {
	l.logger.Panic().Msg(fmt.Sprintf(template, args...))
}

func (l *apiLogger) Fatal(args ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (l *apiLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprintf(template, args...))
}

func (l *apiLogger) Print(args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprint(args...))
}

func (l *apiLogger) Printf(format string, args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprintf(format, args...))
}

func (l *apiLogger) Println(args ...interface{}) {
	l.Print(args...)
}
