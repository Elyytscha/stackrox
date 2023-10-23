//go:build sql_integration

package datastore

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/types"
	checkResultsStorage "github.com/stackrox/rox/central/complianceoperator/v2/checkresults/store/postgres"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/fixtures/fixtureconsts"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sac/resources"
	"github.com/stackrox/rox/pkg/search"
	pgSearch "github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/search/postgres/aggregatefunc"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	scanConfig1 = uuid.NewV4().String()
	scanConfig2 = uuid.NewV4().String()
)

type shrewsTestCheckName struct {
	CheckName string `db:"compliance_check_name"`
}

type resourceCountByResultByCluster struct {
	PassCount          int    `db:"pass_count"`
	FailCount          int    `db:"fail_count"`
	ErrorCount         int    `db:"error_count"`
	InfoCount          int    `db:"info_count"`
	ManualCount        int    `db:"manual_count"`
	NotApplicableCount int    `db:"not_applicable_count"`
	InconsistentCount  int    `db:"inconsistent_count"`
	ClusterID          string `db:"cluster_id"`
	ScanConfigID       string `db:"compliance_scan_config_id"`
}

func TestComplianceCheckResultDataStore(t *testing.T) {
	suite.Run(t, new(complianceCheckResultDataStoreTestSuite))
}

type complianceCheckResultDataStoreTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller

	hasReadCtx  context.Context
	hasWriteCtx context.Context
	noAccessCtx context.Context

	dataStore DataStore
	storage   checkResultsStorage.Store
	db        *pgtest.TestPostgres
}

func (s *complianceCheckResultDataStoreTestSuite) SetupSuite() {
	s.T().Setenv(features.ComplianceEnhancements.EnvVar(), "true")
	if !features.ComplianceEnhancements.Enabled() {
		s.T().Skip("Skip tests when ComplianceEnhancements disabled")
		s.T().SkipNow()
	}
}

func (s *complianceCheckResultDataStoreTestSuite) SetupTest() {
	s.hasReadCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
			sac.ResourceScopeKeys(resources.ComplianceOperator)))
	s.hasWriteCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.ComplianceOperator)))
	s.noAccessCtx = sac.WithGlobalAccessScopeChecker(context.Background(), sac.DenyAllAccessScopeChecker())

	s.mockCtrl = gomock.NewController(s.T())

	s.db = pgtest.ForT(s.T())

	s.storage = checkResultsStorage.New(s.db)
	s.dataStore = New(s.storage, s.db)
}

func (s *complianceCheckResultDataStoreTestSuite) TearDownTest() {
	//s.db.Teardown(s.T())
}

func (s *complianceCheckResultDataStoreTestSuite) TestUpsertResult() {
	// make sure we have nothing
	checkResultIDs, err := s.storage.GetIDs(s.hasReadCtx)
	s.Require().NoError(err)
	s.Require().Empty(checkResultIDs)

	rec1 := getTestRec(fixtureconsts.Cluster1)
	rec2 := getTestRec(fixtureconsts.Cluster2)
	ids := []string{rec1.GetId(), rec2.GetId()}

	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec1))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec2))

	count, err := s.storage.Count(s.hasReadCtx)
	s.Require().NoError(err)
	s.Require().Equal(len(ids), count)

	// upsert with read context
	s.Require().Error(s.dataStore.UpsertResult(s.hasReadCtx, rec2))

	retrieveRec1, found, err := s.storage.Get(s.hasReadCtx, rec1.GetId())
	s.Require().NoError(err)
	s.Require().True(found)
	s.Require().Equal(rec1, retrieveRec1)
}

