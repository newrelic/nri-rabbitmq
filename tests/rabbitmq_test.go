//go:build integration

package tests

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/v3/log"
	"github.com/stretchr/testify/assert"
)

const (
	containerName     = "nri-rabbitmq"
	schema            = "rabbitmq-schema.json"
	connURL           = "amqp://guest:guest@localhost:5672/"
	containerRabbitMQ = "rabbitmq-1"
)

func TestMain(m *testing.M) {
	flag.Parse()
	result := m.Run()
	os.Exit(result)
}

func TestSuccessConnection(t *testing.T) {
	envVars := []string{}
	ports := []string{"5672:5672", "15672:15672"}
	stdout, stderr, err := dockerComposeRunMode(envVars, ports, containerRabbitMQ, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
	if !waitForRabbitMQIsUpAndRunning(20, containerRabbitMQ) {
		t.Fatal("tests cannot be executed, rabbitmq is not running.")
	}
	envVars = []string{
		fmt.Sprintf("HOSTNAME=%s", containerRabbitMQ),
		"USERNAME=guest",
		"PASSWORD=guest",
	}
	response, stderr, err := dockerComposeRun(envVars, containerName)
	fmt.Println(stderr)
	assert.Nil(t, err)
	assert.NotEmpty(t, response)
	err = validateJSONSchema(schema, response)
	assert.NoError(t, err, "The output of rabbitmq integration doesn't have expected format.")
}
