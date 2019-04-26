package inventory

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"
)

// PopulateClusterInventory populates the cluster entity with appropriate inventory data
func PopulateClusterInventory(integrationData *integration.Integration, overviewData *data.OverviewData) {
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
