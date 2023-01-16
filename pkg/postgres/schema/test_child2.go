// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

var (
	// CreateTableTestChild2Stmt holds the create statement for table `test_child2`.
	CreateTableTestChild2Stmt = &postgres.CreateStmts{
		GormModel: (*TestChild2)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// TestChild2Schema is the go schema for table `test_child2`.
	TestChild2Schema = func() *walker.Schema {
		schema := GetSchemaForTable("test_child2")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.TestChild2)(nil)), "test_child2")
		referencedSchemas := map[string]*walker.Schema{
			"storage.TestParent2":     TestParent2Schema,
			"storage.TestGrandparent": TestGrandparentsSchema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(104), "testchild2", (*storage.TestChild2)(nil)))
		RegisterTable(schema, CreateTableTestChild2Stmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory(104), schema)
		return schema
	}()
)

const (
	TestChild2TableName = "test_child2"
)

// TestChild2 holds the Gorm model for Postgres table `test_child2`.
type TestChild2 struct {
	Id             string      `gorm:"column:id;type:uuid;primaryKey"`
	ParentId       string      `gorm:"column:parentid;type:varchar"`
	GrandparentId  string      `gorm:"column:grandparentid;type:varchar"`
	Val            string      `gorm:"column:val;type:varchar"`
	Serialized     []byte      `gorm:"column:serialized;type:bytea"`
	TestParent2Ref TestParent2 `gorm:"foreignKey:parentid;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}
