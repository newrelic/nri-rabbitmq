package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/stretchr/objx"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Hostname         string `default:"localhost" help:"Hostname or IP where RabbitMQ Management Plugin is running."`
	Port             int    `default:"15672" help:"Port on which RabbitMQ Management Plugin is listening."`
	Username         string `default:"" help:"Username for accessing RabbitMQ Management Plugin"`
	Password         string `default:"" help:"Password for the given user."`
	CABundleFile     string `default:"" help:"Alternative Certificate Authority bundle file"`
	CABundleDir      string `default:"" help:"Alternative Certificate Authority bundle directory"`
	NodeNameOverride string `default:"" help:"Overrides the local node name instead of retrieving it from RabbitMQ."`
	ConfigPath       string `default:"" help:"RabbitMQ configuration file path."`
	UseSSL           bool   `default:"false" help:"configure whether to use an SSL connection or not."`
	Queues           string `default:"" help:"JSON array of queue names from which to collect metrics."`
	QueuesRegexes    string `default:"" help:"JSON array of queue name regexes from which to collect metrics."`
	Exchanges        string `default:"" help:"JSON array of exchange names from which to collect metrics."`
	ExchangesRegexes string `default:"" help:"JSON array of exchange name regexes from which to collect metrics."`
	Vhosts           string `default:"" help:"JSON array of vhost names from which to collect metrics."`
	VhostsRegexes    string `default:"" help:"JSON array of vhost name regexes from which to collect metrics."`
}

const (
	integrationName    = "com.newrelic.rabbitmq"
	integrationVersion = "0.1.0"

	defaultExchangeName = "amq.default"

	nodeType     = "node"
	vhostType    = "vhost"
	queueType    = "queue"
	exchangeType = "exchange"
	clusterType  = "cluster"
)

var (
	args   argumentList
	logger log.Logger

	entityMappings = []struct {
		entityType     string
		entityEndpoint string
	}{
		{nodeType, nodesEndpoint},
		{vhostType, vhostsEndpoint},
		{queueType, queuesEndpoint},
		{exchangeType, exchangesEndpoint},
	}

	queueFilter, queueRegexFilter, exchangeFilter, exchangeRegexFilter, vhostFilter, vhostRegexFilter []string
)

func main() {

	// Create Integration
	i, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	panicOnErr(err)
	logger = i.Logger()

	err = checkArgs()
	panicOnErr(err)

	err = parseFilterArgs()
	panicOnErr(err)

	var localNode *integration.Entity

	allEndpointResponses, err := collectEndpoints(getListOfEndpointsToCollect()...)
	panicOnErr(err)

	responses := objx.New(allEndpointResponses)
	if responses == nil {
		panic(errors.New("unexpected management response"))
	}

	if args.All() || args.Metrics {
		logger.Infof("Collecting metrics.")
		collectMetrics(i, &responses)
	}

	if args.All() || args.Inventory {
		logger.Infof("Collecting inventory.")
		collectInventory(i, &responses, localNode)
	}

	if args.All() || (args.Metrics && args.Inventory) {
		populateClusterEntity(i, responses.Get(overviewEndpoint).ObjxMap())
	}

	panicOnErr(i.Publish())
}

func checkArgs() error {
	if args.Metrics && !args.Inventory {
		return errors.New("invalid arguments: can't collect metrics while not collecting inventory")
	} else if !args.All() && !args.Metrics && !args.Inventory {
		return errors.New("invalid arguments: nothing specified to collect")
	}
	return nil
}

func getListOfEndpointsToCollect() (endpoints []string) {
	if args.All() {
		endpoints = allEndpoints
	} else if args.Inventory {
		endpoints = inventoryEndpoints
	}
	return
}

func parseFilterArgs() error {
	argMappings := []struct {
		argString string
		parsed    *[]string
	}{
		{args.Queues, &queueFilter},
		{args.QueuesRegexes, &queueRegexFilter},
		{args.Exchanges, &exchangeFilter},
		{args.ExchangesRegexes, &exchangeRegexFilter},
		{args.Vhosts, &vhostFilter},
		{args.VhostsRegexes, &vhostRegexFilter},
	}

	for _, filterPair := range argMappings {
		// parse arg string, populate parsed field
		if filterPair.argString == "" {
			continue
		}

		var result []string
		err := json.Unmarshal([]byte(filterPair.argString), &result)
		if err != nil {
			return err
		}
		*filterPair.parsed = result
	}
	return nil
}

