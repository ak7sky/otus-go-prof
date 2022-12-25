package logger

import (
	"os"

	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/app"
	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	debugLvl = "DEBUG"
	infoLvl  = "INFO"
	errorLvl = "ERROR"
)

// Zap Based implementation of Logger interface.
type logger struct {
	zapLogger *zap.Logger
}

func (l *logger) Debug(msg string) {
	l.zapLogger.Debug(msg)
}

func (l *logger) Info(msg string) {
	l.zapLogger.Info(msg)
}

func (l *logger) Error(msg string) {
	l.zapLogger.Error(msg)
}

func NewPreconfigured() app.Logger {
	return preconfigured()
}

func NewConfigured(loggerConf config.LoggerConf) app.Logger {
	return configured(loggerConf)
}

// Returns default safety (without error returning) initialized logger.
// Used to log events of app before logger configuration is completed.
// For example to log error during app config initialization.
func preconfigured() *logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "zapLogger",
		TimeKey:        "time",
		StacktraceKey:  "stacktrace",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), os.Stdout, zapcore.ErrorLevel)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	return &logger{zapLogger: zapLogger}
}

func configured(loggerConf config.LoggerConf) *logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "zapLogger",
		TimeKey:        "time",
		StacktraceKey:  "stacktrace",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	if loggerConf.IsJSONEnabled {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var lvlEnabler zapcore.LevelEnabler
	switch loggerConf.Level {
	case debugLvl:
		lvlEnabler = zapcore.DebugLevel
	case infoLvl:
		lvlEnabler = zapcore.InfoLevel
	case errorLvl:
		lvlEnabler = zapcore.ErrorLevel
	}
	core := zapcore.NewCore(encoder, os.Stdout, lvlEnabler)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	return &logger{zapLogger: zapLogger}
}
