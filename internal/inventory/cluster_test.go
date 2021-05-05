package inventory

import (
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/internal/data"
	"github.com/newrelic/nri-rabbitmq/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestPopulateClusterInventory(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	PopulateClusterInventory(i, nil)
	assert.Empty(t, i.Entities)

	overviewData := &data.OverviewData{}
	PopulateClusterInventory(i, overviewData)
	assert.Empty(t, i.Entities)

	testutils.ReadStructFromJSONFile(t, filepath.Join("testdata", "populateClusterEntity.json"), overviewData)

	PopulateClusterInventory(i, overviewData)
	assert.Equal(t, 1, len(i.Entities))
	assert.Equal(t, "my-cluster", i.Entities[0].Metadata.Name)
	assert.Equal(t, "ra-cluster", i.Entities[0].Metadata.Namespace)
	assert.Equal(t, 2, len(i.Entities[0].Inventory.Items()))

	item, ok := i.Entities[0].Inventory.Item("version/rabbitmq")
	assert.True(t, ok)
	assert.Equal(t, "1.0.1", item["value"])

	item, ok = i.Entities[0].Inventory.Item("version/management")
	assert.True(t, ok)
	assert.Equal(t, "2.0.2", item["value"])
}
