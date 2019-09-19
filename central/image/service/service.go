package service

import (
	"context"

	"github.com/stackrox/rox/central/image/datastore"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/expiringcache"
	"github.com/stackrox/rox/pkg/grpc"
	"github.com/stackrox/rox/pkg/images/enricher"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

// Service provides the interface to the microservice that serves alert data.
type Service interface {
	grpc.APIService

	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)

	v1.ImageServiceServer
}

// New returns a new Service instance using the given DataStore.
func New(datastore datastore.DataStore, enricher enricher.ImageEnricher, metadataCache, scanCache expiringcache.Cache) Service {
	return &serviceImpl{
		datastore:     datastore,
		enricher:      enricher,
		metadataCache: metadataCache,
		scanCache:     scanCache,
	}
}
