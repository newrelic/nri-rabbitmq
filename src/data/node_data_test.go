package data

import (
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/data/consts"
	"github.com/newrelic/nri-rabbitmq/src/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_NodeData(t *testing.T) {
	var nodeData NodeData
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "node.json"), &nodeData)
	assert.NotNil(t, nodeData)
	assert.Equal(t, 2, len(nodeData.ConfigFiles))
	assert.Contains(t, nodeData.ConfigFiles, "advanced.config")
	assert.Contains(t, nodeData.ConfigFiles, "rabbit.conf")
	i64 := int64(1024)
	assert.Equal(t, &i64, nodeData.DiskFreeSpace)
	i := int(0)
	assert.Equal(t, &i, nodeData.DiskAlarm)
	assert.Equal(t, &i, nodeData.MemoryAlarm)
	i64 = int64(20)
	assert.Equal(t, &i64, nodeData.FileDescriptorsUsed)
	i64 = int64(2048)
	assert.Equal(t, &i64, nodeData.MemoryUsed)
	assert.Equal(t, "node1", nodeData.Name)
	assert.Equal(t, 2, nodeData.Partitions)
	i64 = int64(3)
	assert.Equal(t, &i64, nodeData.RunQueue)
	i = int(1)
	assert.Equal(t, &i, nodeData.Running)
	i64 = int64(2)
	assert.Equal(t, &i64, nodeData.SocketsUsed)
	assert.Equal(t, "node1", nodeData.EntityName())
	assert.Equal(t, consts.NodeType, nodeData.EntityType())

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
