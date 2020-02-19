package searchbasedpolicies

import (
	"context"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
)

// A PolicyQueryBuilder can build a query for (a part of) a policy.
type PolicyQueryBuilder interface {
	// Query returns a query matching the field(s) that this particular query builder cares about.
	// It may return all nil values -- this is allowed, and means that the given policy fields
	// don't have a value for any field that this particular query builder cares about.
	Query(*storage.PolicyFields, map[search.FieldLabel]*v1.SearchField) (*v1.Query, ViolationPrinter, error)
	// Name returns a human-friendly name for this fieldQueryBuilder.
	Name() string
}

// Violations represents a list of violation sub-objects.
type Violations struct {
	ProcessViolation *storage.Alert_ProcessViolation
	AlertViolations  []*storage.Alert_Violation
}

// A ViolationPrinter knows how to print violation messages from a search result.
type ViolationPrinter func(context.Context, search.Result) Violations

// A ProcessIndicatorGetter knows how to retrieve process indicators given its id.
type ProcessIndicatorGetter interface {
	GetProcessIndicator(ctx context.Context, id string) (*storage.ProcessIndicator, bool, error)
}

// MatcherOpts defines the options for the MatchOne function
type MatcherOpts struct {
	IDFilters map[v1.SearchCategory][]string
}

// Matcher matches objects against a policy.
//go:generate mockgen-wrapper
type Matcher interface {
	// Match matches the policy against all objects, returning a map from object ID to violations.
	Match(ctx context.Context, searcher search.Searcher) (map[string]Violations, error)
	// MatchOne matches the policy against the passed deployment and images
	MatchOne(ctx context.Context, deployment *storage.Deployment, images []*storage.Image, pi *storage.ProcessIndicator) (Violations, error)
	// MatchMany mathes the policy against just the objects with the given ids.
	MatchMany(ctx context.Context, searcher search.Searcher, ids ...string) (map[string]Violations, error)
}
