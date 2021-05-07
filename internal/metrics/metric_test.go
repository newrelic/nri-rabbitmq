package metrics

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/log"

	"github.com/newrelic/nri-rabbitmq/internal/data/consts"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/persist"
	"github.com/newrelic/nri-rabbitmq/internal/data"
	"github.com/newrelic/nri-rabbitmq/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCollectEntityMetrics(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	CollectEntityMetrics(i, []*data.BindingData{},
		"testClusterName",
		&data.NodeData{Name: "node1"},
		&data.QueueData{Name: "queue1"},
	)
	assert.Equal(t, 2, len(i.Entities))
}

func Test_setMetric(t *testing.T) {
	ms := metric.NewSet("TestSample", persist.NewInMemoryStore())
	setMetric(ms, "rate", 0.5, metric.GAUGE)
	assert.Equal(t, float64(0.5), ms.Metrics["rate"])
}

func TestCollectEntityMetrics_Node(t *testing.T) {
	var bindingData []*data.BindingData
	var nodeData []*data.NodeData
	i := testutils.GetTestingIntegration(t)

	sourceFile := filepath.Join("testdata", "populateMetricsTest.node.json")
	testutils.ReadStructFromJSONFile(t, sourceFile, &nodeData)

	entityData := make([]data.EntityData, len(nodeData))
	for i, v := range nodeData {
		entityData[i] = v
	}

	CollectEntityMetrics(i, bindingData, "testClusterName", entityData...)

	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils.Update {
			if err := ioutil.WriteFile(goldenFile, actual, 0o644); err != nil {
				log.Error(err.Error())
			}
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}

func TestCollectEntityMetrics_Queue(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	var queueData []*data.QueueData
	bindingData := []*data.BindingData{
		{
			Vhost:           "test-vhost",
			Source:          "exchange1",
			Destination:     "test-name",
			DestinationType: consts.QueueType,
		},
		{
			Vhost:           "test-vhost",
			Source:          "exchange2",
			Destination:     "test-name",
			DestinationType: consts.QueueType,
		},
	}

	sourceFile := filepath.Join("testdata", "populateMetricsTest.queue.json")
	testutils.ReadStructFromJSONFile(t, sourceFile, &queueData)

	entityData := make([]data.EntityData, len(queueData))
	for i, v := range queueData {
		entityData[i] = v
	}

	CollectEntityMetrics(i, bindingData, "testClusterName", entityData...)

	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils.Update {
			if err := ioutil.WriteFile(goldenFile, actual, 0o644); err != nil {
				log.Error(err.Error())
			}
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}

func TestCollectEntityMetrics_Exchange(t *testing.T) {
	var bindingData []*data.BindingData
	var exchangeData []*data.ExchangeData
	i := testutils.GetTestingIntegration(t)

	sourceFile := filepath.Join("testdata", "populateMetricsTest.exchange.json")
	testutils.ReadStructFromJSONFile(t, sourceFile, &exchangeData)

	entityData := make([]data.EntityData, len(exchangeData))
	for i, v := range exchangeData {
		entityData[i] = v
	}

	CollectEntityMetrics(i, bindingData, "testClusterName", entityData...)

	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils.Update {
			if err := ioutil.WriteFile(goldenFile, actual, 0o644); err != nil {
				log.Error(err.Error())
			}
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}

func TestCollectVhostMetrics(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	var vhostData []*data.VhostData
	var connectionsData []*data.ConnectionData
	sourceFile := filepath.Join("testdata", "populateMetricsTest.vhost.json")
	testutils.ReadStructFromJSONFile(t, sourceFile, &vhostData)
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "populateMetricsTest.connections.json"), &connectionsData)

	CollectVhostMetrics(i, vhostData, connectionsData, "testClusterName")
	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils.Update {
			if err := ioutil.WriteFile(goldenFile, actual, 0o644); err != nil {
				log.Error(err.Error())
			}
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}
