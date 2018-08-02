package metrics

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/newrelic/nri-rabbitmq/client"
	"github.com/newrelic/nri-rabbitmq/consts"
	"github.com/newrelic/nri-rabbitmq/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/objx"
)

func TestCollectMetrics(t *testing.T) {
	testutils.SetTestLogger(t)
	apiResponses := objx.MSI(
		client.VhostsEndpoint, []objx.Map{
			objx.MSI("name", "vhost1"),
		},
		client.QueuesEndpoint, []objx.Map{
			objx.MSI("name", "queue1"),
		},
	)
	i := testutils.GetTestingIntegration(t)
	CollectMetrics(i, &apiResponses)
	assert.Equal(t, 2, len(i.Entities))
}

func TestSetMetric(t *testing.T) {
	l := testutils.SetMockLogger()
	defer func() {
		testutils.SetTestLogger(t)
	}()

	i, _ := integration.New("name", "version", integration.Logger(l))
	l.On("Errorf", mock.Anything, mock.Anything).Once()

	firstQueue, _ := i.Entity("first-queue", "/my-vhost/queue")
	metrics := firstQueue.NewMetricSet("RabbitMQSample")
	setMetric(metrics, "rate", 0.5, metric.RATE)

	l.AssertExpectations(t)
}

func TestPopulateMetrics(t *testing.T) {
	testutils.SetTestLogger(t)
	args.GlobalArgs = args.RabbitMQArguments{}
	actualMetricSet := metric.NewSet("testMetrics", nil)

	sourceFile := filepath.Join("testdata", "populateMetricsTest.json")
	responseString, _ := readFile(sourceFile)
	responseObject, _ := objx.FromJSON(responseString)

	populateMetrics(actualMetricSet, "queue", &responseObject)

	goldenFile := sourceFile + ".golden"
	actual, _ := json.Marshal(actualMetricSet)
	if *testutils.Update {
		ioutil.WriteFile(goldenFile, actual, 0644)
	}
	expected, _ := ioutil.ReadFile(goldenFile)

	if !bytes.Equal(expected, actual) {
		t.Errorf("Actual JSON results do not match expected .golden file for %s", goldenFile)
	}
}

func readFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func TestPopulateVhostMetrics(t *testing.T) {
	testutils.SetTestLogger(t)
	actualVhostMetricSet := metric.NewSet("testVhostMetrics", nil)

	connKeyStruct := connKey{
		Vhost: "/",
		State: "starting",
	}
	var connStats = map[connKey]int{}
	connStats[connKeyStruct] = 7

	connKeyStruct.State = "flow"
	connStats[connKeyStruct] = 2

	connKeyStruct.State = "total"
	connStats[connKeyStruct] = 9

	populateVhostMetrics(connKeyStruct.Vhost, actualVhostMetricSet, connStats)

	goldenFile := filepath.Join("testdata", "populateVhostMetricsTest.json.golden")
	actual, _ := json.Marshal(actualVhostMetricSet)
	if *testutils.Update {
		ioutil.WriteFile(goldenFile, actual, 0644)
	}
	expected, _ := ioutil.ReadFile(goldenFile)

	if !bytes.Equal(expected, actual) {
		t.Errorf("Actual JSON results do not match expected .golden file for %s", goldenFile)
	}
}

func TestPopulateBindingMetrics(t *testing.T) {
	testutils.SetTestLogger(t)
	actualBindingMetricSet := metric.NewSet("testBindingMetrics", nil)

	bindingKeyStruct := bindingKey{
		Vhost:      "/",
		EntityName: "test-queue",
		EntityType: consts.QueueType,
	}
	var bindingStats = map[bindingKey]int{}
	bindingStats[bindingKeyStruct] = 7

	populateBindingMetric(bindingKeyStruct.EntityName, bindingKeyStruct.Vhost, bindingKeyStruct.EntityType, actualBindingMetricSet, bindingStats)

	goldenFile := filepath.Join("testdata", "populateBindingMetricTest.json.golden")
	actual, _ := json.Marshal(actualBindingMetricSet)
	if *testutils.Update {
		ioutil.WriteFile(goldenFile, actual, 0644)
	}
	expected, _ := ioutil.ReadFile(goldenFile)

	if !bytes.Equal(expected, actual) {
		t.Errorf("Actual JSON results do not match expected .golden file for %s", goldenFile)
	}
}

func TestParseJsonFloat(t *testing.T) {
	testutils.SetTestLogger(t)
	jsonString, _ := readFile(filepath.Join("testdata", "parseJsonTest.json"))
	jsonObject, _ := objx.FromJSON(jsonString)
	actual, _ := parseJSON(&jsonObject, "float-test")
	assert.IsType(t, *new(float64), actual)
}

func TestParseJsonInt(t *testing.T) {
	testutils.SetTestLogger(t)
	jsonString, _ := readFile(filepath.Join("testdata", "parseJsonTest.json"))
	jsonObject, _ := objx.FromJSON(jsonString)
	_, err := parseJSON(&jsonObject, "int-test")
	assert.Error(t, err, "output should have an error")
}

func TestParseJsonBool(t *testing.T) {
	testutils.SetTestLogger(t)
	jsonString, _ := readFile(filepath.Join("testdata", "parseJsonTest.json"))
	jsonObject, _ := objx.FromJSON(jsonString)
	actual, _ := parseJSON(&jsonObject, "bool-test")
	assert.IsType(t, *new(int), actual)
}
