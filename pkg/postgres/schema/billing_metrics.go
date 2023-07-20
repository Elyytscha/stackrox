package schema

// Code not generated by pg-bindings generator. DO EDIT.

import (
	"reflect"
	"time"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

// BillingMetricsTableName specifies the name of the table in postgres.
const BillingMetricsTableName = "billing_metrics"

var (
	// createTableBillingMetricsStmt holds the create statement for table `billing_metrics`.
	createTableBillingMetricsStmt = &postgres.CreateStmts{
		GormModel: (*BillingMetrics)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// BillingMetricsSchema is the go schema for table `billingmetrics`.
	BillingMetricsSchema = func() *walker.Schema {
		schema := GetSchemaForTable(BillingMetricsTableName)
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.BillingMetrics)(nil)), BillingMetricsTableName)
		RegisterTable(schema, createTableBillingMetricsStmt)
		return schema
	}()
)

// BillingMetrics holds the Gorm model for Postgres table `billing_metrics`.
type BillingMetrics struct {
	Ts         time.Time `gorm:"column:ts;type:timestamp;primaryKey"`
	Serialized []byte    `gorm:"column:serialized;type:bytea"`
}
