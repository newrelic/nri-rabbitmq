package main

import (
	"os"

	"github.com/newrelic/nri-rabbitmq/metrics"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/newrelic/nri-rabbitmq/client"
	"github.com/newrelic/nri-rabbitmq/data"
	"github.com/newrelic/nri-rabbitmq/inventory"
)

const (
	integrationName    = "com.newrelic.rabbitmq"
	integrationVersion = "0.1.0"
)

func main() {
	var argList args.ArgumentList
	// Create Integration
	rabbitmqIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&argList))
	exitOnError(err)

	log.SetupLogging(args.GlobalArgs.Verbose)

	exitOnError(args.SetGlobalArgs(argList))

	exitOnError(args.GlobalArgs.Validate())

	rabbitData := getNeededData()

	if args.GlobalArgs.All() || (args.GlobalArgs.Metrics && args.GlobalArgs.Inventory) {
		metrics.CollectVhostMetrics(rabbitmqIntegration, rabbitData.vhosts, rabbitData.connections)

		metricEntities := getMetricEntities(rabbitData)
		metrics.CollectEntityMetrics(rabbitmqIntegration, rabbitData.bindings, metricEntities...)

		inventory.PopulateClusterInventory(rabbitmqIntegration, rabbitData.overview)
	}

	if args.GlobalArgs.All() || args.GlobalArgs.Inventory {
		inventory.CollectInventory(rabbitmqIntegration, rabbitData.nodes)
	}

	err = rabbitmqIntegration.Publish()
	if err != nil {
		log.Error("Error publishing integration: %v", err)
		exitOnError(err)
	}
}

type allData struct {
	overview    *data.OverviewData
	vhosts      []*data.VhostData
	nodes       []*data.NodeData
	queues      []*data.QueueData
	exchanges   []*data.ExchangeData
	connections []*data.ConnectionData
	bindings    []*data.BindingData
}

func getNeededData() *allData {
	data := new(allData)
	if args.GlobalArgs.All() || args.GlobalArgs.Metrics {
		warnIfError(client.CollectEndpoint(client.ConnectionsEndpoint, &data.connections), "Error collecting Connections data: %v")
		warnIfError(client.CollectEndpoint(client.BindingsEndpoint, &data.bindings), "Error collecting Bindings data: %v")
		warnIfError(client.CollectEndpoint(client.VhostsEndpoint, &data.vhosts), "Error collecting Vhost data: %v")
		warnIfError(client.CollectEndpoint(client.NodesEndpoint, &data.nodes), "Error collecting Node data: %v")
		warnIfError(client.CollectEndpoint(client.QueuesEndpoint, &data.queues), "Error collecting Queue data: %v")
		warnIfError(client.CollectEndpoint(client.ExchangesEndpoint, &data.exchanges), "Error collecting Exchange data: %v")
	}
	if args.GlobalArgs.All() || args.GlobalArgs.Inventory {
		warnIfError(client.CollectEndpoint(client.OverviewEndpoint, &data.overview), "Error collecting Overview data: %v")
	}
	return data
}

func getMetricEntities(apiData *allData) []data.EntityData {
	i := 0
	dataItems := make([]data.EntityData, len(apiData.nodes)+len(apiData.exchanges)+len(apiData.queues))
	for _, v := range apiData.nodes {
		dataItems[i] = v
		i++
	}
	for _, v := range apiData.exchanges {
		dataItems[i] = v
		i++
	}
	for _, v := range apiData.queues {
		dataItems[i] = v
		i++
	}
	return dataItems
}

func warnIfError(err error, format string, args ...interface{}) {
	if err != nil {
		log.Warn(format, append(args, err))
	}
}

func exitOnError(err error) {
	if err != nil {
		os.Exit(-1)
	}
}
