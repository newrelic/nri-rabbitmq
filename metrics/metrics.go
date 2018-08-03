package metrics

import (
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-rabbitmq/client"
	"github.com/newrelic/nri-rabbitmq/inventory"
	"github.com/newrelic/nri-rabbitmq/logger"
	"github.com/newrelic/nri-rabbitmq/utils"
	"github.com/newrelic/nri-rabbitmq/utils/consts"
	"github.com/stretchr/objx"
)

var metricDefinitions = map[string]map[string]struct {
	name       string
	metricType metric.SourceType
}{
	// map of metric name to json key and metric info
	// these metrics come from /api/nodes
	"node": {
		"node.fileDescriptorsTotalUsed":      {"fd_used", metric.GAUGE},
		"node.diskSpaceFreeInBytes":          {"disk_free", metric.GAUGE},
		"node.totalMemoryUsedInBytes":        {"mem_used", metric.GAUGE},
		"node.averageErlangProcessesWaiting": {"run_queue", metric.GAUGE},
		"node.fileDescriptorsUsedSockets":    {"sockets_used", metric.GAUGE},
		"node.partitionsSeen":                {"partitions", metric.GAUGE},
		"node.running":                       {"running", metric.GAUGE},
		"node.hostMemoryAlarm":               {"mem_alarm", metric.GAUGE},
		"node.diskAlarm":                     {"disk_free_alarm", metric.GAUGE},
	},

	"queue": {
		"queue.consumers":                             {"consumers", metric.GAUGE},
		"queue.consumerMessageUtilizationPerSecond":   {"consumer_utilisation", metric.GAUGE},
		"queue.countActiveConsumersReceiveMessages":   {"active_consumers", metric.GAUGE},
		"queue.erlangBytesConsumedInBytes":            {"memory", metric.GAUGE},
		"queue.totalMessages":                         {"messages", metric.GAUGE},
		"queue.totalMessagesPerSecond":                {"messages_details.rate", metric.GAUGE},
		"queue.messagesReadyDeliveryClients":          {"messages_ready", metric.GAUGE},
		"queue.messagesReadyDeliveryClientsPerSecond": {"messages_ready_details.rate", metric.GAUGE},
		"queue.messagesReadyUnacknowledged":           {"messages_unacknowledged", metric.GAUGE},
		"queue.messagesReadyUnacknowledgedPerSecond":  {"messages_unacknowledged_details.rate", metric.GAUGE},
		"queue.messagesAcknowledged":                  {"message_stats.ack", metric.GAUGE},
		"queue.messagesAcknowledgedPerSecond":         {"message_stats.ack_details.rate", metric.GAUGE},
		"queue.messagesDeliveredAckMode":              {"message_stats.deliver", metric.GAUGE},
		"queue.messagesDeliveredAckModePerSecond":     {"message_stats.deliver_details.rate", metric.GAUGE},
		"queue.sumMessagesDelivered":                  {"message_stats.deliver_get", metric.GAUGE},
		"queue.sumMessagesDeliveredPerSecond":         {"message_stats.deliver_get_details.rate", metric.GAUGE},
		"queue.messagesPublished":                     {"message_stats.publish", metric.GAUGE},
		"queue.messagesPublishedPerSecond":            {"message_stats.publish_details.rate", metric.GAUGE},
		"queue.messagesRedeliverGet":                  {"message_stats.redeliver", metric.GAUGE},
		"queue.messagesRedeliverGetPerSecond":         {"message_stats.redeliver_details.rate", metric.GAUGE},
	},

	"exchange": {
		"exchange.messagesPublishedPerChannel":          {"message_stats.publish_in", metric.GAUGE},
		"exchange.messagesPublishedPerChannelPerSecond": {"message_stats.publish_in_details.rate", metric.GAUGE},
		"exchange.messagesPublishedQueue":               {"message_stats.publish_out", metric.GAUGE},
		"exchange.messagesPublishedQueuePerSecond":      {"message_stats.publish_out_details.rate", metric.GAUGE},
	},

	"vhost": {
		"vhost.connectionsTotal":    {"total", metric.GAUGE},
		"vhost.connectionsStarting": {"starting", metric.GAUGE},
		"vhost.connectionsTuning":   {"tuning", metric.GAUGE},
		"vhost.connectionsOpening":  {"opening", metric.GAUGE},
		"vhost.connectionsRunning":  {"running", metric.GAUGE},
		"vhost.connectionsFlow":     {"flow", metric.GAUGE},
		"vhost.connectionsBlocking": {"blocking", metric.GAUGE},
		"vhost.connectionsBlocked":  {"blocked", metric.GAUGE},
		"vhost.connectionsClosing":  {"closing", metric.GAUGE},
		"vhost.connectionsClosed":   {"closed", metric.GAUGE},
	},
}

