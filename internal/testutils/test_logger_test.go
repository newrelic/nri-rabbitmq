package testutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestLogger(t *testing.T) {
	topFormat := "%v"
	logLevel := "debug"
	x := func(format string, args ...interface{}) {
		switch logLevel {
		case "debug":
			assert.Equal(t, strings.ToUpper(logLevel)+": "+topFormat, format)
			assert.Equal(t, 1, len(args))
			assert.Equal(t, logLevel, args[0])
		}
	}
	testLogger := TestLogger{x}
	testLogger.Debugf(topFormat, logLevel)
	logLevel = "info"
	testLogger.Infof(topFormat, logLevel)
	logLevel = "error"
	testLogger.Errorf(topFormat, logLevel)
	logLevel = "warn"
	testLogger.Warnf(topFormat, logLevel)
}
