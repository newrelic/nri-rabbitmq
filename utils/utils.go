package utils

import (
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/newrelic/nri-rabbitmq/consts"
	"github.com/newrelic/nri-rabbitmq/logger"
)

// CreateEntity will create an entity and metricNamespace attributes with approprate name/namespace values if the entity isn't filtered
func CreateEntity(rabbitmqIntegration *integration.Integration, entityName string, entityType string, vhost string) (entity *integration.Entity, metricNamespace []metric.Attribute, err error) {
	var name, namespace string
	if entityType == consts.ExchangeType && entityName == "" {
		name = consts.DefaultExchangeName
	} else {
		name = entityName
	}
	namespace = entityType

	if !args.GlobalArgs.IncludeEntity(name, entityType, vhost) {
		return nil, nil, nil
	}

	if entityType == consts.QueueType || entityType == consts.ExchangeType {
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

	entity, err = rabbitmqIntegration.Entity(name, namespace)
	return
}

// ConvertBoolToInt converts a boolean to it's metric/inventory representation
func ConvertBoolToInt(val bool) (returnval int) {
	returnval = 0
	if val {
		returnval = 1
	}
	return
}

// CheckErr will invoke the function and the returned error if it's not nil
func CheckErr(f func() error) {
	if err := f(); err != nil {
		logger.Errorf("%v", err)
	}
}

// PanicOnErr pacics if the passed in err is not nil
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
