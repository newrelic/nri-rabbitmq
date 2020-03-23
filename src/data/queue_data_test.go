package data

import (
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/data/consts"
	"github.com/newrelic/nri-rabbitmq/src/testutils"
	"github.com/stretchr/testify/assert"
)

func TestQueueData_UnmarshalJSON_MarshalMetrics(t *testing.T) {
	var queueData QueueData
	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "queue.json"), &queueData)
	assert.NotNil(t, queueData)
	assert.Equal(t, getInt64(1), queueData.Messages)
	assert.Equal(t, getFloat64(0.1), queueData.MessagesDetails.Rate)
	assert.Equal(t, getInt64(2), queueData.MessagesReady)
	assert.Equal(t, getFloat64(0.2), queueData.MessagesReadyDetail.Rate)
	assert.Equal(t, getInt64(3), queueData.MessagesUnacknowledged)
	assert.Equal(t, getFloat64(0.3), queueData.MessagesUnacknowledgedDetail.Rate)
	assert.Equal(t, getInt64(4), queueData.MessageStats.Ack)
	assert.Equal(t, getFloat64(0.4), queueData.MessageStats.AckDetails.Rate)
	assert.Equal(t, getInt64(5), queueData.MessageStats.Deliver)
	assert.Equal(t, getFloat64(0.5), queueData.MessageStats.DeliverDetails.Rate)
	assert.Equal(t, getInt64(6), queueData.MessageStats.DeliverGet)
	assert.Equal(t, getFloat64(0.6), queueData.MessageStats.DeliverGetDetails.Rate)
	assert.Equal(t, getInt64(7), queueData.MessageStats.Publish)
	assert.Equal(t, getFloat64(0.7), queueData.MessageStats.PublishDetails.Rate)
	assert.Equal(t, getInt64(8), queueData.MessageStats.Redeliver)
	assert.Equal(t, getFloat64(0.8), queueData.MessageStats.RedeliverDetails.Rate)
	assert.True(t, queueData.Exclusive)
	assert.False(t, queueData.AutoDelete)
	assert.True(t, queueData.Durable)
	assert.Equal(t, "vhost1", queueData.Vhost)
	assert.Equal(t, "queue1", queueData.Name)
	assert.Equal(t, getInt64(9), queueData.Consumers)
	assert.Equal(t, getInt64(10), queueData.ActiveConsumers)
	assert.Equal(t, getFloat64(11.11), queueData.ConsumerUtilisation)
	assert.Equal(t, getInt64(1024), queueData.Memory)
	assert.Equal(t, "queue1", queueData.EntityName())
	assert.Equal(t, consts.QueueType, queueData.EntityType())

	testIntegration := testutils.GetTestingIntegration(t)
	e, metricAttribs, err := queueData.GetEntity(testIntegration, "testClusterName")
	assert.NotNil(t, e)
	assert.NotEmpty(t, metricAttribs)
	assert.NoError(t, err)

	ms := e.NewMetricSet("TestSample", metricAttribs...)

	assert.NoError(t, ms.MarshalMetrics(queueData))
	expectedMetrics := map[string]interface{}{
		"queue.consumers":                             float64(9),
		"queue.consumerMessageUtilizationPerSecond":   float64(11.11),
		"queue.countActiveConsumersReceiveMessages":   float64(10),
		"queue.erlangBytesConsumedInBytes":            float64(1024),
		"queue.totalMessages":                         float64(1),
		"queue.totalMessagesPerSecond":                float64(0.1),
		"queue.messagesReadyDeliveryClients":          float64(2),
		"queue.messagesReadyDeliveryClientsPerSecond": float64(0.2),
		"queue.messagesReadyUnacknowledged":           float64(3),
		"queue.messagesReadyUnacknowledgedPerSecond":  float64(0.3),
		"queue.messagesAcknowledged":                  float64(4),
		"queue.messagesAcknowledgedPerSecond":         float64(0.4),
		"queue.messagesDeliveredAckMode":              float64(5),
		"queue.messagesDeliveredAckModePerSecond":     float64(0.5),
		"queue.sumMessagesDelivered":                  float64(6),
		"queue.sumMessagesDeliveredPerSecond":         float64(0.6),
		"queue.messagesPublished":                     float64(7),
		"queue.messagesPublishedPerSecond":            float64(0.7),
		"queue.messagesRedeliverGet":                  float64(8),
		"queue.messagesRedeliverGetPerSecond":         float64(0.8),
	}
	assert.Equal(t, 2+len(expectedMetrics)+len(metricAttribs), len(ms.Metrics), "Unexpected metric count for QueueData")
	for k, v := range expectedMetrics {
		assert.Equal(t, v, ms.Metrics[k], k)
	}

	queueData.CollectInventory(e, nil)
	expectedInventory := map[string]interface{}{
		"queue/exclusive":   1,
		"queue/durable":     1,
		"queue/auto_delete": 0,
	}
	assert.Equal(t, 1+len(expectedInventory), len(e.Inventory.Items()), "Unexpected inventory count for QueueData")
	for k, v := range expectedInventory {
		item, exists := e.Inventory.Item(k)
		assert.True(t, exists)
		assert.Equal(t, v, item["value"], k)
	}
	args, exists := e.Inventory.Item("queue/arguments")
	if assert.True(t, exists, "queue/arguments") {
		for k, v := range expectedArguments {
			assert.Equal(t, v, args[k], "queue/arguments/"+k)
		}
	}
}

func getBool(b bool) *bool {
	return &b
}

func getInt(i int) *int {
	return &i
}

func getInt64(i int64) *int64 {
	return &i
}

func getFloat64(f float64) *float64 {
	return &f
}
