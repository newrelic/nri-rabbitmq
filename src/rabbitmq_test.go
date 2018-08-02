package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	origArgs := os.Args
	os.Args = []string{
		"nr-rabbitmq",
		// "-node_name_override", expectedNodeName,
		// "-config_path", testConfigPath,
	}
	defer func() {
		os.Args = origArgs
	}()
	assert.NotPanics(t, func() {
		// TODO: don't test till we fully merge and get other http tests
		// main()
	})
}
