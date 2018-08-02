package logger_test

import (
	"testing"

	"github.com/newrelic/nri-rabbitmq/logger"
	"github.com/newrelic/nri-rabbitmq/testutils"
)

func TestLogger(t *testing.T) {
	mockLogger := new(testutils.MockedLogger)
	mockLogger.On("Debugf", "debug1", []interface{}{"debug2"})
	mockLogger.On("Warnf", "warn1", []interface{}{"warn2"})
	mockLogger.On("Infof", "info1", []interface{}{"info2"})
	mockLogger.On("Errorf", "error1", []interface{}{"error2"})

	logger.SetLogger(mockLogger)
	logger.Debugf("debug1", "debug2")
	logger.Warnf("warn1", "warn2")
	logger.Infof("info1", "info2")
	logger.Errorf("error1", "error2")
	mockLogger.AssertExpectations(t)
}
