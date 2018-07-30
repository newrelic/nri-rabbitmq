# New Relic Infrastructure Integration for RabbitMQ

New Relic Infrastructure Integration for RabbitMQ captures critical performance metrics and inventory reported by the RabbitMQ Management Plugin. Inventory data is obtained from the configuration file, and metrics and additional inventory data is obtained from the management API.

## Requirements

The RabbitMQ integration requires that the [RabbitMQ Management Plugin](https://www.rabbitmq.com/management.html#getting-started) is enabled on the RabbitMQ host being monitored.

## Installation

* download an archive file for the RabbitMQ Integration
* extract `rabbitmq-definition.yml` and `/bin` directory into `/var/db/newrelic-infra/newrelic-integrations`
* add execute permissions for the binary file `nr-rabbitmq` (if necessary)
* extract `rabbitmq-config.yml.sample` into `/etc/newrelic-infra/integrations.d`

## Usage

This is the description about how to run the RabbitMQ Integration with New Relic Infrastructure agent, so it is required to have the agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

In order to use the RabbitMQ Integration it is required to configure `rabbitmq-config.yml.sample` file. Firstly, rename the file to `rabbitmq-config.yml`. Then, depending on your needs, specify all instances that you want to monitor. Once this is done, restart the Infrastructure agent.

You can view your data in Insights by creating your own custom NRQL queries. To do so use the **Rabbitmq*Sample** event types.

## Integration Development usage

The integration can be built and run locally.

For managing external dependencies [govendor tool](https://github.com/kardianos/govendor) is used. It is required to lock all external dependencies to specific version (if possible) into vendor directory.

* Go to the directory of the RabbitMQ integration and build it:
```bash
$ make
```
* The above command will run tests for the integration and build an executable file called `nr-rabbitmq` in the `/bin` directory. This executable can be run by itself:
```bash
$ ./bin/nr-rabbitmq
```
* For additional usage information use the `-help` flag:
```bash
$ ./bin/nr-rabbitmq -help
```