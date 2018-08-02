package metrics

import (
	"testing"

	"github.com/newrelic/nri-rabbitmq/consts"
	"github.com/newrelic/nri-rabbitmq/testutils"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
)

func Test_collectConnectionStats(t *testing.T) {
	testutils.SetTestLogger(t)
	data := []objx.Map{
		objx.MSI("vhost", "/", "state", "running"),
		objx.MSI("vhost", "/", "state", "blocked"),
		objx.MSI("vhost", "/", "state", "running"),
	}

	stats := collectConnectionStats(data)
	assert.NotNil(t, stats)

	assert.Equal(t, 2, stats[connKey{"/", "running"}])
	assert.Equal(t, 1, stats[connKey{"/", "blocked"}])
	assert.Equal(t, 3, stats[connKey{"/", "total"}])
}

func Test_collectBindingStats(t *testing.T) {
	testutils.SetTestLogger(t)
	data := []objx.Map{
		objx.MSI("vhost", "/", "source", "source1", "destination", "dest1", "destination_type", consts.QueueType),
		objx.MSI("vhost", "/", "source", "source1", "destination", "dest2", "destination_type", consts.QueueType),
		objx.MSI("vhost", "/", "source", "source2", "destination", "dest1", "destination_type", consts.QueueType),
		objx.MSI("vhost", "/", "source", "source2", "destination", "dest2", "destination_type", consts.QueueType),

		objx.MSI("vhost", "/", "source", "source1", "destination", "dest1", "destination_type", consts.ExchangeType),
		objx.MSI("vhost", "/", "source", "source1", "destination", "dest2", "destination_type", consts.ExchangeType),
		objx.MSI("vhost", "/", "source", "source2", "destination", "dest1", "destination_type", consts.ExchangeType),
		objx.MSI("vhost", "/", "source", "source2", "destination", "dest2", "destination_type", consts.ExchangeType),
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
