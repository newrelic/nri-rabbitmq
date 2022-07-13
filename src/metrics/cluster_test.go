package metrics

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/src/args"
	"github.com/newrelic/nri-rabbitmq/src/data"
	"github.com/newrelic/nri-rabbitmq/src/testutils"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// TODO remove global args.
	// This test are heavily based on global args to create entities.
	args.GlobalArgs = args.RabbitMQArguments{
		Hostname: "foo",
		Port:     8000,
	}

	os.Exit(m.Run())
}

func TestPopulateClusterInventory(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	PopulateClusterData(i, nil)
	assert.Empty(t, i.Entities)

	overviewData := &data.OverviewData{}
	PopulateClusterData(i, overviewData)
	assert.Empty(t, i.Entities)

	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "populateClusterEntity.json"), overviewData)

	PopulateClusterData(i, overviewData)
	assert.Equal(t, 1, len(i.Entities))

	entityKeyPrefix := fmt.Sprintf("%s:%d:", args.GlobalArgs.Hostname, args.GlobalArgs.Port)
	assert.Equal(t, entityKeyPrefix+"my-cluster", i.Entities[0].Metadata.Name)
	assert.Equal(t, "ra-cluster", i.Entities[0].Metadata.Namespace)
	assert.Equal(t, 2, len(i.Entities[0].Inventory.Items()))

	item, ok := i.Entities[0].Inventory.Item("version/rabbitmq")
	assert.True(t, ok)
	assert.Equal(t, "1.0.1", item["value"])

	item, ok = i.Entities[0].Inventory.Item("version/management")
	assert.True(t, ok)
	assert.Equal(t, "2.0.2", item["value"])
}
