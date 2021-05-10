package main

import (
	"github.com/newrelic/nri-rabbitmq/src/args"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
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

func Test_alivenessTest_SkipCollect(t *testing.T) {
	i := testutils.GetTestingIntegration(t)

	argList := args.ArgumentList{
		Exchanges: "[\"test1\"]",
		Queues:    "[\"test1\"]",
		Vhosts:    "[\"test1\"]",
	}
	err := args.SetGlobalArgs(argList)
	assert.Nil(t, err)

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
	assert.Equal(t, 0, len(i.Entities))
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
