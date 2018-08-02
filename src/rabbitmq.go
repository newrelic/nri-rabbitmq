package main

import (
	"errors"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/newrelic/nri-rabbitmq/client"
	"github.com/newrelic/nri-rabbitmq/inventory"
	"github.com/newrelic/nri-rabbitmq/logger"
	"github.com/newrelic/nri-rabbitmq/metrics"
	"github.com/newrelic/nri-rabbitmq/utils"
	"github.com/stretchr/objx"
)

const (
	integrationName    = "com.newrelic.rabbitmq"
	integrationVersion = "0.1.0"
)

func main() {
	var argList args.ArgumentList
	// Create Integration
	rabbitmqIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&argList))
	utils.PanicOnErr(err)

	logger.SetLogger(rabbitmqIntegration.Logger())

	utils.PanicOnErr(args.SetGlobalArgs(argList))

	utils.PanicOnErr(args.GlobalArgs.Validate())

	endpoints := client.GetEndpointsToCollect()

	allEndpointResponses, err := client.CollectEndpoints(endpoints...)
	utils.PanicOnErr(err)

	responses := objx.New(allEndpointResponses)
	if responses == nil {
		panic(errors.New("unexpected management response"))
	}

	if args.GlobalArgs.All() || args.GlobalArgs.Metrics {
		logger.Infof("Collecting metrics.")
		metrics.CollectMetrics(rabbitmqIntegration, &responses)
	}

	if args.GlobalArgs.All() || args.GlobalArgs.Inventory {
		logger.Infof("Collecting inventory.")
		inventory.CollectInventory(rabbitmqIntegration, responses.Get(client.NodesEndpoint).ObjxMapSlice())
	}

	if args.GlobalArgs.All() || (args.GlobalArgs.Metrics && args.GlobalArgs.Inventory) {
		inventory.PopulateClusterEntity(rabbitmqIntegration, responses.Get(client.OverviewEndpoint).ObjxMap())
	}

	utils.PanicOnErr(rabbitmqIntegration.Publish())
}

// func collectMetrics(rabbitmqIntegration *integration.Integration, responses *objx.Map) {
// 	connStats := collectConnectionStats(responses.Get(client.ConnectionsEndpoint).ObjxMapSlice())
// 	bindingStats := collectBindingStats(responses.Get(client.BindingsEndpoint).ObjxMapSlice())

// 	for _, mapping := range entityMappings {
// 		responseObjects := responses.Get(mapping.entityEndpoint).ObjxMapSlice()
// 		for _, entityData := range responseObjects {
// 			parseEntityResponse(rabbitmqIntegration, &entityData, mapping.entityType, connStats, bindingStats)
// 		}
// 	}
// }

// func parseEntityResponse(rabbitmqIntegration *integration.Integration, entityData *objx.Map, entityType string, connectionStats map[connKey]int, bindingStats map[bindingKey]int) {
// 	entityName := entityData.Get("name").Str()
// 	vhost := entityData.Get("vhost").Str()
// 	entity, metricNamespace, err := utils.CreateEntity(rabbitmqIntegration, entityName, entityType, vhost)
// 	if entity == nil {
// 		logger.Infof("%v entity %v did not match filter - skipping", entityType, entityName)
// 		return
// 	}
// 	if err != nil {
// 		logger.Errorf("could not create entity [%v, %v, %v]: %v", entityName, entityType, vhost, err)
// 		return
// 	}
// 	sampleName := fmt.Sprintf("Rabbitmq%vSample", strings.Title(entityType))
// 	metricSet := entity.NewMetricSet(sampleName, metricNamespace...)

// 	if entityType == consts.VhostType {
// 		populateVhostMetrics(entityName, metricSet, connectionStats)
// 	} else {
// 		populateMetrics(metricSet, entityType, entityData)
// 		populateEntityInventory(entity, entityType, entityData)

// 		if entityType == consts.QueueType || entityType == consts.ExchangeType {
// 			populateBindingMetric(entityName, vhost, entityType, metricSet, bindingStats)
// 		}
// 	}
// }

// func collectInventory(rabbitmqIntegration *integration.Integration, responses *objx.Map) {
// 	nodeName, err := getNodeName()
// 	utils.PanicOnErr(err)

// 	overviewData := responses.Get(client.OverviewEndpoint).ObjxMap()
// 	var localNode *integration.Entity
// 	nodesList := responses.Get(client.NodesEndpoint).ObjxMapSlice()

// 	if nodesList != nil {
// 		localNode, err = getNodeEntity(nodeName, nodesList, rabbitmqIntegration)
// 	} else {
// 		logger.Warnf("Could not retrieve node data")
// 		localNode, err = getNodeEntity(nodeName, nil, rabbitmqIntegration)
// 	}

// 	if err == nil && localNode != nil {
// 		if err = setInventoryData(localNode, overviewData); err != nil {
// 			logger.Errorf("%v", err)
// 		}
// 	}
// }
