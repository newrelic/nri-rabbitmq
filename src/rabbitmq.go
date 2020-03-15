//go:generate goversioninfo
package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-rabbitmq/src/args"
	"github.com/newrelic/nri-rabbitmq/src/client"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/inventory"
	"github.com/newrelic/nri-rabbitmq/src/metrics"
)

const (
	integrationName    = "com.newrelic.rabbitmq"
	integrationVersion = "2.2.2"
)

func main() {
	var argList args.ArgumentList
	// Create Integration
	rabbitmqIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&argList))
	exitOnError(err)

	exitOnError(args.SetGlobalArgs(argList))

	log.SetupLogging(args.GlobalArgs.Verbose)

	exitOnError(args.GlobalArgs.Validate())

	rabbitData := getNeededData()
	clusterName := rabbitData.overview.ClusterName

	if args.GlobalArgs.HasMetrics() && args.GlobalArgs.HasInventory() {
		metrics.CollectVhostMetrics(rabbitmqIntegration, rabbitData.vhosts, rabbitData.connections, clusterName)

		metricEntities := getMetricEntities(rabbitData)
		metrics.CollectEntityMetrics(rabbitmqIntegration, rabbitData.bindings, clusterName, metricEntities...)

		inventory.PopulateClusterInventory(rabbitmqIntegration, rabbitData.overview)
	}

	if args.GlobalArgs.HasInventory() {
		inventory.CollectInventory(rabbitmqIntegration, rabbitData.nodes, clusterName)
	}

	if args.GlobalArgs.HasEvents() {
		alivenessTest(rabbitmqIntegration, rabbitData.aliveness, clusterName)
		healthcheckTest(rabbitmqIntegration, rabbitData.healthcheck, clusterName)
	}

	if len(rabbitmqIntegration.Entities) > 0 {
		err = rabbitmqIntegration.Publish()
		if err != nil {
			log.Error("Error publishing integration: %v", err)
			exitOnError(err)
		}
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
	healthcheck []*data.NodeTest
	aliveness   []*data.VhostTest
}

func getNeededData() *allData {
	rabbitData := new(allData)
	warnIfError(client.CollectEndpoint(client.NodesEndpoint, &rabbitData.nodes), "Error collecting Node data: %v")
	warnIfError(client.CollectEndpoint(client.OverviewEndpoint, &rabbitData.overview), "Error collecting Overview data: %v")
	if args.GlobalArgs.HasMetrics() {
		warnIfError(client.CollectEndpoint(client.ConnectionsEndpoint, &rabbitData.connections), "Error collecting Connections data: %v")
		warnIfError(client.CollectEndpoint(client.BindingsEndpoint, &rabbitData.bindings), "Error collecting Bindings data: %v")
		warnIfError(client.CollectEndpoint(client.VhostsEndpoint, &rabbitData.vhosts), "Error collecting Vhost data: %v")
		warnIfError(client.CollectEndpoint(client.QueuesEndpoint, &rabbitData.queues), "Error collecting Queue data: %v")
		warnIfError(client.CollectEndpoint(client.ExchangesEndpoint, &rabbitData.exchanges), "Error collecting Exchange data: %v")
	} else if args.GlobalArgs.HasEvents() {
		warnIfError(client.CollectEndpoint(client.VhostsEndpoint, &rabbitData.vhosts), "Error collecting Vhost data: %v")
	}
	if args.GlobalArgs.HasEvents() {
		getEventData(rabbitData)
	}
	return rabbitData
}

func getEventData(rabbitData *allData) {
	var endpoint string
	if len(rabbitData.nodes) > 0 {
		rabbitData.healthcheck = make([]*data.NodeTest, len(rabbitData.nodes))
		for i, node := range rabbitData.nodes {
			nodeTest := &data.NodeTest{
				Node: node,
				Test: new(data.TestData),
			}
			endpoint = fmt.Sprintf(client.HealthCheckEndpoint, url.PathEscape(node.Name))
			if err := client.CollectEndpoint(endpoint, nodeTest.Test); err != nil {
				nodeTest.Test.Status = "error"
				nodeTest.Test.Reason = err.Error()
			}
			rabbitData.healthcheck[i] = nodeTest
		}
	}

	if len(rabbitData.vhosts) > 0 {
		rabbitData.aliveness = make([]*data.VhostTest, len(rabbitData.vhosts))
		for i, vhost := range rabbitData.vhosts {
			vhostTest := &data.VhostTest{
				Vhost: vhost,
				Test:  new(data.TestData),
			}
			endpoint = fmt.Sprintf(client.AlivenessTestEndpoint, url.PathEscape(vhost.Name))
			if err := client.CollectEndpoint(endpoint, vhostTest.Test); err != nil {
				vhostTest.Test.Status = "error"
				vhostTest.Test.Reason = err.Error()
			}
			rabbitData.aliveness[i] = vhostTest
		}
	}
}

// maxQueues is the maximum amount of Queues that can be collect.
// If there are more than this number of Queues then collection of
// Queues will fail.
const maxQueues = 500

func getMetricEntities(apiData *allData) []data.EntityData {
	i := 0
	// Make the length the size of nodes and exchanges but capacity the length + size of queues. This is to accommodate the chance that there are more
	// queues than can be collected.
	dataItems := make([]data.EntityData, len(apiData.nodes)+len(apiData.exchanges), len(apiData.nodes)+len(apiData.exchanges)+len(apiData.queues))

	for _, v := range apiData.nodes {
		dataItems[i] = v
		i++
	}
	for _, v := range apiData.exchanges {
		dataItems[i] = v
		i++
	}

	if queueLength := getFilteredQueueCount(apiData.queues); queueLength > maxQueues {
		log.Error("There are %d queues in collection, the maximum amount of queues to collect is %d. Use the queue whitelist or regex configuration parameter to limit collection size.", queueLength, maxQueues)
		return dataItems
	}

	for _, v := range apiData.queues {
		dataItems = append(dataItems, v)
	}
	return dataItems
}

func getFilteredQueueCount(queuesData []*data.QueueData) int {
	queueCount := 0
	for _, queueData := range queuesData {
		if args.GlobalArgs.IncludeEntity(queueData.Name, "queue", queueData.Vhost) {
			queueCount++
		}
	}

	return queueCount
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
