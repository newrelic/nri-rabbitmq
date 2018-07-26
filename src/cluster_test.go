package main

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
)

func Test_PopulateClusterEntity(t *testing.T) {
	i := getTestingIntegration(t)
	populateClusterEntity(i, nil)
	assert.Empty(t, i.Entities)

	data := objx.MSI("missing", "cluster_name")
	populateClusterEntity(i, data)
	assert.Empty(t, i.Entities)

	data = readObjectFromJSONFile(t, filepath.Join("testdata", "populateClusterEntity.json"))

	populateClusterEntity(i, data)
	assert.Equal(t, 1, len(i.Entities))

	actual, err := json.Marshal(i.Entities[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, actual)
}
