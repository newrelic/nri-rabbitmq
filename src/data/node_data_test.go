package data

import (
	"encoding/json"
	consts2 "github.com/newrelic/nri-rabbitmq/src/data/consts"
	testutils2 "github.com/newrelic/nri-rabbitmq/src/testutils"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeData_UnmarshalJSON_MarshalMetrics(t *testing.T) {
	var nodeData NodeData
	testutils2.ReadStructFromJSONFile(t, filepath.Join("testdata", "node.json"), &nodeData)
	assert.NotNil(t, nodeData)
	assert.Equal(t, 2, len(nodeData.ConfigFiles))
	assert.Contains(t, nodeData.ConfigFiles, "advanced.config")
	assert.Contains(t, nodeData.ConfigFiles, "rabbit.conf")
	assert.Equal(t, getInt64(1024), nodeData.DiskFreeSpace)
	assert.Equal(t, getBool(false), nodeData.DiskAlarm)
	assert.Equal(t, getBool(false), nodeData.MemoryAlarm)
	assert.Equal(t, getInt64(20), nodeData.FileDescriptorsUsed)
	assert.Equal(t, getInt64(65436), nodeData.FileDescriptorsTotal)
	assert.Equal(t, getInt64(1048576), nodeData.ProcessesTotal)
	assert.Equal(t, getInt64(5180), nodeData.ProcessesUsed)
	assert.Equal(t, getInt64(2048), nodeData.MemoryUsed)
	assert.Equal(t, "node1", nodeData.Name)
	assert.Equal(t, 2, nodeData.Partitions)
	assert.Equal(t, getInt64(3), nodeData.RunQueue)
	assert.Equal(t, getBool(true), nodeData.Running)
	assert.Equal(t, getInt64(2), nodeData.SocketsUsed)
	assert.Equal(t, getInt64(58890), nodeData.SocketsTotal)
	assert.Equal(t, "node1", nodeData.EntityName())
	assert.Equal(t, consts2.NodeType, nodeData.EntityType())
	assert.Equal(t, "", nodeData.EntityVhost())

	testIntegration := testutils2.GetTestingIntegration(t)
	e, metricAttribs, err := nodeData.GetEntity(testIntegration, "testClusterName")
	assert.NotNil(t, e)
	assert.NotEmpty(t, metricAttribs)
	assert.NoError(t, err)

	ms := e.NewMetricSet("TestSample", metricAttribs...)

	assert.NoError(t, ms.MarshalMetrics(nodeData))
	assert.Equal(t, float64(20), ms.Metrics["node.fileDescriptorsTotalUsed"])
	assert.Equal(t, float64(65436), ms.Metrics["node.fileDescriptorsTotal"])
	assert.Equal(t, float64(1048576), ms.Metrics["node.processesTotal"])
	assert.Equal(t, float64(5180), ms.Metrics["node.processesUsed"])
	assert.Equal(t, float64(1024), ms.Metrics["node.diskSpaceFreeInBytes"])
	assert.Equal(t, float64(2048), ms.Metrics["node.totalMemoryUsedInBytes"])
	assert.Equal(t, float64(3), ms.Metrics["node.averageErlangProcessesWaiting"])
	assert.Equal(t, float64(2), ms.Metrics["node.fileDescriptorsUsedSockets"])
	assert.Equal(t, float64(58890), ms.Metrics["node.fileDescriptorsTotalSockets"])
	assert.Equal(t, float64(2), ms.Metrics["node.partitionsSeen"])
	assert.Equal(t, float64(1), ms.Metrics["node.running"])
	assert.Equal(t, float64(0), ms.Metrics["node.hostMemoryAlarm"])
	assert.Equal(t, float64(0), ms.Metrics["node.diskAlarm"])
}

func TestNodeData_JSONError(t *testing.T) {
	badJSONData := `{
		"name": "node1",
		"running": "true"
	}`
	var nodeData NodeData
	err := json.Unmarshal([]byte(badJSONData), &nodeData)
	require.Error(t, err)
}
