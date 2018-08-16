package metrics

import (
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"
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

// CollectBindingStats returns a map of BindingKey{vhost,source,dest} -> BindingStats
func collectBindingStats(bindingsData []*data.BindingData) (stats data.BindingStats) {
	stats = make(data.BindingStats)

	for _, binding := range bindingsData {
		srcKey := data.BindingKey{
			Vhost:      binding.Vhost,
			EntityName: binding.Source,
			EntityType: consts.ExchangeType,
		}
		dstKey := data.BindingKey{
			Vhost:      srcKey.Vhost,
			EntityName: binding.Destination,
			EntityType: binding.DestinationType,
		}
		if stat := stats[srcKey]; stat != nil {
			stat.Destination = append(stat.Destination, &dstKey)
		} else {
			stats[srcKey] = &data.Binding{
				Destination: []*data.BindingKey{&dstKey},
				Source:      []*data.BindingKey{},
			}
		}
		if stat := stats[dstKey]; stat != nil {
			stat.Source = append(stat.Source, &srcKey)
		} else {
			stats[dstKey] = &data.Binding{
				Destination: []*data.BindingKey{},
				Source:      []*data.BindingKey{&srcKey},
			}
		}
	}
	return
}
