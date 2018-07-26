package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	update = flag.Bool("update", false, "update .golden files")
)

type mockedLogger struct {
	mock.Mock
}
type testLogger struct {
	f func(format string, args ...interface{})
}

func (l *mockedLogger) Debugf(format string, args ...interface{}) {
	args = append([]interface{}{format}, args...)
	l.Called(args...)
}
func (l *mockedLogger) Infof(format string, args ...interface{}) {
	args = append([]interface{}{format}, args...)
	l.Called(args...)
}
func (l *mockedLogger) Errorf(format string, args ...interface{}) {
	args = append([]interface{}{format}, args...)
	l.Called(args...)
}
func (l *mockedLogger) Warnf(format string, args ...interface{}) {
	args = append([]interface{}{format}, args...)
	l.Called(args...)
}

func (l *testLogger) Debugf(format string, args ...interface{}) {
	l.f("DEBUG: "+format, args...)
}
func (l *testLogger) Infof(format string, args ...interface{}) {
	l.f("INFO : "+format, args...)
}
func (l *testLogger) Errorf(format string, args ...interface{}) {
	l.f("ERROR: "+format, args...)
}
func (l *testLogger) Warnf(format string, args ...interface{}) {
	l.f("WARN: "+format, args...)
}

func getTestingIntegration(t *testing.T) (payload *integration.Integration) {
	payload, err := integration.New("Test", "0.0.1", integration.Logger(&testLogger{t.Logf}))
	require.NoError(t, err)
	require.NotNil(t, payload)
	logger = payload.Logger()
	return
}

func getTestingEntity(t *testing.T, entityArgs ...string) (payload *integration.Integration, entity *integration.Entity) {
	payload = getTestingIntegration(t)
	var err error
	if len(entityArgs) > 1 {
		entity, err = payload.Entity(entityArgs[0], entityArgs[1])
		assert.NoError(t, err)
	} else {
		entity = payload.LocalEntity()
	}
	require.NotNil(t, entity)
	return
}

func readObjectFromJSONFile(t *testing.T, filename string) map[string]interface{} {
	data, err := ioutil.ReadFile(filename)
	assert.NoError(t, err)
	item := map[string]interface{}{}
	err = json.Unmarshal(data, &item)
	assert.NoError(t, err)
	return item
}

func fakeExecCommand(command string, args ...string) (cmd *exec.Cmd) {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd = exec.Command(os.Args[0], cs...)
	// cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess isn't a real test. It's used as a helper process.
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "rabbitmqctl":
		if len(args) == 2 && args[0] == "eval" && args[1] == "path()." {
			if os.Getenv("GET_NODE_NAME_ERROR") == "1" {
				os.Exit(2)
			}
			if os.Getenv("GET_NODE_NAME_EMPTY") == "1" {
				fmt.Fprintf(os.Stdout, "")
			} else {
				fmt.Fprintf(os.Stdout, expectedNodeCmdOutput)
			}
		}
	}
}
