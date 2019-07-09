package sac

import (
	"context"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
)

// SearchHelperFlavor determines how the search helper extracts scope information from search results.
type SearchHelperFlavor int32

const (
	// ClusterIDField instructs the search helper to only look at the `Cluster ID` field. Use this for
	// objects which are tied to a single cluster, and are not namespace-scoped.
	ClusterIDField SearchHelperFlavor = iota
	// ClusterIDAndNamespaceFields instructs the search helper to look at the `Cluster ID`
	// and `Namespace` fields. Use this for objects which are tied to a single cluster and namespace.
	ClusterIDAndNamespaceFields
	// ClusterNSScopesField instructs the search helper to look at the `ClusterNS Scopes`
	// field. Use this for objects that may be associated with multiple clusters or namespaces, such as
	// images.
	ClusterNSScopesField
)

// SearchHelper facilitates applying scoped access control to search operations.
type SearchHelper interface {
	Apply(searchFunc func(*v1.Query) ([]search.Result, error)) func(context.Context, *v1.Query) ([]search.Result, error)
	FilteredSearcher(searcher search.UnsafeSearcher) search.Searcher
}

// searchResultsChecker is responsible for checking whether a single search result is allowed to be seen.
type searchResultsChecker interface {
	TryAllowed(resourceSC ScopeChecker, result search.Result) TryAllowedResult
	SearchFieldLabels() []search.FieldLabel
	PostProcess(result *search.Result)
}

type searchHelper struct {
	resource permissions.Resource

	resultsChecker searchResultsChecker
}

// NewSearchHelper returns a new search helper for the given resource.
func NewSearchHelper(resourceMD permissions.ResourceMetadata, optionsMap search.OptionsMap, flavor SearchHelperFlavor) (SearchHelper, error) {
	// Check that resource scope is consistent with search helper flavor being used.
	switch resourceMD.GetScope() {
	case permissions.GlobalScope:
		return nil, errors.New("search helper cannot be used with globally-scoped resources")
	case permissions.ClusterScope:
		if flavor != ClusterIDField {
			return nil, errors.Errorf("cluster-scoped resource %v need to be used with flavor %v, not %v", resourceMD, ClusterIDField, flavor)
		}
	case permissions.NamespaceScope:
		if flavor == ClusterIDField {
			return nil, errors.Errorf("namespace-scoped resource %v must not be used with flavor %v", resourceMD, flavor)
		}
	default:
		return nil, errors.Errorf("unknown resource scope %v", resourceMD.GetScope())
	}

	optMap := optionsMap.Original()

	var resultsChecker searchResultsChecker
	var err error

	switch flavor {
	case ClusterNSScopesField:
		resultsChecker, err = newClusterNSScopesBasedResultsChecker(optMap)
	case ClusterIDField, ClusterIDAndNamespaceFields:
		resultsChecker, err = newClusterNSFieldBaseResultsChecker(optMap, flavor == ClusterIDAndNamespaceFields)
	default:
		err = errors.Errorf("unknown search helper flavor %v", flavor)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "creating search helper for resource %v", resourceMD)
	}

	return &searchHelper{
		resource:       resourceMD.GetResource(),
		resultsChecker: resultsChecker,
	}, nil
}

// Apply takes in a context-less search function, and returns a search function taking in a context and applying
// scoped access control checks for result filtering.
func (h *searchHelper) Apply(rawSearchFunc func(*v1.Query) ([]search.Result, error)) func(context.Context, *v1.Query) ([]search.Result, error) {
	return func(ctx context.Context, q *v1.Query) ([]search.Result, error) {
		return h.execute(ctx, q, search.WrapUnsafeSearchFunc(rawSearchFunc))
	}
}

// FilteredSearcher takes in an unsafe searcher and makes it safe.
func (h *searchHelper) FilteredSearcher(searcher search.UnsafeSearcher) search.Searcher {
	return search.WrapSearchFunc(func(ctx context.Context, q *v1.Query) ([]search.Result, error) {
		return h.execute(ctx, q, searcher)
	})
}

func (h *searchHelper) execute(ctx context.Context, q *v1.Query, searcher search.UnsafeSearcher) ([]search.Result, error) {
	scopeChecker := GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_ACCESS).Resource(h.resource)
	if ok, err := scopeChecker.Allowed(ctx); err != nil {
		return nil, err
	} else if ok {
		return searcher.Search(q)
	}

	// Make sure the cluster and perhaps namespace fields are part of the returned fields.
	fieldQB := search.NewQueryBuilder()
	for _, fieldLabel := range h.resultsChecker.SearchFieldLabels() {
		fieldQB = fieldQB.AddStringsHighlighted(fieldLabel, search.WildcardString)
	}
	queryWithFields := search.NewConjunctionQuery(q, fieldQB.ProtoQuery())

	results, err := searcher.Search(queryWithFields)
	if err != nil {
		return nil, err
	}

	return h.filterResults(ctx, scopeChecker, results)
}

func (h *searchHelper) filterResultsOnce(resourceScopeChecker ScopeChecker, results []search.Result) (allowed []search.Result, maybe []search.Result) {
	for _, result := range results {
		if res := h.resultsChecker.TryAllowed(resourceScopeChecker, result); res == Allow {
			allowed = append(allowed, result)
		} else if res == Unknown {
			maybe = append(maybe, result)
		}
	}
	return
}

