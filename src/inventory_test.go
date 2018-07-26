package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func errorOsOpen(name string) (*os.File, error) {
	return nil, errors.New("an error other than file not found")
}

const (
	expectedNodeCmdOutput = `'node1'
`
	expectedNodeName = "node1"
)

var testConfigPath = filepath.Join("testdata", "sample.conf")

func TestGetNodeName(t *testing.T) {
	args.NodeNameOverride = expectedNodeName
	nodeName, err := getNodeName()
	assert.NoError(t, err)
	assert.Equal(t, expectedNodeName, nodeName)
	args.NodeNameOverride = ""

	prevExec := execCommand
	execCommand = fakeExecCommand
	os.Setenv("GO_WANT_HELPER_PROCESS", "1")
	defer func() {
		execCommand = prevExec
		os.Unsetenv("GO_WANT_HELPER_PROCESS")
	}()

	os.Setenv("GET_NODE_NAME_ERROR", "1")
	_, err = getNodeName()
	assert.Error(t, err)
	os.Unsetenv("GET_NODE_NAME_ERROR")

	os.Setenv("GET_NODE_NAME_EMPTY", "1")
	_, err = getNodeName()
	assert.EqualError(t, err, "could not determine the local node name")
	os.Unsetenv("GET_NODE_NAME_EMPTY")

	nodeName, err = getNodeName()
	assert.NoError(t, err)
	assert.Equal(t, expectedNodeName, nodeName)
}

func TestGetNodeEntity(t *testing.T) {
	i := getTestingIntegration(t)
	nodeEntity, err := getNodeEntity("node1", nil, nil)
	assert.EqualError(t, err, "node name [node1] not found in cluster")

	nodeData := []objx.Map{
		objx.MSI("name", "node1"),
		objx.MSI("name", "node2"),
	}

	nodeEntity, err = getNodeEntity("node1", nodeData, i)
	assert.NoError(t, err)
	assert.NotNil(t, nodeEntity)
	assert.Equal(t, "node1", nodeEntity.Metadata.Name)
	assert.Equal(t, nodeType, nodeEntity.Metadata.Namespace)

	nodeEntity, err = getNodeEntity("node2", nodeData, i)
	assert.NoError(t, err)
	assert.NotNil(t, nodeEntity)
	assert.Equal(t, "node2", nodeEntity.Metadata.Name)
	assert.Equal(t, nodeType, nodeEntity.Metadata.Namespace)
}

func TestConfigNotExist(t *testing.T) {
	args.ConfigPath = filepath.Join("testdata", "file-not_found.config")
	_, e := getTestingEntity(t)
	err := setInventoryData(e, nil)
	assert.Nil(t, err)
}

func TestConfigOpenError(t *testing.T) {
	args.ConfigPath = filepath.Join("testdata")
	_, e := getTestingEntity(t)
	err := setInventoryData(e, nil)
	assert.Error(t, err)

	prevOpen := osOpen
	osOpen = errorOsOpen
	defer func() {
		osOpen = prevOpen
	}()

	err = setInventoryData(e, nil)
	assert.Error(t, err)
}

func TestSetInventoryData(t *testing.T) {
	i, e := getTestingEntity(t)

	args.ConfigPath = ""
	err := setInventoryData(e, nil)
	assert.NoError(t, err)

	args.ConfigPath = testConfigPath
	golden := args.ConfigPath + ".golden"
	err = setInventoryData(e, nil)
	require.NoError(t, err)

	actual, err := i.MarshalJSON()
	require.NoError(t, err)

	if *update {
		t.Log("Writing .golden file")
		err = ioutil.WriteFile(golden, actual, 0644)
		require.NoError(t, err)
	}

	expected, _ := ioutil.ReadFile(golden)

	if !bytes.Equal(expected, actual) {
		t.Errorf("Actual JSON results do not match expected .golden file for %s", args.ConfigPath)
	}
}

func TestGetConfigPath(t *testing.T) {
	args.ConfigPath = testConfigPath
	actual := getConfigPath(nil)
	assert.Equal(t, testConfigPath, actual)
	args.ConfigPath = ""

	actual = getConfigPath(nil)
	assert.Empty(t, actual)

	overviewData := map[string]interface{}{
		"config_files": []string{
			"/etc/rabbitmq/rabbitmq.config",
			"/etc/rabbitmq/advanced.config",
		},
	}
	actual = getConfigPath(overviewData)
	assert.Empty(t, actual)

	overviewData = map[string]interface{}{
		"config_files": []string{
			"/etc/rabbitmq/rabbitmq.conf",
			"/etc/rabbitmq/advanced.config",
		},
	}
	actual = getConfigPath(overviewData)
	assert.Equal(t, "/etc/rabbitmq/rabbitmq.conf", actual)
}

func Test_PopulateEntityInventory(t *testing.T) {
	sourceFile := filepath.Join("testdata", "populateEntityInventory.json")
	exchangeGoldenFile := sourceFile + ".exchange.golden"
	queueGoldenFile := sourceFile + ".queue.golden"
	data := objx.New(readObjectFromJSONFile(t, sourceFile))

	i, exchangeEntity := getTestingEntity(t, "my-exchange", exchangeType)
	l := new(mockedLogger)
	logger = l
	l.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	// exchange test
	populateEntityInventory(exchangeEntity, exchangeType, &data)
	setInventoryItem(exchangeEntity, "this-key-is-longer-than-375-to-force-an-error-lorem-ipsum-dolor-sit-amet-consectetur-adipiscing-elit-sed-do-eiusmod-tempor-incididunt-ut-labore-et-dolore-magna-aliqua-ut-enim-ad-minim-veniam-quis-nostrud-exercitation-ullamco-laboris-nisi-ut-aliquip-ex-ea-commodo-consequat-duis-aute-irure-dolor-in-reprehenderit-in-voluptate-velit-esse-cillum-dolore-eu-fugiat-nulla-pariatures", "nope", "false")

	exchangeActual, err := json.Marshal(exchangeEntity.Inventory)
	assert.NoError(t, err)

	if *update {
		t.Log("Writing .golden file")
		err := ioutil.WriteFile(exchangeGoldenFile, exchangeActual, 0644)
		assert.NoError(t, err)
	}
	exchangeExpected, err := ioutil.ReadFile(exchangeGoldenFile)
	assert.NoError(t, err)

	assert.Equal(t, exchangeExpected, exchangeActual, "Expected doesn't match .golden file %v", exchangeGoldenFile)

	// queue test
	queueEntity, _, err := createEntity(i, "my-queue", queueType, "/")
	assert.NoError(t, err)
	populateEntityInventory(queueEntity, queueType, &data)
	l.AssertExpectations(t)

	queueActual, err := json.Marshal(queueEntity.Inventory)
	assert.NoError(t, err)

	if *update {
		t.Log("Writing .golden file")
		err := ioutil.WriteFile(queueGoldenFile, queueActual, 0644)
		assert.NoError(t, err)
	}
	queueExpected, err := ioutil.ReadFile(queueGoldenFile)
	assert.NoError(t, err)

	assert.Equal(t, queueExpected, queueActual, "Excpected doesn't match .golden file %v", queueGoldenFile)
}
