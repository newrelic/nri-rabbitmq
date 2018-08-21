package data

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/data/consts"
	"github.com/newrelic/nri-rabbitmq/src/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeData(t *testing.T) {
	var nodeData NodeData
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "node.json"), &nodeData)
	assert.NotNil(t, nodeData)
	assert.Equal(t, 2, len(nodeData.ConfigFiles))
	assert.Contains(t, nodeData.ConfigFiles, "advanced.config")
	assert.Contains(t, nodeData.ConfigFiles, "rabbit.conf")
	assert.Equal(t, getInt64(1024), nodeData.DiskFreeSpace)
	assert.Equal(t, getInt(0), nodeData.DiskAlarm)
	assert.Equal(t, getInt(0), nodeData.MemoryAlarm)
	assert.Equal(t, getInt64(20), nodeData.FileDescriptorsUsed)
	assert.Equal(t, getInt64(2048), nodeData.MemoryUsed)
	assert.Equal(t, "node1", nodeData.Name)
	assert.Equal(t, 2, nodeData.Partitions)
	assert.Equal(t, getInt64(3), nodeData.RunQueue)
	assert.Equal(t, getInt(1), nodeData.Running)
	assert.Equal(t, getInt64(2), nodeData.SocketsUsed)
	assert.Equal(t, "node1", nodeData.EntityName())
	assert.Equal(t, consts.NodeType, nodeData.EntityType())
	assert.Equal(t, "", nodeData.EntityVhost())

	testIntegration := testutils.GetTestingIntegration(t)
	e, metricAttribs, err := nodeData.GetEntity(testIntegration)
	assert.NotNil(t, e)
	assert.NotEmpty(t, metricAttribs)
	assert.NoError(t, err)

	ms := e.NewMetricSet("TestSample", metricAttribs...)

	assert.NoError(t, ms.MarshalMetrics(nodeData))
	assert.Equal(t, float64(20), ms.Metrics["node.fileDescriptorsTotalUsed"])
	assert.Equal(t, float64(1024), ms.Metrics["node.diskSpaceFreeInBytes"])
	assert.Equal(t, float64(2048), ms.Metrics["node.totalMemoryUsedInBytes"])
	assert.Equal(t, float64(3), ms.Metrics["node.averageErlangProcessesWaiting"])
	assert.Equal(t, float64(2), ms.Metrics["node.fileDescriptorsUsedSockets"])
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
