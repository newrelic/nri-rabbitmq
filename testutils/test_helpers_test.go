package testutils

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/stretchr/testify/assert"
)

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

type subStruct struct {
	SubString string `json:"sub-string"`
}
type topStruct struct {
	String string
	Number int64
	Bool   bool
	Obj    *subStruct
	Array  []bool
}

var expectedTopStruct = &topStruct{
	"value",
	int64(1),
	true,
	&subStruct{"sub-value"},
	[]bool{true, false},
}

func Test_ReadStructFromJSONFile(t *testing.T) {
	actual := new(topStruct)
	ReadStructFromJSONFile(t, filepath.Join("testdata", "sample.json"), actual)
	assert.Equal(t, expectedTopStruct, actual)
}

func Test_ReadStructFromJSONString(t *testing.T) {
	actual := new(topStruct)
	data, _ := ioutil.ReadFile(filepath.Join("testdata", "sample.json"))
	ReadStructFromJSONString(t, string(data), actual)
	assert.Equal(t, expectedTopStruct, actual)
}

func Test_GetTestServer(t *testing.T) {
	mux, closer := GetTestServer(false)
	assert.NotNil(t, mux)
	assert.NotNil(t, closer)
	assert.NotEmpty(t, args.GlobalArgs.Hostname)
	assert.True(t, args.GlobalArgs.Port > 0)
	closer()
}
