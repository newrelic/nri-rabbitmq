package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/stretchr/testify/assert"
)

func TestGetListOfEndpointsToCollectAllArgs(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	actual := GetEndpointsToCollect()
	assert.Equal(t, AllEndpoints, actual, "should publish all endpoints")

	args.GlobalArgs.Metrics = true
	args.GlobalArgs.Inventory = true
	actual = GetEndpointsToCollect()
	assert.Equal(t, AllEndpoints, actual, "should publish all endpoints")
}

func TestGetListOfEndpointsToCollectNoEndpoints(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	args.GlobalArgs.Metrics = true
	actual := GetEndpointsToCollect()
	assert.Empty(t, actual, "should be empty slice")
}

func TestGetListOfEndpointsToCollectInventory(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	args.GlobalArgs.Inventory = true
	actual := GetEndpointsToCollect()
	assert.Equal(t, InventoryEndpoints, actual, "should publish inventory endpoints")
}

func TestCollectEndpointsConnectionsNoSSL(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	mux, teardown := setupTestServer(false)
	defer teardown()
	mux.HandleFunc(ConnectionsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
				"destination": "test-destination",
				"source": "test-source",
				"vhost": "test-vhost"
			}
		]`)
	})
	mux.HandleFunc(QueuesEndpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
				"message_details": {
					"rate": 0.1
				},
				"messages": 10,
				"node": "rabbit@rabbitmq-1"
			}
		]`)
	})

	actualResultMap, _ := CollectEndpoints(ConnectionsEndpoint, QueuesEndpoint)

	expectedResultMap := map[string]interface{}{
		ConnectionsEndpoint: []interface{}{
			map[string]interface{}{
				"destination": "test-destination",
				"source":      "test-source",
				"vhost":       "test-vhost",
			},
		},
		QueuesEndpoint: []interface{}{
			map[string]interface{}{
				"message_details": map[string]interface{}{
					"rate": 0.1,
				},
				"messages": 10.0,
				"node":     "rabbit@rabbitmq-1",
			},
		},
	}

	assert.Equal(t, expectedResultMap, actualResultMap)
}

func setupTestServer(tls bool) (mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	var server *httptest.Server
	if tls {
		args.GlobalArgs.UseSSL = true
		server = httptest.NewTLSServer(mux)
	} else {
		args.GlobalArgs.UseSSL = false
		server = httptest.NewServer(mux)
	}
	url, _ := url.Parse(server.URL)

	port, _ := strconv.Atoi(url.Port())
	args.GlobalArgs.Hostname, args.GlobalArgs.Port = url.Hostname(), port
	return mux, server.Close
}

func TestCollectEndpointErrors(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	_, err := collectEndpoint(nil, nil)
	assert.Error(t, err, "An error should be returned when passing a nil http.Client")
	err = nil

	_, err = collectEndpoint(http.DefaultClient, nil)
	assert.Error(t, err, "An error should be returned when passing a nil http.Request")
	err = nil

	req := new(http.Request)
	_, err = collectEndpoint(http.DefaultClient, req)
	assert.Error(t, err, "An error should be returned when a bad request is provided")
	err = nil

	origCAFile := args.GlobalArgs.CABundleFile
	args.GlobalArgs.CABundleFile = "not-found"

	_, err = CollectEndpoints("/bad")
	assert.Error(t, err, "An error should be returned when bad CA files are supplied")
	err = nil
	args.GlobalArgs.CABundleFile = origCAFile

	mux, closer := setupTestServer(false)
	defer func() {
		closer()
	}()

	mux.HandleFunc("/bad", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "{[]}")
	})

	_, err = CollectEndpoints("/bad")
	assert.Error(t, err, "An error should be returned when a bad JSON response is returned")
	err = nil
}

func TestCollectEndpointsEmptyArray(t *testing.T) {
	actualResultMap, _ := CollectEndpoints()
	assert.Empty(t, actualResultMap)
}

func TestCreateRequestURL(t *testing.T) {
	args.GlobalArgs = args.RabbitMQArguments{}
	args.GlobalArgs.UseSSL = true
	args.GlobalArgs.Hostname = "test-hostname"
	args.GlobalArgs.Port = 3000
	endpoint := "test-endpoint"
	r, err := createRequest(endpoint)
	actualURL := fmt.Sprintf("https://%v", r.Host)
	expectedURL := fmt.Sprintf("https://%v:%v%v", args.GlobalArgs.Hostname, args.GlobalArgs.Port, endpoint)
	assert.Equal(t, expectedURL, actualURL, "expect url to use https")
	assert.NoError(t, err)
	if r.Method != http.MethodGet {
		t.Error("Expected GET method, got POST method.")
	}
}
