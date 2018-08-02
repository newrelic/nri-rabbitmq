package inventory

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/testutils"
	"github.com/newrelic/nri-rabbitmq/utils"
	"github.com/newrelic/nri-rabbitmq/utils/consts"

	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func Test_CollectInventory(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	args.GlobalArgs = args.RabbitMQArguments{}
	prevExec := execCommand
	execCommand = fakeExecCommand
	os.Setenv("GO_WANT_HELPER_PROCESS", "1")
	defer func() {
		execCommand = prevExec
		os.Unsetenv("GO_WANT_HELPER_PROCESS")
	}()

	nodesData := []objx.Map{}
	assert.Panics(t, func() {
		CollectInventory(i, nodesData)
	}, "CollectInventory should fail with empty nodeData")

	nodesData = []objx.Map{
		objx.MSI("name", "node2"),
	}
	assert.Panics(t, func() {
		CollectInventory(i, nodesData)
	}, "CollectInventory should fail with mismatched nodeData")

	nodeData, _ := objx.FromJSON(fmt.Sprintf(`{
		"name": %q,
		"config_files": [
			%q
		]
	}`, expectedNodeName, testConfigPath))
	nodesData = []objx.Map{
		nodeData,
	}

	prevOsOpen := osOpen
	osOpen = errorOsOpen
	require.Panics(t, func() {
		CollectInventory(i, nodesData)
	})
	osOpen = prevOsOpen

	CollectInventory(i, nodesData)
	assert.Equal(t, 1, len(i.Entities), "CollectInventory should create one Entity")
	actual, _ := i.Entities[0].Inventory.MarshalJSON()

	golden := testConfigPath + ".golden"
	if *testutils.Update {
		t.Log("Writing .golden file")
		ioutil.WriteFile(golden, actual, 0644)
	}

	expected, _ := ioutil.ReadFile(golden)
	assert.Equal(t, expected, actual, "CollectInventory does not have the expected output")
}

func Test_getLocalNodeName(t *testing.T) {
	testutils.SetTestLogger(t)
	args.GlobalArgs = args.RabbitMQArguments{
		NodeNameOverride: expectedNodeName,
	}

	nodeName, err := getLocalNodeName()
	assert.NoError(t, err)
	assert.Equal(t, expectedNodeName, nodeName)
	args.GlobalArgs.NodeNameOverride = ""

	prevExec := execCommand
	execCommand = fakeExecCommand
	os.Setenv("GO_WANT_HELPER_PROCESS", "1")
	defer func() {
		execCommand = prevExec
		os.Unsetenv("GO_WANT_HELPER_PROCESS")
	}()

	os.Setenv("GET_NODE_NAME_ERROR", "1")
	_, err = getLocalNodeName()
	assert.Error(t, err)
	os.Unsetenv("GET_NODE_NAME_ERROR")

	os.Setenv("GET_NODE_NAME_EMPTY", "1")
	_, err = getLocalNodeName()
	assert.EqualError(t, err, "could not determine the local node name")
	os.Unsetenv("GET_NODE_NAME_EMPTY")

	nodeName, err = getLocalNodeName()
	assert.NoError(t, err)
	assert.Equal(t, expectedNodeName, nodeName)
}

func Test_getNodeEntity(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	nodeEntity, actualNodeData, err := getNodeEntity("node1", nil, nil)
	assert.EqualError(t, err, "node name [node1] not found in cluster")

	nodeData := []objx.Map{
		objx.MSI("name", "node1"),
		objx.MSI("name", "node2"),
	}

	nodeEntity, actualNodeData, err = getNodeEntity("node1", nodeData, i)
	assert.NoError(t, err)
	assert.NotNil(t, nodeEntity)
	assert.Equal(t, nodeData[0], actualNodeData)
	assert.Equal(t, "node1", nodeEntity.Metadata.Name)
	assert.Equal(t, consts.NodeType, nodeEntity.Metadata.Namespace)

	nodeEntity, actualNodeData, err = getNodeEntity("node2", nodeData, i)
	assert.NoError(t, err)
	assert.NotNil(t, nodeEntity)
	assert.Equal(t, nodeData[1], actualNodeData)
	assert.Equal(t, "node2", nodeEntity.Metadata.Name)
	assert.Equal(t, consts.NodeType, nodeEntity.Metadata.Namespace)
}

func TestConfigNotExist(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{
		ConfigPath: filepath.Join("testdata", "file-not_found.config"),
	}
	_, e := testutils.GetTestingEntity(t)
	err := setInventoryData(e, nil)
	assert.Nil(t, err)
}

