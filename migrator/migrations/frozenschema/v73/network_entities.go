// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableNetworkEntitiesStmt holds the create statement for table `network_entities`.
	CreateTableNetworkEntitiesStmt = &postgres.CreateStmts{
		GormModel: (*NetworkEntities)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// NetworkEntitiesSchema is the go schema for table `network_entities`.
	NetworkEntitiesSchema = func() *walker.Schema {
		schema := walker.Walk(reflect.TypeOf((*storage.NetworkEntity)(nil)), "network_entities")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_NETWORK_ENTITY, "networkentity", (*storage.NetworkEntity)(nil)))
		return schema
	}()
)

const (
	NetworkEntitiesTableName = "network_entities"
)

// NetworkEntities holds the Gorm model for Postgres table `network_entities`.
type NetworkEntities struct {
	InfoId                    string `gorm:"column:info_id;type:varchar;primaryKey"`
	InfoExternalSourceDefault bool   `gorm:"column:info_externalsource_default;type:bool"`
	Serialized                []byte `gorm:"column:serialized;type:bytea"`
}
