package manager

import (
	"context"
	"testing"
	"time"

	"github.com/etcd-io/bbolt"
	"github.com/stackrox/rox/central/telemetry/manager/internal/store"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

var (
	withAccessCtx = sac.WithAllAccess(context.Background())
)

func TestManager(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(managerSuite))
}

type managerSuite struct {
	suite.Suite

	envIsolator *testutils.EnvIsolator

	runCtx       context.Context
	runCtxCancel context.CancelFunc

	db    *bbolt.DB
	store store.Store
}

func (s *managerSuite) createManager(ctx context.Context) *manager {
	return newManager(ctx, s.store, nil, nil)
}

func (s *managerSuite) SetupTest() {
	s.envIsolator = testutils.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(env.OfflineModeEnv.EnvVar(), "true")

	s.runCtx, s.runCtxCancel = context.WithCancel(sac.WithAllAccess(context.Background()))
	s.db = testutils.DBForT(s.T())
	dbStore, err := store.New(s.db)
	s.Require().NoError(err)
	s.store = dbStore
}

func (s *managerSuite) TearDownTest() {
	s.runCtxCancel()
	_ = s.db.Close()
	s.envIsolator.RestoreAll()
}

func (s *managerSuite) TestInitConfig_Unset() {
	s.envIsolator.Unsetenv(env.InitialTelemetryEnabledEnv.EnvVar())

	mgr := newManager(context.Background(), s.store, nil, nil)

	cfg, err := mgr.GetTelemetryConfig(withAccessCtx)
	s.NoError(err)
	s.False(cfg.GetEnabled())
}

func (s *managerSuite) TestInitConfig_False() {
	s.envIsolator.Setenv(env.InitialTelemetryEnabledEnv.EnvVar(), "false")

	mgr := s.createManager(context.Background())

	cfg, err := mgr.GetTelemetryConfig(withAccessCtx)
	s.NoError(err)
	s.False(cfg.GetEnabled())
}

func (s *managerSuite) TestInitConfig_True() {
	s.envIsolator.Setenv(env.InitialTelemetryEnabledEnv.EnvVar(), "true")

	mgr := s.createManager(context.Background())

	cfg, err := mgr.GetTelemetryConfig(withAccessCtx)
	s.NoError(err)
	s.True(cfg.GetEnabled())
}

func (s *managerSuite) TestReadConfig_WithoutAccess() {
	mgr := s.createManager(context.Background())

	_, err := mgr.GetTelemetryConfig(context.Background())
	s.Error(err)
}

func (s *managerSuite) TestUpdateConfig_WithAccess() {
	mgr := s.createManager(context.Background())

	cfg, err := mgr.GetTelemetryConfig(withAccessCtx)
	s.NoError(err)

	cfg.Enabled = !cfg.GetEnabled()

	err = mgr.UpdateTelemetryConfig(withAccessCtx, cfg)
	s.NoError(err)

	updatedCfg, err := mgr.GetTelemetryConfig(withAccessCtx)
	s.NoError(err)

	s.Equal(cfg.GetEnabled(), updatedCfg.GetEnabled())
}

func (s *managerSuite) TestUpdateConfig_WithoutAccess() {
	mgr := s.createManager(context.Background())

	cfg := &storage.TelemetryConfiguration{}

	err := mgr.UpdateTelemetryConfig(context.Background(), cfg)
	s.Error(err)
}

func (s *managerSuite) TestUpdateConfig_AfterCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	mgr := s.createManager(ctx)
	cancel()
	time.Sleep(100 * time.Millisecond)

	cfg := &storage.TelemetryConfiguration{}

	err := mgr.UpdateTelemetryConfig(withAccessCtx, cfg)
	s.Error(err)
}
