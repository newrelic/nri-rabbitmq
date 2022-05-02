package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	nrHttp "github.com/newrelic/infra-integrations-sdk/http"
	"github.com/newrelic/infra-integrations-sdk/log"

	args2 "github.com/newrelic/nri-rabbitmq/src/args"
)

const (
	// OverviewEndpoint path
	OverviewEndpoint = "/api/overview"
	// NodesEndpoint path
	NodesEndpoint = "/api/nodes"
	// QueuesEndpoint path
	QueuesEndpoint = "/api/queues"
	// ExchangesEndpoint path
	ExchangesEndpoint = "/api/exchanges"
	// VhostsEndpoint path
	VhostsEndpoint = "/api/vhosts"
	// ConnectionsEndpoint path
	ConnectionsEndpoint = "/api/connections"
	// BindingsEndpoint path
	BindingsEndpoint = "/api/bindings"
	// AlivenessTestEndpoint path, this is formatted with the vhost name
	AlivenessTestEndpoint = "/api/aliveness-test/%s"
	// HealthCheckEndpoint path, this is formatted with the node name
	HealthCheckEndpoint = "/api/healthchecks/node/%s"
	// Default client timeout
	DefaultTimout = time.Second * 30
)

var defaultClient *http.Client

// CollectEndpoint calls the endpoint and populates its response into result
func CollectEndpoint(endpoint string, result interface{}) error {
	if endpoint == "" {
		err := errors.New("endpoint cannot be empty")
		log.Error("Error collecting endpoint: %v", err)
		return err
	}
	if result == nil {
		err := errors.New("the result destination for the endpoint cannot be nil")
		log.Error("Error collecting endpoint: %v", err)
		return err
	}
	request, err := createRequest(endpoint)
	if err != nil {
		log.Error("Error creating request to Management API: %v", err)
		return err
	}
	if err := collectEndpoint(request, result); err != nil {
		return err
	}
	return nil
}

func collectEndpoint(req *http.Request, jsonResult interface{}) error {
	ensureClient()
	if req == nil {
		return errors.New("an http request was not specified")
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 || !strings.HasPrefix(resp.Header.Get("content-type"), "application/json") {
		err := fmt.Errorf("unexpected http response from [%s]: %s", req.URL, resp.Status)
		log.Error("Error making API call: %v", err)
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("Error closing response body: %v", err)
		}
	}()

	if err = json.NewDecoder(resp.Body).Decode(jsonResult); err != nil {
		return err
	}

	return nil
}

func ensureClient() {
	if defaultClient == nil {
		clientOptions := []nrHttp.ClientOption{
			nrHttp.WithTimeout(DefaultTimout),
		}
		if args2.GlobalArgs.CABundleDir != "" {
			clientOptions = append(clientOptions, nrHttp.WithCABundleDir(args2.GlobalArgs.CABundleDir))
		}
		if args2.GlobalArgs.CABundleFile != "" {
			clientOptions = append(clientOptions, nrHttp.WithCABundleFile(args2.GlobalArgs.CABundleFile))
		}

		client, err := nrHttp.New(clientOptions...)
		if err != nil {
			log.Error("Unable to create HTTP Client: %v", err)
			os.Exit(2)
		}

		defaultClient = client
	}
}

func createRequest(endpoint string) (*http.Request, error) {
	var fullURL string
	if args2.GlobalArgs.UseSSL {
		fullURL = fmt.Sprintf("https://%s:%d%s%s", args2.GlobalArgs.Hostname, args2.GlobalArgs.Port, args2.GlobalArgs.ManagementPathPrefix, endpoint)
	} else {
		fullURL = fmt.Sprintf("http://%s:%d%s%s", args2.GlobalArgs.Hostname, args2.GlobalArgs.Port, args2.GlobalArgs.ManagementPathPrefix, endpoint)
	}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(args2.GlobalArgs.Username, args2.GlobalArgs.Password)
	return req, nil
}
