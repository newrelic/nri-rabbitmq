package main

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
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
		"queue.consumerMessageUtilizationPerSecond":   {"consumer_utilisation", metric.RATE},
		"queue.countActiveConsumersReceiveMessages":   {"active_consumers", metric.GAUGE},
		"queue.erlangBytesConsumedInBytes":            {"memory", metric.GAUGE},
		"queue.totalMessages":                         {"messages", metric.GAUGE},
		"queue.totalMessagesPerSecond":                {"messages_details.rate", metric.RATE},
		"queue.messagesReadyDeliveryClients":          {"messages_ready", metric.GAUGE},
		"queue.messagesReadyDeliveryClientsPerSecond": {"messages_ready_details.rate", metric.RATE},
		"queue.messagesReadyUnacknowledged":           {"messages_unacknowledged", metric.GAUGE},
		"queue.messagesReadyUnacknowledgedPerSecond":  {"messages_unacknowledged_details.rate", metric.RATE},
		"queue.messagesAcknowledged":                  {"message_stats.ack", metric.GAUGE},
		"queue.messagesAcknowledgedPerSecond":         {"message_stats.ack_details.rate", metric.RATE},
		"queue.messagesDeliveredAckMode":              {"message_stats.deliver", metric.GAUGE},
		"queue.messagesDeliveredAckModePerSecond":     {"message_stats.deliver_details.rate", metric.RATE},
		"queue.sumMessagesDelivered":                  {"message_stats.deliver_get", metric.GAUGE},
		"queue.sumMessagesDeliveredPerSecond":         {"message_stats.deliver_get_details.rate", metric.RATE},
		"queue.messagesPublished":                     {"message_stats.publish", metric.GAUGE},
		"queue.messagesPublishedPerSecond":            {"message_stats.publish_details.rate", metric.RATE},
		"queue.messagesRedeliverGet":                  {"message_stats.redeliver", metric.GAUGE},
		"queue.messagesRedeliverGetPerSecond":         {"message_stats.redeliver_details.rate", metric.RATE},
	},

	"exchange": {
		"exchange.messagesPublishedQueue":          {"message_stats.publish_out", metric.GAUGE},
		"exchange.messagesPublishedQueuePerSecond": {"message_stats.publish_out_details.rate", metric.RATE},
		"exchange.messagesPublished":               {"message_stats.publish_in", metric.GAUGE},
		"exchange.messagesPublishedPerSecond":      {"message_stats.publish_in_details.rate", metric.RATE},
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
	for metricKey, metricInfo := range metricDefinitions[vhostType] {
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

	if value.IsFloat64() {
		return value.Float64(), nil
	} else if value.IsBool() {
		return convertBoolToInt(value.Bool()), nil
	}

	return nil, fmt.Errorf("could not parse json value for key [%v], unexpected data type", key)
}

func convertBoolToInt(val bool) (returnval int) {
	returnval = 0
	if val {
		returnval = 1
	}
	return
}
