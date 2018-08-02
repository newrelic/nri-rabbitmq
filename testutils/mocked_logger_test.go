package testutils

import "testing"

func Test_MockedLogger(t *testing.T) {
	topFormat := "%v"
	mockedLogger := MockedLogger{}
	mockedLogger.On("Debugf", topFormat, "debug").Once()
	mockedLogger.On("Infof", topFormat, "info").Once()
	mockedLogger.On("Errorf", topFormat, "error").Once()
	mockedLogger.On("Warnf", topFormat, "warn").Once()

	mockedLogger.Debugf(topFormat, "debug")
	mockedLogger.Infof(topFormat, "info")
	mockedLogger.Errorf(topFormat, "error")
	mockedLogger.Warnf(topFormat, "warn")

	mockedLogger.AssertExpectations(t)
}
