package metrics

import (
	"github.com/newrelic/nri-rabbitmq/data"
	"github.com/newrelic/nri-rabbitmq/data/consts"
)

// connKey is used to uniquely identify a connection by Vhost and State
type connKey struct {
	Vhost, State string
}

// collectConnectionStats returns a map of vhost -> connection totals by status and an overall status
func collectConnectionStats(connectionsData []*data.ConnectionData) (stats map[connKey]int) {
	stats = map[connKey]int{}

	for _, connection := range connectionsData {
		key := connKey{
			connection.Vhost,
			connection.State,
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
func collectBindingStats(bindingsData []*data.BindingData) (stats map[bindingKey]int) {
	stats = map[bindingKey]int{}

	for _, binding := range bindingsData {
		srcKey := bindingKey{
			binding.Vhost,
			binding.Source,
			consts.ExchangeType,
		}
		dstKey := bindingKey{
			srcKey.Vhost,
			binding.Destination,
			binding.DestinationType,
		}
		stats[srcKey]++
		stats[dstKey]++
	}
	return
}
