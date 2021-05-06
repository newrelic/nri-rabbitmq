// +build integration

package tests

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const (
	containerName = "nri-rabbitmq"
	schema        = "rabbitmq-jsonschema.json"
)

func TestMain(m *testing.M) {
	flag.Parse()
	result := m.Run()
	os.Exit(result)
}

func TestSuccessConnection(t *testing.T) {
	hostname := "rabbitmq"
	envVars := []string{
		fmt.Sprintf("HOSTNAME=%s", hostname),
		"USERNAME=guest",
		"PASSWORD=guest",
	}
	response, _, err := executeDockerCompose(containerName, envVars)
	assert.Nil(t, err)
	assert.NotEmpty(t, response)

	assert.Equal(t, "com.newrelic.rabbitmq", gjson.Get(response, "name").String())
	assert.Equal(t, "3", gjson.Get(response, "protocol_version").String())
	assert.Equal(t, "/", gjson.Get(response, "data.0.entity.name").String())
	evt := gjson.Get(
		response, "data.0.metrics.#(event_type==\"RabbitmqVhostSample\")",
	).Map()
	assert.Equal(t, "vhost:/", evt["entityName"].String())
	assert.Equal(t, "/", evt["displayName"].String())
	fmt.Println(response)
}
