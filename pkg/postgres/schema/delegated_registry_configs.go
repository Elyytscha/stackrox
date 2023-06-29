// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

var (
	// CreateTableDelegatedRegistryConfigsStmt holds the create statement for table `delegated_registry_configs`.
	CreateTableDelegatedRegistryConfigsStmt = &postgres.CreateStmts{
		GormModel: (*DelegatedRegistryConfigs)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// DelegatedRegistryConfigsSchema is the go schema for table `delegated_registry_configs`.
	DelegatedRegistryConfigsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("delegated_registry_configs")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.DelegatedRegistryConfig)(nil)), "delegated_registry_configs")
		RegisterTable(schema, CreateTableDelegatedRegistryConfigsStmt)
		return schema
	}()
)

const (
	// DelegatedRegistryConfigsTableName specifies the name of the table in postgres.
	DelegatedRegistryConfigsTableName = "delegated_registry_configs"
)

// DelegatedRegistryConfigs holds the Gorm model for Postgres table `delegated_registry_configs`.
type DelegatedRegistryConfigs struct {
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}