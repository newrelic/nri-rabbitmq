package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/mock"
)

func TestSetMetric(t *testing.T) {
	l := new(mockedLogger)
	i, _ := integration.New("name", "version", integration.Logger(l))
	logger = l
	l.On("Errorf", mock.Anything, mock.Anything).Once()

	firstQueue, _ := i.Entity("first-queue", "/my-vhost/queue")
	metrics := firstQueue.NewMetricSet("RabbitMQSample")
	setMetric(metrics, "rate", 0.5, metric.RATE)

	l.AssertExpectations(t)
}

func TestPopulateMetrics(t *testing.T) {
	logger = log.NewStdErr(true)
	actualMetricSet := metric.NewSet("testMetrics", nil)

	responseString, _ := readFile("testdata/populateMetricsTest.json")
	responseObject, _ := objx.FromJSON(responseString)

	populateMetrics(actualMetricSet, "queue", &responseObject)

	actual, _ := json.Marshal(actualMetricSet)
	if *update {
		ioutil.WriteFile("testdata/populateMetricsTest.json.golden", actual, 0644)
	}
	expected, _ := ioutil.ReadFile("testdata/populateMetricsTest.json.golden")

	if !bytes.Equal(expected, actual) {
		t.Errorf("Actual JSON results do not match expected .golden file for %s", args.ConfigPath)
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
	logger = &testLogger{t.Logf}
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

	actual, _ := json.Marshal(actualVhostMetricSet)
	if *update {
		ioutil.WriteFile("testdata/populateVhostMetricsTest.json.golden", actual, 0644)
	}
	expected, _ := ioutil.ReadFile("testdata/populateVhostMetricsTest.json.golden")

	if !bytes.Equal(expected, actual) {
		t.Errorf("Actual JSON results do not match expected .golden file for %s", args.ConfigPath)
	}
}

func TestPopulateBindingMetrics(t *testing.T) {
	logger = log.NewStdErr(true)
	actualBindingMetricSet := metric.NewSet("testBindingMetrics", nil)

	bindingKeyStruct := bindingKey{
		Vhost:      "/",
		EntityName: "test-queue",
		EntityType: queueType,
	}
	var bindingStats = map[bindingKey]int{}
	bindingStats[bindingKeyStruct] = 7

	populateBindingMetric(bindingKeyStruct.EntityName, bindingKeyStruct.Vhost, bindingKeyStruct.EntityType, actualBindingMetricSet, bindingStats)

	actual, _ := json.Marshal(actualBindingMetricSet)
	if *update {
		ioutil.WriteFile("testdata/populateBindingMetricTest.json.golden", actual, 0644)
	}
	expected, _ := ioutil.ReadFile("testdata/populateBindingMetricTest.json.golden")

	if !bytes.Equal(expected, actual) {
		t.Errorf("Actual JSON results do not match expected .golden file for %s", args.ConfigPath)
	}
}

func TestParseJsonFloat(t *testing.T) {
	jsonString, _ := readFile("testdata/parseJsonTest.json")
	jsonObject, _ := objx.FromJSON(jsonString)
	actual, _ := parseJSON(&jsonObject, "float-test")
	assert.IsType(t, *new(float64), actual)
}

func TestParseJsonInt(t *testing.T) {
	jsonString, _ := readFile("testdata/parseJsonTest.json")
	jsonObject, _ := objx.FromJSON(jsonString)
	_, err := parseJSON(&jsonObject, "int-test")
	assert.Error(t, err, "output should have an error")
}

func TestParseJsonBool(t *testing.T) {
	jsonString, _ := readFile("testdata/parseJsonTest.json")
	jsonObject, _ := objx.FromJSON(jsonString)
	actual, _ := parseJSON(&jsonObject, "bool-test")
	assert.IsType(t, *new(int), actual)
}

func TestConvertBoolToIntTrue(t *testing.T) {
	actual := convertBoolToInt(true)
	expected := 1
	assert.Equal(t, expected, actual)
}

func TestConvertBoolToIntFalse(t *testing.T) {
	actual := convertBoolToInt(false)
	expected := 0
	assert.Equal(t, expected, actual)
}
