package main

import (
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/stretchr/testify/assert"

	"github.com/newrelic/nri-rabbitmq/src/testutils"
)

func Test_alivenessTest_Pass(t *testing.T) {
	i := testutils.GetTestingIntegration(t)

	vhostTests := []*data.VhostTest{
		{
			Vhost: &data.VhostData{Name: "vhost1"},
			Test: &data.TestData{
				Status: "ok",
			},
		},
	}
	alivenessTest(i, vhostTests, "testClusterName")
	assert.Empty(t, i.Entities)
}

func Test_alivenessTest_FailCreateEntity(t *testing.T) {
	i := testutils.GetTestingIntegration(t)

	vhostTests := []*data.VhostTest{
		{
			Vhost: &data.VhostData{},
			Test:  &data.TestData{},
		},
	}
	alivenessTest(i, vhostTests, "testClusterName")
	assert.Empty(t, i.Entities)
}

func Test_alivenessTest_FailAliveness(t *testing.T) {
	i := testutils.GetTestingIntegration(t)

	vhostTests := []*data.VhostTest{
		{
			Vhost: &data.VhostData{Name: "vhost1"},
			Test: &data.TestData{
				Status: "failed",
				Reason: "nodedown",
			},
		},
	}
	alivenessTest(i, vhostTests, "testClusterName")
	assert.Equal(t, 1, len(i.Entities))
	assert.Equal(t, 1, len(i.Entities[0].Events))
}

func Test_healthcheckTest_Pass(t *testing.T) {
	i := testutils.GetTestingIntegration(t)

	nodeTests := []*data.NodeTest{
		{
			Node: &data.NodeData{Name: "node1"},
			Test: &data.TestData{
				Status: "ok",
			},
		},
	}
	healthcheckTest(i, nodeTests, "testClusterName")
	assert.Empty(t, i.Entities)
}

func Test_healthcheckTest_FailCreateEntity(t *testing.T) {
	i := testutils.GetTestingIntegration(t)

	nodeTests := []*data.NodeTest{
		{
			Node: &data.NodeData{},
			Test: &data.TestData{},
		},
	}
	healthcheckTest(i, nodeTests, "testClusterName")
	assert.Empty(t, i.Entities)
}

func Test_healthcheckTest_FailAliveness(t *testing.T) {
	i := testutils.GetTestingIntegration(t)

	nodeTests := []*data.NodeTest{
		{
			Node: &data.NodeData{Name: "vhost1"},
			Test: &data.TestData{
				Status: "failed",
				Reason: "nodedown",
			},
		},
	}
	healthcheckTest(i, nodeTests, "testClusterName")
	assert.Equal(t, 1, len(i.Entities))
	assert.Equal(t, 1, len(i.Entities[0].Events))
}
