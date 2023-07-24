// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"
	"time"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

var (
	// CreateTableBillingMetricsStmt holds the create statement for table `billing_metrics`.
	CreateTableBillingMetricsStmt = &postgres.CreateStmts{
		GormModel: (*BillingMetrics)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// BillingMetricsSchema is the go schema for table `billing_metrics`.
	BillingMetricsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("billing_metrics")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.BillingMetrics)(nil)), "billing_metrics")
		RegisterTable(schema, CreateTableBillingMetricsStmt)
		return schema
	}()
)

const (
	// BillingMetricsTableName specifies the name of the table in postgres.
	BillingMetricsTableName = "billing_metrics"
)

// BillingMetrics holds the Gorm model for Postgres table `billing_metrics`.
type BillingMetrics struct {
	Ts         *time.Time `gorm:"column:ts;type:timestamp;primaryKey"`
	Serialized []byte     `gorm:"column:serialized;type:bytea"`
}