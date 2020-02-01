package pruning

import (
	"context"
	"testing"
	"time"

	protoTypes "github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	alertDatastore "github.com/stackrox/rox/central/alert/datastore"
	alertDatastoreMocks "github.com/stackrox/rox/central/alert/datastore/mocks"
	clusterDatastoreMocks "github.com/stackrox/rox/central/cluster/datastore/mocks"
	configDatastore "github.com/stackrox/rox/central/config/datastore"
	configDatastoreMocks "github.com/stackrox/rox/central/config/datastore/mocks"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	"github.com/stackrox/rox/central/globalindex"
	imageDatastore "github.com/stackrox/rox/central/image/datastore"
	imageDatastoreMocks "github.com/stackrox/rox/central/image/datastore/mocks"
	networkFlowDatastoreMocks "github.com/stackrox/rox/central/networkflow/datastore/mocks"
	processIndicatorDatastoreMocks "github.com/stackrox/rox/central/processindicator/datastore/mocks"
	processWhitelistDatastoreMocks "github.com/stackrox/rox/central/processwhitelist/datastore/mocks"
	riskDatastoreMocks "github.com/stackrox/rox/central/risk/datastore/mocks"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/alert/convert"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/images/types"
	filterMocks "github.com/stackrox/rox/pkg/process/filter/mocks"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRetentionResolvedDeploy = 7
	testRetentionAllRuntime     = 6
	testRetentionDeletedRuntime = 3
)

var (
	testConfig = &storage.Config{
		PrivateConfig: &storage.PrivateConfig{
			AlertRetention: &storage.PrivateConfig_AlertConfig{
				AlertConfig: &storage.AlertRetentionConfig{
					AllRuntimeRetentionDurationDays:     testRetentionAllRuntime,
					DeletedRuntimeRetentionDurationDays: testRetentionDeletedRuntime,
					ResolvedDeployRetentionDurationDays: testRetentionResolvedDeploy,
				},
			},
			ImageRetentionDurationDays: configDatastore.DefaultImageRetention,
		},
	}
)

func newAlertInstance(id string, daysOld int, stage storage.LifecycleStage, state storage.ViolationState) *storage.Alert {
	return newAlertInstanceWithDeployment(id, daysOld, stage, state, nil)
}
func newAlertInstanceWithDeployment(id string, daysOld int, stage storage.LifecycleStage, state storage.ViolationState, deployment *storage.Deployment) *storage.Alert {
	var alertDeployment *storage.Alert_Deployment
	if deployment != nil {
		alertDeployment = convert.ToAlertDeployment(deployment)
	} else {
		alertDeployment = &storage.Alert_Deployment{
			Id:       "inactive",
			Inactive: true,
		}
	}
	return &storage.Alert{
		Id: id,

		LifecycleStage: stage,
		State:          state,
		Deployment:     alertDeployment,
		Time:           protoconv.ConvertTimeToTimestamp(time.Now().Add(-24 * time.Duration(daysOld) * time.Hour)),
	}
}

func newImageInstance(id string, daysOld int) *storage.Image {
	return &storage.Image{
		Id:          id,
		LastUpdated: protoconv.ConvertTimeToTimestamp(time.Now().Add(-24 * time.Duration(daysOld) * time.Hour)),
	}
}

func newDeployment(imageIDs ...string) *storage.Deployment {
	var containers []*storage.Container
	for _, id := range imageIDs {
		digest := types.NewDigest(id).Digest()
		containers = append(containers, &storage.Container{
			Image: &storage.ContainerImage{
				Id: digest,
			},
			Instances: []*storage.ContainerInstance{
				{
					ImageDigest: id,
				},
			},
		})
	}
	return &storage.Deployment{
		Id:         "id",
		Containers: containers,
	}
}

