package service

import (
	"context"
	"testing"

	managerMocks "github.com/stackrox/rox/central/complianceoperator/v2/compliancemanager/mocks"
	scanConfigMocks "github.com/stackrox/rox/central/complianceoperator/v2/scanconfigurations/datastore/mocks"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/grpc/testutils"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

const (
	mockClusterName = "mock-cluster"
)

func TestAuthz(t *testing.T) {
	testutils.AssertAuthzWorks(t, &serviceImpl{})
}

func TestComplianceScanConfigService(t *testing.T) {
	suite.Run(t, new(ComplianceResultsServiceTestSuite))
}

type ComplianceResultsServiceTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller

	ctx                 context.Context
	manager             *managerMocks.MockManager
	scanConfigDatastore *scanConfigMocks.MockDataStore
	service             Service
}

func (s *ComplianceResultsServiceTestSuite) SetupSuite() {
	s.T().Setenv(features.ComplianceEnhancements.EnvVar(), "true")
	if !features.ComplianceEnhancements.Enabled() {
		s.T().Skip("Skip test when compliance enhancements are disabled")
		s.T().SkipNow()
	}

	s.ctx = sac.WithAllAccess(context.Background())
}

func (s *ComplianceResultsServiceTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.manager = managerMocks.NewMockManager(s.mockCtrl)
	s.scanConfigDatastore = scanConfigMocks.NewMockDataStore(s.mockCtrl)

	//s.service = New(s.scanConfigDatastore, s.manager)
}

func (s *ComplianceResultsServiceTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}
