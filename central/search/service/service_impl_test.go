package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	alertMocks "github.com/stackrox/rox/central/alert/datastore/mocks"
	clusterDataStoreMocks "github.com/stackrox/rox/central/cluster/datastore/mocks"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	deploymentMocks "github.com/stackrox/rox/central/deployment/datastore/mocks"
	"github.com/stackrox/rox/central/globalindex"
	imageMocks "github.com/stackrox/rox/central/image/datastore/mocks"
	namespaceMocks "github.com/stackrox/rox/central/namespace/datastore/mocks"
	nodeMocks "github.com/stackrox/rox/central/node/globaldatastore/mocks"
	policyDatastore "github.com/stackrox/rox/central/policy/datastore"
	policyMocks "github.com/stackrox/rox/central/policy/datastore/mocks"
	policyIndex "github.com/stackrox/rox/central/policy/index"
	"github.com/stackrox/rox/central/processindicator/datastore/mocks"
	roleMocks "github.com/stackrox/rox/central/rbac/k8srole/datastore/mocks"
	roleBindingsMocks "github.com/stackrox/rox/central/rbac/k8srolebinding/datastore/mocks"
	riskDatastoreMocks "github.com/stackrox/rox/central/risk/datastore/mocks"
	secretMocks "github.com/stackrox/rox/central/secret/datastore/mocks"
	serviceAccountMocks "github.com/stackrox/rox/central/serviceaccount/datastore/mocks"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchCategoryToResourceMap(t *testing.T) {
	allCategories := set.NewV1SearchCategorySet(GetAllSearchableCategories()...).Union(autocompleteCategories)
	categoryToResource := GetSearchCategoryToResourceMetadata()
	for _, searchCategory := range allCategories.AsSlice() {
		_, ok := categoryToResource[searchCategory]
		// This is a programming error. If you see this, add the new category you've added to the
		// SearchCategoryToResource map!
		assert.True(t, ok, "Please add category %s to the SearchCategoryToResource map used by the authorizer", searchCategory.String())
	}
}

func TestSearchCategoryToOptionsMultiMap(t *testing.T) {
	t.Parallel()

	for _, cat := range autocompleteCategories.AsSlice() {
		_, ok := categoryToOptionsMultimap[cat]
		assert.True(t, ok, "no options multimap for category", cat)
	}
}

func TestSearchFuncs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := NewBuilder().
		WithAlertStore(alertMocks.NewMockDataStore(mockCtrl)).
		WithDeploymentStore(deploymentMocks.NewMockDataStore(mockCtrl)).
		WithImageStore(imageMocks.NewMockDataStore(mockCtrl)).
		WithPolicyStore(policyMocks.NewMockDataStore(mockCtrl)).
		WithSecretStore(secretMocks.NewMockDataStore(mockCtrl)).
		WithServiceAccountStore(serviceAccountMocks.NewMockDataStore(mockCtrl)).
		WithNodeStore(nodeMocks.NewMockGlobalDataStore(mockCtrl)).
		WithNamespaceStore(namespaceMocks.NewMockDataStore(mockCtrl)).
		WithRiskStore(riskDatastoreMocks.NewMockDataStore(mockCtrl)).
		WithRoleStore(roleMocks.NewMockDataStore(mockCtrl)).
		WithRoleBindingStore(roleBindingsMocks.NewMockDataStore(mockCtrl)).
		WithClusterDataStore(clusterDataStoreMocks.NewMockDataStore(mockCtrl)).
		WithAggregator(nil).
		Build()

	searchFuncMap := s.(*serviceImpl).getSearchFuncs()
	for _, searchCategory := range GetAllSearchableCategories() {
		_, ok := searchFuncMap[searchCategory]
		// This is a programming error. If you see this, add the new category you've added to the
		// SearchCategoryToResource map!
		assert.True(t, ok, "Please add category %s to the map in getSearchFuncs()", searchCategory.String())
	}
}

