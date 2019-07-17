package datastore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	searchMock "github.com/stackrox/rox/central/image/datastore/internal/search/mocks"
	storeMock "github.com/stackrox/rox/central/image/datastore/internal/store/mocks"
	indexMock "github.com/stackrox/rox/central/image/index/mocks"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stretchr/testify/suite"
)

func TestImageDataStore(t *testing.T) {
	suite.Run(t, new(ImageDataStoreTestSuite))
}

type ImageDataStoreTestSuite struct {
	suite.Suite

	hasWriteCtx context.Context

	mockIndexer  *indexMock.MockIndexer
	mockSearcher *searchMock.MockSearcher
	mockStore    *storeMock.MockStore

	datastore DataStore

	mockCtrl *gomock.Controller
}

func (suite *ImageDataStoreTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())

	suite.hasWriteCtx = sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowFixedScopes(
		sac.AccessModeScopeKeys(storage.Access_READ_WRITE_ACCESS),
		sac.ResourceScopeKeys(resources.Image)))

	suite.mockIndexer = indexMock.NewMockIndexer(suite.mockCtrl)
	suite.mockIndexer.EXPECT().GetTxnCount().Return(uint64(1))

	suite.mockSearcher = searchMock.NewMockSearcher(suite.mockCtrl)
	suite.mockStore = storeMock.NewMockStore(suite.mockCtrl)
	suite.mockStore.EXPECT().GetTxnCount().Return(uint64(1), nil)

	var err error
	suite.datastore, err = newDatastoreImpl(suite.mockStore, suite.mockIndexer, suite.mockSearcher)
	suite.Require().NoError(err)
}

func (suite *ImageDataStoreTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

// Scenario: We have a new image with a sha and no scan or metadata. And no previously matched registry shas.
// Outcome: Image should be upserted and indexed unchanged.
func (suite *ImageDataStoreTestSuite) TestNewImageAddedWithoutMetadata() {
	image := &storage.Image{
		Id: "sha1",
	}

	suite.mockStore.EXPECT().GetImage("sha1").Return((*storage.Image)(nil), false, nil)

	suite.mockStore.EXPECT().UpsertImage(image).Return(nil)
	suite.mockIndexer.EXPECT().AddImage(image).Return(nil)

	err := suite.datastore.UpsertImage(suite.hasWriteCtx, image)
	suite.NoError(err)
}

// Scenario: We have a new image with metadata, but its sha and the registry sha do not match.
// Outcome: The sha should be changed to the registry sha, and the mapping added to the store.
func (suite *ImageDataStoreTestSuite) TestNewImageAddedWithMetadata() {
	newImage := &storage.Image{
		Id:       "sha1",
		Metadata: &storage.ImageMetadata{},
	}
	upsertedImage := &storage.Image{
		Id:       "sha1",
		Metadata: &storage.ImageMetadata{},
	}

	suite.mockStore.EXPECT().GetImage("sha1").Return((*storage.Image)(nil), false, nil)
	suite.mockStore.EXPECT().UpsertImage(upsertedImage).Return(nil)
	suite.mockIndexer.EXPECT().AddImage(upsertedImage).Return(nil)

	err := suite.datastore.UpsertImage(suite.hasWriteCtx, newImage)
	suite.NoError(err)
}

func (suite *ImageDataStoreTestSuite) TestNewImageAddedWithScanStats() {
	newImage := &storage.Image{
		Id: "sha1",
		Scan: &storage.ImageScan{
			Components: []*storage.ImageScanComponent{
				{
					Vulns: []*storage.Vulnerability{
						{
							Cve: "derp",
							SetFixedBy: &storage.Vulnerability_FixedBy{
								FixedBy: "v1.2",
							},
						},
						{
							Cve: "derp2",
						},
					},
				},
			},
		},
	}
	upsertedImage := &storage.Image{
		Id: "sha1",
		Scan: &storage.ImageScan{
			Components: []*storage.ImageScanComponent{
				{
					Vulns: []*storage.Vulnerability{
						{
							Cve: "derp",
							SetFixedBy: &storage.Vulnerability_FixedBy{
								FixedBy: "v1.2",
							},
						},
						{
							Cve: "derp2",
						},
					},
				},
			},
		},
		SetComponents: &storage.Image_Components{
			Components: 1,
		},
		SetCves: &storage.Image_Cves{
			Cves: 2,
		},
		SetFixable: &storage.Image_FixableCves{
			FixableCves: 1,
		},
	}

	suite.mockStore.EXPECT().GetImage("sha1").Return((*storage.Image)(nil), false, nil)
	suite.mockStore.EXPECT().UpsertImage(upsertedImage).Return(nil)
	suite.mockIndexer.EXPECT().AddImage(upsertedImage).Return(nil)

	err := suite.datastore.UpsertImage(suite.hasWriteCtx, newImage)
	suite.NoError(err)
}
