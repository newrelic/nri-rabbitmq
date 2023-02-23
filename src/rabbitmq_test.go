package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/args"
	"github.com/newrelic/nri-rabbitmq/src/client"
	"github.com/newrelic/nri-rabbitmq/src/testutils"

	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/stretchr/testify/assert"
)

const timeout = 30

func Test_main(t *testing.T) {
	mux, closer := testutils.GetTestServer(false)
	defer closer()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		switch r.RequestURI {
		case client.NodesEndpoint:
			fmt.Fprintf(w, `[{ "name": "node1", "running": false }]`)
		case client.OverviewEndpoint:
			fmt.Fprint(w, "{}")
		default:
			fmt.Fprint(w, "[]")
		}
	})
	origArgs := os.Args
	os.Args = []string{
		"nri-rabbitmq",
		"-node_name_override", "node1",
		"-config_path", "",
		"-hostname", args.GlobalArgs.Hostname,
		"-port", strconv.Itoa(args.GlobalArgs.Port),
	}
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Args = origArgs
		os.Stdout = origStdout
	}()
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			log.Error(err.Error())
		}
		outC <- buf.String()
	}()

	assert.NotPanics(t, func() {
		main()
	})
	w.Close()
	os.Stdout = origStdout
	out := <-outC

	assert.Equal(t, fmt.Sprintf(`{"name":%q,"protocol_version":"3","integration_version":%q,"data":[{"entity":{"name":"%s:%d:node1","type":"ra-node","id_attributes":[{"Key":"clusterName","Value":""}]},"metrics":[{"displayName":"node1","entityName":"node:node1","event_type":"RabbitmqNodeSample","node.partitionsSeen":0,"node.running":0,"rabbitmqClusterName":"","reportingEndpoint":"127.0.0.1:%d"}],"inventory":{"config/nodeName":{"value":"node1"}},"events":[{"summary":"Response is [%s] for node [node1] running status","category":"integration","attributes":{"reportingEndpoint":"127.0.0.1:%d"}}]}]}%s`, integrationName, integrationVersion, args.GlobalArgs.Hostname, args.GlobalArgs.Port, args.GlobalArgs.Port, NotRunning, args.GlobalArgs.Port, "\n"), out)
}

func Test_getNeededData(t *testing.T) {
	mux, closer := testutils.GetTestServer(false)
	defer closer()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		if r.RequestURI == client.OverviewEndpoint ||
			r.RequestURI == fmt.Sprintf(client.AlivenessTestEndpoint, "") ||
			r.RequestURI == fmt.Sprintf(client.HealthCheckEndpoint, "node1") {
			fmt.Fprint(w, "{}")
		} else {
			fmt.Fprint(w, "[{}]")
		}
	})

	rabbitData := getNeededData(timeout)
	assert.NotNil(t, rabbitData)
	assert.NotNil(t, rabbitData.overview)
	assert.Equal(t, 1, len(rabbitData.bindings))
	assert.Equal(t, 1, len(rabbitData.connections))
	assert.Equal(t, 1, len(rabbitData.exchanges))
	assert.Equal(t, 1, len(rabbitData.nodes))
	assert.Equal(t, 1, len(rabbitData.queues))
	assert.Equal(t, 1, len(rabbitData.vhosts))

	metricData := getMetricEntities(rabbitData)
	assert.Equal(t, 3, len(metricData))
}
