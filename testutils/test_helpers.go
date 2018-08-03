package testutils

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-rabbitmq/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// Update flag will update .golden files to the current actual
	Update     = flag.Bool("update", false, "update .golden files")
	testLogger log.Logger
)

// SetTestLogger creates a logger that logs to the testing framework and returns that logger
func SetTestLogger(t *testing.T) log.Logger {
	if testLogger == nil {
		testLogger = &TestLogger{F: t.Logf}
	}
	logger.SetLogger(testLogger)
	return testLogger
}

// SetMockLogger creates a new MockedLogger, sets it to be used, and returns it so that it can be setup
func SetMockLogger() *MockedLogger {
	l := new(MockedLogger)
	logger.SetLogger(l)
	return l
}

// GetTestingIntegration creates an Integration used for testing and sets the logger to the integration's logger
func GetTestingIntegration(t *testing.T) (payload *integration.Integration) {
	testLogger := SetTestLogger(t)
	payload, err := integration.New("Test", "0.0.1", integration.Logger(testLogger))
	require.NoError(t, err)
	require.NotNil(t, payload)
	logger.SetLogger(payload.Logger())
	return
}

// GetTestingEntity creates an Entity used for testing
func GetTestingEntity(t *testing.T, entityArgs ...string) (payload *integration.Integration, entity *integration.Entity) {
	payload = GetTestingIntegration(t)
	var err error
	if len(entityArgs) > 1 {
		entity, err = payload.Entity(entityArgs[0], entityArgs[1])
		assert.NoError(t, err)
	} else {
		entity = payload.LocalEntity()
	}
	require.NotNil(t, entity)
	return
}

// ReadObjectFromJSONFile reads a generic map[string]interface{} from a file, typically used for reading JSON
func ReadObjectFromJSONFile(t *testing.T, filename string) map[string]interface{} {
	data, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	item := map[string]interface{}{}
	err = json.Unmarshal(data, &item)
	require.NoError(t, err)
	return item
}

// ReadObjectFromJSONString reads a generic map[string]interface{} from a json string
func ReadObjectFromJSONString(t *testing.T, rawJSON string) map[string]interface{} {
	item := map[string]interface{}{}
	err := json.Unmarshal([]byte(rawJSON), &item)
	require.NoError(t, err)
	return item
}
