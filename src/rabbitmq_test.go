package main

import (
	"bytes"
	"fmt"
	"github.com/newrelic/nri-rabbitmq/src"
	args2 "github.com/newrelic/nri-rabbitmq/src/args"
	client2 "github.com/newrelic/nri-rabbitmq/src/client"
	testutils2 "github.com/newrelic/nri-rabbitmq/src/testutils"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/log"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	mux, closer := testutils2.GetTestServer(false)
	defer closer()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		if r.RequestURI == fmt.Sprintf(client2.HealthCheckEndpoint, "node1") {
			fmt.Fprint(w, `{"status":"ok"}`)
		} else if r.RequestURI == client2.NodesEndpoint {
			fmt.Fprintf(w, `[{ "name": "node1" }]`)
		} else if r.RequestURI == client2.OverviewEndpoint {
			fmt.Fprint(w, "{}")
		} else {
			fmt.Fprint(w, "[]")
		}
	})
	origArgs := os.Args
	os.Args = []string{
		"nri-rabbitmq",
		"-node_name_override", "node1",
		"-config_path", "",
		"-hostname", args2.GlobalArgs.Hostname,
		"-port", strconv.Itoa(args2.GlobalArgs.Port),
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
		src.main()
	})
	w.Close()
	os.Stdout = origStdout
	out := <-outC

	assert.Equal(t, fmt.Sprintf(`{"name":%q,"protocol_version":"3","integration_version":%q,"data":[{"entity":{"name":"node1","type":"ra-node","id_attributes":[{"Key":"clusterName","Value":""}]},"metrics":[{"clusterName":"","displayName":"node1","entityName":"node:node1","event_type":"RabbitmqNodeSample","node.partitionsSeen":0,"reportingEndpoint":"127.0.0.1:%d"}],"inventory":{"config/nodeName":{"value":"node1"}},"events":[]}]}%s`, src.integrationName, src.integrationVersion, args2.GlobalArgs.Port, "\n"), out)
}

func Test_getNeededData(t *testing.T) {
	mux, closer := testutils2.GetTestServer(false)
	defer closer()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		if r.RequestURI == client2.OverviewEndpoint ||
			r.RequestURI == fmt.Sprintf(client2.AlivenessTestEndpoint, "") ||
			r.RequestURI == fmt.Sprintf(client2.HealthCheckEndpoint, "node1") {
			fmt.Fprint(w, "{}")
		} else {
			fmt.Fprint(w, "[{}]")
		}
	})

	rabbitData := src.getNeededData()
	assert.NotNil(t, rabbitData)
	assert.NotNil(t, rabbitData.overview)
	assert.Equal(t, 1, len(rabbitData.bindings))
	assert.Equal(t, 1, len(rabbitData.connections))
	assert.Equal(t, 1, len(rabbitData.exchanges))
	assert.Equal(t, 1, len(rabbitData.nodes))
	assert.Equal(t, 1, len(rabbitData.queues))
	assert.Equal(t, 1, len(rabbitData.vhosts))

	metricData := src.getMetricEntities(rabbitData)
	assert.Equal(t, 3, len(metricData))
}
