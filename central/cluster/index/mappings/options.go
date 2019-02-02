package mappings

import (
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search/blevesearch"
)

// OptionsMap describes the options for Clusters
var OptionsMap = blevesearch.Walk(v1.SearchCategory_CLUSTERS, "cluster", (*storage.Cluster)(nil))
