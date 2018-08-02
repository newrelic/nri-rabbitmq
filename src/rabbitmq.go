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
