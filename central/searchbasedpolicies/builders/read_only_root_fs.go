package builders

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/searchbasedpolicies"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
)

// ReadOnlyRootFSQueryBuilder checks for read only root fs in containers.
type ReadOnlyRootFSQueryBuilder struct {
}

// Query implements the PolicyQueryBuilder interface.
func (p ReadOnlyRootFSQueryBuilder) Query(fields *storage.PolicyFields, optionsMap map[search.FieldLabel]*v1.SearchField) (q *v1.Query, v searchbasedpolicies.ViolationPrinter, err error) {
	// We don't match on readonlyrootfs = true, because that seems pointless.
	if fields.GetSetReadOnlyRootFs() == nil {
		return
	}
	if fields.GetReadOnlyRootFs() {
		return nil, nil, errors.New("Policy can only check for non read-only root filesystems")
	}
	searchField, err := getSearchField(search.ReadOnlyRootFilesystem, optionsMap)
	if err != nil {
		err = errors.Wrapf(err, "%s", p.Name())
		return
	}

	q = search.NewQueryBuilder().AddBoolsHighlighted(search.ReadOnlyRootFilesystem, false).ProtoQuery()
	v = violationPrinterForField(searchField.GetFieldPath(), func(match string) string {
		if match != "false" {
			return ""
		}
		return "Container using read-write root filesystem found"
	})
	return
}

// Name implements the PolicyQueryBuilder interface.
func (ReadOnlyRootFSQueryBuilder) Name() string {
	return "Query builder for read-write filesystem containers"
}
