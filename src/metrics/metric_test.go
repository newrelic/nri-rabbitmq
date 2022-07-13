package metrics

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	data2 "github.com/newrelic/nri-rabbitmq/src/data"
	consts2 "github.com/newrelic/nri-rabbitmq/src/data/consts"
	testutils2 "github.com/newrelic/nri-rabbitmq/src/testutils"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/persist"
	"github.com/stretchr/testify/assert"
)

func TestCollectEntityMetrics(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)
	CollectEntityMetrics(i, []*data2.BindingData{},
		"testClusterName",
		&data2.NodeData{Name: "node1"},
		&data2.QueueData{Name: "queue1"},
	)
	assert.Equal(t, 2, len(i.Entities))
}

func Test_setMetric(t *testing.T) {
	ms := metric.NewSet("TestSample", persist.NewInMemoryStore())
	setMetric(ms, "rate", 0.5, metric.GAUGE)
	assert.Equal(t, float64(0.5), ms.Metrics["rate"])
}

func TestCollectEntityMetrics_Node(t *testing.T) {
	var bindingData []*data2.BindingData
	var nodeData []*data2.NodeData
	i := testutils2.GetTestingIntegration(t)

	sourceFile := filepath.Join("testdata", "populateMetricsTest.node.json")
	testutils2.ReadStructFromJSONFile(t, sourceFile, &nodeData)

	entityData := make([]data2.EntityData, len(nodeData))
	for i, v := range nodeData {
		entityData[i] = v
	}

	CollectEntityMetrics(i, bindingData, "testClusterName", entityData...)

	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils2.Update {
			if err := ioutil.WriteFile(goldenFile, actual, 0o644); err != nil {
				log.Error(err.Error())
			}
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}

func TestCollectEntityMetrics_Queue(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)
	var queueData []*data2.QueueData
	bindingData := []*data2.BindingData{
		{
			Vhost:           "test-vhost",
			Source:          "exchange1",
			Destination:     "test-name",
			DestinationType: consts2.QueueType,
		},
		{
			Vhost:           "test-vhost",
			Source:          "exchange2",
			Destination:     "test-name",
			DestinationType: consts2.QueueType,
		},
	}

	sourceFile := filepath.Join("testdata", "populateMetricsTest.queue.json")
	testutils2.ReadStructFromJSONFile(t, sourceFile, &queueData)

	entityData := make([]data2.EntityData, len(queueData))
	for i, v := range queueData {
		entityData[i] = v
	}

	CollectEntityMetrics(i, bindingData, "testClusterName", entityData...)

	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils2.Update {
			if err := ioutil.WriteFile(goldenFile, actual, 0o644); err != nil {
				log.Error(err.Error())
			}
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}

func TestCollectEntityMetrics_Exchange(t *testing.T) {
	var bindingData []*data2.BindingData
	var exchangeData []*data2.ExchangeData
	i := testutils2.GetTestingIntegration(t)

	sourceFile := filepath.Join("testdata", "populateMetricsTest.exchange.json")
	testutils2.ReadStructFromJSONFile(t, sourceFile, &exchangeData)

	entityData := make([]data2.EntityData, len(exchangeData))
	for i, v := range exchangeData {
		entityData[i] = v
	}

	CollectEntityMetrics(i, bindingData, "testClusterName", entityData...)

	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils2.Update {
			if err := ioutil.WriteFile(goldenFile, actual, 0o644); err != nil {
				log.Error(err.Error())
			}
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}

func TestCollectVhostMetrics(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)
	var vhostData []*data2.VhostData
	var connectionsData []*data2.ConnectionData
	sourceFile := filepath.Join("testdata", "populateMetricsTest.vhost.json")
	testutils2.ReadStructFromJSONFile(t, sourceFile, &vhostData)
	testutils2.ReadStructFromJSONFile(t, filepath.Join("testdata", "populateMetricsTest.connections.json"), &connectionsData)

	CollectVhostMetrics(i, vhostData, connectionsData, "testClusterName")
	if assert.Equal(t, 1, len(i.Entities)) && assert.Equal(t, 1, len(i.Entities[0].Metrics)) {
		goldenFile := sourceFile + ".golden"
		actual, _ := i.Entities[0].Metrics[0].MarshalJSON()
		if *testutils2.Update {
			if err := ioutil.WriteFile(goldenFile, actual, 0o644); err != nil {
				log.Error(err.Error())
			}
		}
		expected, _ := ioutil.ReadFile(goldenFile)
		assert.Equal(t, string(expected), string(actual))
	}
}
