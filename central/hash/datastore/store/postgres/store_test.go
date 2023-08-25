// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

type HashesStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestHashesStore(t *testing.T) {
	suite.Run(t, new(HashesStoreSuite))
}

func (s *HashesStoreSuite) SetupSuite() {

	s.T().Setenv(features.StoreEventHashes.EnvVar(), "true")
	if !features.StoreEventHashes.Enabled() {
		s.T().Skip("Skip postgres store tests because feature flag is off")
		s.T().SkipNow()
	}

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *HashesStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE hashes CASCADE")
	s.T().Log("hashes", tag)
	s.NoError(err)
}

func (s *HashesStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *HashesStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	hash := &storage.Hash{}
	s.NoError(testutils.FullInit(hash, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundHash, exists, err := store.Get(ctx, hash.GetClusterId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundHash)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, hash))
	foundHash, exists, err = store.Get(ctx, hash.GetClusterId())
	s.NoError(err)
	s.True(exists)
	s.Equal(hash, foundHash)

	hashCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, hashCount)
	hashCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(hashCount)

	hashExists, err := store.Exists(ctx, hash.GetClusterId())
	s.NoError(err)
	s.True(hashExists)
	s.NoError(store.Upsert(ctx, hash))
	s.ErrorIs(store.Upsert(withNoAccessCtx, hash), sac.ErrResourceAccessDenied)

	foundHash, exists, err = store.Get(ctx, hash.GetClusterId())
	s.NoError(err)
	s.True(exists)
	s.Equal(hash, foundHash)

	s.NoError(store.Delete(ctx, hash.GetClusterId()))
	foundHash, exists, err = store.Get(ctx, hash.GetClusterId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundHash)
	s.ErrorIs(store.Delete(withNoAccessCtx, hash.GetClusterId()), sac.ErrResourceAccessDenied)

	var hashs []*storage.Hash
	var hashIDs []string
	for i := 0; i < 200; i++ {
		hash := &storage.Hash{}
		s.NoError(testutils.FullInit(hash, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		hashs = append(hashs, hash)
		hashIDs = append(hashIDs, hash.GetClusterId())
	}

	s.NoError(store.UpsertMany(ctx, hashs))

	hashCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, hashCount)

	s.NoError(store.DeleteMany(ctx, hashIDs))

	hashCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, hashCount)
}
