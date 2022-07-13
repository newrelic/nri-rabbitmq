package metrics

import (
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// PopulateClusterData populates the cluster entity with appropriate data
func PopulateClusterData(integrationData *integration.Integration, overviewData *data.OverviewData) {
	if overviewData == nil || overviewData.ClusterName == "" {
		return
	}
	e, _, err := data.CreateEntity(integrationData, overviewData.ClusterName, consts.ClusterType, "", overviewData.ClusterName)
	if err != nil {
		log.Error("Error creating cluster entity: %v", err)
		return
	}
	data.SetInventoryItem(e, "version", "rabbitmq", overviewData.RabbitMQVersion)
	data.SetInventoryItem(e, "version", "management", overviewData.ManagementVersion)
}