func generateImageDataStructures(ctx context.Context, t *testing.T) (alertDatastore.DataStore, configDatastore.DataStore, imageDatastore.DataStore, deploymentDatastore.DataStore) {
	db := testutils.BadgerDBForT(t)

	bleveIndex, err := globalindex.MemOnlyIndex()
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	mockRiskDatastore := riskDatastoreMocks.NewMockDataStore(ctrl)
	mockRiskDatastore.EXPECT().SearchRawRisks(gomock.Any(), gomock.Any()).AnyTimes()
	mockRiskDatastore.EXPECT().RemoveRisk(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	dacky, err := dackbox.NewDackBox(db, nil, []byte("graph"), []byte("dirty"), []byte("valid"))
	require.NoError(t, err)

	// Initialize real datastore
	images, err := imageDatastore.NewBadger(dacky, concurrency.NewKeyFence(), db, bleveIndex, true, mockRiskDatastore)
	require.NoError(t, err)

	mockProcessDataStore := processIndicatorDatastoreMocks.NewMockDataStore(ctrl)
	mockProcessDataStore.EXPECT().RemoveProcessIndicatorsOfStaleContainers(gomock.Any(), gomock.Any()).Return((error)(nil))

	mockWhitelistDataStore := processWhitelistDatastoreMocks.NewMockDataStore(ctrl)

	mockConfigDatastore := configDatastoreMocks.NewMockDataStore(ctrl)
	mockConfigDatastore.EXPECT().GetConfig(ctx).Return(testConfig, nil)

	mockAlertDatastore := alertDatastoreMocks.NewMockDataStore(ctrl)

	mockFilter := filterMocks.NewMockFilter(ctrl)
	mockFilter.EXPECT().Update(gomock.Any()).AnyTimes()

	deployments, err := deploymentDatastore.NewBadger(dacky, concurrency.NewKeyFence(), db, bleveIndex, nil, mockProcessDataStore, mockWhitelistDataStore, nil, mockRiskDatastore, nil, mockFilter)
	require.NoError(t, err)

	return mockAlertDatastore, mockConfigDatastore, images, deployments
}

func generateAlertDataStructures(ctx context.Context, t *testing.T) (alertDatastore.DataStore, configDatastore.DataStore, imageDatastore.DataStore, deploymentDatastore.DataStore) {
	db := testutils.BadgerDBForT(t)

	bleveIndex, err := globalindex.MemOnlyIndex()
	require.NoError(t, err)

	dacky, err := dackbox.NewDackBox(db, nil, []byte("graph"), []byte("dirty"), []byte("valid"))
	require.NoError(t, err)

	// Initialize real datastore
	alerts := alertDatastore.NewWithDb(db, bleveIndex)

	ctrl := gomock.NewController(t)
	mockProcessDataStore := processIndicatorDatastoreMocks.NewMockDataStore(ctrl)
	mockProcessDataStore.EXPECT().RemoveProcessIndicatorsOfStaleContainers(gomock.Any(), gomock.Any()).Return((error)(nil))

	mockWhitelistDataStore := processWhitelistDatastoreMocks.NewMockDataStore(ctrl)

	mockImageDatastore := imageDatastoreMocks.NewMockDataStore(ctrl)
	mockConfigDatastore := configDatastoreMocks.NewMockDataStore(ctrl)
	mockConfigDatastore.EXPECT().GetConfig(ctx).Return(testConfig, nil)

	mockRiskDatastore := riskDatastoreMocks.NewMockDataStore(ctrl)
	mockRiskDatastore.EXPECT().SearchRawRisks(gomock.Any(), gomock.Any()).AnyTimes()

	mockFilter := filterMocks.NewMockFilter(ctrl)
	mockFilter.EXPECT().Update(gomock.Any()).AnyTimes()

	deployments, err := deploymentDatastore.NewBadger(dacky, concurrency.NewKeyFence(), db, bleveIndex, nil, mockProcessDataStore, mockWhitelistDataStore, nil, mockRiskDatastore, nil, mockFilter)
	require.NoError(t, err)

	return alerts, mockConfigDatastore, mockImageDatastore, deployments
}

func TestImagePruning(t *testing.T) {
	var cases = []struct {
		name        string
		images      []*storage.Image
		deployment  *storage.Deployment
		expectedIDs []string
	}{
		{
			name: "No pruning",
			images: []*storage.Image{
				newImageInstance("id1", 1),
				newImageInstance("id2", 1),
			},
			expectedIDs: []string{"id1", "id2"},
		},
		{
			name: "one old and one new - no deployments",
			images: []*storage.Image{
				newImageInstance("id1", 1),
				newImageInstance("id2", configDatastore.DefaultImageRetention+1),
			},
			expectedIDs: []string{"id1"},
		},
		{
			name: "one old and one new - 1 deployment with new",
			images: []*storage.Image{
				newImageInstance("id1", 1),
				newImageInstance("id2", configDatastore.DefaultImageRetention+1),
			},
			deployment:  newDeployment("id1"),
			expectedIDs: []string{"id1"},
		},
		{
			name: "one old and one new - 1 deployment with old",
			images: []*storage.Image{
				newImageInstance("id1", 1),
				newImageInstance("id2", configDatastore.DefaultImageRetention+1),
			},
			deployment:  newDeployment("id2"),
			expectedIDs: []string{"id1", "id2"},
		},
		{
			name: "two old - 1 deployment with old",
			images: []*storage.Image{
				newImageInstance("id1", configDatastore.DefaultImageRetention+1),
				newImageInstance("id2", configDatastore.DefaultImageRetention+1),
			},
			deployment:  newDeployment("id2"),
			expectedIDs: []string{"id2"},
		},
		{
			name: "two old - 1 deployment with old, but has reference to old",
			images: []*storage.Image{
				newImageInstance("id1", configDatastore.DefaultImageRetention+1),
				newImageInstance("id2", configDatastore.DefaultImageRetention+1),
			},
			deployment: &storage.Deployment{
				Id: "d1",
				Containers: []*storage.Container{
					{
						Image: &storage.ContainerImage{
							Id: "id1",
						},
						Instances: []*storage.ContainerInstance{
							{
								ImageDigest: "id2",
							},
						},
					},
				},
			},
			expectedIDs: []string{},
		},
	}

	scc := sac.OneStepSCC{
		sac.AccessModeScopeKey(storage.Access_READ_ACCESS): sac.AllowFixedScopes(
			sac.ResourceScopeKeys(resources.Alert, resources.Config, resources.Deployment, resources.Image, resources.Risk)),
		sac.AccessModeScopeKey(storage.Access_READ_WRITE_ACCESS): sac.AllowFixedScopes(
			sac.ResourceScopeKeys(resources.Alert, resources.Image, resources.Deployment, resources.Risk)),
	}

	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), scc)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Get all of the image constructs because I update the time within the store
			// So to test need to update them separately
			alerts, config, images, deployments := generateImageDataStructures(ctx, t)

			gc := newGarbageCollector(alerts, images, nil, deployments, nil, nil, nil, config).(*garbageCollectorImpl)

			// Add images and deployments into the datastores
			if c.deployment != nil {
				require.NoError(t, deployments.UpsertDeployment(ctx, c.deployment))
			}
			for _, image := range c.images {
				image.Id = types.NewDigest(image.Id).Digest()
				require.NoError(t, images.UpsertImage(ctx, image))
			}

			conf, err := config.GetConfig(ctx)
			require.NoError(t, err, "failed to get config")
			// Garbage collect all of the images
			gc.collectImages(conf.GetPrivateConfig())

			// Grab the  actual remaining images and make sure they match the images expected to be remaining
			remainingImages, err := images.SearchListImages(ctx, search.EmptyQuery())
			require.NoError(t, err)

			var ids []string
			for _, i := range remainingImages {
				ids = append(ids, i.GetId())
			}
			for i, eid := range c.expectedIDs {
				c.expectedIDs[i] = types.NewDigest(eid).Digest()
			}

			assert.ElementsMatch(t, c.expectedIDs, ids)
		})
	}
}