var entityMappings = []struct {
	entityType     string
	entityEndpoint string
}{
	{consts.NodeType, client.NodesEndpoint},
	{consts.VhostType, client.VhostsEndpoint},
	{consts.QueueType, client.QueuesEndpoint},
	{consts.ExchangeType, client.ExchangesEndpoint},
}

// CollectMetrics collects the metrics present in the apiResponses
func CollectMetrics(rabbitmqIntegration *integration.Integration, apiResponses *objx.Map) {
	connStats := collectConnectionStats(apiResponses.Get(client.ConnectionsEndpoint).ObjxMapSlice())
	bindingStats := collectBindingStats(apiResponses.Get(client.BindingsEndpoint).ObjxMapSlice())

	for _, mapping := range entityMappings {
		responseObjects := apiResponses.Get(mapping.entityEndpoint).ObjxMapSlice()
		for _, entityData := range responseObjects {
			parseEntityResponse(rabbitmqIntegration, &entityData, mapping.entityType, connStats, bindingStats)
		}
	}
}

func parseEntityResponse(rabbitmqIntegration *integration.Integration, entityData *objx.Map, entityType string, connectionStats map[connKey]int, bindingStats map[bindingKey]int) {
	entityName := entityData.Get("name").Str()
	vhost := entityData.Get("vhost").Str()
	entity, metricNamespace, err := utils.CreateEntity(rabbitmqIntegration, entityName, entityType, vhost)
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

	if entityType == consts.VhostType {
		populateVhostMetrics(entityName, metricSet, connectionStats)
	} else {
		populateMetrics(metricSet, entityType, entityData)
		inventory.PopulateEntityInventory(entity, entityType, entityData)

		if entityType == consts.QueueType || entityType == consts.ExchangeType {
			populateBindingMetric(entityName, vhost, entityType, metricSet, bindingStats)
		}
	}
}

func setMetric(metricSet *metric.Set, metricName string, metricValue interface{}, metricType metric.SourceType) {
	if err := metricSet.SetMetric(metricName, metricValue, metricType); err != nil {
		logger.Errorf("There was an error when trying to set metric value: %s", err)
	}
}

func populateMetrics(metricSet *metric.Set, entityType string, response *objx.Map) {
	notFoundMetrics := make([]string, 0)
	for metricKey, metricInfo := range metricDefinitions[entityType] {
		metricInfoValue, err := parseJSON(response, metricInfo.name)
		if err != nil {
			notFoundMetrics = append(notFoundMetrics, metricKey)
		}
		if metricInfoValue != nil {
			setMetric(metricSet, metricKey, metricInfoValue, metricInfo.metricType)
		}
	}
	if len(notFoundMetrics) > 0 {
		logger.Debugf("Can't find raw metrics in results for keys: %v", notFoundMetrics)
	}
}

func populateVhostMetrics(vhostName string, metricSet *metric.Set, connStats map[connKey]int) {
	for metricKey, metricInfo := range metricDefinitions[consts.VhostType] {
		connKey := connKey{vhostName, metricInfo.name}
		setMetric(metricSet, metricKey, connStats[connKey], metricInfo.metricType)
	}
}

func populateBindingMetric(entityName, vhost, entityType string, metricSet *metric.Set, bindingsStats map[bindingKey]int) {
	bindingKey := bindingKey{
		Vhost:      vhost,
		EntityName: entityName,
		EntityType: entityType,
	}
	setMetric(metricSet, entityType+".bindings", bindingsStats[bindingKey], metric.GAUGE)
}

func parseJSON(jsonData *objx.Map, key string) (interface{}, error) {
	value := jsonData.Get(key)

	if value.IsInterSlice() {
		return len(value.InterSlice()), nil
	} else if value.IsFloat64() {
		return value.Float64(), nil
	} else if value.IsBool() {
		return utils.ConvertBoolToInt(value.Bool()), nil
	}

	return nil, fmt.Errorf("could not parse json value for key [%v], unexpected data type", key)
}
