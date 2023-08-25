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

type TestChild1StoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestTestChild1Store(t *testing.T) {
	suite.Run(t, new(TestChild1StoreSuite))
}

func (s *TestChild1StoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *TestChild1StoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE test_child1 CASCADE")
	s.T().Log("test_child1", tag)
	s.NoError(err)
}

func (s *TestChild1StoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *TestChild1StoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	testChild1 := &storage.TestChild1{}
	s.NoError(testutils.FullInit(testChild1, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundTestChild1, exists, err := store.Get(ctx, testChild1.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTestChild1)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, testChild1))
	foundTestChild1, exists, err = store.Get(ctx, testChild1.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(testChild1, foundTestChild1)

	testChild1Count, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, testChild1Count)
	testChild1Count, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(testChild1Count)

	testChild1Exists, err := store.Exists(ctx, testChild1.GetId())
	s.NoError(err)
	s.True(testChild1Exists)
	s.NoError(store.Upsert(ctx, testChild1))
	s.ErrorIs(store.Upsert(withNoAccessCtx, testChild1), sac.ErrResourceAccessDenied)

	foundTestChild1, exists, err = store.Get(ctx, testChild1.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(testChild1, foundTestChild1)

	s.NoError(store.Delete(ctx, testChild1.GetId()))
	foundTestChild1, exists, err = store.Get(ctx, testChild1.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTestChild1)
	s.NoError(store.Delete(withNoAccessCtx, testChild1.GetId()))

	var testChild1s []*storage.TestChild1
	var testChild1IDs []string
	for i := 0; i < 200; i++ {
		testChild1 := &storage.TestChild1{}
		s.NoError(testutils.FullInit(testChild1, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		testChild1s = append(testChild1s, testChild1)
		testChild1IDs = append(testChild1IDs, testChild1.GetId())
	}

	s.NoError(store.UpsertMany(ctx, testChild1s))

	testChild1Count, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, testChild1Count)

	s.NoError(store.DeleteMany(ctx, testChild1IDs))

	testChild1Count, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, testChild1Count)
}
