package data

import (
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_OverviewData_JSON(t *testing.T) {
	var overviewData *OverviewData
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "overview.json"), &overviewData)
	assert.NotNil(t, overviewData)
	assert.Equal(t, "cluster1", overviewData.ClusterName)
	assert.Equal(t, "1.0.1", overviewData.RabbitMQVersion)
	assert.Equal(t, "2.0.2", overviewData.ManagementVersion)
}

func Test_ConnectionData_JSON(t *testing.T) {
	var connectionData []*ConnectionData
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "connections.json"), &connectionData)
	assert.Equal(t, 1, len(connectionData))
	assert.Equal(t, "running", connectionData[0].State)
	assert.Equal(t, "vhost1", connectionData[0].Vhost)
}

func Test_BindingData_JSON(t *testing.T) {
	var bindingData []*BindingData
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "bindings.json"), &bindingData)
	assert.Equal(t, 1, len(bindingData))
	assert.Equal(t, "source1", bindingData[0].Source)
	assert.Equal(t, "vhost1", bindingData[0].Vhost)
	assert.Equal(t, "dest1", bindingData[0].Destination)
	assert.Equal(t, "queue", bindingData[0].DestinationType)
}

func Test_VhostData_JSON(t *testing.T) {
	var vhostData []*VhostData
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "vhosts.json"), &vhostData)
	assert.Equal(t, 1, len(vhostData))
	assert.Equal(t, "vhost1", vhostData[0].Name)
}
