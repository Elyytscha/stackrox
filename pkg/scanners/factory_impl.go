package scanners

import (
	"fmt"

	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/scanners/types"
)

type factoryImpl struct {
	creators map[string]Creator
}

func (e *factoryImpl) CreateScanner(source *v1.ImageIntegration) (types.ImageScanner, error) {
	creator, exists := e.creators[source.GetType()]
	if !exists {
		return nil, fmt.Errorf("Scanner with type '%s' does not exist", source.GetType())
	}
	return creator(source)
}
