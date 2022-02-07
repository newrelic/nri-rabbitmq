package metrics

import (
	data2 "github.com/newrelic/nri-rabbitmq/src/data"
	testutils2 "github.com/newrelic/nri-rabbitmq/src/testutils"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateClusterInventory(t *testing.T) {
	i := testutils2.GetTestingIntegration(t)
	PopulateClusterData(i, nil)
	assert.Empty(t, i.Entities)

	overviewData := &data2.OverviewData{}
	PopulateClusterData(i, overviewData)
	assert.Empty(t, i.Entities)

	testutils2.ReadStructFromJSONFile(t, filepath.Join("testdata", "populateClusterEntity.json"), overviewData)

	PopulateClusterData(i, overviewData)
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
