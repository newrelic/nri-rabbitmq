package main

import (
	"os"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/args"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/testutils"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// TODO remove global args.
	// This test are heavily based on global args to filter entities on creation.
	args.GlobalArgs.Vhosts = []string{"vhost1"}

	os.Exit(m.Run())
}

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
			Vhost: &data.VhostData{Name: "bar"},
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
	running := true

	nodeTests := []*data.NodeData{
		{Name: "node1", Running: &running},
	}
	healthcheckTest(i, nodeTests, "testClusterName")
	assert.Empty(t, i.Entities)
}

func Test_healthcheckTest_NoRunningField(t *testing.T) {
	i := testutils.GetTestingIntegration(t)

	nodeTests := []*data.NodeData{
		{Name: "node1"},
	}
	healthcheckTest(i, nodeTests, "testClusterName")
	assert.Equal(t, 1, len(i.Entities))
	assert.Equal(t, 1, len(i.Entities[0].Events))
	assert.Contains(t, i.Entities[0].Events[0].Summary, RunningUnknown)
}

func Test_healthcheckTest_FailAliveness(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	running := false

	nodeTests := []*data.NodeData{
		{Name: "node1", Running: &running},
	}
	healthcheckTest(i, nodeTests, "testClusterName")
	assert.Equal(t, 1, len(i.Entities))
	assert.Equal(t, 1, len(i.Entities[0].Events))
	assert.Contains(t, i.Entities[0].Events[0].Summary, NotRunning)
}
