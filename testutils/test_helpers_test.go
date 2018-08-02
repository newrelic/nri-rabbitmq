package testutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetTestLogger(t *testing.T) {
	testLogger = nil
	newLogger := SetTestLogger(t)
	assert.NotNil(t, newLogger)
	assert.Equal(t, newLogger, testLogger)
}

func Test_SetMockLogger(t *testing.T) {
	testLogger = nil
	newLogger := SetMockLogger()
	assert.NotNil(t, newLogger)
	assert.IsType(t, &MockedLogger{}, newLogger)
}

func Test_GetIntegrationEntity(t *testing.T) {
	testIntegration, testEntity := GetTestingEntity(t)
	assert.NotNil(t, testIntegration)
	assert.Equal(t, "Test", testIntegration.Name)
	assert.Equal(t, "0.0.1", testIntegration.IntegrationVersion)
	assert.NotNil(t, testEntity)
	assert.True(t, testEntity.Metadata == nil || testEntity.Metadata.Name == "")

	_, testEntity = GetTestingEntity(t, "name", "namespace")
	assert.NotNil(t, testEntity)
	assert.Equal(t, "name", testEntity.Metadata.Name)
	assert.Equal(t, "namespace", testEntity.Metadata.Namespace)
}

func Test_ReadObjectFromJSONFile(t *testing.T) {
	actual := ReadObjectFromJSONFile(t, filepath.Join("testdata", "sample.json"))
	expected := map[string]interface{}{
		"string": "value",
		"number": 1.0,
		"bool":   true,
		"obj": map[string]interface{}{
			"sub-string": "sub-value",
		},
		"array": []interface{}{true, false},
	}
	assert.Equal(t, expected, actual)
}
