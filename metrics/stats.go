package metrics

import (
	"github.com/newrelic/nri-rabbitmq/utils/consts"
	"github.com/stretchr/objx"
)

// connKey is used to uniquely identify a connection by Vhost and State
type connKey struct {
	Vhost, State string
}

// collectConnectionStats returns a map of vhost -> connection totals by status and an overall status
func collectConnectionStats(connectionsData []objx.Map) (stats map[connKey]int) {
	stats = map[connKey]int{}

	for _, connection := range connectionsData {
		key := connKey{
			connection.Get("vhost").Str(),
			connection.Get("state").Str(),
		}

		stats[key]++
		stats[connKey{key.Vhost, "total"}]++
	}
	return
}

// bindingKey is used to uniquely identify a binding by Vhost, EntityName, and EntityType
type bindingKey struct {
	Vhost, EntityName, EntityType string
}

// CollectBindingStats returns a map of bindingKey{vhost,source,dest} -> count
func collectBindingStats(bindingsData []objx.Map) (stats map[bindingKey]int) {
	stats = map[bindingKey]int{}

	for _, binding := range bindingsData {
		srcKey := bindingKey{
			binding.Get("vhost").Str(),
			binding.Get("source").Str(),
			consts.ExchangeType,
		}
		dstKey := bindingKey{
			srcKey.Vhost,
			binding.Get("destination").Str(),
			binding.Get("destination_type").Str(),
		}
		stats[srcKey]++
		stats[dstKey]++
	}
	return
}
