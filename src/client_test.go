package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectEndpointsConnectionsNoSSL(t *testing.T) {
	mux, teardown := setupTestServer(false)
	defer teardown()
	mux.HandleFunc(connectionsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
				"destination": "test-destination",
				"source": "test-source",
				"vhost": "test-vhost"
			}
		]`)
	})
	mux.HandleFunc(queuesEndpoint, func(w http.ResponseWriter, r *http.Request) {
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

	actualResultMap, _ := collectEndpoints(connectionsEndpoint, queuesEndpoint)

	expectedResultMap := map[string]interface{}{
		connectionsEndpoint: []interface{}{
			map[string]interface{}{
				"destination": "test-destination",
				"source":      "test-source",
				"vhost":       "test-vhost",
			},
		},
		queuesEndpoint: []interface{}{
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
		args.UseSSL = true
		server = httptest.NewTLSServer(mux)
	} else {
		args.UseSSL = false
		server = httptest.NewServer(mux)
	}
	url, _ := url.Parse(server.URL)

	port, _ := strconv.Atoi(url.Port())
	args.Hostname, args.Port = url.Hostname(), port
	return mux, server.Close
}

func TestCollectEndpointErrors(t *testing.T) {
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

	origCAFile := args.CABundleFile
	args.CABundleFile = "not-found"

	_, err = collectEndpoints("/bad")
	assert.Error(t, err, "An error should be returned when bad CA files are supplied")
	err = nil
	args.CABundleFile = origCAFile

	mux, closer := setupTestServer(false)
	defer func() {
		closer()
	}()

	mux.HandleFunc("/bad", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "{[]}")
	})

	_, err = collectEndpoints("/bad")
	assert.Error(t, err, "An error should be returned when a bad JSON response is returned")
	err = nil
}

func TestCollectEndpointsEmptyArray(t *testing.T) {
	actualResultMap, _ := collectEndpoints()
	assert.Empty(t, actualResultMap)
}

func TestCreateRequestURL(t *testing.T) {
	args.UseSSL = true
	args.Hostname = "test-hostname"
	args.Port = 3000
	endpoint := "test-endpoint"
	r, err := createRequest(endpoint)
	actualURL := fmt.Sprintf("https://%v", r.Host)
	expectedURL := fmt.Sprintf("https://%v:%v%v", args.Hostname, args.Port, endpoint)
	assert.Equal(t, expectedURL, actualURL, "expect url to use https")
	assert.NoError(t, err)
	if r.Method != http.MethodGet {
		t.Error("Expected GET method, got POST method.")
	}
}
