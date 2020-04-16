package filtered

import (
	"context"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/sac"
)

type globalFilterImpl struct {
	resourceHelper sac.ForResourceHelper
	access         storage.Access
}

func (f *globalFilterImpl) Apply(ctx context.Context, from ...string) ([]string, error) {
	if ok, err := f.resourceHelper.AccessAllowed(ctx, f.access); err != nil || !ok {
		return nil, err
	}
	return from, nil
}
