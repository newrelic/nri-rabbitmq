package metrics

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/data/consts"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/persist"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCollectMetrics(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	CollectEntityMetrics(i, []*data.BindingData{},
		&data.NodeData{Name: "node1"},
		&data.QueueData{Name: "queue1"},
	)
	assert.Equal(t, 2, len(i.Entities))
}

func TestSetMetric(t *testing.T) {
	ms := metric.NewSet("TestSample", persist.NewInMemoryStore())
	setMetric(ms, "rate", 0.5, metric.GAUGE)
	assert.Equal(t, float64(0.5), ms.Metrics["rate"])
}

func Test_PopulateMetrics_Node(t *testing.T) {
	var bindingData []*data.BindingData
	var nodeData []*data.NodeData
	i := testutils.GetTestingIntegration(t)

	sourceFile := filepath.Join("testdata", "populateMetricsTest.node.json")
	testutils.ReadStructFromJSONFile(t, sourceFile, &nodeData)

	entityData := make([]data.EntityData, len(nodeData))
	for i, v := range nodeData {
		entityData[i] = v
	}

	CollectEntityMetrics(i, bindingData, entityData...)

	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils.Update {
			ioutil.WriteFile(goldenFile, actual, 0644)
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}

func Test_PopulateMetrics_Queue(t *testing.T) {
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

	CollectEntityMetrics(i, bindingData, entityData...)

	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils.Update {
			ioutil.WriteFile(goldenFile, actual, 0644)
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}

func Test_PopulateMetrics_Exchange(t *testing.T) {
	var bindingData []*data.BindingData
	var exchangeData []*data.ExchangeData
	i := testutils.GetTestingIntegration(t)

	sourceFile := filepath.Join("testdata", "populateMetricsTest.exchange.json")
	testutils.ReadStructFromJSONFile(t, sourceFile, &exchangeData)

	entityData := make([]data.EntityData, len(exchangeData))
	for i, v := range exchangeData {
		entityData[i] = v
	}

	CollectEntityMetrics(i, bindingData, entityData...)

	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils.Update {
			ioutil.WriteFile(goldenFile, actual, 0644)
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}

func TestPopulateVhostMetrics(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	var vhostData []*data.VhostData
	var connectionsData []*data.ConnectionData
	sourceFile := filepath.Join("testdata", "populateMetricsTest.vhost.json")
	testutils.ReadStructFromJSONFile(t, sourceFile, &vhostData)
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "populateMetricsTest.connections.json"), &connectionsData)

	CollectVhostMetrics(i, vhostData, connectionsData)
	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils.Update {
			ioutil.WriteFile(goldenFile, actual, 0644)
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}
