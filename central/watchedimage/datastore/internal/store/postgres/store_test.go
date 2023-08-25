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

type WatchedImagesStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestWatchedImagesStore(t *testing.T) {
	suite.Run(t, new(WatchedImagesStoreSuite))
}

func (s *WatchedImagesStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *WatchedImagesStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE watched_images CASCADE")
	s.T().Log("watched_images", tag)
	s.NoError(err)
}

func (s *WatchedImagesStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *WatchedImagesStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	watchedImage := &storage.WatchedImage{}
	s.NoError(testutils.FullInit(watchedImage, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundWatchedImage, exists, err := store.Get(ctx, watchedImage.GetName())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundWatchedImage)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, watchedImage))
	foundWatchedImage, exists, err = store.Get(ctx, watchedImage.GetName())
	s.NoError(err)
	s.True(exists)
	s.Equal(watchedImage, foundWatchedImage)

	watchedImageCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, watchedImageCount)
	watchedImageCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(watchedImageCount)

	watchedImageExists, err := store.Exists(ctx, watchedImage.GetName())
	s.NoError(err)
	s.True(watchedImageExists)
	s.NoError(store.Upsert(ctx, watchedImage))
	s.ErrorIs(store.Upsert(withNoAccessCtx, watchedImage), sac.ErrResourceAccessDenied)

	foundWatchedImage, exists, err = store.Get(ctx, watchedImage.GetName())
	s.NoError(err)
	s.True(exists)
	s.Equal(watchedImage, foundWatchedImage)

	s.NoError(store.Delete(ctx, watchedImage.GetName()))
	foundWatchedImage, exists, err = store.Get(ctx, watchedImage.GetName())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundWatchedImage)
	s.ErrorIs(store.Delete(withNoAccessCtx, watchedImage.GetName()), sac.ErrResourceAccessDenied)

	var watchedImages []*storage.WatchedImage
	var watchedImageIDs []string
	for i := 0; i < 200; i++ {
		watchedImage := &storage.WatchedImage{}
		s.NoError(testutils.FullInit(watchedImage, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		watchedImages = append(watchedImages, watchedImage)
		watchedImageIDs = append(watchedImageIDs, watchedImage.GetName())
	}

	s.NoError(store.UpsertMany(ctx, watchedImages))

	watchedImageCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, watchedImageCount)

	s.NoError(store.DeleteMany(ctx, watchedImageIDs))

	watchedImageCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, watchedImageCount)
}