func TestConfigOpenError(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{
		ConfigPath: filepath.Join("testdata"),
	}
	_, e := testutils.GetTestingEntity(t)
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
	args.GlobalArgs = args.RabbitMQArguments{}
	i, e := testutils.GetTestingEntity(t)

	err := setInventoryData(e, nil)
	assert.NoError(t, err)

	args.GlobalArgs.ConfigPath = testConfigPath
	golden := args.GlobalArgs.ConfigPath + ".golden"
	err = setInventoryData(e, nil)
	require.NoError(t, err)

	require.Equal(t, 1, len(i.Entities))
	actual, err := i.Entities[0].Inventory.MarshalJSON()
	require.NoError(t, err)

	if *testutils.Update {
		t.Log("Writing .golden file")
		ioutil.WriteFile(golden, actual, 0644)
	}

	expected, _ := ioutil.ReadFile(golden)

	if !bytes.Equal(expected, actual) {
		t.Errorf("Actual JSON results do not match expected .golden file for %s", args.GlobalArgs.ConfigPath)
	}
}

func TestGetConfigPath(t *testing.T) {
	testutils.SetTestLogger(t)
	args.GlobalArgs = args.RabbitMQArguments{
		ConfigPath: testConfigPath,
	}

	actual := getConfigPath(nil)
	assert.Equal(t, testConfigPath, actual)
	args.GlobalArgs.ConfigPath = ""

	actual = getConfigPath(nil)
	assert.Empty(t, actual)

	nodeData, _ := objx.FromJSON(`{
		"config_files": [
			"/etc/rabbitmq/rabbitmq.config",
			"/etc/rabbitmq/advanced.config"
		]
	}`)
	actual = getConfigPath(nodeData)
	assert.Empty(t, actual)

	nodeData, _ = objx.FromJSON(`{
		"config_files": [
			"/etc/rabbitmq/rabbitmq.conf",
			"/etc/rabbitmq/advanced.config"
		]
	}`)
	actual = getConfigPath(nodeData)
	assert.Equal(t, "/etc/rabbitmq/rabbitmq.conf", actual)
}

func Test_PopulateEntityInventory(t *testing.T) {
	sourceFile := filepath.Join("testdata", "populateEntityInventory.json")
	exchangeGoldenFile := sourceFile + ".exchange.golden"
	queueGoldenFile := sourceFile + ".queue.golden"
	data := objx.New(testutils.ReadObjectFromJSONFile(t, sourceFile))

	i, exchangeEntity := testutils.GetTestingEntity(t, "my-exchange", consts.ExchangeType)
	l := testutils.SetMockLogger()
	l.On("Infof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	// exchange test
	PopulateEntityInventory(exchangeEntity, consts.ExchangeType, &data)
	setInventoryItem(exchangeEntity, "this-key-is-longer-than-375-to-force-an-error-lorem-ipsum-dolor-sit-amet-consectetur-adipiscing-elit-sed-do-eiusmod-tempor-incididunt-ut-labore-et-dolore-magna-aliqua-ut-enim-ad-minim-veniam-quis-nostrud-exercitation-ullamco-laboris-nisi-ut-aliquip-ex-ea-commodo-consequat-duis-aute-irure-dolor-in-reprehenderit-in-voluptate-velit-esse-cillum-dolore-eu-fugiat-nulla-pariatures", "nope", "false")

	exchangeActual, err := exchangeEntity.Inventory.MarshalJSON()
	assert.NoError(t, err)

	if *testutils.Update {
		t.Log("Writing .golden file")
		ioutil.WriteFile(exchangeGoldenFile, exchangeActual, 0644)
	}
	exchangeExpected, _ := ioutil.ReadFile(exchangeGoldenFile)

	assert.Equal(t, exchangeExpected, exchangeActual, "Expected doesn't match .golden file %v", exchangeGoldenFile)

	// queue test
	queueEntity, _, err := utils.CreateEntity(i, "my-queue", consts.QueueType, "/")
	assert.NoError(t, err)
	PopulateEntityInventory(queueEntity, consts.QueueType, &data)
	l.AssertExpectations(t)

	queueActual, err := json.Marshal(queueEntity.Inventory)
	assert.NoError(t, err)

	if *testutils.Update {
		t.Log("Writing .golden file")
		ioutil.WriteFile(queueGoldenFile, queueActual, 0644)
	}
	queueExpected, _ := ioutil.ReadFile(queueGoldenFile)

	assert.Equal(t, queueExpected, queueActual, "Excpected doesn't match .golden file %v", queueGoldenFile)
}

func fakeExecCommand(command string, args ...string) (cmd *exec.Cmd) {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd = exec.Command(os.Args[0], cs...)
	return cmd
}

// TestHelperProcess isn't a real test. It's used as a helper process.
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "rabbitmqctl":
		if len(args) == 2 && args[0] == "eval" && args[1] == "node()." {
			if os.Getenv("GET_NODE_NAME_ERROR") == "1" {
				os.Exit(2)
			}
			if os.Getenv("GET_NODE_NAME_EMPTY") == "1" {
				fmt.Fprintf(os.Stdout, "")
			} else {
				fmt.Fprintf(os.Stdout, `'node1'
`)
			}
		}
	}
}
