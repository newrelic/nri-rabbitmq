package main

import (
	"errors"
	"os"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func resetGlobalState() {
	args = argumentList{}
	logger = log.NewStdErr(false)
	queueFilter = make([]string, 0)
	queueRegexFilter = make([]string, 0)
	exchangeFilter = make([]string, 0)
	exchangeRegexFilter = make([]string, 0)
}

func TestMain(t *testing.T) {
	origArgs := os.Args
	os.Args = []string{
		"nr-rabbitmq",
		"-node_name_override", expectedNodeName,
		"-config_path", testConfigPath,
	}
	defer func() {
		os.Args = origArgs
	}()
	assert.NotPanics(t, func() {
		// TODO: don't test till we fully merge and get other http tests
		// main()
	})
}

func TestCheckArgsBadInventory(t *testing.T) {
	resetGlobalState()
	args.Metrics = true
	args.Inventory = false
	err := checkArgs()
	assert.Error(t, err, "should have an error for not having valid inventory")
}

func TestCheckArgsBadMetrics(t *testing.T) {
	resetGlobalState()
	args.Metrics = false
	args.Inventory = false
	args.Events = true
	err := checkArgs()
	assert.Error(t, err, "should have an error for not having anything specified to collect")
}

func TestCheckArgsGoodArgs(t *testing.T) {
	resetGlobalState()
	err := checkArgs()
	assert.NoError(t, err, "err should be nil")
}

func TestGetListOfEndpointsToCollectAllArgs(t *testing.T) {
	actual := getListOfEndpointsToCollect()
	expected := allEndpoints
	assert.Equal(t, expected, actual, "should publish all endpoints")
}

func TestGetListOfEndpointsToCollectNoEndpoints(t *testing.T) {
	resetGlobalState()
	args.Metrics = true
	actual := getListOfEndpointsToCollect()
	assert.Empty(t, actual, "should be empty slice")
}

func TestGetListOfEndpointsToCollectInventory(t *testing.T) {
	args.Inventory = true
	actual := getListOfEndpointsToCollect()
	expected := inventoryEndpoints
	assert.Equal(t, expected, actual, "should publish inventory endpoints")
}

func TestParseFilterArgsDefault(t *testing.T) {
	resetGlobalState()
	err := parseFilterArgs()
	assert.NoError(t, err, "err should be nil")
}

func TestParseFilterArgsBadArg(t *testing.T) {
	resetGlobalState()
	args.Queues = "bad-argument"
	err := parseFilterArgs()
	assert.Error(t, err, "should have error from bad arg when unmarshaling")

}

func TestParseFilterArgsValidJson(t *testing.T) {
	resetGlobalState()
	args.Queues = `["test-1", "test-2", "test-3"]`
	parseFilterArgs()
}

func TestCreateEntity(t *testing.T) {
	resetGlobalState()
	i := getTestingIntegration(t)

	expectedEntityName := "/firstEntity"
	expectedMetricEntityName := "queue:" + expectedEntityName
	e1, metricNS, err := createEntity(i, "firstEntity", queueType, "/")
	assert.NoError(t, err)
	assert.NotNil(t, e1)
	assert.Equal(t, expectedEntityName, e1.Metadata.Name)
	assert.Equal(t, queueType, e1.Metadata.Namespace)
	assert.Equal(t, 2, len(metricNS))
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "displayName", Value: expectedEntityName,
	})
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "entityName", Value: expectedMetricEntityName,
	})

	expectedEntityName = "/vhost2/" + defaultExchangeName
	expectedMetricEntityName = "exchange:" + expectedEntityName
	e2, metricNS, err := createEntity(i, "", exchangeType, "/vhost2")
	assert.NoError(t, err)
	assert.NotNil(t, e2)
	assert.NotNil(t, metricNS)
	assert.Equal(t, expectedEntityName, e2.Metadata.Name)
	assert.Equal(t, exchangeType, e2.Metadata.Namespace)
	assert.Equal(t, 2, len(metricNS))
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "displayName", Value: expectedEntityName,
	})
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "entityName", Value: expectedMetricEntityName,
	})

	existingQueueFilter := queueFilter
	defer func() {
		queueFilter = existingQueueFilter
	}()
	queueFilter = []string{"missing-queue"}

	e3, metricNS, err := createEntity(i, "actual-queue", queueType, "/")
	assert.Nil(t, e3)
	assert.Nil(t, metricNS)
	assert.Nil(t, err)
}

func TestCheckFiltersEmptyLists(t *testing.T) {
	entityName := ""
	var nameList []string
	var nameRegexList []string
	actual := checkFilters(entityName, nameList, nameRegexList)
	assert.True(t, actual, "should return true")
}

func TestCheckFilterNameList(t *testing.T) {
	entityName := "test-entity"
	nameList := []string{"test-entity", "name2", "name3"}
	var nameRegexList []string
	actual := checkFilters(entityName, nameList, nameRegexList)
	assert.True(t, actual, "should return true")
}

func TestCheckErr(t *testing.T) {
	l := new(mockedLogger)
	logger = l
	l.On("Errorf", mock.Anything, mock.Anything).Once()
	checkErr(func() error {
		return nil
	})
	checkErr(func() error {
		return errors.New("Test Error")
	})
	l.AssertExpectations(t)
}

func TestPanicErr(t *testing.T) {
	assert.NotPanics(t, func() {
		panicOnErr(nil)
	}, "panicOnEror should not panic when not given an error")

	assert.Panics(t, func() {
		panicOnErr(errors.New("Panic Error"))
	}, "panicOnEror should panic when given an error")
}
