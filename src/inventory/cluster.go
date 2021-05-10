package inventory

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	data2 "github.com/newrelic/nri-rabbitmq/src/data"
	consts2 "github.com/newrelic/nri-rabbitmq/src/data/consts"
)

// PopulateClusterInventory populates the cluster entity with appropriate inventory data
func PopulateClusterInventory(integrationData *integration.Integration, overviewData *data2.OverviewData) {
	if overviewData == nil || overviewData.ClusterName == "" {
		return
	}
	e, _, err := data2.CreateEntity(integrationData, overviewData.ClusterName, consts2.ClusterType, "", overviewData.ClusterName)
	if err != nil {
		log.Error("Error creating cluster entity: %v", err)
		return
	}
	data2.SetInventoryItem(e, "version", "rabbitmq", overviewData.RabbitMQVersion)
	data2.SetInventoryItem(e, "version", "management", overviewData.ManagementVersion)
}
