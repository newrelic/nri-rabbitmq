package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	nrHttp "github.com/newrelic/infra-integrations-sdk/http"
)

const (
	overviewEndpoint    = "/api/overview"
	nodesEndpoint       = "/api/nodes"
	queuesEndpoint      = "/api/queues"
	exchangesEndpoint   = "/api/exchanges"
	vhostsEndpoint      = "/api/vhosts"
	connectionsEndpoint = "/api/connections"
	bindingsEndpoint    = "/api/bindings"
)

var (
	allEndpoints = []string{
		overviewEndpoint,
		nodesEndpoint,
		queuesEndpoint,
		exchangesEndpoint,
		vhostsEndpoint,
		connectionsEndpoint,
		bindingsEndpoint,
	}
	inventoryEndpoints = []string{
		overviewEndpoint,
		nodesEndpoint,
		queuesEndpoint,
		exchangesEndpoint,
	}
)

func collectEndpoints(endpoints ...string) (resultMap map[string]interface{}, err error) {
	resultMap = make(map[string]interface{})

	if len(endpoints) == 0 {
		return
	}

	client, err := nrHttp.New(args.CABundleFile, args.CABundleDir, time.Second*30)
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
	defer checkErr(resp.Body.Close)

	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func createRequest(endpoint string) (*http.Request, error) {
	var fullURL string
	if args.UseSSL {
		fullURL = fmt.Sprintf("https://%v:%v%v", args.Hostname, args.Port, endpoint)
	} else {
		fullURL = fmt.Sprintf("http://%v:%v%v", args.Hostname, args.Port, endpoint)
	}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(args.Username, args.Password)

	return req, nil
}
