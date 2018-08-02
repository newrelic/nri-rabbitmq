package inventory

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-rabbitmq/utils"
	"github.com/newrelic/nri-rabbitmq/utils/consts"
	"github.com/stretchr/objx"
)

// PopulateClusterEntity populates the cluse entity with appropriate inventory data
func PopulateClusterEntity(integrationData *integration.Integration, overviewData objx.Map) {
	if overviewData != nil {
		clusterName := overviewData.Get("cluster_name").Str()
		if clusterName == "" {
			return
		}
		e, _, _ := utils.CreateEntity(integrationData, clusterName, consts.ClusterType, "")
		if e != nil {
			setInventoryValue(e, "version", "rabbitmq", overviewData.Get("rabbitmq_version").Str())
			setInventoryValue(e, "version", "management", overviewData.Get("management_version").Str())
		}
	}
}
