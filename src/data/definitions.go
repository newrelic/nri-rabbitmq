package data

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

// EntityData is capabile of reporting it's own data to inventory
type EntityData interface {
	GetEntity(integration *integration.Integration) (*integration.Entity, []metric.Attribute, error)
	EntityName() string
	EntityType() string
}

// OverviewData is the representation of the overview endpoint
type OverviewData struct {
	ClusterName       string `json:"cluster_name"`
	RabbitMQVersion   string `json:"rabbitmq_version"`
	ManagementVersion string `json:"management_version"`
}

// ConnectionData is the representation of the connections endpoint
type ConnectionData struct {
	Vhost string
	State string
}

// BindingData is the representation of the bindings endpoint
type BindingData struct {
	Vhost           string
	Source          string
	Destination     string
	DestinationType string `json:"destination_type"`
}

// VhostData is the representation of the vhosts endpoint
type VhostData struct {
	Name string
}
