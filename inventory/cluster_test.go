package inventory

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-rabbitmq/testutils"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
)

func Test_PopulateClusterEntity(t *testing.T) {
	i := testutils.GetTestingIntegration(t)
	PopulateClusterEntity(i, nil)
	assert.Empty(t, i.Entities)

	data := objx.MSI("missing", "cluster_name")
	PopulateClusterEntity(i, data)
	assert.Empty(t, i.Entities)

	data = testutils.ReadObjectFromJSONFile(t, filepath.Join("testdata", "populateClusterEntity.json"))

	PopulateClusterEntity(i, data)
	assert.Equal(t, 1, len(i.Entities))

	actual, err := json.Marshal(i.Entities[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, actual)
}
