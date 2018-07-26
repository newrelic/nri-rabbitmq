package main

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/objx"
)

func populateClusterEntity(integrationData *integration.Integration, overviewData objx.Map) {
	if overviewData != nil {
		clusterName := overviewData.Get("cluster_name").Str()
		if clusterName == "" {
			return
		}
		e, _, _ := createEntity(integrationData, clusterName, clusterType, "")
		if e != nil {
			setInventoryValue(e, "version", "rabbitmq", overviewData.Get("rabbitmq_version").Str())
			setInventoryValue(e, "version", "management", overviewData.Get("management_version").Str())
		}
	}
}
