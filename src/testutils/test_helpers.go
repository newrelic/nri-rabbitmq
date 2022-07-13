package testutils

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/args"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Update flag will update .golden files to the current actual
var Update = flag.Bool("update", false, "update .golden files")

// GetTestingIntegration creates an Integration used for testing and sets the logger to the integration's logger
func GetTestingIntegration(t *testing.T) (payload *integration.Integration) {
	testLogger := &TestLogger{F: t.Logf}
	payload, err := integration.New("Test", "0.0.1", integration.Logger(testLogger))
	require.NoError(t, err)
	require.NotNil(t, payload)
	return
}

// GetTestingEntity creates an Entity used for testing
func GetTestingEntity(t *testing.T, entityArgs ...string) (payload *integration.Integration, entity *integration.Entity) {
	payload = GetTestingIntegration(t)
	var err error
	if len(entityArgs) > 1 {
		entity, err = payload.Entity(entityArgs[0], entityArgs[1])
		assert.NoError(t, err)
	} else {
		entity = payload.LocalEntity()
	}
	require.NotNil(t, entity)
	return
}

// ReadStructFromJSONFile Unmarshals the json file into the specified object
func ReadStructFromJSONFile(t *testing.T, filename string, object interface{}) {
	data, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	err = json.Unmarshal(data, object)
	require.NoError(t, err)
}

// ReadStructFromJSONString reads a generic map[string]interface{} from a json string
func ReadStructFromJSONString(t *testing.T, rawJSON string, object interface{}) {
	err := json.Unmarshal([]byte(rawJSON), &object)
	require.NoError(t, err)
}

// GetTestServer creates a test server
func GetTestServer(tls bool) (mux *http.ServeMux, teardown func()) {
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