func (s *complianceCheckResultDataStoreTestSuite) TestDeleteResult() {
	// make sure we have nothing
	checkResultIDs, err := s.storage.GetIDs(s.hasReadCtx)
	s.Require().NoError(err)
	s.Require().Empty(checkResultIDs)

	rec1 := getTestRec(fixtureconsts.Cluster1)
	rec2 := getTestRec(fixtureconsts.Cluster2)
	ids := []string{rec1.GetId(), rec2.GetId()}

	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec1))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec2))

	count, err := s.storage.Count(s.hasReadCtx)
	s.Require().NoError(err)
	s.Require().Equal(len(ids), count)

	// Try to delete with wrong context
	s.Require().Error(s.dataStore.DeleteResult(s.hasReadCtx, rec1.GetId()))

	// Now delete rec1
	s.Require().NoError(s.dataStore.DeleteResult(s.hasWriteCtx, rec1.GetId()))
	retrieveRec1, found, err := s.storage.Get(s.hasReadCtx, rec1.GetId())
	s.Require().NoError(err)
	s.Require().False(found)
	s.Require().Nil(retrieveRec1)
}

func (s *complianceCheckResultDataStoreTestSuite) TestSearchCheckResults() {
	// make sure we have nothing
	checkResultIDs, err := s.storage.GetIDs(s.hasReadCtx)
	s.Require().NoError(err)
	s.Require().Empty(checkResultIDs)

	_, err = s.db.DB.Exec(context.Background(), "insert into compliance_operator_scan_configuration_v2 (id, scanname) values ($1, $2)", scanConfig1, "scan 1")
	s.Require().NoError(err)
	_, err = s.db.DB.Exec(context.Background(), "insert into compliance_operator_scan_configuration_v2 (id, scanname) values ($1, $2)", scanConfig2, "scan 2")
	s.Require().NoError(err)

	rec1 := getTestRec(fixtureconsts.Cluster1)
	rec2 := getTestRec(fixtureconsts.Cluster2)
	rec3 := getTestRec(fixtureconsts.Cluster2)
	rec4 := getTestRec(fixtureconsts.Cluster2)
	rec5 := getTestRec(fixtureconsts.Cluster3)
	rec6 := getTestRec2(fixtureconsts.Cluster3)
	ids := []string{rec1.GetId(), rec2.GetId(), rec3.GetId(), rec4.GetId()}

	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec1))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec2))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec3))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec4))

	count, err := s.storage.Count(s.hasReadCtx)
	s.Require().NoError(err)
	s.Require().Equal(len(ids), count)

	// Search results of fake cluster, should return 0
	searchResults, err := s.dataStore.SearchCheckResults(s.hasReadCtx, search.NewQueryBuilder().
		AddExactMatches(search.ClusterID, fixtureconsts.ClusterFake2).ProtoQuery())
	s.Require().NoError(err)
	s.Require().Equal(0, len(searchResults))

	// Search results of cluster 2 should return 3 records.
	searchResults, err = s.dataStore.SearchCheckResults(s.hasReadCtx, search.NewQueryBuilder().
		AddExactMatches(search.ClusterID, fixtureconsts.Cluster2).ProtoQuery())
	s.Require().NoError(err)
	s.Require().Equal(3, len(searchResults))

	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec5))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec6))

	log.Infof("SHREWS=========================================")
	//selects := []*v1.QuerySelect{
	//	search.NewQuerySelect(search.ComplianceOperatorCheckName).Distinct().Filter("compliance_check_name",
	//		search.NewQueryBuilder().
	//			AddExactMatches(search.ClusterID, fixtureconsts.Cluster2).
	//			AddExactMatches(search.ClusterID, fixtureconsts.Cluster3).ProtoQuery(),
	//	).Proto(),
	//}
	//query := &v1.Query{
	//	Selects: selects,
	//}
	// Get distinct check name
	q := search.NewQueryBuilder().
		AddSelectFields(
			search.NewQuerySelect(search.ComplianceOperatorCheckName).Distinct(),
		).
		AddExactMatches(search.ClusterID, fixtureconsts.Cluster2).
		AddExactMatches(search.ClusterID, fixtureconsts.Cluster1).ProtoQuery()

	log.Infof("SHREWS -- query %v", q)
	searchResults, err = s.dataStore.SearchCheckResults(s.hasReadCtx, q)
	s.Require().NoError(err)

	for i, result := range searchResults {
		log.Infof("SHREWS -- check shape %d", i)
		log.Infof("%v", result)
	}

	//var results []*shrewsTestCheckName
	results, err := pgSearch.RunSelectRequestForSchema[shrewsTestCheckName](s.hasReadCtx, s.db, schema.ComplianceOperatorCheckResultV2Schema, q)
	log.Infof("err %v", err)
	log.Infof("%v", results)
	for _, result := range results {
		log.Infof("result %v", result)
	}

	query := search.NewQueryBuilder().
		AddSelectFields(
			search.NewQuerySelect(search.ClusterID),
			search.NewQuerySelect(search.ComplianceOperatorScanConfig),
		).
		AddExactMatches(search.ClusterID, fixtureconsts.Cluster2).
		AddExactMatches(search.ClusterID, fixtureconsts.Cluster3).
		AddExactMatches(search.ClusterID, fixtureconsts.Cluster1).
		AddGroupBy(search.ClusterID, search.ComplianceOperatorScanConfig).ProtoQuery()
	countQuery := withCountByResultSelectQuery(query, search.ClusterID)
	countResults, err := pgSearch.RunSelectRequestForSchema[resourceCountByResultByCluster](s.hasReadCtx, s.db, schema.ComplianceOperatorCheckResultV2Schema, countQuery)
	log.Infof("err %v", err)
	log.Infof("%v", results)
	for _, result := range countResults {
		log.Infof("result %v", result)
	}

	log.Infof("END SHREWS=========================================")
}

