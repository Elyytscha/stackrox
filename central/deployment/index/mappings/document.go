package mappings

import (
	"github.com/blevesearch/bleve/mapping"
	"github.com/stackrox/rox/pkg/search/blevesearch"
)

// DocumentMap provides the document mapping for the indexer to use.
var DocumentMap = func() *mapping.DocumentMapping {
	documentMap := blevesearch.FieldsToDocumentMapping(OptionsMap)
	blevesearch.AddDefaultTypeField(documentMap)
	return documentMap
}()
