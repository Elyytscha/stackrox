package service

import (
	"context"

	"github.com/stackrox/rox/central/license/manager"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/grpc"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

// Service is the GRPC service interface that provides the entry point for processing deployment events.
type Service interface {
	grpc.APIService

	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)

	GetMetadata(context.Context, *v1.Empty) (*v1.Metadata, error)
}

// New returns a new instance of service.
func New(restarting *concurrency.Flag, licenseMgr manager.LicenseManager) Service {
	return &serviceImpl{
		restarting: restarting,
		licenseMgr: licenseMgr,
	}
}
