package data

import (
	"github.com/newrelic/infra-integrations-sdk/v3/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/v3/integration"
)

// EntityData is capable of reporting it's own data to inventory
type EntityData interface {
	GetEntity(integration *integration.Integration, entityName string) (*integration.Entity, []attribute.Attribute, error)
	EntityVhost() string
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

// BindingKey is used to uniquely identify a binding by Vhost, EntityName, and EntityType
type BindingKey struct {
	Vhost, EntityName, EntityType string
}

// Binding contains a list of Source and Destination BindingKeys
type Binding struct {
	Source      []*BindingKey
	Destination []*BindingKey
}

// BindingStats contains the calculation for Source/Destination Binding for each entity/BindingKey
type BindingStats map[BindingKey]*Binding

// VhostData is the representation of the vhosts endpoint
type VhostData struct {
	Name string
}

// VhostTest holds data around a test against a Vhost
type VhostTest struct {
	Vhost *VhostData
	Test  *TestData
}

// NodeTest holds data around a test against a Node
type NodeTest struct {
	Node *NodeData
	Test *TestData
}

// TestData is the representation of both the AlivenessTest and Healthchecks endpoints
type TestData struct {
	Status string
	Reason string
}
