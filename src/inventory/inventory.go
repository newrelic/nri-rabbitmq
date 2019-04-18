package inventory

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-rabbitmq/src/args"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/data/consts"
)

var (
	execCommand = exec.Command
	osOpen      = os.Open
)

type inventoryKey struct {
	category string
	key      string
}

// CollectInventory collects the inventory items (config file values) from the apiResponses
func CollectInventory(rabbitmqIntegration *integration.Integration, nodesData []*data.NodeData, clusterName string) {
	if len(nodesData) == 0 {
		log.Warn("No node data available to collect inventory")
		return
	}

	nodeName, err := getLocalNodeName()
	if err != nil {
		log.Error("Error getting local node name: %v", err)
		return
	}

	nodeData, err := findNodeData(nodeName, nodesData)
	if err != nil {
		log.Error("Error finding node: %v", err)
		return
	}

	if config := getConfigData(nodeData); len(config) > 0 {
		localNode, _, err := data.CreateEntity(rabbitmqIntegration, nodeName, strings.TrimPrefix(consts.NodeType, "ra-"), "", clusterName)
		if err != nil {
			log.Error("Error creating local node entity: %v", err)
		}

		for k, v := range config {
			data.SetInventoryItem(localNode, k.category, k.key, v)
		}
	}
}

func getLocalNodeName() (string, error) {
	if len(args.GlobalArgs.NodeNameOverride) > 0 {
		return args.GlobalArgs.NodeNameOverride, nil
	}
	output, err := execCommand("rabbitmqctl", "eval", "node().").Output()
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

func findNodeData(nodeName string, nodesData []*data.NodeData) (nodeData *data.NodeData, err error) {
	for _, node := range nodesData {
		if node.Name == nodeName {
			return node, err
		}
	}
	return nil, fmt.Errorf("node name [%v] not found in RabbitMQ", nodeName)
}

func getConfigData(nodeData *data.NodeData) map[inventoryKey]string {
	configPath := getConfigPath(nodeData)
	if len(configPath) > 0 {
		file, err := osOpen(configPath)
		if os.IsNotExist(err) {
			log.Error("The specified configuration file does not exist: %v", args.GlobalArgs.ConfigPath)
			return nil
		}
		if err != nil {
			log.Error("Could not open the specified configuration file: %v", err)
			return nil
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Error("Error closing config file [%s]: %v", configPath, err)
			}
		}()

		if config, err := parseConfigInventory(file); err != nil {
			log.Error("Error parsing config file: %v", err)
		} else {
			return config
		}

	}
	return nil
}

func getConfigPath(nodeData *data.NodeData) string {
	if len(args.GlobalArgs.ConfigPath) > 0 {
		return args.GlobalArgs.ConfigPath
	}
	if nodeData != nil {
		for _, config := range nodeData.ConfigFiles {
			if strings.HasSuffix(config, ".conf") {
				return config
			}
		}
	}
	return ""
}

func parseConfigInventory(reader io.Reader) (map[inventoryKey]string, error) {
	values := make(map[inventoryKey]string)
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
		if len(line) >= 2 {
			if eqIndex = bytes.IndexByte(line, '='); eqIndex >= 1 {
				key = string(bytes.TrimSpace(line[0:eqIndex]))
				value = string(bytes.TrimSpace(line[eqIndex+1:]))
				values[inventoryKey{"config", key}] = value
			}
		}
	}
	return values, scanner.Err()
}
