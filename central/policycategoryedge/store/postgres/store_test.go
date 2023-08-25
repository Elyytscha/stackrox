// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

type PolicyCategoryEdgesStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestPolicyCategoryEdgesStore(t *testing.T) {
	suite.Run(t, new(PolicyCategoryEdgesStoreSuite))
}

func (s *PolicyCategoryEdgesStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *PolicyCategoryEdgesStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE policy_category_edges CASCADE")
	s.T().Log("policy_category_edges", tag)
	s.NoError(err)
}

func (s *PolicyCategoryEdgesStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *PolicyCategoryEdgesStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	policyCategoryEdge := &storage.PolicyCategoryEdge{}
	s.NoError(testutils.FullInit(policyCategoryEdge, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundPolicyCategoryEdge, exists, err := store.Get(ctx, policyCategoryEdge.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundPolicyCategoryEdge)

}