func (s *complianceCheckResultDataStoreTestSuite) TestCheckResultStats() {
	_, err := s.db.DB.Exec(context.Background(), "insert into compliance_operator_scan_configuration_v2 (id, scanname) values ($1, $2)", scanConfig1, "scan 1")
	s.Require().NoError(err)
	_, err = s.db.DB.Exec(context.Background(), "insert into compliance_operator_scan_configuration_v2 (id, scanname) values ($1, $2)", scanConfig2, "scan 2")
	s.Require().NoError(err)

	_, err = s.db.DB.Exec(context.Background(), "insert into clusters (id, name) values ($1, $2)", fixtureconsts.Cluster1, "cluster1")
	s.Require().NoError(err)
	_, err = s.db.DB.Exec(context.Background(), "insert into clusters (id, name) values ($1, $2)", fixtureconsts.Cluster2, "cluster2")
	s.Require().NoError(err)
	_, err = s.db.DB.Exec(context.Background(), "insert into clusters (id, name) values ($1, $2)", fixtureconsts.Cluster3, "cluster3")
	s.Require().NoError(err)

	rec1 := getTestRec(fixtureconsts.Cluster1)
	rec2 := getTestRec(fixtureconsts.Cluster2)
	rec3 := getTestRec(fixtureconsts.Cluster2)
	rec4 := getTestRec(fixtureconsts.Cluster2)
	rec5 := getTestRec(fixtureconsts.Cluster3)
	rec6 := getTestRec2(fixtureconsts.Cluster3)

	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec1))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec2))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec3))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec4))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec5))
	s.Require().NoError(s.dataStore.UpsertResult(s.hasWriteCtx, rec6))

	// Counts by Scan Config by Cluster
	query := search.NewQueryBuilder().
		AddExactMatches(search.ClusterID, fixtureconsts.Cluster2).
		AddExactMatches(search.ClusterID, fixtureconsts.Cluster3).ProtoQuery()

	results, err := s.dataStore.CheckResultStats(s.hasReadCtx, query)
	log.Infof("err %v", err)
	s.Require().NoError(err)
	log.Infof("%v", results)
	for _, result := range results {
		log.Infof("result %v", result)
	}

	log.Infof("END SHREWS=========================================")
}

