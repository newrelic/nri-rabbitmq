package data

import (
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/data/consts"
	"github.com/newrelic/nri-rabbitmq/src/testutils"
	"github.com/stretchr/testify/assert"
)

var expectedArguments = map[string]interface{}{
	"one":   float64(1),
	"two":   "two",
	"three": true,
	"four":  "[true,false]",
}

func TestExchangeData_UnmarshalJSON_MarshalMetrics(t *testing.T) {
	var exchangeData ExchangeData
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "exchange.json"), &exchangeData)
	assert.NotNil(t, exchangeData)
	assert.Equal(t, "exchange1", exchangeData.Name)
	assert.Equal(t, "vhost1", exchangeData.Vhost)
	assert.Equal(t, "direct", exchangeData.Type)
	assert.True(t, exchangeData.Durable)
	assert.False(t, exchangeData.AutoDelete)
	assert.NotNil(t, exchangeData.MessageStats)
	assert.Equal(t, getInt64(1), exchangeData.MessageStats.PublishIn)
	assert.Equal(t, getFloat64(1.1), exchangeData.MessageStats.PublishInDetails.Rate)
	assert.Equal(t, getInt64(2), exchangeData.MessageStats.PublishOut)
	assert.Equal(t, getFloat64(2.2), exchangeData.MessageStats.PublishOutDetails.Rate)
	assert.Equal(t, "exchange1", exchangeData.EntityName())
	assert.Equal(t, consts.ExchangeType, exchangeData.EntityType())
	assert.Equal(t, "vhost1", exchangeData.EntityVhost())

	testIntegration := testutils.GetTestingIntegration(t)
	e, metricAttribs, err := exchangeData.GetEntity(testIntegration, "testClusterName")
	assert.NotNil(t, e)
	assert.NotEmpty(t, metricAttribs)
	assert.NoError(t, err)

	ms := e.NewMetricSet("TestSample", metricAttribs...)

	assert.NoError(t, ms.MarshalMetrics(exchangeData))
	expectedMetrics := map[string]interface{}{
		"exchange.messagesPublishedPerChannel":          float64(1),
		"exchange.messagesPublishedPerChannelPerSecond": float64(1.1),
		"exchange.messagesPublishedQueue":               float64(2),
		"exchange.messagesPublishedQueuePerSecond":      float64(2.2),
	}
	assert.Equal(t, 2+len(expectedMetrics)+len(metricAttribs), len(ms.Metrics), "Unexpected metric count for ExchangeData")
	for k, v := range expectedMetrics {
		assert.Equal(t, v, ms.Metrics[k], k)
	}

	exchangeData.CollectInventory(e, nil)
	expectedInventory := map[string]interface{}{
		"exchange/type":        "direct",
		"exchange/durable":     1,
		"exchange/auto_delete": 0,
	}
	assert.Equal(t, 1+len(expectedInventory), len(e.Inventory.Items()), "Unexpected inventory count for ExchangeData")
	for k, v := range expectedInventory {
		item, exists := e.Inventory.Item(k)
		assert.True(t, exists)
		assert.Equal(t, v, item["value"], k)
	}
	args, exists := e.Inventory.Item("exchange/arguments")
	if assert.True(t, exists, "exchange/arguments") {
		for k, v := range expectedArguments {
			assert.Equal(t, v, args[k], "exchange/arguments/"+k)
		}
	}
}
