package data

import (
	"encoding/json"
	"math/big"

	"github.com/newrelic/nri-rabbitmq/src/data/consts"

	"github.com/newrelic/infra-integrations-sdk/v3/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/v3/integration"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
)

// NodeData is the representation of the nodes endpoint
type NodeData struct {
	Name                 string
	ConfigFiles          []string `json:"config_files"`
	DiskAlarm            *bool    `json:"disk_free_alarm" metric_name:"node.diskAlarm" source_type:"gauge"`
	DiskFreeSpace        *int64   `json:"disk_free" metric_name:"node.diskSpaceFreeInBytes" source_type:"gauge"`
	FileDescriptorsUsed  *int64   `json:"fd_used" metric_name:"node.fileDescriptorsTotalUsed" source_type:"gauge"`
	FileDescriptorsTotal *int64   `json:"fd_total" metric_name:"node.fileDescriptorsTotal" source_type:"gauge"`
	ProcessesTotal       *int64   `json:"proc_total" metric_name:"node.processesTotal" source_type:"gauge"`
	ProcessesUsed        *int64   `json:"proc_used" metric_name:"node.processesUsed" source_type:"gauge"`
	MemoryAlarm          *bool    `json:"mem_alarm" metric_name:"node.hostMemoryAlarm" source_type:"gauge"`
	MemoryUsed           *int64   `json:"mem_used" metric_name:"node.totalMemoryUsedInBytes" source_type:"gauge"`
	Partitions           int      `json:"-" metric_name:"node.partitionsSeen" source_type:"gauge"`
	Running              *bool    `metric_name:"node.running" source_type:"gauge"`
	RunQueue             *int64   `json:"run_queue" metric_name:"node.averageErlangProcessesWaiting" source_type:"gauge"`
	SocketsTotal         *int64   `json:"sockets_total" metric_name:"node.fileDescriptorsTotalSockets" source_type:"gauge"`
	SocketsUsed          *int64   `json:"sockets_used" metric_name:"node.fileDescriptorsUsedSockets" source_type:"gauge"`
}

// GetEntity creates an integration.Entity for this NodeData
func (n *NodeData) GetEntity(integration *integration.Integration, clusterName string) (*integration.Entity, []attribute.Attribute, error) {
	return CreateEntity(integration, n.Name, consts.NodeType, "", clusterName)
}

// EntityType returns the type of this entity
func (n *NodeData) EntityType() string {
	return consts.NodeType
}

// EntityName returns the main name of this entity
func (n *NodeData) EntityName() string {
	return n.Name
}

// EntityVhost returns the vhost of this entity
func (n *NodeData) EntityVhost() string {
	return ""
}

// UnmarshalJSON handles custom JSON Unmarshaling in order to convert values to metrics
func (n *NodeData) UnmarshalJSON(data []byte) error {
	// take from: http://choly.ca/post/go-json-marshalling/
	type Alias NodeData
	aux := &struct {
		Partitions []interface{} `json:"partitions"`
		*Alias
		DiskFreeSpace *big.Int `json:"disk_free" metric_name:"node.diskSpaceFreeInBytes" source_type:"gauge"`
	}{
		Alias: (*Alias)(n),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	n.Partitions = len(aux.Partitions)
	if aux.DiskFreeSpace == nil {
		return nil
	}
	if aux.DiskFreeSpace.IsInt64() {
		diskFree := aux.DiskFreeSpace.Int64()
		n.DiskFreeSpace = &diskFree
	} else {
		log.Warn("Node's disk_free value is too high to be reported (%v), ignoring it", aux.DiskFreeSpace)
	}
	return nil
}
