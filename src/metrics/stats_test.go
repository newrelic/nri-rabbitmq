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
	data := []*data.BindingData{
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
	}

	stats := collectBindingStats(data)
	assert.NotNil(t, stats)

	assert.Equal(t, 0, stats[bindingKey{}], "A missing key should return a 0 count")

	assert.Equal(t, 4, stats[bindingKey{"/", "source1", consts.ExchangeType}], "Source exchange [source1] should have 4 bindings")
	assert.Equal(t, 4, stats[bindingKey{"/", "source2", consts.ExchangeType}], "Source exchange [source2] should have 4 bindings")

	assert.Equal(t, 2, stats[bindingKey{"/", "dest1", consts.QueueType}])
	assert.Equal(t, 2, stats[bindingKey{"/", "dest2", consts.QueueType}])
	assert.Equal(t, 2, stats[bindingKey{"/", "dest1", consts.ExchangeType}])
	assert.Equal(t, 2, stats[bindingKey{"/", "dest2", consts.ExchangeType}])
}
