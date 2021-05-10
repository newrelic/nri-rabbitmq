package metrics

import (
	data2 "github.com/newrelic/nri-rabbitmq/src/data"
	consts2 "github.com/newrelic/nri-rabbitmq/src/data/consts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_collectConnectionStats(t *testing.T) {
	data := []*data2.ConnectionData{
		{Vhost: "/", State: "running"},
		{Vhost: "/", State: "blocked"},
		{Vhost: "/", State: "running"},
	}

	stats := collectConnectionStats(data)
	assert.NotNil(t, stats)

	assert.Equal(t, 2, stats[connKey{"/", "running"}])
	assert.Equal(t, 1, stats[connKey{"/", "blocked"}])
	assert.Equal(t, 3, stats[connKey{"/", "total"}])
}

func Test_collectBindingStats(t *testing.T) {
	bindingData := []*data2.BindingData{
		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "dest1",
			DestinationType: consts2.QueueType,
		},
		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "dest2",
			DestinationType: consts2.QueueType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "dest1",
			DestinationType: consts2.QueueType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "dest2",
			DestinationType: consts2.QueueType,
		},

		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "dest1",
			DestinationType: consts2.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "dest2",
			DestinationType: consts2.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "dest1",
			DestinationType: consts2.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "dest2",
			DestinationType: consts2.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "source2",
			DestinationType: consts2.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "source1",
			DestinationType: consts2.ExchangeType,
		},
	}

	stats := collectBindingStats(bindingData)
	assert.NotNil(t, stats)

	stat := stats[data2.BindingKey{}]
	assert.Nil(t, stat)

	stat = stats[data2.BindingKey{Vhost: "/", EntityName: "source1", EntityType: consts2.ExchangeType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 1, len(stat.Source))
		assert.Equal(t, 5, len(stat.Destination))
	}

	stat = stats[data2.BindingKey{Vhost: "/", EntityName: "source2", EntityType: consts2.ExchangeType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 1, len(stat.Source))
		assert.Equal(t, 5, len(stat.Destination))
	}

	stat = stats[data2.BindingKey{Vhost: "/", EntityName: "dest1", EntityType: consts2.QueueType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 2, len(stat.Source))
		assert.Equal(t, 0, len(stat.Destination))
	}

	stat = stats[data2.BindingKey{Vhost: "/", EntityName: "dest2", EntityType: consts2.QueueType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 2, len(stat.Source))
		assert.Equal(t, 0, len(stat.Destination))
	}

	stat = stats[data2.BindingKey{Vhost: "/", EntityName: "dest1", EntityType: consts2.ExchangeType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 2, len(stat.Source))
		assert.Equal(t, 0, len(stat.Destination))
	}

	stat = stats[data2.BindingKey{Vhost: "/", EntityName: "dest2", EntityType: consts2.ExchangeType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 2, len(stat.Source))
		assert.Equal(t, 0, len(stat.Destination))
	}
}
