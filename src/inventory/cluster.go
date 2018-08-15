package inventory

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"
)

// PopulateClusterInventory populates the cluse entity with appropriate inventory data
func PopulateClusterInventory(integrationData *integration.Integration, overviewData *data.OverviewData) {
	if overviewData == nil || overviewData.ClusterName == "" {
		return
	}
	e, _, _ := data.CreateEntity(integrationData, overviewData.ClusterName, consts.ClusterType, "")
	if e != nil {
		data.SetInventoryItem(e, "version", "rabbitmq", overviewData.RabbitMQVersion)
		data.SetInventoryItem(e, "version", "management", overviewData.ManagementVersion)
	}
}
