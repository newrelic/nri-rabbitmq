// +build integration

package tests

import (
	"flag"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

const (
	containerName     = "nri-rabbitmq"
	schema            = "rabbitmq-schema.json"
	connURL           = "amqp://guest:guest@localhost:5672/"
	containerRabbitMQ = "rabbitmq"
)

func TestMain(m *testing.M) {
	flag.Parse()
	result := m.Run()
	os.Exit(result)
}

func TestSuccessConnection(t *testing.T) {
	if !waitForRabbitMQIsUpAndRunning(20) {
		t.Fatal("tests cannot be executed")
	}
	hostname := "rabbitmq"
	envVars := []string{
		fmt.Sprintf("HOSTNAME=%s", hostname),
		"USERNAME=guest",
		"PASSWORD=guest",
	}
	response, _, err := dockerComposeRun(envVars, containerName)
	assert.Nil(t, err)
	assert.NotEmpty(t, response)
	validateJSONSchema(response, schema)
}

func waitForRabbitMQIsUpAndRunning(maxTries int) bool {
	envVars := []string{}
	ports := []string{"5672:5672", "15672:15672"}
	stdout, stderr, err := dockerComposeRunMode(envVars, ports, containerRabbitMQ, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stdout)
	fmt.Println(stderr)
	for ; maxTries > 0; maxTries-- {
		log.Info("try to establish de connection with the rabbitmq...")
		conn, err := amqp.Dial(connURL)
		if err != nil {
			log.Warn(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		if conn != nil {
			conn.Close()
			log.Info("rabbitmq is up & running!")
			return true
		}
	}
	return false
}
