package service

import (
	"context"

	"github.com/stackrox/rox/central/externalbackups/datastore"
	"github.com/stackrox/rox/central/externalbackups/manager"
	integrationHealthDS "github.com/stackrox/rox/central/integrationhealth/datastore"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once

	as Service
)

func initialize() {
	as = New(datastore.Singleton(), integrationHealthDS.Singleton(), initializeManager())
}

func initializeManager() manager.Manager {
	ctx := sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.BackupPlugins)))

	backups, err := datastore.Singleton().ListBackups(ctx)
	if err != nil {
		panic(err)
	}
	mgr := manager.New()
	for _, b := range backups {
		if err := mgr.Upsert(ctx, b); err != nil {
			log.Errorf("error initializing backup: %v", err)
		}
	}
	return mgr
}

// Singleton provides the instance of the Service interface to register.
func Singleton() Service {
	once.Do(initialize)
	return as
}