func TestAlertPruning(t *testing.T) {
	existsDeployment := &storage.Deployment{
		Id:        "deploymentId1",
		Name:      "test deployment",
		Namespace: "ns",
		ClusterId: "clusterid",
	}

	var cases = []struct {
		name                 string
		alerts               []*storage.Alert
		expectedIDsRemaining []string
		deployments          []*storage.Deployment
	}{
		{
			name: "No pruning",
			alerts: []*storage.Alert{
				newAlertInstance("id1", 1, storage.LifecycleStage_RUNTIME, storage.ViolationState_ACTIVE),
				newAlertInstance("id2", 1, storage.LifecycleStage_RUNTIME, storage.ViolationState_RESOLVED),
			},
			expectedIDsRemaining: []string{"id1", "id2"},
		},
		{
			name: "One old alert, and one new alert",
			alerts: []*storage.Alert{
				newAlertInstance("id1", 1, storage.LifecycleStage_RUNTIME, storage.ViolationState_ACTIVE),
				newAlertInstance("id2", testRetentionAllRuntime+1, storage.LifecycleStage_RUNTIME, storage.ViolationState_RESOLVED),
			},
			expectedIDsRemaining: []string{"id1"},
		},
		{
			name: "One old runtime alert, and one old deploy time unresolved alert",
			alerts: []*storage.Alert{
				newAlertInstance("id1", testRetentionAllRuntime+1, storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newAlertInstance("id2", testRetentionAllRuntime+1, storage.LifecycleStage_RUNTIME, storage.ViolationState_RESOLVED),
			},
			expectedIDsRemaining: []string{"id1"},
		},
		{
			name: "one old deploy time alert resolved",
			alerts: []*storage.Alert{
				newAlertInstance("id1", testRetentionResolvedDeploy+1, storage.LifecycleStage_DEPLOY, storage.ViolationState_RESOLVED),
			},
			expectedIDsRemaining: []string{},
		},
		{
			name: "two old-ish runtime alerts, one with no deployment",
			alerts: []*storage.Alert{
				newAlertInstanceWithDeployment("id1", testRetentionDeletedRuntime+1, storage.LifecycleStage_RUNTIME, storage.ViolationState_RESOLVED, nil),
				newAlertInstanceWithDeployment("id2", testRetentionDeletedRuntime+1, storage.LifecycleStage_RUNTIME, storage.ViolationState_RESOLVED, existsDeployment),
			},
			expectedIDsRemaining: []string{"id2"},
			deployments: []*storage.Deployment{
				existsDeployment,
			},
		},
		{
			name: "expired runtime alert with no deployment",
			alerts: []*storage.Alert{
				newAlertInstanceWithDeployment("id1", testRetentionDeletedRuntime+1, storage.LifecycleStage_RUNTIME, storage.ViolationState_ACTIVE, nil),
			},
			expectedIDsRemaining: []string{},
		},
	}
	scc := sac.OneStepSCC{
		sac.AccessModeScopeKey(storage.Access_READ_ACCESS): sac.AllowFixedScopes(
			sac.ResourceScopeKeys(resources.Alert, resources.Config, resources.Deployment, resources.Image)),
		sac.AccessModeScopeKey(storage.Access_READ_WRITE_ACCESS): sac.AllowFixedScopes(
			sac.ResourceScopeKeys(resources.Alert, resources.Image, resources.Deployment)),
	}

	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), scc)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Get all of the image constructs because I update the time within the store
			// So to test need to update them separately
			alerts, config, images, deployments := generateAlertDataStructures(ctx, t)

			gc := newGarbageCollector(alerts, images, nil, deployments, nil, nil, nil, config).(*garbageCollectorImpl)

			// Add alerts into the datastores
			for _, alert := range c.alerts {
				require.NoError(t, alerts.UpsertAlert(ctx, alert))
			}
			for _, deployment := range c.deployments {
				require.NoError(t, deployments.UpsertDeployment(ctx, deployment))
			}
			all, err := alerts.Search(ctx, search.NewQueryBuilder().AddStrings(search.ViolationState, storage.ViolationState_RESOLVED.String()).ProtoQuery())
			if err != nil {
				t.Error(err)
			}
			log.Infof("All query returns %d objects: %v", len(all), search.ResultsToIDs(all))

			conf, err := config.GetConfig(ctx)
			require.NoError(t, err, "failed to get config")

			// Garbage collect all of the alerts
			gc.collectAlerts(conf.GetPrivateConfig())

			q := search.NewQueryBuilder().AddStrings(search.ViolationState,
				storage.ViolationState_ACTIVE.String(), storage.ViolationState_RESOLVED.String()).ProtoQuery()
			// Grab the actual remaining alerts and make sure they match the alerts expected to be remaining
			remainingAlerts, err := alerts.SearchListAlerts(ctx, q)
			require.NoError(t, err)

			log.Infof("Remaining alerts: %v", remainingAlerts)
			var ids []string
			for _, i := range remainingAlerts {
				ids = append(ids, i.GetId())
			}

			assert.ElementsMatch(t, c.expectedIDsRemaining, ids)
		})
	}
}

