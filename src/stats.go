package main

import (
	"github.com/stretchr/objx"
)

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

type bindingKey struct {
	Vhost, EntityName, EntityType string
}

func collectBindingStats(bindingsData []objx.Map) (stats map[bindingKey]int) {
	stats = map[bindingKey]int{}

	for _, binding := range bindingsData {
		srcKey := bindingKey{
			binding.Get("vhost").Str(),
			binding.Get("source").Str(),
			exchangeType,
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
