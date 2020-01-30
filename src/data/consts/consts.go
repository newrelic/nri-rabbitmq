package consts

const (
	// DefaultExchangeName is the common name to give the exchange with an empty name
	DefaultExchangeName = "amq.default"

	// NodeType name
	NodeType = "node"
	// VhostType name
	VhostType = "vhost"
	// QueueType name
	QueueType = "queue"
	// ExchangeType name
	ExchangeType = "exchange"
	// ClusterType name
	ClusterType = "cluster"
)

var ValidTypes = map[string]struct{}{
	NodeType:     {},
	VhostType:    {},
	QueueType:    {},
	ExchangeType: {},
	ClusterType:  {},
}