func (h *searchHelper) filterResults(ctx context.Context, resourceScopeChecker ScopeChecker, results []search.Result) ([]search.Result, error) {
	allowed, maybe := h.filterResultsOnce(resourceScopeChecker, results)
	if len(maybe) > 0 {
		if err := resourceScopeChecker.PerformChecks(ctx); err != nil {
			return nil, err
		}
		extraAllowed, maybe := h.filterResultsOnce(resourceScopeChecker, maybe)
		if len(maybe) > 0 {
			errorhelpers.PanicOnDevelopmentf("still %d maybe results after PerformChecks", len(maybe))
		}
		allowed = append(allowed, extraAllowed...)
	}
	for i := range allowed {
		h.resultsChecker.PostProcess(&allowed[i])
	}

	return allowed, nil
}

// searchHelper implementations

// clusterNSScopesBasedResultsChecker inspects the `ClusterNS Scopes` map field of a result object, to
// determine whether the object is used in ANY scope that the principal performing the search is allowed
// to see.
type clusterNSScopesBasedResultsChecker struct {
	clusterNSScopesValuePath string
}

func newClusterNSScopesBasedResultsChecker(opts map[search.FieldLabel]*v1.SearchField) (searchResultsChecker, error) {
	clusterNSScopesField := opts[search.ClusterNSScopes]
	if clusterNSScopesField == nil {
		return nil, errors.Errorf("field %v not found", search.ClusterNSScopes)
	}
	if !clusterNSScopesField.GetStore() {
		return nil, errors.Errorf("field %s is not stored, which is a requirement for access scope enforcement", clusterNSScopesField.GetFieldPath())
	}

	return &clusterNSScopesBasedResultsChecker{
		clusterNSScopesValuePath: blevesearch.ToMapValuePath(clusterNSScopesField.GetFieldPath()),
	}, nil
}

func (c *clusterNSScopesBasedResultsChecker) TryAllowed(resourceSC ScopeChecker, result search.Result) TryAllowedResult {
	clusterNSScopeStrs := blevesearch.GetValuesFromFields(c.clusterNSScopesValuePath, result.Fields)

	tryAllowedRes := Deny
	for _, clusterNSScopeStr := range clusterNSScopeStrs {
		if clusterNSScopeStr == "" {
			continue
		}

		scopeKey := ParseClusterNSScopeString(clusterNSScopeStr)
		switch resourceSC.TryAllowed(scopeKey...) {
		case Allow:
			return Allow
		case Unknown:
			tryAllowedRes = Unknown
		}
	}
	return tryAllowedRes
}

func (c *clusterNSScopesBasedResultsChecker) SearchFieldLabels() []search.FieldLabel {
	return []search.FieldLabel{search.ClusterNSScopes}
}

func (c *clusterNSScopesBasedResultsChecker) PostProcess(result *search.Result) {
	// Make sure this doesn't get leaked via search results.
	delete(result.Fields, c.clusterNSScopesValuePath)
}

// clusterNSFieldBasedResultsChecker inspects the `Cluster ID` and optionally the `Namespace`
// field of search results, to determine whether the principal performing the search is allowed
// to see an object.
type clusterNSFieldBasedResultsChecker struct {
	clusterIDFieldPath string
	namespaceFieldPath string
}

func newClusterNSFieldBaseResultsChecker(opts map[search.FieldLabel]*v1.SearchField, namespaceScoped bool) (searchResultsChecker, error) {
	clusterIDField := opts[search.ClusterID]
	if clusterIDField == nil {
		return nil, errors.Errorf("field %v not found", search.ClusterID)
	}
	if !clusterIDField.GetStore() {
		return nil, errors.Errorf("field %s is not stored, which is a requirement for access scope enforcement", clusterIDField.GetFieldPath())
	}

	var nsField *v1.SearchField
	if namespaceScoped {
		nsField = opts[search.Namespace]
		if nsField == nil {
			return nil, errors.Errorf("field %v not found", search.Namespace)
		}
		if !nsField.GetStore() {
			return nil, errors.Errorf("field %s is not stored, which is a requirement for access scope enforcement", nsField.GetFieldPath())
		}
	}

	return &clusterNSFieldBasedResultsChecker{
		clusterIDFieldPath: clusterIDField.GetFieldPath(),
		namespaceFieldPath: nsField.GetFieldPath(),
	}, nil
}

func (c *clusterNSFieldBasedResultsChecker) TryAllowed(resourceSC ScopeChecker, result search.Result) TryAllowedResult {
	key := make([]ScopeKey, 0, 2)
	clusterID, _ := result.Fields[c.clusterIDFieldPath].(string)
	key = append(key, ClusterScopeKey(clusterID))
	if c.namespaceFieldPath != "" {
		namespace, _ := result.Fields[c.namespaceFieldPath].(string)
		key = append(key, NamespaceScopeKey(namespace))
	}
	return resourceSC.TryAllowed(key...)
}

func (c *clusterNSFieldBasedResultsChecker) SearchFieldLabels() []search.FieldLabel {
	fieldLabels := make([]search.FieldLabel, 0, 2)
	fieldLabels = append(fieldLabels, search.ClusterID)
	if c.namespaceFieldPath != "" {
		fieldLabels = append(fieldLabels, search.Namespace)
	}
	return fieldLabels
}

func (c *clusterNSFieldBasedResultsChecker) PostProcess(result *search.Result) {

}
