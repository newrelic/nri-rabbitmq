package metrics

import (
	data2 "github.com/newrelic/nri-rabbitmq/src/data"
	consts2 "github.com/newrelic/nri-rabbitmq/src/data/consts"
)

// connKey is used to uniquely identify a connection by Vhost and State
type connKey struct {
	Vhost, State string
}

// collectConnectionStats returns a map of vhost -> connection totals by status and an overall status
func collectConnectionStats(connectionsData []*data2.ConnectionData) (stats map[connKey]int) {
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
func collectBindingStats(bindingsData []*data2.BindingData) (stats data2.BindingStats) {
	stats = make(data2.BindingStats)

	for _, binding := range bindingsData {
		srcKey := data2.BindingKey{
			Vhost:      binding.Vhost,
			EntityName: binding.Source,
			EntityType: consts2.ExchangeType,
		}
		dstKey := data2.BindingKey{
			Vhost:      srcKey.Vhost,
			EntityName: binding.Destination,
			EntityType: binding.DestinationType,
		}
		if stat := stats[srcKey]; stat != nil {
			stat.Destination = append(stat.Destination, &dstKey)
		} else {
			stats[srcKey] = &data2.Binding{
				Destination: []*data2.BindingKey{&dstKey},
				Source:      []*data2.BindingKey{},
			}
		}
		if stat := stats[dstKey]; stat != nil {
			stat.Source = append(stat.Source, &srcKey)
		} else {
			stats[dstKey] = &data2.Binding{
				Destination: []*data2.BindingKey{},
				Source:      []*data2.BindingKey{&srcKey},
			}
		}
	}
	return
}