func timestampNowMinus(t time.Duration) *protoTypes.Timestamp {
	return protoconv.ConvertTimeToTimestamp(time.Now().Add(-t))
}

func newListAlertWithDeployment(id string, age time.Duration, deploymentID string, stage storage.LifecycleStage, state storage.ViolationState) *storage.ListAlert {
	return &storage.ListAlert{
		Id: id,
		Deployment: &storage.ListAlertDeployment{
			Id: deploymentID,
		},
		State:          state,
		LifecycleStage: stage,
		Time:           timestampNowMinus(age),
	}
}

func newIndicatorWithDeployment(id string, age time.Duration, deploymentID string) *storage.ProcessIndicator {
	return &storage.ProcessIndicator{
		Id:                id,
		DeploymentId:      deploymentID,
		DeploymentStateTs: 0,
		ContainerName:     "",
		PodId:             "",
		Signal: &storage.ProcessSignal{
			Time: timestampNowMinus(age),
		},
	}
}

func TestRemoveOrphanedProcesses(t *testing.T) {
	cases := []struct {
		name              string
		initialProcesses  []*storage.ProcessIndicator
		deployments       set.FrozenStringSet
		expectedDeletions []string
	}{
		{
			name: "no deployments - remove all old indicators",
			initialProcesses: []*storage.ProcessIndicator{
				newIndicatorWithDeployment("pi1", 1*time.Hour, "dep1"),
				newIndicatorWithDeployment("pi2", 1*time.Hour, "dep2"),
				newIndicatorWithDeployment("pi3", 1*time.Hour, "dep3"),
			},
			deployments:       set.NewFrozenStringSet(),
			expectedDeletions: []string{"pi1", "pi2", "pi3"},
		},
		{
			name: "no deployments - remove no new orphaned indicators",
			initialProcesses: []*storage.ProcessIndicator{
				newIndicatorWithDeployment("pi1", 20*time.Minute, "dep1"),
				newIndicatorWithDeployment("pi2", 20*time.Minute, "dep2"),
				newIndicatorWithDeployment("pi3", 20*time.Minute, "dep3"),
			},
			deployments:       set.NewFrozenStringSet(),
			expectedDeletions: []string{},
		},
		{
			name: "all deployments - remove no indicators",
			initialProcesses: []*storage.ProcessIndicator{
				newIndicatorWithDeployment("pi1", 1*time.Hour, "dep1"),
				newIndicatorWithDeployment("pi2", 1*time.Hour, "dep2"),
				newIndicatorWithDeployment("pi3", 1*time.Hour, "dep3"),
			},
			deployments:       set.NewFrozenStringSet("dep1", "dep2", "dep3"),
			expectedDeletions: []string{},
		},
		{
			name: "some deployments - remove some indicators",
			initialProcesses: []*storage.ProcessIndicator{
				newIndicatorWithDeployment("pi1", 1*time.Hour, "dep1"),
				newIndicatorWithDeployment("pi2", 20*time.Minute, "dep2"),
				newIndicatorWithDeployment("pi3", 1*time.Hour, "dep3"),
			},
			deployments:       set.NewFrozenStringSet("dep3"),
			expectedDeletions: []string{"pi1"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			processes := processIndicatorDatastoreMocks.NewMockDataStore(ctrl)
			gci := &garbageCollectorImpl{
				processes: processes,
			}
			processes.EXPECT().WalkAll(pruningCtx, gomock.Any()).DoAndReturn(
				func(ctx context.Context, fn func(pi *storage.ProcessIndicator) error) error {
					for _, a := range c.initialProcesses {
						assert.NoError(t, fn(a))
					}
					return nil
				})
			processes.EXPECT().RemoveProcessIndicators(pruningCtx, testutils.AssertionMatcher(assert.ElementsMatch, c.expectedDeletions))
			gci.removeOrphanedProcesses(c.deployments)
		})
	}
}

