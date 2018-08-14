package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/newrelic/nri-rabbitmq/data/consts"
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
			name = fmt.Sprintf("%s/%s", vhost, name)
		}
	}
	metricNamespace = []metric.Attribute{
		{Key: "displayName", Value: name},
		{Key: "entityName", Value: fmt.Sprintf("%s:%s", namespace, name)},
	}

	entity, err = rabbitmqIntegration.Entity(name, namespace)
	return
}

// SetInventoryItem sets an inventory item in a consistent way
func SetInventoryItem(entity *integration.Entity, category, key string, value interface{}) {
	if entity != nil && key != "" && value != nil {
		if category != "" {
			key = category + "/" + key
		}
		if err := entity.SetInventoryItem(key, "value", value); err != nil {
			if entity.Metadata == nil {
				log.Warn("Error setting inventory [%s] on LocalEntity: %v", key, err)
			} else {
				log.Warn("Error setting inventory [%s] on [%s]: %v", key, entity.Metadata.Name, err)
			}
		}
	}
}

// setInventoryMap sets an inventory map in a consistent way
func setInventoryMap(entity *integration.Entity, category, key string, value map[string]interface{}) {
	if entity != nil && key != "" && value != nil {
		if category != "" {
			key = category + "/" + key
		}
		for k, v := range value {
			if arrayVal, ok := v.([]interface{}); ok {
				setInventoryArray(entity, key, k, arrayVal)
			} else {
				if err := entity.SetInventoryItem(key, k, v); err != nil {
					logInventoryErr(entity, err, key)
				}
			}
		}
	}
}

func setInventoryArray(entity *integration.Entity, key, field string, v []interface{}) {
	if arrayJSON, err := json.Marshal(v); err != nil {
		logInventoryErr(entity, err, key)
	} else {
		if err := entity.SetInventoryItem(key, field, string(arrayJSON)); err != nil {
			logInventoryErr(entity, err, key)
		}
	}
}

func logInventoryErr(entity *integration.Entity, err error, key string) {
	if entity.Metadata == nil {
		log.Warn("Error setting inventory [%s] on LocalEntity: %v", key, err)
	} else {
		log.Warn("Error setting inventory [%s] on [%s]: %v", key, entity.Metadata.Name, err)
	}
}
