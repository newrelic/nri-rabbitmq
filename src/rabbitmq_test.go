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
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	mux, closer := testutils.GetTestServer(false)
	defer closer()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		if r.RequestURI == client.OverviewEndpoint {
			fmt.Fprint(w, "{}")
		} else {
			fmt.Fprint(w, "[]")
		}
	})
	origArgs := os.Args
	os.Args = []string{
		"nr-rabbitmq",
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
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	assert.NotPanics(t, func() {
		main()
	})
	w.Close()
	os.Stdout = origStdout
	out := <-outC

	assert.Equal(t, fmt.Sprintf(`{"name":%q,"protocol_version":"2","integration_version":%q,"data":[]}`, integrationName, integrationVersion), out)
}

func TestGetNeededData(t *testing.T) {
	mux, closer := testutils.GetTestServer(false)
	defer closer()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		if r.RequestURI == client.OverviewEndpoint {
			fmt.Fprint(w, "{}")
		} else {
			fmt.Fprint(w, "[{}]")
		}
	})

	rabbitData := getNeededData()
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
