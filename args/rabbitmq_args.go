package args

import (
	"encoding/json"
	"errors"
	"regexp"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/nri-rabbitmq/utils/consts"
)

// GlobalArgs are the global set of arguments
var GlobalArgs RabbitMQArguments

// RabbitMQArguments is the fully parsed arguments, converting the JSON string into actual types
type RabbitMQArguments struct {
	sdkArgs.DefaultArgumentList
	Hostname         string
	Port             int
	Username         string
	Password         string
	CABundleFile     string
	CABundleDir      string
	NodeNameOverride string
	ConfigPath       string
	UseSSL           bool
	Queues           []string
	QueuesRegexes    []*regexp.Regexp
	Exchanges        []string
	ExchangesRegexes []*regexp.Regexp
	Vhosts           []string
	VhostsRegexes    []*regexp.Regexp
}

// Validate checks that valid collection arguments were specified
func (args *RabbitMQArguments) Validate() error {
	if args.Metrics && !args.Inventory {
		return errors.New("invalid arguments: can't collect metrics while not collecting inventory")
	} else if !args.All() && !args.Metrics && !args.Inventory {
		return errors.New("invalid arguments: nothing specified to collect")
	}
	return nil
}

// IncludeEntity returns true if the entity should be included; false otherwise
func (args *RabbitMQArguments) IncludeEntity(entityName string, entityType string, vhostName string) bool {
	if !args.IncludeVhost(vhostName) {
		return false
	}
	if entityType == consts.QueueType {
		return args.IncludeQueue(entityName)
	} else if entityType == consts.ExchangeType {
		return args.IncludeExchange(entityName)
	} else {
		return true
	}
}

// IncludeExchange returns true if exchage should be included; false otherwise
func (args *RabbitMQArguments) IncludeExchange(exchangeName string) bool {
	return includeName(exchangeName, args.Exchanges, args.ExchangesRegexes)
}

// IncludeQueue returns true if queue should be included; false otherwise
func (args *RabbitMQArguments) IncludeQueue(queueName string) bool {
	return includeName(queueName, args.Queues, args.QueuesRegexes)
}

// IncludeVhost returns true if vhost should be included; false otherwise
func (args *RabbitMQArguments) IncludeVhost(vhostName string) bool {
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
		CABundleDir:         args.CABundleDir,
		CABundleFile:        args.CABundleFile,
		ConfigPath:          args.ConfigPath,
		DefaultArgumentList: args.DefaultArgumentList,
		Hostname:            args.Hostname,
		NodeNameOverride:    args.NodeNameOverride,
		Password:            args.Password,
		Port:                args.Port,
		Username:            args.Username,
		UseSSL:              args.UseSSL,
	}
	var err error
	if err = parseStrings(args.Exchanges, &rabbitArgs.Exchanges); err != nil {
		return err
	}
	if err = parseStrings(args.Queues, &rabbitArgs.Queues); err != nil {
		return err
	}
	if err = parseStrings(args.Vhosts, &rabbitArgs.Queues); err != nil {
		return err
	}

	if rabbitArgs.ExchangesRegexes, err = parseRegexes(args.ExchangesRegexes); err != nil {
		return err
	}
	if rabbitArgs.QueuesRegexes, err = parseRegexes(args.QueuesRegexes); err != nil {
		return err
	}
	if rabbitArgs.VhostsRegexes, err = parseRegexes(args.VhostsRegexes); err != nil {
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
