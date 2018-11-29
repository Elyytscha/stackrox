package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	alertMocks "github.com/stackrox/rox/central/alert/datastore/mocks"
	deploymentMocks "github.com/stackrox/rox/central/deployment/datastore/mocks"
	imageMocks "github.com/stackrox/rox/central/image/datastore/mocks"
	policyMocks "github.com/stackrox/rox/central/policy/datastore/mocks"
	secretMocks "github.com/stackrox/rox/central/secret/datastore/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSearchCategoryToResourceMap(t *testing.T) {
	for _, searchCategory := range GetAllSearchableCategories() {
		_, ok := searchCategoryToResource[searchCategory]
		// This is a programming error. If you see this, add the new category you've added to the
		// searchCategoryToResource map!
		assert.True(t, ok, "Please add category %s to the searchCategoryToResource map used by the authorizer", searchCategory.String())
	}
}

func TestSearchFuncs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := New(alertMocks.NewMockDataStore(mockCtrl), deploymentMocks.NewMockDataStore(mockCtrl), imageMocks.NewMockDataStore(mockCtrl), policyMocks.NewMockDataStore(gomock.NewController(t)), secretMocks.NewMockDataStore(mockCtrl))
	searchFuncMap := s.(*serviceImpl).getSearchFuncs()
	for _, searchCategory := range GetAllSearchableCategories() {
		_, ok := searchFuncMap[searchCategory]
		// This is a programming error. If you see this, add the new category you've added to the
		// searchCategoryToResource map!
		assert.True(t, ok, "Please add category %s to the map in getSearchFuncs()", searchCategory.String())
	}
}
