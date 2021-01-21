package args

import sdkArgs "github.com/newrelic/infra-integrations-sdk/args"

// ArgumentList is the raw arguments passed into the integration via YAML, CLI args, or ENV variables
type ArgumentList struct {
	sdkArgs.DefaultArgumentList
	Hostname             string `default:"localhost" help:"Hostname or IP where RabbitMQ Management Plugin is running."`
	Port                 int    `default:"15672" help:"Port on which RabbitMQ Management Plugin is listening."`
	Username             string `default:"" help:"Username for accessing RabbitMQ Management Plugin"`
	Password             string `default:"" help:"Password for the given user."`
	ManagementPathPrefix string `default:"" help:"RabbitMQ Management Prefix."`
	CABundleFile         string `default:"" help:"Alternative Certificate Authority bundle file"`
	CABundleDir          string `default:"" help:"Alternative Certificate Authority bundle directory"`
	NodeNameOverride     string `default:"" help:"Overrides the local node name instead of retrieving it from RabbitMQ."`
	ConfigPath           string `default:"" help:"RabbitMQ configuration file path."`
	UseSSL               bool   `default:"false" help:"configure whether to use an SSL connection or not."`
	Queues               string `default:"" help:"JSON array of queue names from which to collect metrics."`
	QueuesRegexes        string `default:"" help:"JSON array of queue name regexes from which to collect metrics."`
	Exchanges            string `default:"" help:"JSON array of exchange names from which to collect metrics."`
	ExchangesRegexes     string `default:"" help:"JSON array of exchange name regexes from which to collect metrics."`
	Vhosts               string `default:"" help:"JSON array of vhost names from which to collect metrics."`
	VhostsRegexes        string `default:"" help:"JSON array of vhost name regexes from which to collect metrics."`
	ShowVersion          bool   `default:"false" help:"Print build information and exit"`
}
