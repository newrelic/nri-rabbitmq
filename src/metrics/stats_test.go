package metrics

import (
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"
	"github.com/stretchr/testify/assert"
)

func Test_collectConnectionStats(t *testing.T) {
	data := []*data.ConnectionData{
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
	bindingData := []*data.BindingData{
		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "dest1",
			DestinationType: consts.QueueType,
		},
		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "dest2",
			DestinationType: consts.QueueType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "dest1",
			DestinationType: consts.QueueType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "dest2",
			DestinationType: consts.QueueType,
		},

		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "dest1",
			DestinationType: consts.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "dest2",
			DestinationType: consts.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "dest1",
			DestinationType: consts.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "dest2",
			DestinationType: consts.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source1",
			Destination:     "source2",
			DestinationType: consts.ExchangeType,
		},
		{
			Vhost:           "/",
			Source:          "source2",
			Destination:     "source1",
			DestinationType: consts.ExchangeType,
		},
	}

	stats := collectBindingStats(bindingData)
	assert.NotNil(t, stats)

	stat := stats[data.BindingKey{}]
	assert.Nil(t, stat)

	stat = stats[data.BindingKey{Vhost: "/", EntityName: "source1", EntityType: consts.ExchangeType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 1, len(stat.Source))
		assert.Equal(t, 5, len(stat.Destination))
	}

	stat = stats[data.BindingKey{Vhost: "/", EntityName: "source2", EntityType: consts.ExchangeType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 1, len(stat.Source))
		assert.Equal(t, 5, len(stat.Destination))
	}

	stat = stats[data.BindingKey{Vhost: "/", EntityName: "dest1", EntityType: consts.QueueType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 2, len(stat.Source))
		assert.Equal(t, 0, len(stat.Destination))
	}

	stat = stats[data.BindingKey{Vhost: "/", EntityName: "dest2", EntityType: consts.QueueType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 2, len(stat.Source))
		assert.Equal(t, 0, len(stat.Destination))
	}

	stat = stats[data.BindingKey{Vhost: "/", EntityName: "dest1", EntityType: consts.ExchangeType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 2, len(stat.Source))
		assert.Equal(t, 0, len(stat.Destination))
	}

	stat = stats[data.BindingKey{Vhost: "/", EntityName: "dest2", EntityType: consts.ExchangeType}]
	if assert.NotNil(t, stat) {
		assert.Equal(t, 2, len(stat.Source))
		assert.Equal(t, 0, len(stat.Destination))
	}
}
