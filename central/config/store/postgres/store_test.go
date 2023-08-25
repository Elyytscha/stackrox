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

type ConfigsStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestConfigsStore(t *testing.T) {
	suite.Run(t, new(ConfigsStoreSuite))
}

func (s *ConfigsStoreSuite) SetupTest() {
	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *ConfigsStoreSuite) TearDownTest() {
	s.testDB.Teardown(s.T())
}

func (s *ConfigsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	config := &storage.Config{}
	s.NoError(testutils.FullInit(config, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundConfig, exists, err := store.Get(ctx)
	s.NoError(err)
	s.False(exists)
	s.Nil(foundConfig)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, config))
	foundConfig, exists, err = store.Get(ctx)
	s.NoError(err)
	s.True(exists)
	s.Equal(config, foundConfig)

	foundConfig, exists, err = store.Get(ctx)
	s.NoError(err)
	s.True(exists)
	s.Equal(config, foundConfig)

	s.NoError(store.Delete(ctx))
	foundConfig, exists, err = store.Get(ctx)
	s.NoError(err)
	s.False(exists)
	s.Nil(foundConfig)

	s.ErrorIs(store.Delete(withNoAccessCtx), sac.ErrResourceAccessDenied)

	config = &storage.Config{}
	s.NoError(testutils.FullInit(config, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))
	s.NoError(store.Upsert(ctx, config))

	foundConfig, exists, err = store.Get(ctx)
	s.NoError(err)
	s.True(exists)
	s.Equal(config, foundConfig)

	config = &storage.Config{}
	s.NoError(testutils.FullInit(config, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))
	s.NoError(store.Upsert(ctx, config))

	foundConfig, exists, err = store.Get(ctx)
	s.NoError(err)
	s.True(exists)
	s.Equal(config, foundConfig)
}
