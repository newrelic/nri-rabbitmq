package args

import (
	"encoding/json"
	"regexp"

	"github.com/newrelic/nri-rabbitmq/src/data/consts"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// GlobalArgs are the global set of arguments
var GlobalArgs RabbitMQArguments

// RabbitMQArguments is the fully parsed arguments, converting the JSON string into actual types
type RabbitMQArguments struct {
	sdkArgs.DefaultArgumentList
	Hostname             string
	Port                 int
	Username             string
	Password             string
	ManagementPathPrefix string
	CABundleFile         string
	CABundleDir          string
	NodeNameOverride     string
	ConfigPath           string
	UseSSL               bool
	Timeout              int
	DisableEntities      bool
	QueuesMaxLimit       int
	Queues               []string
	QueuesRegexes        []*regexp.Regexp
	Exchanges            []string
	ExchangesRegexes     []*regexp.Regexp
	Vhosts               []string
	VhostsRegexes        []*regexp.Regexp
}

// IncludeEntity returns true if the entity should be included; false otherwise
func (args *RabbitMQArguments) IncludeEntity(entityName string, entityType string, vhostName string) bool {
	if entityType == consts.NodeType {
		return true
	}

	if !args.includeVhost(vhostName) {
		return false
	}

	if entityType == consts.QueueType {
		return args.includeQueue(entityName)
	} else if entityType == consts.ExchangeType {
		return args.includeExchange(entityName)
	} else {
		return true
	}
}

// includeExchange returns true if exchange should be included; false otherwise
func (args *RabbitMQArguments) includeExchange(exchangeName string) bool {
	return includeName(exchangeName, args.Exchanges, args.ExchangesRegexes)
}

// includeQueue returns true if queue should be included; false otherwise
func (args *RabbitMQArguments) includeQueue(queueName string) bool {
	return includeName(queueName, args.Queues, args.QueuesRegexes)
}

// includeVhost returns true if vhost should be included; false otherwise
func (args *RabbitMQArguments) includeVhost(vhostName string) bool {
	return includeName(vhostName, args.Vhosts, args.VhostsRegexes)
}

func includeName(itemName string, names []string, namesRegex []*regexp.Regexp) bool {
	for _, name := range names {
		if name == itemName {
			return true
		}
	}
	for _, reg := range namesRegex {
		if reg.MatchString(itemName) {
			return true
		}
	}
	if len(names) > 0 || len(namesRegex) > 0 {
		return false
	}
	return true
}

// SetGlobalArgs validates the arguments in ArgumentList and sets GlobalArgs to the result
func SetGlobalArgs(args ArgumentList) error {
	rabbitArgs := RabbitMQArguments{
		ManagementPathPrefix: args.ManagementPathPrefix,
		CABundleDir:          args.CABundleDir,
		CABundleFile:         args.CABundleFile,
		ConfigPath:           args.ConfigPath,
		DefaultArgumentList:  args.DefaultArgumentList,
		Hostname:             args.Hostname,
		NodeNameOverride:     args.NodeNameOverride,
		Password:             args.Password,
		Port:                 args.Port,
		Username:             args.Username,
		UseSSL:               args.UseSSL,
		Timeout:              args.Timeout,
		DisableEntities:      args.DisableEntities,
		QueuesMaxLimit:       args.QueuesMaxLimit,
	}
	var err error
	if err = parseStrings(args.Exchanges, &rabbitArgs.Exchanges); err != nil {
		log.Error("Error parsing arguments [Exchanges]: %v", err)
		return err
	}
	if err = parseStrings(args.Queues, &rabbitArgs.Queues); err != nil {
		log.Error("Error parsing arguments [Queues]: %v", err)
		return err
	}
	if err = parseStrings(args.Vhosts, &rabbitArgs.Vhosts); err != nil {
		log.Error("Error parsing arguments [Vhosts]: %v", err)
		return err
	}

	if rabbitArgs.ExchangesRegexes, err = parseRegexes(args.ExchangesRegexes); err != nil {
		log.Error("Error parsing arguments [ExchangesRegexes]: %v", err)
		return err
	}
	if rabbitArgs.QueuesRegexes, err = parseRegexes(args.QueuesRegexes); err != nil {
		log.Error("Error parsing arguments [QueuesRegexes]: %v", err)
		return err
	}
	if rabbitArgs.VhostsRegexes, err = parseRegexes(args.VhostsRegexes); err != nil {
		log.Error("Error parsing arguments [VhostsRegexes]: %v", err)
		return err
	}
	GlobalArgs = rabbitArgs
	return nil
}

func parseStrings(argValue string, value *[]string) error {
	if argValue != "" {
		return json.Unmarshal([]byte(argValue), value)
	}
	return nil
}

func parseRegexes(argValue string) (regexes []*regexp.Regexp, err error) {
	if argValue == "" {
		return nil, nil
	}
	var values []string
	if err = json.Unmarshal([]byte(argValue), &values); err != nil {
		return
	}
	for _, item := range values {
		regex, err := regexp.Compile(item)
		if err != nil {
			return nil, err
		}
		regexes = append(regexes, regex)
	}
	return
}
