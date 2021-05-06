// +build integration

package tests

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/streadway/amqp"

	"github.com/newrelic/infra-integrations-sdk/log"

	"github.com/xeipuuv/gojsonschema"
)

const (
	connURL           = "amqp://guest:guest@localhost:5672/"
	containerRabbitMQ = "rabbitmq"
)

func executeDockerCompose(containerName string, envVars []string) (string, string, error) {
	cmdLine := []string{"run"}
	for i := range envVars {
		cmdLine = append(cmdLine, "-e")
		cmdLine = append(cmdLine, envVars[i])
	}
	cmdLine = append(cmdLine, containerName)
	fmt.Printf("executing: docker-compose %s\n", strings.Join(cmdLine, " "))
	cmd := exec.Command("docker-compose", cmdLine...)
	if err := checkRabbitMQIsUpAndRunning(10); err != nil {
		return "", "", err
	}
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	return outbuf.String(), errbuf.String(), nil
}

func validateJSONSchema(fileName string, input string) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	schemaURI := fmt.Sprintf("file://%s", filepath.Join(pwd, "testdata", fileName))
	log.Info("loading schema from %s", schemaURI)
	schemaLoader := gojsonschema.NewReferenceLoader(schemaURI)
	documentLoader := gojsonschema.NewStringLoader(input)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("Error loading JSON schema, error: %v", err)
	}
	if result.Valid() {
		return nil
	}
	fmt.Printf("Errors for JSON schema: '%s'\n", schemaURI)
	for _, desc := range result.Errors() {
		fmt.Printf("\t- %s\n", desc)
	}
	fmt.Printf("\n")
	return fmt.Errorf("The output of the integration doesn't have expected JSON format")
}

func checkRabbitMQIsUpAndRunning(maxTries int) error {
	cmdLine := []string{"up", "-d"}
	cmdLine = append(cmdLine, containerRabbitMQ)
	fmt.Printf("executing: docker-compose %s\n", strings.Join(cmdLine, " "))
	cmd := exec.Command("docker-compose", cmdLine...)
	cmd.Run()
	for ; maxTries > 0; maxTries-- {
		log.Info("check rabbitmq connection")
		conn, err := amqp.Dial(connURL)
		if err != nil {
			log.Warn(err.Error())
			if maxTries == 0 {
				return errors.New("rabbitmq connection cannot be established!. Tests won't be executed")
			}

			time.Sleep(3 * time.Second)
		} else {
			conn.Close()
			break
		}
	}
	log.Info("rabbit is up & running")
	return nil
}
