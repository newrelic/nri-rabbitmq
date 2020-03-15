package client

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/args"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCollectEndpoint(t *testing.T) {
	defaultClient = nil
	args.GlobalArgs = args.RabbitMQArguments{}
	mux, teardown := testutils.GetTestServer(false)
	defer teardown()
	mux.HandleFunc(ConnectionsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		fmt.Fprint(w, `[
			{
				"state": "running",
				"vhost": "test-vhost"
			}
		]`)
	})
	mux.HandleFunc(QueuesEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		fmt.Fprint(w, `[
			{
				"messages_details": {
					"rate": 0.1
				},
				"messages": 10,
				"name": "test-queue",
				"vhost": "test-vhost"
			}
		]`)
	})

	err := CollectEndpoint("", nil)
	assert.Error(t, err)

	err = CollectEndpoint(ConnectionsEndpoint, nil)
	assert.Error(t, err)
	var actualConnections []data.ConnectionData

	err = CollectEndpoint(ConnectionsEndpoint, &actualConnections)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actualConnections), "Connection length should be 1")
	assert.Equal(t, "test-vhost", actualConnections[0].Vhost, "Vhost is different")
	assert.Equal(t, "running", actualConnections[0].State, "State is different")
	var actualQueues []data.QueueData

	err = CollectEndpoint(QueuesEndpoint, &actualQueues)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actualQueues))
	f := float64(0.1)
	assert.Equal(t, &f, actualQueues[0].MessagesDetails.Rate, "MessageDetails.Rate is different")
	i := int64(10)
	assert.Equal(t, &i, actualQueues[0].Messages, "Messages is different")
	assert.Equal(t, "test-queue", actualQueues[0].Name, "Name is different")
	assert.Equal(t, "test-vhost", actualQueues[0].Vhost, "Vhost is different")
}

func TestCollectEndpoint_Errors(t *testing.T) {
	defaultClient = nil
	args.GlobalArgs = args.RabbitMQArguments{}
	mux, teardown := testutils.GetTestServer(false)
	defer teardown()
	mux.HandleFunc(ConnectionsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	err := CollectEndpoint(ConnectionsEndpoint, &struct{}{})
	assert.Error(t, err)

	defaultClient = nil
	args.GlobalArgs.Hostname = "[" + args.GlobalArgs.Hostname

	err = CollectEndpoint("/missing", &struct{}{})
	assert.Error(t, err)
}

func Test_ensureClient_CannotCreateClient(t *testing.T) {
	defaultClient = nil
	args.GlobalArgs = args.RabbitMQArguments{}

	if os.Getenv("INVALID_CLIENT") == "1" {
		// If this test was called to execute the invalid client test, then call collectEndpoint
		// This is only called in this fashion below
		args.GlobalArgs.CABundleFile = filepath.Join("not-found")
		defer func() {
			args.GlobalArgs.CABundleDir = ""
		}()

		ensureClient()
		return
	}

	// If this is the first time this test method is ran, re-execute it telling it to perform the actual method call.
	// Then test the result of that, which should be an os.Exit(2).
	// The downside to this is the ensureClient() will not show full coverage when it actually does (since it's ran as a sub-test)
	cmd := exec.Command(os.Args[0], "-test.run=Test_ensureClient_CannotCreateClient")
	cmd.Env = append(os.Environ(), "INVALID_CLIENT=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok {
		assert.Equal(t, "exit status 2", e.Error(), "Exit status of invalid HTTP client should be 2")
		return
	}
	if err != nil {
		t.Fatalf("Unexpected error type from invalid HTTP client %v", err)
	} else {
		t.Fatalf("Expected ensureClient() to os.Exit() and it did not")
	}
}

func Test_collectEndpoint_Errors(t *testing.T) {
	defaultClient = nil
	args.GlobalArgs = args.RabbitMQArguments{}

	mux, close := testutils.GetTestServer(false)
	defer close()
	mux.HandleFunc(ConnectionsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		fmt.Fprint(w, `[
			{}
		]`)
	})
	mux.HandleFunc(QueuesEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	err := collectEndpoint(nil, &struct{}{})
	assert.Error(t, err)

	req, err := createRequest(ConnectionsEndpoint)
	assert.NoError(t, err)
	err = collectEndpoint(req, &struct{}{})
	assert.Error(t, err)

	req, err = createRequest(QueuesEndpoint)
	assert.NoError(t, err)
	err = collectEndpoint(req, &struct{}{})
	assert.Error(t, err)

	req.URL = nil
	err = collectEndpoint(req, &struct{}{})
	assert.Error(t, err)
}

func Test_createRequest(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	args.GlobalArgs.UseSSL = true
	args.GlobalArgs.Hostname = "test-hostname"
	args.GlobalArgs.Port = 3000
	args.GlobalArgs.ManagementPathPrefix = "/test-management-prefix"
	endpoint := "/test-endpoint"
	r, err := createRequest(endpoint)
	assert.NoError(t, err)

	actualURL := r.URL.String()
	expectedURL := fmt.Sprintf("https://%v:%v%v%v", args.GlobalArgs.Hostname, args.GlobalArgs.Port, args.GlobalArgs.ManagementPathPrefix, endpoint)
	assert.Equal(t, expectedURL, actualURL, "expect url to use https")
	if r.Method != http.MethodGet {
		t.Error("Expected GET method, got POST method.")
	}
}

func Test_createRequest_Error(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	args.GlobalArgs.Hostname = "[test-hostname"
	endpoint := "/test-endpoint"
	_, err := createRequest(endpoint)
	assert.Error(t, err)
}
