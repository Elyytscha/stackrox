// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"
	"time"

	"github.com/lib/pq"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

var (
	// CreateTableTestSingleUuidKeyStructsStmt holds the create statement for table `test_single_uuid_key_structs`.
	CreateTableTestSingleUuidKeyStructsStmt = &postgres.CreateStmts{
		GormModel: (*TestSingleUuidKeyStructs)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// TestSingleUuidKeyStructsSchema is the go schema for table `test_single_uuid_key_structs`.
	TestSingleUuidKeyStructsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("test_single_uuid_key_structs")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.TestSingleUUIDKeyStruct)(nil)), "test_single_uuid_key_structs")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(115), "testsingleuuidkeystruct", (*storage.TestSingleUUIDKeyStruct)(nil)))
		RegisterTable(schema, CreateTableTestSingleUuidKeyStructsStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory(115), schema)
		return schema
	}()
)

const (
	TestSingleUuidKeyStructsTableName = "test_single_uuid_key_structs"
)

// TestSingleUuidKeyStructs holds the Gorm model for Postgres table `test_single_uuid_key_structs`.
type TestSingleUuidKeyStructs struct {
	Key         string                               `gorm:"column:key;type:uuid;primaryKey;index:testsingleuuidkeystructs_key,type:hash"`
	Name        string                               `gorm:"column:name;type:varchar;unique"`
	StringSlice *pq.StringArray                      `gorm:"column:stringslice;type:text[]"`
	Bool        bool                                 `gorm:"column:bool;type:bool"`
	Uint64      uint64                               `gorm:"column:uint64;type:bigint"`
	Int64       int64                                `gorm:"column:int64;type:bigint"`
	Float       float32                              `gorm:"column:float;type:numeric"`
	Labels      map[string]string                    `gorm:"column:labels;type:jsonb"`
	Timestamp   *time.Time                           `gorm:"column:timestamp;type:timestamp"`
	Enum        storage.TestSingleUUIDKeyStruct_Enum `gorm:"column:enum;type:integer"`
	Enums       *pq.Int32Array                       `gorm:"column:enums;type:int[]"`
	Serialized  []byte                               `gorm:"column:serialized;type:bytea"`
}