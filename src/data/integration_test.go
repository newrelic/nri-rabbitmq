package data

import (
	"fmt"
	"testing"

	args2 "github.com/newrelic/nri-rabbitmq/src/args"
	consts2 "github.com/newrelic/nri-rabbitmq/src/data/consts"
	testutils2 "github.com/newrelic/nri-rabbitmq/src/testutils"

	"github.com/newrelic/infra-integrations-sdk/data/attribute"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// TODO remove global args.
	// This test are heavily based on global args to create entities.
	args2.GlobalArgs = args2.RabbitMQArguments{
		Hostname: "foo",
		Port:     8000,
	}
}

func TestCreateEntity(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)

	expectedEntityName := "/firstEntity"
	expectedEntityKey := fmt.Sprintf("%s:%d:%s", args2.GlobalArgs.Hostname, args2.GlobalArgs.Port, expectedEntityName)
	expectedMetricEntityName := "queue:" + expectedEntityName
	e1, metricNS, err := CreateEntity(i, "firstEntity", consts2.QueueType, "/", "testClusterName")
	assert.NoError(t, err)
	assert.NotNil(t, e1)
	assert.Equal(t, expectedEntityKey, e1.Metadata.Name)
	assert.Equal(t, "ra-queue", e1.Metadata.Namespace)
	assert.Equal(t, 3, len(metricNS))
	assert.Contains(t, metricNS, attribute.Attribute{
		Key: "displayName", Value: expectedEntityName,
	})
	assert.Contains(t, metricNS, attribute.Attribute{
		Key: "entityName", Value: expectedMetricEntityName,
	})

	expectedEntityName = "/vhost2/" + consts2.DefaultExchangeName
	expectedEntityKey = fmt.Sprintf("%s:%d:%s", args2.GlobalArgs.Hostname, args2.GlobalArgs.Port, expectedEntityName)
	expectedMetricEntityName = "exchange:" + expectedEntityName
	e2, metricNS, err := CreateEntity(i, "", consts2.ExchangeType, "/vhost2", "testClusterName")
	assert.NoError(t, err)
	assert.NotNil(t, e2)
	assert.NotNil(t, metricNS)
	assert.Equal(t, expectedEntityKey, e2.Metadata.Name)
	assert.Equal(t, "ra-exchange", e2.Metadata.Namespace)
	assert.Equal(t, 3, len(metricNS))
	assert.Contains(t, metricNS, attribute.Attribute{
		Key: "displayName", Value: expectedEntityName,
	})
	assert.Contains(t, metricNS, attribute.Attribute{
		Key: "entityName", Value: expectedMetricEntityName,
	})

	existingArgs := args2.GlobalArgs
	defer func() {
		args2.GlobalArgs = existingArgs
	}()
	args2.GlobalArgs = args2.RabbitMQArguments{
		Queues: []string{"missing-queue"},
	}

	e3, metricNS, err := CreateEntity(i, "actual-queue", consts2.QueueType, "/", "testClusterName")
	assert.Nil(t, e3)
	assert.Nil(t, metricNS)
	assert.Nil(t, err)
}

const errorKey = "this-key-is-longer-than-375-to-force-an-error-lorem-ipsum-dolor-sit-amet-consectetur-adipiscing-elit-sed-do-eiusmod-tempor-incididunt-ut-labore-et-dolore-magna-aliqua-ut-enim-ad-minim-veniam-quis-nostrud-exercitation-ullamco-laboris-nisi-ut-aliquip-ex-ea-commodo-consequat-duis-aute-irure-dolor-in-reprehenderit-in-voluptate-velit-esse-cillum-dolore-eu-fugiat-nulla-pariatures"

func TestSetInventoryItem_ErrorKeyToLong(t *testing.T) {
	_, e := testutils2.GetTestingEntity(t)
	SetInventoryItem(e, errorKey, "nope", "false")
	assert.Empty(t, e.Inventory.Items())

	_, e = testutils2.GetTestingEntity(t, "name", "namespace")
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

func TestCollectInventory_Exchange(t *testing.T) {
	data := &ExchangeData{
		Vhost:      "vhost1",
		Name:       "exchange1",
		Type:       "test-type",
		Durable:    true,
		AutoDelete: false,
		Arguments:  testArgs,
	}
	bindingStats := BindingStats{
		BindingKey{
			Vhost:      "vhost1",
			EntityName: "exchange1",
			EntityType: consts2.ExchangeType,
		}: &Binding{
			Source: []*BindingKey{
				{
					Vhost:      "vhost1",
					EntityName: "exchange2",
					EntityType: consts2.ExchangeType,
				},
				{
					Vhost:      "vhost1",
					EntityName: "exchange3",
					EntityType: consts2.ExchangeType,
				},
			},
			Destination: []*BindingKey{
				{
					Vhost:      "vhost1",
					EntityName: "queue2",
					EntityType: consts2.QueueType,
				},
				{
					Vhost:      "vhost1",
					EntityName: "exchange4",
					EntityType: consts2.ExchangeType,
				},
			},
		},
	}
	_, e := testutils2.GetTestingEntity(t)
	data.CollectInventory(e, bindingStats)

	item, ok := e.Inventory.Item("exchange/bindings.source")
	assert.True(t, ok)
	assert.Equal(t, "exchange:vhost1/exchange2, exchange:vhost1/exchange3", item["value"])

	item, ok = e.Inventory.Item("exchange/bindings.destination")
	assert.True(t, ok)
	assert.Equal(t, "queue:vhost1/queue2, exchange:vhost1/exchange4", item["value"])

	item, ok = e.Inventory.Item("exchange/type")
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

func TestCollectInventory_Queue(t *testing.T) {
	data := &QueueData{
		Vhost:      "vhost1",
		Name:       "queue1",
		Exclusive:  true,
		Durable:    false,
		AutoDelete: true,
		Arguments:  testArgs,
	}
	bindingStats := BindingStats{
		BindingKey{
			Vhost:      "vhost1",
			EntityName: "queue1",
			EntityType: consts2.QueueType,
		}: &Binding{
			Source: []*BindingKey{
				{
					Vhost:      "vhost1",
					EntityName: "exchange1",
					EntityType: consts2.ExchangeType,
				},
				{
					Vhost:      "vhost1",
					EntityName: "exchange2",
					EntityType: consts2.ExchangeType,
				},
			},
			Destination: []*BindingKey{},
		},
	}
	_, e := testutils2.GetTestingEntity(t)
	data.CollectInventory(e, bindingStats)

	item, ok := e.Inventory.Item("queue/bindings.source")
	assert.True(t, ok)
	assert.Equal(t, "exchange:vhost1/exchange1, exchange:vhost1/exchange2", item["value"])

	_, ok = e.Inventory.Item("queue/bindings.destination")
	assert.False(t, ok)

	item, ok = e.Inventory.Item("queue/exclusive")
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
