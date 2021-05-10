package metrics

import (
	"fmt"
	data2 "github.com/newrelic/nri-rabbitmq/src/data"
	consts2 "github.com/newrelic/nri-rabbitmq/src/data/consts"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

var vhostMetrics = []struct {
	metricName string
	state      string
	sourceType metric.SourceType
}{
	{"vhost.connectionsTotal", "total", metric.GAUGE},
	{"vhost.connectionsStarting", "starting", metric.GAUGE},
	{"vhost.connectionsTuning", "tuning", metric.GAUGE},
	{"vhost.connectionsOpening", "opening", metric.GAUGE},
	{"vhost.connectionsRunning", "running", metric.GAUGE},
	{"vhost.connectionsFlow", "flow", metric.GAUGE},
	{"vhost.connectionsBlocking", "blocking", metric.GAUGE},
	{"vhost.connectionsBlocked", "blocked", metric.GAUGE},
	{"vhost.connectionsClosing", "closing", metric.GAUGE},
	{"vhost.connectionsClosed", "closed", metric.GAUGE},
}

// CollectEntityMetrics ...
func CollectEntityMetrics(rabbitmqIntegration *integration.Integration, bindings []*data2.BindingData, clusterName string, dataItems ...data2.EntityData) {
	bindingStats := collectBindingStats(bindings)
	for _, dataItem := range dataItems {
		if entity, metricNamespace, err := dataItem.GetEntity(rabbitmqIntegration, clusterName); err != nil {
			log.Error("Could not create %s entity [%s]: %v", dataItem.EntityType(), dataItem.EntityName(), err)
		} else if entity != nil {
			metricSet := entity.NewMetricSet(getSampleName(dataItem.EntityType()), metricNamespace...)
			warnIfError(metricSet.MarshalMetrics(dataItem), "Error collecting metrics for [%s:%s]", dataItem.EntityType(), dataItem.EntityName())

			if queue, ok := dataItem.(*data2.QueueData); ok {
				populateBindingMetric(queue.Name, queue.Vhost, consts2.QueueType, metricSet, bindingStats)
				queue.CollectInventory(entity, bindingStats)
			} else if exchange, ok := dataItem.(*data2.ExchangeData); ok {
				populateBindingMetric(exchange.Name, exchange.Vhost, consts2.ExchangeType, metricSet, bindingStats)
				exchange.CollectInventory(entity, bindingStats)
			}
		}
	}
}

// CollectVhostMetrics collects the metrics for VHost entities
func CollectVhostMetrics(rabbitmqIntegration *integration.Integration, vhosts []*data2.VhostData, connections []*data2.ConnectionData, clusterName string) {
	connStats := collectConnectionStats(connections)
	for _, vhost := range vhosts {
		if entity, metricNamespace, err := data2.CreateEntity(rabbitmqIntegration, vhost.Name, consts2.VhostType, vhost.Name, clusterName); err != nil {
			log.Error("Could not create vhost entity [%s]: %v", vhost.Name, err)
		} else if entity != nil {
			metricSet := entity.NewMetricSet(getSampleName(consts2.VhostType), metricNamespace...)
			for _, connStatus := range vhostMetrics {
				connKey := connKey{vhost.Name, connStatus.state}
				setMetric(metricSet, connStatus.metricName, connStats[connKey], connStatus.sourceType)
			}
		}
	}
}

func getSampleName(entityType string) string {
	namespace := entityType
	return fmt.Sprintf("Rabbitmq%sSample", strings.Title(namespace))
}

func warnIfError(err error, format string, args ...interface{}) {
	if err != nil {
		log.Warn(format, append(args, err))
	}
}

func setMetric(metricSet *metric.Set, metricName string, metricValue interface{}, metricType metric.SourceType) {
	if err := metricSet.SetMetric(metricName, metricValue, metricType); err != nil {
		log.Error("There was an error when trying to set metric value: %s", err)
	}
}

func populateBindingMetric(entityName, vhost, entityType string, metricSet *metric.Set, bindingsStats data2.BindingStats) {
	count := 0
	if bindingsStats != nil {
		key := data2.BindingKey{Vhost: vhost, EntityName: entityName, EntityType: entityType}
		if stat := bindingsStats[key]; stat != nil {
			count = len(stat.Destination) + len(stat.Source)
		}
	}
	setMetric(metricSet, entityType+".bindings", count, metric.GAUGE)
}
