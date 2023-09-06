// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
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
	// CreateTableNodeCvesStmt holds the create statement for table `node_cves`.
	CreateTableNodeCvesStmt = &postgres.CreateStmts{
		GormModel: (*NodeCves)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// NodeCvesSchema is the go schema for table `node_cves`.
	NodeCvesSchema = func() *walker.Schema {
		schema := GetSchemaForTable("node_cves")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.NodeCVE)(nil)), "node_cves")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_NODE_VULNERABILITIES, "nodecve", (*storage.NodeCVE)(nil)))
		schema.SetSearchScope([]v1.SearchCategory{
			v1.SearchCategory_NODE_VULNERABILITIES,
			v1.SearchCategory_NODE_COMPONENT_CVE_EDGE,
			v1.SearchCategory_NODE_COMPONENTS,
			v1.SearchCategory_NODE_COMPONENT_EDGE,
			v1.SearchCategory_NODES,
			v1.SearchCategory_CLUSTERS,
		}...)
		schema.ScopingResource = resources.Node
		RegisterTable(schema, CreateTableNodeCvesStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory_NODE_VULNERABILITIES, schema)
		return schema
	}()
)

const (
	// NodeCvesTableName specifies the name of the table in postgres.
	NodeCvesTableName = "node_cves"
)

// NodeCves holds the Gorm model for Postgres table `node_cves`.
type NodeCves struct {
	ID                     string                        `gorm:"column:id;type:varchar;primaryKey"`
	CveBaseInfoCve         string                        `gorm:"column:cvebaseinfo_cve;type:varchar;index:nodecves_cvebaseinfo_cve,type:hash"`
	CveBaseInfoPublishedOn *time.Time                    `gorm:"column:cvebaseinfo_publishedon;type:timestamp"`
	CveBaseInfoCreatedAt   *time.Time                    `gorm:"column:cvebaseinfo_createdat;type:timestamp"`
	OperatingSystem        string                        `gorm:"column:operatingsystem;type:varchar"`
	Cvss                   float32                       `gorm:"column:cvss;type:numeric"`
	Severity               storage.VulnerabilitySeverity `gorm:"column:severity;type:integer"`
	ImpactScore            float32                       `gorm:"column:impactscore;type:numeric"`
	Snoozed                bool                          `gorm:"column:snoozed;type:bool"`
	SnoozeExpiry           *time.Time                    `gorm:"column:snoozeexpiry;type:timestamp"`
	Serialized             []byte                        `gorm:"column:serialized;type:bytea"`
}
