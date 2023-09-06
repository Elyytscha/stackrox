// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"
	"time"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/sac/resources"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

var (
	// CreateTableImageCveEdgesStmt holds the create statement for table `image_cve_edges`.
	CreateTableImageCveEdgesStmt = &postgres.CreateStmts{
		GormModel: (*ImageCveEdges)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ImageCveEdgesSchema is the go schema for table `image_cve_edges`.
	ImageCveEdgesSchema = func() *walker.Schema {
		schema := GetSchemaForTable("image_cve_edges")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ImageCVEEdge)(nil)), "image_cve_edges")
		referencedSchemas := map[string]*walker.Schema{
			"storage.Image":    ImagesSchema,
			"storage.ImageCVE": ImageCvesSchema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_IMAGE_VULN_EDGE, "imagecveedge", (*storage.ImageCVEEdge)(nil)))
		schema.SetSearchScope([]v1.SearchCategory{
			v1.SearchCategory_IMAGE_VULNERABILITIES,
			v1.SearchCategory_COMPONENT_VULN_EDGE,
			v1.SearchCategory_IMAGE_COMPONENTS,
			v1.SearchCategory_IMAGE_COMPONENT_EDGE,
			v1.SearchCategory_IMAGE_VULN_EDGE,
			v1.SearchCategory_IMAGES,
			v1.SearchCategory_DEPLOYMENTS,
			v1.SearchCategory_NAMESPACES,
			v1.SearchCategory_CLUSTERS,
		}...)
		schema.ScopingResource = resources.Image
		RegisterTable(schema, CreateTableImageCveEdgesStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory_IMAGE_VULN_EDGE, schema)
		return schema
	}()
)

const (
	// ImageCveEdgesTableName specifies the name of the table in postgres.
	ImageCveEdgesTableName = "image_cve_edges"
)

// ImageCveEdges holds the Gorm model for Postgres table `image_cve_edges`.
type ImageCveEdges struct {
	ID                   string                     `gorm:"column:id;type:varchar;primaryKey"`
	FirstImageOccurrence *time.Time                 `gorm:"column:firstimageoccurrence;type:timestamp"`
	State                storage.VulnerabilityState `gorm:"column:state;type:integer"`
	ImageID              string                     `gorm:"column:imageid;type:varchar;index:imagecveedges_imageid,type:hash"`
	ImageCveID           string                     `gorm:"column:imagecveid;type:varchar;index:imagecveedges_imagecveid,type:hash"`
	Serialized           []byte                     `gorm:"column:serialized;type:bytea"`
	ImagesRef            Images                     `gorm:"foreignKey:imageid;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}
