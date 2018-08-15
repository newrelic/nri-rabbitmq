package data

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/newrelic/nri-rabbitmq/data/consts"
	"github.com/newrelic/nri-rabbitmq/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_CreateEntity(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	i := testutils.GetTestingIntegration(t)

	expectedEntityName := "/firstEntity"
	expectedMetricEntityName := "queue:" + expectedEntityName
	e1, metricNS, err := CreateEntity(i, "firstEntity", consts.QueueType, "/")
	assert.NoError(t, err)
	assert.NotNil(t, e1)
	assert.Equal(t, expectedEntityName, e1.Metadata.Name)
	assert.Equal(t, consts.QueueType, e1.Metadata.Namespace)
	assert.Equal(t, 2, len(metricNS))
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "displayName", Value: expectedEntityName,
	})
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "entityName", Value: expectedMetricEntityName,
	})

	expectedEntityName = "/vhost2/" + consts.DefaultExchangeName
	expectedMetricEntityName = "exchange:" + expectedEntityName
	e2, metricNS, err := CreateEntity(i, "", consts.ExchangeType, "/vhost2")
	assert.NoError(t, err)
	assert.NotNil(t, e2)
	assert.NotNil(t, metricNS)
	assert.Equal(t, expectedEntityName, e2.Metadata.Name)
	assert.Equal(t, consts.ExchangeType, e2.Metadata.Namespace)
	assert.Equal(t, 2, len(metricNS))
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "displayName", Value: expectedEntityName,
	})
	assert.Contains(t, metricNS, metric.Attribute{
		Key: "entityName", Value: expectedMetricEntityName,
	})

	existingArgs := args.GlobalArgs
	defer func() {
		args.GlobalArgs = existingArgs
	}()
	args.GlobalArgs = args.RabbitMQArguments{
		Queues: []string{"missing-queue"},
	}

	e3, metricNS, err := CreateEntity(i, "actual-queue", consts.QueueType, "/")
	assert.Nil(t, e3)
	assert.Nil(t, metricNS)
	assert.Nil(t, err)
}

const errorKey = "this-key-is-longer-than-375-to-force-an-error-lorem-ipsum-dolor-sit-amet-consectetur-adipiscing-elit-sed-do-eiusmod-tempor-incididunt-ut-labore-et-dolore-magna-aliqua-ut-enim-ad-minim-veniam-quis-nostrud-exercitation-ullamco-laboris-nisi-ut-aliquip-ex-ea-commodo-consequat-duis-aute-irure-dolor-in-reprehenderit-in-voluptate-velit-esse-cillum-dolore-eu-fugiat-nulla-pariatures"

func Test_SetInventoryItem_ErrorKeyToLong(t *testing.T) {
	_, e := testutils.GetTestingEntity(t)
	SetInventoryItem(e, errorKey, "nope", "false")
	assert.Empty(t, e.Inventory.Items())

	_, e = testutils.GetTestingEntity(t, "name", "namespace")
	SetInventoryItem(e, errorKey, "nope", "false")
	assert.Empty(t, e.Inventory.Items())
}

var testArgs = map[string]interface{}{
	"string":  "test-string",
	"number":  123.456,
	"boolean": true,
	"array": []interface{}{
		"sub-string",
		654.321,
		false,
		[]interface{}{
			"sub-array",
			987,
		},
	},
}

const testArgsArrayValue = `["sub-string",654.321,false,["sub-array",987]]`

func Test_CollectInventory_Exchange(t *testing.T) {
	data := &ExchangeData{
		Type:       "test-type",
		Durable:    true,
		AutoDelete: false,
		Arguments:  testArgs,
	}
	_, e := testutils.GetTestingEntity(t)
	data.CollectInventory(e)

	item, ok := e.Inventory.Item("exchange/type")
	assert.True(t, ok)
	assert.Equal(t, "test-type", item["value"])

	item, ok = e.Inventory.Item("exchange/durable")
	assert.True(t, ok)
	assert.Equal(t, 1, item["value"])

	item, ok = e.Inventory.Item("exchange/auto_delete")
	assert.True(t, ok)
	assert.Equal(t, 0, item["value"])

	item, ok = e.Inventory.Item("exchange/arguments")
	if assert.True(t, ok) {
		assert.Equal(t, len(testArgs), len(item))
		for k, v := range testArgs {
			if k == "array" {
				assert.Equal(t, testArgsArrayValue, item[k])
			} else {
				assert.Equal(t, v, item[k])
			}
		}
	}
}

func Test_CollectInventory_Queue(t *testing.T) {
	data := &QueueData{
		Exclusive:  true,
		Durable:    false,
		AutoDelete: true,
		Arguments:  testArgs,
	}
	_, e := testutils.GetTestingEntity(t)
	data.CollectInventory(e)

	item, ok := e.Inventory.Item("queue/exclusive")
	assert.True(t, ok)
	assert.Equal(t, 1, item["value"])

	item, ok = e.Inventory.Item("queue/durable")
	assert.True(t, ok)
	assert.Equal(t, 0, item["value"])

	item, ok = e.Inventory.Item("queue/auto_delete")
	assert.True(t, ok)
	assert.Equal(t, 1, item["value"])

	item, ok = e.Inventory.Item("queue/arguments")
	if assert.True(t, ok) {
		assert.Equal(t, len(testArgs), len(item))
		for k, v := range testArgs {
			if k == "array" {
				assert.Equal(t, testArgsArrayValue, item[k])
			} else {
				assert.Equal(t, v, item[k])
			}
		}
	}
}
