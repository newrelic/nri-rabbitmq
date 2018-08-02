package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	nrHttp "github.com/newrelic/infra-integrations-sdk/http"
	"github.com/newrelic/nri-rabbitmq/args"
	"github.com/newrelic/nri-rabbitmq/utils"
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

var (
	// AllEndpoints is the list of all RabbitMQ Management API endpoints that are needed to be collected for both Inventory and Metrics
	AllEndpoints = []string{
		OverviewEndpoint,
		NodesEndpoint,
		QueuesEndpoint,
		ExchangesEndpoint,
		VhostsEndpoint,
		ConnectionsEndpoint,
		BindingsEndpoint,
	}
	// InventoryEndpoints is the list of all RabbitMQ Management API endpoints that are needed to be collected for Inventory
	InventoryEndpoints = []string{
		OverviewEndpoint,
		NodesEndpoint,
		QueuesEndpoint,
		ExchangesEndpoint,
	}
)

// GetEndpointsToCollect returns the endpoints that are needed to collect based on what's being collected
func GetEndpointsToCollect() (endpoints []string) {
	if args.GlobalArgs.All() || (args.GlobalArgs.Metrics && args.GlobalArgs.Inventory) {
		endpoints = AllEndpoints
	} else if !args.GlobalArgs.Metrics && args.GlobalArgs.Inventory {
		endpoints = InventoryEndpoints
	}
	return
}

// CollectEndpoints calls each endpoint and returns it's result in the map by endpoint path
func CollectEndpoints(endpoints ...string) (resultMap map[string]interface{}, err error) {
	resultMap = make(map[string]interface{})

	if len(endpoints) == 0 {
		return
	}

	client, err := nrHttp.New(args.GlobalArgs.CABundleFile, args.GlobalArgs.CABundleDir, time.Second*30)
	if err != nil {
		return
	}

	for _, endpoint := range endpoints {
		var request *http.Request
		var response interface{}
		request, err = createRequest(endpoint)
		if err == nil {
			response, err = collectEndpoint(client, request)
			if err == nil {
				resultMap[endpoint] = response
			}
		}
		if err != nil {
			return nil, fmt.Errorf("could not collect required endpoint [%v]: %v", endpoint, err)
		}
	}
	return
}

func collectEndpoint(client *http.Client, req *http.Request) (interface{}, error) {
	if client == nil {
		return nil, errors.New("an http client was not specified")
	}
	if req == nil {
		return nil, errors.New("an http request was not specified")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer utils.CheckErr(resp.Body.Close)

	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func createRequest(endpoint string) (*http.Request, error) {
	var fullURL string
	if args.GlobalArgs.UseSSL {
		fullURL = fmt.Sprintf("https://%v:%v%v", args.GlobalArgs.Hostname, args.GlobalArgs.Port, endpoint)
	} else {
		fullURL = fmt.Sprintf("http://%v:%v%v", args.GlobalArgs.Hostname, args.GlobalArgs.Port, endpoint)
	}
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth(args.GlobalArgs.Username, args.GlobalArgs.Password)

	return req, nil
}
