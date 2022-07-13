package inventory

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/args"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/testutils"

	"github.com/newrelic/infra-integrations-sdk/log"
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

func TestCollectInventory(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	args.GlobalArgs = args.RabbitMQArguments{}

	var nodesData []*data.NodeData
	CollectInventory(i, nodesData, "testClusterName")
	assert.Empty(t, i.Entities, "CollectInventory shouldn't create anything with empty NodeData")

	prevExec := execCommand
	execCommand = fakeExecCommand
	os.Setenv("GO_WANT_HELPER_PROCESS", "1")
	defer func() {
		execCommand = prevExec
		os.Unsetenv("GO_WANT_HELPER_PROCESS")
	}()

	nodesData = []*data.NodeData{
		{
			Name: "node2",
		},
	}
	os.Setenv("GET_NODE_NAME_ERROR", "1")
	CollectInventory(i, nodesData, "testClusterName")
	assert.Empty(t, i.Entities, "CollectInventory shouldn't create anything with error getting node name")
	os.Unsetenv("GET_NODE_NAME_ERROR")

	CollectInventory(i, nodesData, "testClustreName")
	assert.Empty(t, i.Entities, "CollectInventory shouldn't create anything with mismatched nodeData")

	nodesData = []*data.NodeData{
		{
			Name:        expectedNodeName,
			ConfigFiles: []string{testConfigPath},
		},
	}

	prevOsOpen := osOpen
	osOpen = errorOsOpen
	CollectInventory(i, nodesData, "testClusterName")
	assert.NotEmpty(t, i.Entities, "CollectInventory should create nodeName when config file fails to open")
	osOpen = prevOsOpen

	CollectInventory(i, nodesData, "testClusterName")
	assert.Equal(t, 1, len(i.Entities), "CollectInventory should create one Entity")
	actual, _ := i.Entities[0].Inventory.MarshalJSON()

	golden := testConfigPath + ".golden"
	if *testutils.Update {
		t.Log("Writing .golden file")
		if err := ioutil.WriteFile(golden, actual, 0o644); err != nil {
			log.Error(err.Error())
		}
	}

	expected, _ := ioutil.ReadFile(golden)
	assert.Equal(t, expected, actual, "CollectInventory does not have the expected output")
}

func TestCollectInventory_Errors(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{
		NodeNameOverride: "node1",
	}
	i := testutils.GetTestingIntegration(t)

	nodesData := []*data.NodeData{
		{Name: "node2"},
	}

	CollectInventory(i, nodesData, "testClusterName")
}

func Test_getLocalNodeName(t *testing.T) {
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

func Test_findNodeData(t *testing.T) {
	_, err := findNodeData("node1", nil)
	assert.EqualError(t, err, "node name [node1] not found in RabbitMQ")

	nodeData := []*data.NodeData{
		{Name: "node1"},
		{Name: "node2"},
	}

	actualNodeData, err := findNodeData("node1", nodeData)
	assert.NoError(t, err)
	assert.Equal(t, nodeData[0], actualNodeData)

	actualNodeData, err = findNodeData("node2", nodeData)
	assert.NoError(t, err)
	assert.Equal(t, nodeData[1], actualNodeData)
}

func Test_getConfigData_ConfigNotExist(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{
		ConfigPath: filepath.Join("testdata", "file-not_found.config"),
	}
	config := getConfigData(nil)
	assert.Empty(t, config)
}

func Test_getConfigData_ConfigOpenError(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{
		ConfigPath: filepath.Join("testdata"),
	}
	config := getConfigData(nil)
	assert.Empty(t, config)

	prevOpen := osOpen
	osOpen = errorOsOpen
	defer func() {
		osOpen = prevOpen
	}()

	config = getConfigData(nil)
	assert.Empty(t, config)
}

func Test_getConfigData(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}

	config := getConfigData(nil)
	assert.Empty(t, config)

	args.GlobalArgs.ConfigPath = testConfigPath
	config = getConfigData(nil)
	require.NotEmpty(t, config)
}

func Test_getConfigPath(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{
		ConfigPath: testConfigPath,
	}

	actual := getConfigPath(nil)
	assert.Equal(t, testConfigPath, actual)
	args.GlobalArgs.ConfigPath = ""

	actual = getConfigPath(nil)
	assert.Empty(t, actual)

	nodeData := new(data.NodeData)
	testutils.ReadStructFromJSONString(t, `{
		"config_files": [
			"/etc/rabbitmq/rabbitmq.config",
			"/etc/rabbitmq/advanced.config"
		]
	}`, nodeData)
	actual = getConfigPath(nodeData)
	assert.Empty(t, actual)

	testutils.ReadStructFromJSONString(t, `{
		"config_files": [
			"/etc/rabbitmq/rabbitmq.conf",
			"/etc/rabbitmq/advanced.config"
		]
	}`, nodeData)
	actual = getConfigPath(nodeData)
	assert.Equal(t, "/etc/rabbitmq/rabbitmq.conf", actual)
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
		fakeRabbitmqctl(args)
	}
}

func fakeRabbitmqctl(args []string) {
	if len(args) == 2 && args[0] == "eval" && args[1] == "node()." {
		if os.Getenv("GET_NODE_NAME_ERROR") == "1" {
			os.Exit(2)
		}
		if os.Getenv("GET_NODE_NAME_EMPTY") == "1" {
			fmt.Fprintf(os.Stdout, "")
		} else {
			// nolint
			fmt.Fprintf(os.Stdout, expectedNodeCmdOutput)
		}
	}
}
