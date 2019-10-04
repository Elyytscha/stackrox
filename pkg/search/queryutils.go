package search

import (
	"github.com/gogo/protobuf/proto"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/errorhelpers"
)

// ApplyFnToAllBaseQueries walks recursively over the query, applying fn to all the base queries.
func ApplyFnToAllBaseQueries(q *v1.Query, fn func(*v1.BaseQuery)) {
	if q.GetQuery() == nil {
		return
	}

	switch typedQ := q.GetQuery().(type) {
	case *v1.Query_Disjunction:
		for _, subQ := range typedQ.Disjunction.GetQueries() {
			ApplyFnToAllBaseQueries(subQ, fn)
		}
	case *v1.Query_Conjunction:
		for _, subQ := range typedQ.Conjunction.GetQueries() {
			ApplyFnToAllBaseQueries(subQ, fn)
		}
	case *v1.Query_BaseQuery:
		fn(typedQ.BaseQuery)
	default:
		log.Errorf("Unhandled query type: %T; query was %s", q, proto.MarshalTextString(q))
	}
}

// ValidateQuery validates that the input query only contains fields in the given option map.
func ValidateQuery(q *v1.Query, optionsMap OptionsMap) error {
	errList := errorhelpers.NewErrorList("query for contains unrecognized fields:")
	ApplyFnToAllBaseQueries(q, func(bq *v1.BaseQuery) {
		matchFieldQuery, ok := bq.GetQuery().(*v1.BaseQuery_MatchFieldQuery)
		if !ok {
			return
		}
		if _, isValid := optionsMap.Get(matchFieldQuery.MatchFieldQuery.GetField()); !isValid {
			errList.AddString(matchFieldQuery.MatchFieldQuery.GetField())
		}
	})
	return errList.ToError()
}

// FilterQuery applies the given function on every base query, and returns a new
// query that has only the sub-queries that the function returns true for.
// It will NOT mutate q unless the function passed mutates its argument.
func FilterQuery(q *v1.Query, fn func(*v1.BaseQuery) bool) (*v1.Query, bool) {
	if q.GetQuery() == nil {
		return nil, false
	}
	switch typedQ := q.GetQuery().(type) {
	case *v1.Query_Disjunction:
		filteredQueries := filterQueriesByFunction(typedQ.Disjunction.GetQueries(), fn)
		if len(filteredQueries) == 0 {
			return nil, false
		}
		return DisjunctionQuery(filteredQueries...), true
	case *v1.Query_Conjunction:
		filteredQueries := filterQueriesByFunction(typedQ.Conjunction.GetQueries(), fn)
		if len(filteredQueries) == 0 {
			return nil, false
		}
		return ConjunctionQuery(filteredQueries...), true
	case *v1.Query_BaseQuery:
		if fn(typedQ.BaseQuery) {
			return q, true
		}
		return nil, false
	default:
		log.Errorf("Unhandled query type: %T; query was %s", q, proto.MarshalTextString(q))
		return nil, false
	}
}

// Helper function used by FilterQuery.
func filterQueriesByFunction(qs []*v1.Query, fn func(*v1.BaseQuery) bool) (filteredQueries []*v1.Query) {
	for _, q := range qs {
		filteredQuery, found := FilterQuery(q, fn)
		if found {
			filteredQueries = append(filteredQueries, filteredQuery)
		}
	}
	return
}
