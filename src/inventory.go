package main

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
	"github.com/stretchr/objx"
)

var (
	execCommand     = exec.Command
	osOpen          = os.Open
	entityInventory = map[string][]string{
		exchangeType: {
			"type",
			"durable",
			"auto_delete",
			"arguments",
		},
		queueType: {
			"exclusive",
			"durable",
			"auto_delete",
			"arguments",
		},
	}
)

func getNodeName() (string, error) {
	if len(args.NodeNameOverride) > 0 {
		return args.NodeNameOverride, nil
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

func getNodeEntity(nodeName string, nodeData [](objx.Map), integration *integration.Integration) (entity *integration.Entity, err error) {
	if integration != nil {
		for _, node := range nodeData {
			if node.Get("name").Str() == nodeName {
				e, _, err := createEntity(integration, nodeName, nodeType, "")
				return e, err
			}
		}
	}
	return nil, fmt.Errorf("node name [%v] not found in cluster", nodeName)
}

func setInventoryData(nodeEntity *integration.Entity, overviewData objx.Map) error {
	configPath := getConfigPath(overviewData)
	if len(configPath) > 0 {
		file, err := osOpen(configPath)
		if os.IsNotExist(err) {
			logger.Infof("The specified configuration file does not exist: %v", args.ConfigPath)
			return nil
		}
		if err != nil {
			logger.Errorf("Could not open the specified configuration file: %v", err)
			return err
		}
		defer checkErr(file.Close)

		logger.Debugf("Reading configuration data from %q", args.ConfigPath)
		return populateConfigInventory(file, nodeEntity)
	}
	return nil
}

func getConfigPath(overviewData objx.Map) string {
	if len(args.ConfigPath) > 0 {
		return args.ConfigPath
	}
	if overviewData != nil {
		configs := overviewData.Get("config_files").StrSlice()
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

func populateEntityInventory(entity *integration.Entity, entityType string, entityData *objx.Map) {
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
		setInventoryItem(entity, typeName, key, convertBoolToInt(value.(bool)), keyPrefix...)
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
