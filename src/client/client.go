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
	"github.com/newrelic/nri-rabbitmq/src/args"
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
	request := createRequest(endpoint)
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
		client, err := nrHttp.New(args.GlobalArgs.CABundleFile, args.GlobalArgs.CABundleDir, time.Second*30)
		if err != nil {
			log.Error("Unable to create HTTP Client: %v", err)
			os.Exit(2)
		}
		defaultClient = client
	}
}

func createRequest(endpoint string) *http.Request {
	var fullURL string
	if args.GlobalArgs.UseSSL {
		fullURL = fmt.Sprintf("https://%s:%d%s", args.GlobalArgs.Hostname, args.GlobalArgs.Port, endpoint)
	} else {
		fullURL = fmt.Sprintf("http://%s:%d%s", args.GlobalArgs.Hostname, args.GlobalArgs.Port, endpoint)
	}
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth(args.GlobalArgs.Username, args.GlobalArgs.Password)

	return req
}
