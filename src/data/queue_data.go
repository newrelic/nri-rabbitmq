package data

import (
  "strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"
)

// QueueData is the representation of the queues endpoint
type QueueData struct {
	Name                string
	Vhost               string
	Exclusive           bool
	Durable             bool
	Arguments           map[string]interface{}
	AutoDelete          bool     `json:"auto_delete"`
	Consumers           *int64   `metric_name:"queue.consumers" source_type:"gauge"`
	ConsumerUtilisation *float64 `json:"consumer_utilisation" metric_name:"queue.consumerMessageUtilizationPerSecond" source_type:"gauge"`
	ActiveConsumers     *int64   `json:"active_consumers" metric_name:"queue.countActiveConsumersReceiveMessages" source_type:"gauge"`
	Memory              *int64   `metric_name:"queue.erlangBytesConsumedInBytes" source_type:"gauge"`
	Messages            *int64   `metric_name:"queue.totalMessages" source_type:"gauge"`
	MessagesDetails     struct {
		Rate *float64 `metric_name:"queue.totalMessagesPerSecond" source_type:"gauge"`
	} `json:"messages_details"`
	MessagesReady       *int64 `json:"messages_ready" metric_name:"queue.messagesReadyDeliveryClients" source_type:"gauge"`
	MessagesReadyDetail struct {
		Rate *float64 `metric_name:"queue.messagesReadyDeliveryClientsPerSecond" source_type:"gauge"`
	} `json:"messages_ready_details"`
	MessagesUnacknowledged       *int64 `json:"messages_unacknowledged" metric_name:"queue.messagesReadyUnacknowledged" source_type:"gauge"`
	MessagesUnacknowledgedDetail struct {
		Rate *float64 `metric_name:"queue.messagesReadyUnacknowledgedPerSecond" source_type:"gauge"`
	} `json:"messages_unacknowledged_details"`
	MessageStats struct {
		Ack        *int64 `metric_name:"queue.messagesAcknowledged" source_type:"gauge"`
		AckDetails struct {
			Rate *float64 `metric_name:"queue.messagesAcknowledgedPerSecond" source_type:"gauge"`
		} `json:"ack_details"`
		Deliver        *int64 `json:"deliver" metric_name:"queue.messagesDeliveredAckMode" source_type:"gauge"`
		DeliverDetails struct {
			Rate *float64 `metric_name:"queue.messagesDeliveredAckModePerSecond" source_type:"gauge"`
		} `json:"deliver_details"`
		DeliverGet        *int64 `json:"deliver_get" metric_name:"queue.sumMessagesDelivered" source_type:"gauge"`
		DeliverGetDetails struct {
			Rate *float64 `metric_name:"queue.sumMessagesDeliveredPerSecond" source_type:"gauge"`
		} `json:"deliver_get_details"`
		Publish        *int64 `metric_name:"queue.messagesPublished" source_type:"gauge"`
		PublishDetails struct {
			Rate *float64 `metric_name:"queue.messagesPublishedPerSecond" source_type:"gauge"`
		} `json:"publish_details"`
		Redeliver        *int64 `metric_name:"queue.messagesRedeliverGet" source_type:"gauge"`
		RedeliverDetails struct {
			Rate *float64 `metric_name:"queue.messagesRedeliverGetPerSecond" source_type:"gauge"`
		} `json:"redeliver_details"`
	} `json:"message_stats"`
}

// CollectInventory collects inventory data and reports it to the integration.Entity
func (q *QueueData) CollectInventory(entity *integration.Entity, bindingStats BindingStats) {
	SetInventoryItem(entity, strings.TrimPrefix(consts.QueueType, "ra-"), "exclusive", ConvertBoolToInt(q.Exclusive))
	SetInventoryItem(entity, strings.TrimPrefix(consts.QueueType, "ra-"), "durable", ConvertBoolToInt(q.Durable))
	SetInventoryItem(entity, strings.TrimPrefix(consts.QueueType, "ra-"), "auto_delete", ConvertBoolToInt(q.AutoDelete))
	setInventoryMap(entity, strings.TrimPrefix(consts.QueueType, "ra-"), "arguments", q.Arguments)
	setInventoryBindings(entity, q, bindingStats)
}

// GetEntity creates an integration.Entity for this QueueData
func (q *QueueData) GetEntity(integration *integration.Integration, clusterName string) (*integration.Entity, []metric.Attribute, error) {
	return CreateEntity(integration, q.Name, strings.TrimPrefix(consts.QueueType, "ra-"), q.Vhost, clusterName)
}

// EntityType returns the type of this entity
func (q *QueueData) EntityType() string {
	return consts.QueueType
}

// EntityName returns the main name of this entity
func (q *QueueData) EntityName() string {
	return q.Name
}

// EntityVhost returns the vhost of this entity
func (q *QueueData) EntityVhost() string {
	return q.Vhost
}
