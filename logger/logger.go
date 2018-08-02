package logger

import "github.com/newrelic/infra-integrations-sdk/log"

var defaultLogger = log.NewStdErr(false)
var integrationLogger = defaultLogger

// SetLogger sets the package logger
func SetLogger(newLogger log.Logger) {
	if newLogger == nil && integrationLogger != defaultLogger {
		integrationLogger = defaultLogger
	} else if integrationLogger != newLogger {
		integrationLogger = newLogger
	}
}

// Debugf is a wrapper for log.Logger.Debugf
func Debugf(format string, args ...interface{}) {
	integrationLogger.Debugf(format, args)
}

// Warnf is a wrapper for log.Logger.Warnf
func Warnf(format string, args ...interface{}) {
	integrationLogger.Warnf(format, args)
}

// Infof is a wrapper for log.Logger.Infof
func Infof(format string, args ...interface{}) {
	integrationLogger.Infof(format, args)
}

// Errorf is a wrapper for log.Logger.Errorf
func Errorf(format string, args ...interface{}) {
	integrationLogger.Errorf(format, args)
}