func TestAutocomplete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create Deployment Indexer
	idx, err := globalindex.MemOnlyIndex()
	require.NoError(t, err)

	testDB := testutils.BadgerDBForT(t)
	defer utils.IgnoreError(testDB.Close)

	mockIndicators := mocks.NewMockDataStore(mockCtrl)
	// This gets called as a side effect of `UpsertDeployment`.
	mockIndicators.EXPECT().RemoveProcessIndicatorsOfStaleContainers(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	mockRiskDatastore := riskDatastoreMocks.NewMockDataStore(mockCtrl)
	mockRiskDatastore.EXPECT().SearchRawRisks(gomock.Any(), gomock.Any())
	mockRiskDatastore.EXPECT().GetRisk(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	deploymentDS, err := deploymentDatastore.NewBadger(testDB, idx, nil, mockIndicators, nil, nil, mockRiskDatastore, nil, nil)
	require.NoError(t, err)

	allAccessCtx := sac.WithAllAccess(context.Background())

	deploymentNameOneOff := fixtures.GetDeployment()
	require.NoError(t, deploymentDS.UpsertDeployment(allAccessCtx, deploymentNameOneOff))

	deploymentName1 := fixtures.GetDeployment()
	deploymentName1.Id = "name1"
	deploymentName1.Name = "name1"
	require.NoError(t, deploymentDS.UpsertDeployment(allAccessCtx, deploymentName1))

	deploymentName1Duplicate := fixtures.GetDeployment()
	deploymentName1Duplicate.Id = "name1Dup"
	deploymentName1Duplicate.Name = "name1"
	require.NoError(t, deploymentDS.UpsertDeployment(allAccessCtx, deploymentName1Duplicate))

	deploymentName2 := fixtures.GetDeployment()
	deploymentName2.Id = "name12"
	deploymentName2.Name = "name12"
	deploymentName2.Labels = map[string]string{"hello": "hi", "hey": "ho"}
	require.NoError(t, deploymentDS.UpsertDeployment(allAccessCtx, deploymentName2))

	service := NewBuilder().
		WithAlertStore(alertMocks.NewMockDataStore(mockCtrl)).
		WithDeploymentStore(deploymentDS).
		WithImageStore(imageMocks.NewMockDataStore(mockCtrl)).
		WithPolicyStore(policyMocks.NewMockDataStore(mockCtrl)).
		WithSecretStore(secretMocks.NewMockDataStore(mockCtrl)).
		WithServiceAccountStore(serviceAccountMocks.NewMockDataStore(mockCtrl)).
		WithNodeStore(nodeMocks.NewMockGlobalDataStore(mockCtrl)).
		WithNamespaceStore(namespaceMocks.NewMockDataStore(mockCtrl)).
		WithRiskStore(riskDatastoreMocks.NewMockDataStore(mockCtrl)).
		WithRoleStore(roleMocks.NewMockDataStore(mockCtrl)).
		WithRoleBindingStore(roleBindingsMocks.NewMockDataStore(mockCtrl)).
		WithClusterDataStore(clusterDataStoreMocks.NewMockDataStore(mockCtrl)).
		WithAggregator(nil).
		Build().(*serviceImpl)

	for _, testCase := range []struct {
		query           string
		expectedResults []string
		ignoreOrder     bool
	}{
		{
			query:           search.NewQueryBuilder().AddStrings(search.DeploymentName, deploymentNameOneOff.Name).Query(),
			expectedResults: []string{deploymentNameOneOff.GetName()},
		},
		{
			query: search.NewQueryBuilder().AddStrings(search.DeploymentName, "name").Query(),
			// This is odd, but this is correct. Bleve scores name12 higher than name1
			expectedResults: []string{"name12", "name1"},
		},
		{
			query:           fmt.Sprintf("%s:", search.DeploymentName),
			expectedResults: []string{"name12", "nginx_server", "name1"},
		},
		{
			query:           fmt.Sprintf("%s:name12,", search.DeploymentName),
			expectedResults: []string{"name12", "nginx_server", "name1"},
		},
		{
			query:           fmt.Sprintf("%s:he=h", search.Label),
			expectedResults: []string{"hello=hi", "hey=ho"},
			ignoreOrder:     true,
		},
		{
			query:           fmt.Sprintf("%s:hey=", search.Label),
			expectedResults: []string{"hey=ho"},
			ignoreOrder:     true,
		},
		{
			query:           fmt.Sprintf("%s:%s+%s:", search.DeploymentName, deploymentName2.Name, search.Label),
			expectedResults: []string{"hello=hi", "hey=ho"},
			ignoreOrder:     true,
		},
	} {
		t.Run(fmt.Sprintf("Test case %q", testCase.query), func(t *testing.T) {
			results, err := service.autocomplete(allAccessCtx, testCase.query, []v1.SearchCategory{v1.SearchCategory_DEPLOYMENTS})
			require.NoError(t, err)
			if testCase.ignoreOrder {
				assert.ElementsMatch(t, testCase.expectedResults, results)
			} else {
				assert.Equal(t, testCase.expectedResults, results)
			}
		})
	}
}

func TestAutocompleteForEnums(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create Deployment Indexer
	idx, err := globalindex.MemOnlyIndex()
	require.NoError(t, err)
	policyIndexer := policyIndex.New(idx)
	require.NoError(t, policyIndexer.AddPolicy(fixtures.GetPolicy()))

	ds := policyDatastore.New(nil, policyIndexer, nil)

	service := NewBuilder().
		WithAlertStore(alertMocks.NewMockDataStore(mockCtrl)).
		WithDeploymentStore(deploymentMocks.NewMockDataStore(mockCtrl)).
		WithImageStore(imageMocks.NewMockDataStore(mockCtrl)).
		WithPolicyStore(ds).
		WithSecretStore(secretMocks.NewMockDataStore(mockCtrl)).
		WithSecretStore(secretMocks.NewMockDataStore(mockCtrl)).
		WithServiceAccountStore(serviceAccountMocks.NewMockDataStore(mockCtrl)).
		WithNodeStore(nodeMocks.NewMockGlobalDataStore(mockCtrl)).
		WithNamespaceStore(namespaceMocks.NewMockDataStore(mockCtrl)).
		WithRoleStore(roleMocks.NewMockDataStore(mockCtrl)).
		WithRoleBindingStore(roleBindingsMocks.NewMockDataStore(mockCtrl)).
		WithClusterDataStore(clusterDataStoreMocks.NewMockDataStore(mockCtrl)).
		WithAggregator(nil).
		Build().(*serviceImpl)

	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowAllAccessScopeChecker())
	results, err := service.autocomplete(ctx, fmt.Sprintf("%s:", search.Severity), []v1.SearchCategory{v1.SearchCategory_POLICIES})
	require.NoError(t, err)
	assert.Equal(t, []string{fixtures.GetPolicy().GetSeverity().String()}, results)
}
