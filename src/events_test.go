package main

import (
	"github.com/newrelic/nri-rabbitmq/src"
	args2 "github.com/newrelic/nri-rabbitmq/src/args"
	data2 "github.com/newrelic/nri-rabbitmq/src/data"
	testutils2 "github.com/newrelic/nri-rabbitmq/src/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_alivenessTest_Pass(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)

	vhostTests := []*data2.VhostTest{
		{
			Vhost: &data2.VhostData{Name: "vhost1"},
			Test: &data2.TestData{
				Status: "ok",
			},
		},
	}
	src.alivenessTest(i, vhostTests, "testClusterName")
	assert.Empty(t, i.Entities)
}

func Test_alivenessTest_FailCreateEntity(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)

	vhostTests := []*data2.VhostTest{
		{
			Vhost: &data2.VhostData{},
			Test:  &data2.TestData{},
		},
	}
	src.alivenessTest(i, vhostTests, "testClusterName")
	assert.Empty(t, i.Entities)
}

func Test_alivenessTest_FailAliveness(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)

	vhostTests := []*data2.VhostTest{
		{
			Vhost: &data2.VhostData{Name: "vhost1"},
			Test: &data2.TestData{
				Status: "failed",
				Reason: "nodedown",
			},
		},
	}
	src.alivenessTest(i, vhostTests, "testClusterName")
	assert.Equal(t, 1, len(i.Entities))
	assert.Equal(t, 1, len(i.Entities[0].Events))
}

func Test_alivenessTest_SkipCollect(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)

	argList := args2.ArgumentList{
		Exchanges: "[\"test1\"]",
		Queues:    "[\"test1\"]",
		Vhosts:    "[\"test1\"]",
	}
	err := args2.SetGlobalArgs(argList)
	assert.Nil(t, err)

	vhostTests := []*data2.VhostTest{
		{
			Vhost: &data2.VhostData{Name: "vhost1"},
			Test: &data2.TestData{
				Status: "failed",
				Reason: "nodedown",
			},
		},
	}
	src.alivenessTest(i, vhostTests, "testClusterName")
	assert.Equal(t, 0, len(i.Entities))
}

func Test_healthcheckTest_Pass(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)

	nodeTests := []*data2.NodeTest{
		{
			Node: &data2.NodeData{Name: "node1"},
			Test: &data2.TestData{
				Status: "ok",
			},
		},
	}
	src.healthcheckTest(i, nodeTests, "testClusterName")
	assert.Empty(t, i.Entities)
}

func Test_healthcheckTest_FailCreateEntity(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)

	nodeTests := []*data2.NodeTest{
		{
			Node: &data2.NodeData{},
			Test: &data2.TestData{},
		},
	}
	src.healthcheckTest(i, nodeTests, "testClusterName")
	assert.Empty(t, i.Entities)
}

func Test_healthcheckTest_FailAliveness(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)
	nodeTests := []*data2.NodeTest{
		{
			Node: &data2.NodeData{Name: "vhost1"},
			Test: &data2.TestData{
				Status: "failed",
				Reason: "nodedown",
			},
		},
	}
	src.healthcheckTest(i, nodeTests, "testClusterName")
	assert.Equal(t, 1, len(i.Entities))
	assert.Equal(t, 1, len(i.Entities[0].Events))
}