func getTestRec(clusterID string) *storage.ComplianceOperatorCheckResultV2 {
	return &storage.ComplianceOperatorCheckResultV2{
		Id:           uuid.NewV4().String(),
		CheckId:      uuid.NewV4().String(),
		CheckName:    "test-check",
		ClusterId:    clusterID,
		Status:       storage.ComplianceOperatorCheckResultV2_INCONSISTENT,
		Severity:     storage.RuleSeverity_HIGH_RULE_SEVERITY,
		Description:  "this is a test",
		Instructions: "this is a test",
		Labels:       nil,
		Annotations:  nil,
		CreatedTime:  types.TimestampNow(),
		ScanId:       uuid.NewV4().String(),
		ScanConfigId: scanConfig1,
	}
}

func getTestRec2(clusterID string) *storage.ComplianceOperatorCheckResultV2 {
	return &storage.ComplianceOperatorCheckResultV2{
		Id:           uuid.NewV4().String(),
		CheckId:      uuid.NewV4().String(),
		CheckName:    "test-check-2",
		ClusterId:    clusterID,
		Status:       storage.ComplianceOperatorCheckResultV2_INFO,
		Severity:     storage.RuleSeverity_INFO_RULE_SEVERITY,
		Description:  "this is a test",
		Instructions: "this is a test",
		Labels:       nil,
		Annotations:  nil,
		CreatedTime:  types.TimestampNow(),
		ScanId:       uuid.NewV4().String(),
		ScanConfigId: scanConfig2,
	}
}

func withCountByResultSelectQuery(q *v1.Query, countOn search.FieldLabel) *v1.Query {
	cloned := q.Clone()
	cloned.Selects = append(cloned.Selects,
		search.NewQuerySelect(countOn).
			AggrFunc(aggregatefunc.Count).
			Filter("pass_count",
				search.NewQueryBuilder().
					AddExactMatches(
						search.ComplianceOperatorCheckStatus,
						storage.ComplianceOperatorCheckResultV2_PASS.String(),
					).ProtoQuery(),
			).Proto(),
		search.NewQuerySelect(countOn).
			AggrFunc(aggregatefunc.Count).
			Filter("fail_count",
				search.NewQueryBuilder().
					AddExactMatches(
						search.ComplianceOperatorCheckStatus,
						storage.ComplianceOperatorCheckResultV2_FAIL.String(),
					).ProtoQuery(),
			).Proto(),
		search.NewQuerySelect(countOn).
			AggrFunc(aggregatefunc.Count).
			Filter("error_count",
				search.NewQueryBuilder().
					AddExactMatches(
						search.ComplianceOperatorCheckStatus,
						storage.ComplianceOperatorCheckResultV2_ERROR.String(),
					).ProtoQuery(),
			).Proto(),
		search.NewQuerySelect(countOn).
			AggrFunc(aggregatefunc.Count).
			Filter("info_count",
				search.NewQueryBuilder().
					AddExactMatches(
						search.ComplianceOperatorCheckStatus,
						storage.ComplianceOperatorCheckResultV2_INFO.String(),
					).ProtoQuery(),
			).Proto(),
		search.NewQuerySelect(countOn).
			AggrFunc(aggregatefunc.Count).
			Filter("manual_count",
				search.NewQueryBuilder().
					AddExactMatches(
						search.ComplianceOperatorCheckStatus,
						storage.ComplianceOperatorCheckResultV2_MANUAL.String(),
					).ProtoQuery(),
			).Proto(),
		search.NewQuerySelect(countOn).
			AggrFunc(aggregatefunc.Count).
			Filter("not_applicable_count",
				search.NewQueryBuilder().
					AddExactMatches(
						search.ComplianceOperatorCheckStatus,
						storage.ComplianceOperatorCheckResultV2_NOT_APPLICABLE.String(),
					).ProtoQuery(),
			).Proto(),
		search.NewQuerySelect(countOn).
			AggrFunc(aggregatefunc.Count).
			Filter("inconsistent_count",
				search.NewQueryBuilder().
					AddExactMatches(
						search.ComplianceOperatorCheckStatus,
						storage.ComplianceOperatorCheckResultV2_INCONSISTENT.String(),
					).ProtoQuery(),
			).Proto(),
	)
	return cloned
}
