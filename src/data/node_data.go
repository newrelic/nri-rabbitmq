package data

import (
	"encoding/json"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"
)

// NodeData is the representation of the nodes endpoint
type NodeData struct {
	Name                string
	ConfigFiles         []string `json:"config_files"`
	DiskAlarm           *int     `json:"-" metric_name:"node.diskAlarm" source_type:"gauge"`
	DiskFreeSpace       *int64   `json:"disk_free" metric_name:"node.diskSpaceFreeInBytes" source_type:"gauge"`
	FileDescriptorsUsed *int64   `json:"fd_used" metric_name:"node.fileDescriptorsTotalUsed" source_type:"gauge"`
	MemoryAlarm         *int     `json:"-" metric_name:"node.hostMemoryAlarm" source_type:"gauge"`
	MemoryUsed          *int64   `json:"mem_used" metric_name:"node.totalMemoryUsedInBytes" source_type:"gauge"`
	Partitions          int      `json:"-" metric_name:"node.partitionsSeen" source_type:"gauge"`
	Running             *int     `json:"-" metric_name:"node.running" source_type:"gauge"`
	RunQueue            *int64   `json:"run_queue" metric_name:"node.averageErlangProcessesWaiting" source_type:"gauge"`
	SocketsUsed         *int64   `json:"sockets_used" metric_name:"node.fileDescriptorsUsedSockets" source_type:"gauge"`
}

// GetEntity creates an integration.Entity for this NodeData
func (n *NodeData) GetEntity(integration *integration.Integration) (*integration.Entity, []metric.Attribute, error) {
	return CreateEntity(integration, n.Name, consts.NodeType, "")
}

// EntityType returns the type of this entity
func (n *NodeData) EntityType() string {
	return consts.NodeType
}

// EntityName returns the main name of this entity
func (n *NodeData) EntityName() string {
	return n.Name
}

// UnmarshalJSON handles custom JSON Unmarshaling in order to convert values to metrics
func (n *NodeData) UnmarshalJSON(data []byte) error {
	// take from: http://choly.ca/post/go-json-marshalling/
	type Alias NodeData
	aux := &struct {
		DiskAlarm   *bool         `json:"disk_free_alarm"`
		MemoryAlarm *bool         `json:"mem_alarm"`
		Running     *bool         `json:"running"`
		Partitions  []interface{} `json:"partitions"`
		*Alias
	}{
		Alias: (*Alias)(n),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.DiskAlarm != nil {
		x := ConvertBoolToInt(*aux.DiskAlarm)
		n.DiskAlarm = &x
	}
	if aux.MemoryAlarm != nil {
		x := ConvertBoolToInt(*aux.MemoryAlarm)
		n.MemoryAlarm = &x
	}
	if aux.Running != nil {
		x := ConvertBoolToInt(*aux.Running)
		n.Running = &x
	}
	n.Partitions = len(aux.Partitions)
	return nil
}
