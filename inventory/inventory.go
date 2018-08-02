package inventory

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/newrelic/nri-rabbitmq/logger"
	"github.com/newrelic/nri-rabbitmq/utils"
	"github.com/newrelic/nri-rabbitmq/utils/consts"
	"github.com/stretchr/objx"
)

var (
	execCommand     = exec.Command
	osOpen          = os.Open
	entityInventory = map[string][]string{
		consts.ExchangeType: {
			"type",
			"durable",
			"auto_delete",
			"arguments",
		},
		consts.QueueType: {
			"exclusive",
			"durable",
			"auto_delete",
			"arguments",
		},
	}
)

// CollectInventory collects the inventory items (config file values) from the apiResponses
func CollectInventory(rabbitmqIntegration *integration.Integration, nodesData []objx.Map) {
	if len(nodesData) == 0 {
		panic(errors.New("could not retrieve the list of nodes"))
	}

	nodeName, err := getLocalNodeName()
	utils.PanicOnErr(err)

	localNode, nodeData, err := getNodeEntity(nodeName, nodesData, rabbitmqIntegration)
	utils.PanicOnErr(err)

	utils.PanicOnErr(setInventoryData(localNode, nodeData))
}

func getLocalNodeName() (string, error) {
	if len(args.GlobalArgs.NodeNameOverride) > 0 {
		return args.GlobalArgs.NodeNameOverride, nil
	}
	output, err := execCommand("rabbitmqctl", "eval", "path().").Output()
	if err != nil {
		return "", err
	}
	output = bytes.TrimFunc(output, trimNodeName)
	if len(output) == 0 {
		return "", errors.New("could not determine the local node name")
	}
	return string(output), nil
}

func trimNodeName(r rune) bool {
	return unicode.IsSpace(r) || r == '\''
}

func getNodeEntity(nodeName string, nodesData []objx.Map, rabbitIntegration *integration.Integration) (entity *integration.Entity, nodeData objx.Map, err error) {
	if rabbitIntegration != nil {
		for _, node := range nodesData {
			if node.Get("name").Str() == nodeName {
				e, _, err := utils.CreateEntity(rabbitIntegration, nodeName, consts.NodeType, "")
				return e, node, err
			}
		}
	}
	return nil, nil, fmt.Errorf("node name [%v] not found in cluster", nodeName)
}

func setInventoryData(nodeEntity *integration.Entity, nodeData objx.Map) error {
	configPath := getConfigPath(nodeData)
	if len(configPath) > 0 {
		file, err := osOpen(configPath)
		if os.IsNotExist(err) {
			logger.Infof("The specified configuration file does not exist: %v", args.GlobalArgs.ConfigPath)
			return nil
		}
		if err != nil {
			logger.Errorf("Could not open the specified configuration file: %v", err)
			return err
		}
		defer utils.CheckErr(file.Close)

		return populateConfigInventory(file, nodeEntity)
	}
	return nil
}

func getConfigPath(nodeData objx.Map) string {
	if len(args.GlobalArgs.ConfigPath) > 0 {
		return args.GlobalArgs.ConfigPath
	}
	if nodeData != nil {
		configs := nodeData.Get("config_files").StrSlice()
		for _, config := range configs {
			if strings.HasSuffix(config, ".conf") {
				return config
			}
		}
	}
	return ""
}

func populateConfigInventory(reader io.Reader, nodeEntity *integration.Entity) error {
	scanner := bufio.NewScanner(reader)
	var line []byte
	var commentIndex int
	var eqIndex int
	var key string
	var value string

	for scanner.Scan() {
		line = scanner.Bytes()
		commentIndex = bytes.IndexByte(line, '#')
		if commentIndex >= 0 {
			line = line[:commentIndex]
		}
		if len(line) > 3 {
			if eqIndex = bytes.IndexByte(line, '='); eqIndex > 1 {
				key = string(bytes.TrimSpace(line[0:eqIndex]))
				value = string(bytes.TrimSpace(line[eqIndex+1:]))
				setInventoryItem(nodeEntity, "config", key, value)
			}
		}
	}
	return scanner.Err()
}

// PopulateEntityInventory adds inventory items from the entityData to the entity
func PopulateEntityInventory(entity *integration.Entity, entityType string, entityData *objx.Map) {
	for _, key := range entityInventory[entityType] {
		setInventoryMap(entity, entityType, key, entityData)
	}
}

func setInventoryMap(entity *integration.Entity, typeName, dataKey string, data *objx.Map, keyPrefix ...string) {
	val := data.Get(dataKey)
	setInventoryValue(entity, typeName, dataKey, val.Data(), keyPrefix...)
}

func setInventoryValue(entity *integration.Entity, typeName, key string, value interface{}, keyPrefix ...string) {
	switch value.(type) {
	case bool:
		setInventoryItem(entity, typeName, key, utils.ConvertBoolToInt(value.(bool)), keyPrefix...)
	case []interface{}:
		arrayVal := value.([]interface{})
		for i := range arrayVal {
			setInventoryValue(entity, typeName, strconv.Itoa(i+1), arrayVal[i], append(keyPrefix, key)...)
		}
	case map[string]interface{}:
		var ok bool
		var objxMap objx.Map
		if objxMap, ok = value.(objx.Map); !ok {
			objxMap = objx.New(value)
		}
		for mapKey := range objxMap {
			setInventoryMap(entity, typeName, mapKey, &objxMap, append(keyPrefix, key)...)
		}
	case nil:
		// skip
	default:
		setInventoryItem(entity, typeName, key, value, keyPrefix...)
	}
}

func setInventoryItem(entity *integration.Entity, typeName, key string, value interface{}, keyPrefix ...string) {
	actualKey := strings.Join(append(keyPrefix, key), ".")
	if err := entity.SetInventoryItem(typeName, actualKey, value); err != nil {
		logger.Infof("Error setting inventory [%s.%s] on [%s]: %v", typeName, actualKey, entity.Metadata.Name, err)
	}
}
