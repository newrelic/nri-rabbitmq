package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-rabbitmq/src/args"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"
)

// CreateEntity will create an entity and metricNamespace attributes with appropriate name/namespace values if the entity isn't filtered
func CreateEntity(rabbitmqIntegration *integration.Integration, entityName, entityType, vhost, clusterName string) (entity *integration.Entity, metricNamespace []metric.Attribute, err error) {
	name := cleanEntityName(entityName, entityType)
	namespace := entityType

	if !args.GlobalArgs.IncludeEntity(name, entityType, vhost) {
		return nil, nil, nil
	}

	if entityType == consts.QueueType || entityType == consts.ExchangeType {
		name = joinVhostName(vhost, name)
	}
	metricNamespace = []metric.Attribute{
		{Key: "displayName", Value: name},
		{Key: "entityName", Value: fmt.Sprintf("%s:%s", strings.TrimPrefix(namespace, "ra-"), name)},
	}

	clusterNameAttribute := integration.IDAttribute{Key: "clusterName", Value: clusterName}
	endpoint := fmt.Sprintf("%s:%d", args.GlobalArgs.Hostname, args.GlobalArgs.Port)
	if err != nil {
		return nil, nil, err
	}

	entity, err = rabbitmqIntegration.EntityReportedVia(endpoint, name, namespace, clusterNameAttribute)
	return
}

func cleanEntityName(entityName, entityType string) string {
	if entityType == consts.ExchangeType && entityName == "" {
		return consts.DefaultExchangeName
	}
	return entityName
}

func joinVhostName(vhost, name string) string {
	if strings.HasSuffix(vhost, "/") {
		return vhost + name
	}
	return vhost + "/" + name
}

// SetInventoryItem sets an inventory item in a consistent way
func SetInventoryItem(entity *integration.Entity, category, key string, value interface{}) {
	if entity == nil || key == "" || value == nil {
		return
	}
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

func setInventoryBindings(entity *integration.Entity, data EntityData, bindingStats BindingStats) {
	if bindingStats != nil {
		if stat := bindingStats[BindingKey{data.EntityVhost(), data.EntityName(), data.EntityType()}]; stat != nil {
			if len(stat.Source) > 0 {
				SetInventoryItem(entity, strings.TrimPrefix(data.EntityType(), "ra-"), "bindings.source", getKeyList(stat.Source))
			}
			if len(stat.Destination) > 0 {
				SetInventoryItem(entity, strings.TrimPrefix(data.EntityType(), "ra-"), "bindings.destination", getKeyList(stat.Destination))
			}
		}
	}
}

func getKeyList(keys []*BindingKey) string {
	names := []string{}
	for _, v := range keys {
		name := v.EntityType + ":" + joinVhostName(v.Vhost, cleanEntityName(v.EntityName, v.EntityType))
		names = append(names, name)
	}
	return strings.Join(names, ", ")
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
