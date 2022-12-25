package logger

import (
	"testing"

	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/app"
	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
	"tideland.dev/go/audit/capture"
)

var levels = map[string]int{debugLvl: 0, infoLvl: 1, errorLvl: 2}

func TestNewConfigured(t *testing.T) {
	testCases := []struct {
		name         string
		conf         config.LoggerConf
		msgLvl       string
		msg          string
		expLogFields []string
	}{
		{
			name:   "logLvl-DEBUG;json:-true;msgLvl-DEBUG",
			conf:   config.LoggerConf{Level: debugLvl, IsJSONEnabled: true},
			msgLvl: debugLvl,
			msg:    "test debug msg",
			expLogFields: []string{
				`"level":"debug"`,
				`"time":`,
				`"msg":"test debug msg"`,
				`"caller":"logger/logger_test.go:108"`,
			},
		},
		{
			name:         "logLvl-DEBUG;json:-false;msgLvl-DEBUG",
			conf:         config.LoggerConf{Level: debugLvl, IsJSONEnabled: false},
			msgLvl:       debugLvl,
			msg:          "test debug msg",
			expLogFields: []string{"debug\tlogger/logger_test.go:108\ttest debug msg"},
		},
		{
			name:   "logLvl-INFO;json:-true;msgLvl-INFO",
			conf:   config.LoggerConf{Level: infoLvl, IsJSONEnabled: true},
			msgLvl: infoLvl,
			msg:    "test info msg",
			expLogFields: []string{
				`"level":"info"`,
				`"time":`,
				`"msg":"test info msg"`,
				`"caller":"logger/logger_test.go:108"`,
			},
		},
		{
			name:         "logLvl-INFO;json:-false;msgLvl-INFO",
			conf:         config.LoggerConf{Level: infoLvl, IsJSONEnabled: false},
			msgLvl:       infoLvl,
			msg:          "test info msg",
			expLogFields: []string{"info\tlogger/logger_test.go:108\ttest info msg"},
		},
		{
			name:   "logLvl-INFO;json:-true;msgLvl-DEBUG",
			conf:   config.LoggerConf{Level: infoLvl, IsJSONEnabled: true},
			msgLvl: debugLvl,
			msg:    "test debug msg",
		},
		{
			name:   "logLvl-INFO;json:-true;msgLvl-ERROR",
			conf:   config.LoggerConf{Level: infoLvl, IsJSONEnabled: true},
			msgLvl: errorLvl,
			msg:    "test error msg",
			expLogFields: []string{
				`"level":"error"`,
				`"time":`,
				`"msg":"test error msg"`,
				`"caller":"logger/logger_test.go:108"`,
				`"stacktrace":`,
			},
		},
		{
			name:   "logLvl-ERROR;json:-true;msgLvl-DEBUG",
			conf:   config.LoggerConf{Level: errorLvl, IsJSONEnabled: true},
			msgLvl: debugLvl,
			msg:    "test debug msg",
		},
		{
			name:   "logLvl-ERROR;json:-true;msgLvl-INFO",
			conf:   config.LoggerConf{Level: errorLvl, IsJSONEnabled: true},
			msgLvl: infoLvl,
			msg:    "test info msg",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var loggerMethod func(app.Logger, string)
			switch tc.msgLvl {
			case debugLvl:
				loggerMethod = app.Logger.Debug
			case infoLvl:
				loggerMethod = app.Logger.Info
			case errorLvl:
				loggerMethod = app.Logger.Error
			}

			log := capture.Stdout(func() {
				logger := NewConfigured(tc.conf)
				loggerMethod(logger, tc.msg)
			})

			if levels[tc.msgLvl] < levels[tc.conf.Level] {
				require.Equal(t, "", log.String(),
					"unexpected log, msgLvl (%s) < loggerLvl (%s)", tc.msgLvl, tc.conf.Level)
				return
			}

			for _, expLogField := range tc.expLogFields {
				require.Contains(t, log.String(), expLogField)
			}
		})
	}
}

func TestNewPreconfigured(t *testing.T) {
	testCases := []struct {
		name               string
		msgLvl             string
		msg                string
		isEmptyLogExpected bool
		expLogFields       []string
	}{
		{
			name:               "msgLvl-DEBUG",
			msgLvl:             debugLvl,
			msg:                "test debug msg",
			isEmptyLogExpected: true,
		},
		{
			name:               "msgLvl-INFO",
			msgLvl:             infoLvl,
			msg:                "test info msg",
			isEmptyLogExpected: true,
		},
		{
			name:         "msgLvl-ERROR",
			msgLvl:       errorLvl,
			msg:          "test error msg",
			expLogFields: []string{"error\tlogger/logger_test.go:167\ttest error msg"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var loggerMethod func(app.Logger, string)
			switch tc.msgLvl {
			case debugLvl:
				loggerMethod = app.Logger.Debug
			case infoLvl:
				loggerMethod = app.Logger.Info
			case errorLvl:
				loggerMethod = app.Logger.Error
			}

			log := capture.Stdout(func() {
				logger := NewPreconfigured()
				loggerMethod(logger, tc.msg)
			})

			if tc.isEmptyLogExpected {
				require.Equal(t, "", log.String(), "unexpected log, msgLvl < loggerLvl")
				return
			}

			for _, expLogField := range tc.expLogFields {
				require.Contains(t, log.String(), expLogField)
			}
		})
	}
}
