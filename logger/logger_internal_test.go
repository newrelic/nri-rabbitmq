package logger

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {
	newLogger := log.NewStdErr(false)
	SetLogger(newLogger)
	assert.Equal(t, newLogger, integrationLogger)
	SetLogger(nil)
	assert.Equal(t, defaultLogger, integrationLogger)
}