func TestMarkOrphanedAlerts(t *testing.T) {
	cases := []struct {
		name              string
		initialAlerts     []*storage.ListAlert
		deployments       set.FrozenStringSet
		expectedDeletions []string
	}{
		{
			name: "no deployments - remove all old alerts",
			initialAlerts: []*storage.ListAlert{
				newListAlertWithDeployment("alert1", 1*time.Hour, "dep1", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert2", 1*time.Hour, "dep2", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
			},
			deployments:       set.NewFrozenStringSet(),
			expectedDeletions: []string{"alert1", "alert2"},
		},
		{
			name: "no deployments - remove no new orphaned alerts",
			initialAlerts: []*storage.ListAlert{
				newListAlertWithDeployment("alert1", 20*time.Minute, "dep1", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert2", 20*time.Minute, "dep2", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert3", 20*time.Minute, "dep3", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
			},
			deployments:       set.NewFrozenStringSet(),
			expectedDeletions: []string{},
		},
		{
			name: "all deployments - remove no alerts",
			initialAlerts: []*storage.ListAlert{
				newListAlertWithDeployment("alert1", 1*time.Hour, "dep1", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert2", 1*time.Hour, "dep2", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert3", 1*time.Hour, "dep3", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
			},
			deployments:       set.NewFrozenStringSet("dep1", "dep2", "dep3"),
			expectedDeletions: []string{},
		},
		{
			name: "some deployments - remove some alerts",
			initialAlerts: []*storage.ListAlert{
				newListAlertWithDeployment("alert1", 1*time.Hour, "dep1", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert2", 20*time.Minute, "dep2", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert3", 1*time.Hour, "dep3", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
			},
			deployments:       set.NewFrozenStringSet("dep3"),
			expectedDeletions: []string{"alert1"},
		},
		{
			name: "some deployments - remove some alerts due to stages",
			initialAlerts: []*storage.ListAlert{
				newListAlertWithDeployment("alert1", 1*time.Hour, "dep1", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert2", 1*time.Hour, "dep2", storage.LifecycleStage_BUILD, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert3", 1*time.Hour, "dep3", storage.LifecycleStage_RUNTIME, storage.ViolationState_ACTIVE),
			},
			deployments:       set.NewFrozenStringSet("dep3"),
			expectedDeletions: []string{"alert1"},
		},
		{
			name: "some deployments - remove some alerts due to state",
			initialAlerts: []*storage.ListAlert{
				newListAlertWithDeployment("alert1", 1*time.Hour, "dep1", storage.LifecycleStage_DEPLOY, storage.ViolationState_ACTIVE),
				newListAlertWithDeployment("alert2", 1*time.Hour, "dep2", storage.LifecycleStage_DEPLOY, storage.ViolationState_RESOLVED),
				newListAlertWithDeployment("alert3", 1*time.Hour, "dep3", storage.LifecycleStage_DEPLOY, storage.ViolationState_SNOOZED),
			},
			deployments:       set.NewFrozenStringSet("dep3"),
			expectedDeletions: []string{"alert1"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			alerts := alertDatastoreMocks.NewMockDataStore(ctrl)
			gci := &garbageCollectorImpl{
				alerts: alerts,
			}
			alerts.EXPECT().WalkAll(pruningCtx, gomock.Any()).DoAndReturn(
				func(ctx context.Context, fn func(la *storage.ListAlert) error) error {
					for _, a := range c.initialAlerts {
						assert.NoError(t, fn(a))
					}
					return nil
				})
			for _, a := range c.expectedDeletions {
				alerts.EXPECT().MarkAlertStale(pruningCtx, a)
			}
			gci.markOrphanedAlertsAsResolved(c.deployments)
		})
	}
}

func TestRemoveOrphanedNetworkFlows(t *testing.T) {
	cases := []struct {
		name             string
		flows            []*storage.NetworkFlow
		deployments      set.FrozenStringSet
		expectedDeletion bool
	}{
		{
			name: "no deployments - remove all flows",
			flows: []*storage.NetworkFlow{
				{
					LastSeenTimestamp: timestampNowMinus(1 * time.Hour),
					Props: &storage.NetworkFlowProperties{
						SrcEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep1",
						},
						DstEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep2",
						},
					},
				},
			},
			deployments:      set.NewFrozenStringSet(),
			expectedDeletion: true,
		},
		{
			name: "no deployments - but no flows with deployments",
			flows: []*storage.NetworkFlow{
				{
					LastSeenTimestamp: timestampNowMinus(1 * time.Hour),
					Props: &storage.NetworkFlowProperties{
						SrcEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_INTERNET,
							Id:   "i1",
						},
						DstEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_INTERNET,
							Id:   "i2",
						},
					},
				},
			},
			deployments:      set.NewFrozenStringSet(),
			expectedDeletion: false,
		},
		{
			name: "no deployments - but flows too recent",
			flows: []*storage.NetworkFlow{
				{
					LastSeenTimestamp: timestampNowMinus(20 * time.Minute),
					Props: &storage.NetworkFlowProperties{
						SrcEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep1",
						},
						DstEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep2",
						},
					},
				},
			},
			deployments:      set.NewFrozenStringSet(),
			expectedDeletion: false,
		},
		{
			name: "some deployments with matching flows",
			flows: []*storage.NetworkFlow{
				{
					LastSeenTimestamp: timestampNowMinus(1 * time.Hour),
					Props: &storage.NetworkFlowProperties{
						SrcEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep1",
						},
						DstEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep2",
						},
					},
				},
			},
			deployments:      set.NewFrozenStringSet("dep1", "dep2"),
			expectedDeletion: false,
		},
		{
			name: "some deployments with matching src",
			flows: []*storage.NetworkFlow{
				{
					LastSeenTimestamp: timestampNowMinus(1 * time.Hour),
					Props: &storage.NetworkFlowProperties{
						SrcEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep1",
						},
						DstEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep2",
						},
					},
				},
			},
			deployments:      set.NewFrozenStringSet("dep1"),
			expectedDeletion: true,
		},
		{
			name: "some deployments with matching dst",
			flows: []*storage.NetworkFlow{
				{
					LastSeenTimestamp: timestampNowMinus(1 * time.Hour),
					Props: &storage.NetworkFlowProperties{
						SrcEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep1",
						},
						DstEntity: &storage.NetworkEntityInfo{
							Type: storage.NetworkEntityInfo_DEPLOYMENT,
							Id:   "dep2",
						},
					},
				},
			},
			deployments:      set.NewFrozenStringSet("dep2"),
			expectedDeletion: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			clusters := clusterDatastoreMocks.NewMockDataStore(ctrl)
			clusterFlows := networkFlowDatastoreMocks.NewMockClusterDataStore(ctrl)
			flows := networkFlowDatastoreMocks.NewMockFlowDataStore(ctrl)

			clusters.EXPECT().GetClusters(pruningCtx).Return([]*storage.Cluster{{Id: "cluster"}}, nil)
			clusterFlows.EXPECT().GetFlowStore(pruningCtx, "cluster").Return(flows)

			flows.EXPECT().RemoveMatchingFlows(pruningCtx, gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, keyFn func(props *storage.NetworkFlowProperties) bool, valueFn func(flow *storage.NetworkFlow) bool) error {
					var deleted bool
					for _, f := range c.flows {
						if !keyFn(f.Props) || !valueFn(f) {
							continue
						}
						deleted = true
					}
					assert.Equal(t, c.expectedDeletion, deleted)
					return nil
				})

			gci := &garbageCollectorImpl{
				clusters:     clusters,
				networkflows: clusterFlows,
			}
			gci.removeOrphanedNetworkFlows(c.deployments)
		})
	}
}