func collectMetrics(integration *integration.Integration, responses *objx.Map) {
	connStats := collectConnectionStats(responses.Get(connectionsEndpoint).ObjxMapSlice())
	bindingStats := collectBindingStats(responses.Get(bindingsEndpoint).ObjxMapSlice())

	for _, mapping := range entityMappings {
		responseObjects := responses.Get(mapping.entityEndpoint).ObjxMapSlice()
		for _, entityData := range responseObjects {
			parseEntityResponse(integration, &entityData, mapping.entityType, connStats, bindingStats)
		}
	}
}

func parseEntityResponse(integration *integration.Integration, entityData *objx.Map, entityType string, connectionStats map[connKey]int, bindingStats map[bindingKey]int) {
	entityName := entityData.Get("name").Str()
	vhost := entityData.Get("vhost").Str()
	entity, metricNamespace, err := createEntity(integration, entityName, entityType, vhost)
	if entity == nil {
		logger.Infof("%v entity %v did not match filter - skipping", entityType, entityName)
		return
	}
	if err != nil {
		logger.Errorf("could not create entity [%v, %v, %v]: %v", entityName, entityType, vhost, err)
		return
	}
	sampleName := fmt.Sprintf("Rabbitmq%vSample", strings.Title(entityType))
	metricSet := entity.NewMetricSet(sampleName, metricNamespace...)

	if entityType == vhostType {
		populateVhostMetrics(entityName, metricSet, connectionStats)
	} else {
		populateMetrics(metricSet, entityType, entityData)
		populateEntityInventory(entity, entityType, entityData)

		if entityType == queueType || entityType == exchangeType {
			populateBindingMetric(entityName, vhost, entityType, metricSet, bindingStats)
		}
	}
}

func collectInventory(integration *integration.Integration, responses *objx.Map, localNode *integration.Entity) {
	nodeName, err := getNodeName()
	panicOnErr(err)

	overviewData := responses.Get(overviewEndpoint).ObjxMap()
	if localNode == nil {
		nodesList := responses.Get(nodesEndpoint).ObjxMapSlice()

		if nodesList != nil {
			localNode, err = getNodeEntity(nodeName, nodesList, integration)
		} else {
			logger.Warnf("Could not retrieve node data")
			localNode, err = getNodeEntity(nodeName, nil, integration)
		}
	}
	if err == nil && localNode != nil {
		if err = setInventoryData(localNode, overviewData); err != nil {
			logger.Errorf("%v", err)
		}
	}
}

func createEntity(integration *integration.Integration, entityName string, entityType string, vhost string) (entity *integration.Entity, metricNamespace []metric.Attribute, err error) {
	var name, namespace string
	if entityType == exchangeType && entityName == "" {
		name = defaultExchangeName
	} else {
		name = entityName
	}
	namespace = entityType

	if !includeEntity(name, entityType, vhost) {
		return nil, nil, nil
	}

	if entityType == queueType || entityType == exchangeType {
		if strings.HasSuffix(vhost, "/") {
			name = vhost + name
		} else {
			name = fmt.Sprintf("%v/%v", vhost, name)
		}
	}
	metricNamespace = []metric.Attribute{
		{Key: "displayName", Value: name},
		{Key: "entityName", Value: fmt.Sprintf("%v:%v", namespace, name)},
	}

	entity, err = integration.Entity(name, namespace)
	return
}

func includeEntity(entityName, entityType, vhost string) bool {
	vhostIncluded := vhost == "" || checkFilters(vhost, vhostFilter, vhostRegexFilter)
	if !vhostIncluded {
		return false
	}

	if entityType == queueType {
		return checkFilters(entityName, queueFilter, queueRegexFilter)
	} else if entityType == exchangeType {
		return checkFilters(entityName, exchangeFilter, exchangeRegexFilter)
	} else {
		return true
	}
}

func checkFilters(entityName string, nameList, nameRegexList []string) bool {
	if len(nameList) == 0 && len(nameRegexList) == 0 {
		return true
	}

	for _, name := range nameList {
		if entityName == name {
			return true
		}
	}
	for _, nameRegex := range nameRegexList {
		matched, err := regexp.MatchString(nameRegex, entityName)
		if matched && err == nil {
			return true
		}
	}

	return false
}

func checkErr(f func() error) {
	if err := f(); err != nil {
		logger.Errorf("%v", err)
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
